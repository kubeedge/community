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
	kerrors "errors"
	glog "log"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	cloudrobotkubeedgev1beta1 "github.com/ospp2023/api/v1beta1"
)

// TaskReconciler reconciles a Task object
type TaskReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

const (
	expireTime = 60 * time.Second
)

//+kubebuilder:rbac:groups=cloudrobot.kubeedge.cloudrobot.kubeedge,resources=tasks,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=cloudrobot.kubeedge.cloudrobot.kubeedge,resources=tasks/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=cloudrobot.kubeedge.cloudrobot.kubeedge,resources=tasks/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Task object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.4/pkg/reconcile
func (r *TaskReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	glog.Println("Start Task Reconcile")

	// TODO(user): your logic here
	// GET TASK
	// task := cloudrobotkubeedgev1beta1.Task{}
	var task cloudrobotkubeedgev1beta1.Task
	if err := r.Get(ctx, req.NamespacedName, &task); err != nil {
		if errors.IsNotFound(err) {
			// resource not found, it might be deleted
			return ctrl.Result{}, nil
		}
		// Error fetching Task resource
		logger.Error(err, "Failed to fetch Task")
		return ctrl.Result{}, err
	}

	glog.Println("Start Task Reconcile")

	// DELETION
	finalizerName := "cloudrobot.kubeedge.io/finalizer"
	// check DeletionTimestamp to determine if object is under deletion
	if task.ObjectMeta.DeletionTimestamp.IsZero() {
		glog.Println("Task is not being deleted")
		if !controllerutil.ContainsFinalizer(&task, finalizerName) {
			glog.Println("Adding finalizer for task")
			controllerutil.AddFinalizer(&task, finalizerName)
			if err := r.Update(ctx, &task); err != nil {
				return ctrl.Result{}, err
			}
			glog.Println("Finish Task Reconcile")
		}
	} else {
		// task is being deleted
		glog.Println("Task is being deleted")
		if controllerutil.ContainsFinalizer(&task, finalizerName) {
			if err := r.handleTaskDelete(ctx, &task); err != nil {
				return ctrl.Result{}, err
			}
			controllerutil.RemoveFinalizer(&task, finalizerName)
			if err := r.Update(ctx, &task); err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	// ADD && UPDATE
	var allocatedRobotID uint
	var err error
	// Get robotSync
	robotSyncList := &cloudrobotkubeedgev1beta1.RobotSyncList{}
	if err := r.List(ctx, robotSyncList); err != nil {
		glog.Println("Failed to list the robotSync", err)
		return ctrl.Result{}, err
	}
	robotSync := &robotSyncList.Items[0]
	// Get registed robot list
	if task.Spec.Allocated < 0 {
		glog.Println("start matching robot for a new task")
		if allocatedRobotID, err = r.matchTask(ctx, robotSync, &task); err != nil {
			glog.Println("Failed to match task to robot", err)
			return ctrl.Result{}, nil
		}
		glog.Println("newly allocate to robotID: ", allocatedRobotID)
		task.Spec.Allocated = int(allocatedRobotID)
		r.Update(ctx, &task)
	} else {
		glog.Println("task is already allocated to: ", task.Spec.Allocated)
		allocatedRobotID = uint(task.Spec.Allocated)
	}
	robot, err := r.getRobotByID(ctx, allocatedRobotID)
	if err != nil {
		glog.Println("Failed to get robot by ID", err)
		return ctrl.Result{}, nil
	}
	r.updateRobot(ctx, robot, &task) // update taskinfo and underTask

	glog.Println("Finish Task Reconcile")

	return ctrl.Result{}, nil
}

func (r *TaskReconciler) handleTaskDelete(ctx context.Context, task *cloudrobotkubeedgev1beta1.Task) error {
	if task.Spec.Allocated <= 0 {
		glog.Println("not alloc yet, delete directly")
		r.Delete(ctx, task)
		return nil
	}
	robot, err := r.getRobotByID(ctx, uint(task.Spec.Allocated))
	if err != nil {
		glog.Println("Failed to get robot")
		return err
	}
	// 2. update robot Undertask
	robot.Status.UnderTask = false
	r.Status().Update(ctx, robot)

	glog.Println("clean up, deleting task")

	// 3. delete
	if err := r.Delete(ctx, task); err != nil {
		glog.Println("Failed to delete task", err)
		return err
	}
	

	return nil
}

func (r *TaskReconciler) getRobotByID(ctx context.Context, robotId uint) (*cloudrobotkubeedgev1beta1.Robot, error) {
	robotList := &cloudrobotkubeedgev1beta1.RobotList{}
	if err := r.List(ctx, robotList); err != nil {
		glog.Println("Failed to list the robot", err)
		return nil, err
	}
	robot := cloudrobotkubeedgev1beta1.Robot{}
	for _, robot = range robotList.Items {
		if robot.Spec.RobotID == robotId {
			err := r.Get(ctx, types.NamespacedName{Namespace: NamespacedName, Name: robot.ObjectMeta.Name}, &robot)
			if err != nil {
				glog.Println("Failed to match robotID in task controller")
				return nil, err
			}
		}
	}

	return &robot, nil
}

func (r *TaskReconciler) updateRobot(ctx context.Context, robot *cloudrobotkubeedgev1beta1.Robot, task *cloudrobotkubeedgev1beta1.Task) error {
	// TODO: update taskinfo and underTask of robot
	robot.Status.UnderTask = true
	robot.Status.TaskInfo.OrderID = task.Spec.OrderID
	robot.Status.TaskInfo.TaskID = task.Spec.TaskID
	robot.Status.TaskInfo.PointStateSequence = task.Spec.PointStateSequence
	robot.Status.TaskInfo.SegmentStateSequence = task.Spec.SegmentStateSequence
	robot.Status.TaskInfo.RequiredSensors = task.Spec.RequiredSensors
	r.Status().Update(ctx, robot)

	return nil
}

func (r *TaskReconciler) matchTask(ctx context.Context, robotSync *cloudrobotkubeedgev1beta1.RobotSync, task *cloudrobotkubeedgev1beta1.Task) (robotID uint, err error) {
	// TODO: return a uint if match successfully, otherwise err
	registedList := robotSync.Status.RegistedRobots
	for _, robotID := range registedList {
		glog.Println("checking robotId: ", robotID)
		registedRobot, err := r.getRobotByID(ctx, robotID)
		if err != nil {
			glog.Println("Failed to get robot by ID", err)
			return 0, err
		}
		// check if the node is alive
		lastHeartbeat := robotSync.Status.LastHeartbeat[registedRobot.ObjectMeta.Name]
		glog.Println("checking lastheartbeat")
		if time.Since(lastHeartbeat.Time) > expireTime {
			glog.Println("robot not alive")
			if err := r.removeFromRegistedList(ctx, registedRobot); err != nil {
				glog.Println("Failed to remove robot from registedList", err)
				return 0, err
			}
			continue
		}
		glog.Println("checking lastheartbeat successfully")
		glog.Println("checking robot underTask")
		if registedRobot.Status.UnderTask {
			continue
		}
		glog.Println("robot is not underTask")
		glog.Println("checking sensors")
		// if task.Spec.RequiredSensors in robot.Spec.Sensors, match task to robot
		requiredSensors := task.Spec.RequiredSensors
		sensors := registedRobot.Status.Sensors
		if isSubset(requiredSensors, sensors) {
			glog.Println("checking sensor successfully returning")
			return registedRobot.Spec.RobotID, nil
		}
	}
	return 0, kerrors.New("no robot match the task")
}

func (r *TaskReconciler) removeFromRegistedList(ctx context.Context, robot *cloudrobotkubeedgev1beta1.Robot) error {
	glog.Println("removing from registedList")
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

	index := -1
	for i, e := range robotSync.Status.RegistedRobots {
		if e == robot.Spec.RobotID {
			index = i
			break
		}
	}
	if index != -1 {
		robotSync.Status.RegistedRobots = append(robotSync.Status.RegistedRobots[:index], robotSync.Status.RegistedRobots[index+1:]...)
	}

	r.Status().Update(ctx, &robotSync)

	glog.Println("update robotSync status in remove from registed list successfully")

	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *TaskReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cloudrobotkubeedgev1beta1.Task{}).
		Complete(r)
}
