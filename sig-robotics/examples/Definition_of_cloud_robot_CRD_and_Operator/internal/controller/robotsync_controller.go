/*
Copyright 2022 The KubeEdge Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	glog "log"
	"strconv"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	cloudrobotkubeedgev1beta1 "github.com/ospp2023/api/v1beta1"
)

// RobotSyncReconciler reconciles a RobotSync object
type RobotSyncReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=cloudrobot.kubeedge.cloudrobot.kubeedge,resources=robotsyncs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=cloudrobot.kubeedge.cloudrobot.kubeedge,resources=robotsyncs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=cloudrobot.kubeedge.cloudrobot.kubeedge,resources=robotsyncs/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the RobotSync object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.4/pkg/reconcile
func (r *RobotSyncReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	glog.Println("Start RobotSync Reconcile")

	// TODO(user): your logic here
	// GET RobotSync
	robotSyncList := &cloudrobotkubeedgev1beta1.RobotSyncList{}
	if err := r.List(ctx, robotSyncList); err != nil {
		glog.Println("Failed to list the robotSync", err)
		return ctrl.Result{}, err
	}
	robotSync := &robotSyncList.Items[0]

	// GET NODE
	node := &corev1.Node{}
	if err := r.Get(ctx, req.NamespacedName, node); err != nil {
		if errors.IsNotFound(err) {
			// resource not found, it might be deleted
			return ctrl.Result{}, nil
		}
		// Error fetching Node resource
		logger.Error(err, "Failed to GET node")
		return ctrl.Result{}, err
	}

	// ADD & UPDATE
	annotations := node.GetAnnotations()
	nodeName := node.GetName()

	// Determine whether there is this robotId in the registered list to prevent repeated creation of robots
	var robotId uint
	var ok = false
	if robotId_str := node.Labels["robotId"]; robotId_str != "" {
		// If it has been registered, get the robotId.
		robotId_int, err := strconv.Atoi(robotId_str)
		if err != nil {
			glog.Println("Failed to convert robotId", err)
			return ctrl.Result{}, err
		} else {
			robotId = uint(robotId_int)
			ok = true // Mark robot is registered
		}
	}
	// Regardless of whether the node that triggers Reconcile is registered or not, the update heartbeat is not processed here.
	if !ok || !ContainsElement(robotSync.Status.RegistedRobots, robotId) {
		// Not registered
		// Create Robot resource
		robot, err := r.applyRobot(ctx, annotations)
		if err != nil {
			glog.Println(err, "Failed to create Robot resource")
			return ctrl.Result{}, err
		}
		// Add to Registed Robot list + add robotId into node labels
		err = r.updateAllRobotSync(ctx, node, robotSync, robot, nodeName, annotations)
		if err != nil {
			glog.Println("Failed to update RobotSync resource", err)
			return ctrl.Result{}, err
		}
	}

	glog.Println("finish RobotSync Reconcile")

	return ctrl.Result{}, nil
}

func (r *RobotSyncReconciler) handleNodeDelete(ctx context.Context, node *corev1.Node) error {
	// Delete Robot CR && remove from Registed Robot list
	robotID, err := strconv.Atoi(node.GetLabels()["robotId"])
	if err != nil {
		glog.Println("Failed to convert robotId", err)
		return err
	}

	robot, err := r.getRobotByID(ctx, uint(robotID))
	if err != nil {
		glog.Println("Failed to get robot", err)
		return err
	}

	r.removeFromRegistedList(ctx, robot)

	if err := r.Delete(ctx, robot); err != nil {
		glog.Println("Failed to delete robot", err)
		return err
	}

	return nil
}

func (r *RobotSyncReconciler) removeFromRegistedList(ctx context.Context, robot *cloudrobotkubeedgev1beta1.Robot) error {
	robotSyncList := &cloudrobotkubeedgev1beta1.RobotSyncList{}
	// get robotSync CR
	if err := r.List(ctx, robotSyncList); err != nil {
		glog.Println("Failed to list the robotSync", err)
		return err
	}
	// There will only be one custom resource object of type robotSync globally.
	robotSync := robotSyncList.Items[0]
	robotName := robot.ObjectMeta.Name
	delete(robotSync.Status.LastHeartbeat, robotName)

	r.Update(ctx, &robotSync)

	return nil
}

func (r *RobotSyncReconciler) getRobotByID(ctx context.Context, robotId uint) (*cloudrobotkubeedgev1beta1.Robot, error) {
	robotList := &cloudrobotkubeedgev1beta1.RobotList{}
	if err := r.List(ctx, robotList); err != nil {
		glog.Println("Failed to list the robot", err)
		return nil, err
	}
	robot := cloudrobotkubeedgev1beta1.Robot{}
	for _, robot = range robotList.Items {
		if robot.Spec.RobotID == robotId {
			if err := r.Get(ctx, types.NamespacedName{Namespace: NamespacedName, Name: robot.ObjectMeta.Name}, &robot); err == nil {
				return &robot, nil
			}
		}
	}
	glog.Println("Failed to match robotID in robotsync controller")
	return nil, nil
}

func (r *RobotSyncReconciler) createOrPatch(ctx context.Context, obj client.Object) error {
	// TODO(user): your logic here
	glog.Println("Applying", obj.GetName())
	if err := r.Create(ctx, obj); err != nil && r.isAlreadyExistError(err) {
		glog.Println(obj.GetName(), "Already Exist, Patching")
		if err := r.Patch(ctx, obj, client.Merge); err != nil {
			glog.Println(err, "Failed to PATCH ", obj.GetName())
		}
	}
	return nil
}

func (r *RobotSyncReconciler) isAlreadyExistError(err error) bool {
	return strings.Contains(err.Error(), "already exists") || strings.Contains(err.Error(), "port is already allocated")
}

func (r *RobotSyncReconciler) applyRobot(ctx context.Context, annotations map[string]string) (*cloudrobotkubeedgev1beta1.Robot, error) {
	// TODO(user): call getBase()  && r.createOrPatch()
	var robot *cloudrobotkubeedgev1beta1.Robot
	var err error
	robot, err = getBase()
	if err != nil {
		glog.Println("Failed to get base", err)
		return nil, err
	}

	r.createOrPatch(ctx, robot)
	return robot, nil
}

func (r *RobotSyncReconciler) updateAllRobotSync(ctx context.Context, node *corev1.Node, robotSync *cloudrobotkubeedgev1beta1.RobotSync, robot *cloudrobotkubeedgev1beta1.Robot, nodeName string, annotations map[string]string) error {
	// TODO: update RegistedList, and add robotId into node labels

	if robotSync.Status.RegistedRobots == nil {
		robotSync.Status.RegistedRobots = []uint{}
	}
	// check if robot.Spec.RobotID is in robotSync.Status.RegistedRobots
	if !ContainsElement(robotSync.Status.RegistedRobots, robot.Spec.RobotID) {
		robotSync.Status.RegistedRobots = append(robotSync.Status.RegistedRobots, robot.Spec.RobotID)
	}

	err := r.Status().Update(ctx, robotSync)
	if err != nil {
		glog.Println("update robotSync failed", err)
		return err
	}

	node.Labels["robotId"] = strconv.FormatUint(uint64(robot.Spec.RobotID), 10)
	err = r.Update(ctx, node)
	if err != nil {
		glog.Println("update node failed", err)
		return err
	}

	glog.Println("Register successfilly")

	return nil
}

func (r *RobotSyncReconciler) updateHeartbeatInfo(ctx context.Context, robotSync *cloudrobotkubeedgev1beta1.RobotSync, nodeName string, annotations map[string]string) error {
	// TODO: only update LastHeatbeat
	glog.Println("timestamp in annotations: ", annotations["lastHeartbeat"])
	lastHeartbeat, err := time.Parse("2006-01-02 15:04:05", annotations["lastHeartbeat"])
	if err != nil {
		glog.Println("Failed to parse time", err)
		return err
	}
	robotSync.Status.LastHeartbeat[nodeName] = metav1.NewTime(lastHeartbeat)

	r.Status().Update(ctx, robotSync)

	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *RobotSyncReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		// For(&cloudrobotkubeedgev1beta1.RobotSync{}).
		For(&corev1.Node{}).
		WithEventFilter(nodeLabelFilter()).
		Complete(r)
}
