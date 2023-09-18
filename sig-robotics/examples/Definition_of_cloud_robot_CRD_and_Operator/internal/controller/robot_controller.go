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
	"encoding/json"
	glog "log"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	cloudrobotkubeedgev1beta1 "github.com/ospp2023/api/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// RobotReconciler reconciles a Robot object
type RobotReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

type Message struct {
	MsgType        string // Indicates whether it is a heartbeat message or feedback on task completion
	Urgent         bool   // Indicates whether it is an exception message
	Time           time.Time
	Sensors        []cloudrobotkubeedgev1beta1.Sensor
	BatteryStatus  cloudrobotkubeedgev1beta1.BatteryStatus
	ResourceStatus cloudrobotkubeedgev1beta1.ResourceStatus
	Position       cloudrobotkubeedgev1beta1.Position
	RunningStatus  cloudrobotkubeedgev1beta1.RunningStatus
	AbnormalEvents []cloudrobotkubeedgev1beta1.AbnormalEvents
}

const (
	NamespacedName = "default"
)

// Special task information for low battery handling
var ChargingStation = cloudrobotkubeedgev1beta1.PointStateSequence{
	Position_x: float32Ptr(0.0),
	Position_y: float32Ptr(2.0),
}
var ChargingTask = cloudrobotkubeedgev1beta1.TaskInfo{
	// reserved OrderID and TaskID
	OrderID:            uintPtr(114514),
	TaskID:             uintPtr(114514),
	PointStateSequence: []cloudrobotkubeedgev1beta1.PointStateSequence{ChargingStation},
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var HttpRegisted = false

//+kubebuilder:rbac:groups=cloudrobot.kubeedge.cloudrobot.kubeedge,resources=robots,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=cloudrobot.kubeedge.cloudrobot.kubeedge,resources=robots/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=cloudrobot.kubeedge.cloudrobot.kubeedge,resources=robots/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Robot object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.15.0/pkg/reconcile
func (r *RobotReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	glog.Println("Start Robot Reconcile")

	var err error
	var robotNotFound bool = false

	// TODO(user): your logic here
	// GET ROBOT
	robot := &cloudrobotkubeedgev1beta1.Robot{}
	if err = r.Get(ctx, req.NamespacedName, robot); err != nil && !errors.IsNotFound(err) {
		// if errors.IsNotFound(err) {
		// 	// resource not found, it might be deleted
		// 	return ctrl.Result{}, nil
		// }
		// Error fetching Robot resource
		logger.Error(err, "Failed to fetch Robot")
		return ctrl.Result{}, err
	}
	if errors.IsNotFound(err) {
		if !robotNotFound {
			robotNotFound = true
		}
	}

	robotId := robot.Spec.RobotID

	// GET NODE
	node, err := r.getNodeByRobotId(ctx, robotId)
	if err != nil {
		glog.Println("Failed to fetch Node by robotId")
		return ctrl.Result{}, err
	}

	// DELETION
	finalizerName := "cloudrobot.kubeedge.io/finalizer"
	// check DeletionTimestamp to determine if object is under deletion
	if !robotNotFound && robot.ObjectMeta.DeletionTimestamp.IsZero() {
		glog.Println("Robot is not being deleted")
		if !controllerutil.ContainsFinalizer(robot, finalizerName) {
			glog.Println("Adding finalizer for robot")
			controllerutil.AddFinalizer(robot, finalizerName)
			if err := r.Update(ctx, robot); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else if !robotNotFound && !robot.ObjectMeta.DeletionTimestamp.IsZero() {
		// robot is being deleted
		glog.Println("Robot is being deleted")
		if controllerutil.ContainsFinalizer(robot, finalizerName) {
			if err := r.handleRobotDelete(ctx, robot, node); err != nil {
				return ctrl.Result{}, err
			}
			controllerutil.RemoveFinalizer(robot, finalizerName)
			if err := r.Update(ctx, robot); err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	// ADD && UPDATE
	if robotNotFound {
		glog.Println("Robot not found, skip")
		glog.Println("Finish Robot Reconcile 1")
		return ctrl.Result{}, nil
	}

	// node trigger: Listen to 11451 port to receive message from edge node
	server := &http.Server{
		Addr:        ":11451",
		Handler:     http.DefaultServeMux, // Set default handler
		IdleTimeout: 5 * time.Second,      // Set idle timeout to 5 seconds
	}
	if !HttpRegisted {
		http.HandleFunc("/path", r.handleWebSocket(ctx, robotId, robotNotFound))
		HttpRegisted = true
	}

	glog.Println("WebSocket server started on : 11451")

	// ListenAndServe will block until the server is shut down
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		glog.Println("HTTP server ListenAndServe error: ", err)
	} else if err == http.ErrServerClosed {
		server.Shutdown(context.Background())
		glog.Println("HTTP server ListenAndServe closed")
	}

	err = server.Shutdown(context.Background())
	if err != nil {
		if err == http.ErrServerClosed {
			// Server has been shut down successfully
			glog.Println("Server has been shut down")
		} else {
			// Other errors
			glog.Printf("Error during server shutdown: %v\n", err)
		}
	}

	glog.Println("Port released. Finish Robot Reconcile 2")

	return ctrl.Result{}, nil
}

func (r *RobotReconciler) handleWebSocket(ctx context.Context, robotId uint, robotNotFound bool) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		glog.Println("robotId: ", robotId)
		robot, err := r.getRobotByID(ctx, robotId)
		if err != nil {
			glog.Println("Failed to get robot by id")
			return
		}
		conn, err := upgrader.Upgrade(w, req, nil)
		if err != nil {
			glog.Println("Failed to upgrade the connection: ", err)
			return
		}
		// defer conn.Close()

		if robotNotFound {
			glog.Println("Robot not found")
			conn.Close()
			return
		}

		for {
			messageType, receivedMessage, err := conn.ReadMessage()
			glog.Println("messageType", messageType)
			if err != nil {
				// Handle read error
				glog.Println("Failed to read message or connection closed: ", err)

				return
			}
			// Handle closing messages
			if messageType == websocket.CloseMessage {
				glog.Println("Close message received")
				if err := closeWebSocketConnection(conn); err != nil {
					glog.Println("Failed to close websocket connection: ", err)
				}
			}
			// Deserialization
			var message Message
			err = json.Unmarshal(receivedMessage, &message)
			if err != nil {
				// Handle unmarshal error
				glog.Println("Failed to unmarshal message", err)
				return
			}
			glog.Println("unmarshal message successfully!")

			if message.MsgType == "heartbeat" {
				glog.Println("into: ", message.MsgType)
				// Process the heartbeat message
				// Update robot information && update robotsync's lastheartbeat
				if err := r.updateBasicInfo(ctx, robot, message, robotNotFound); err != nil {
					glog.Println("Failed to update basic info", err)
					return
				}
				if message.Urgent {
					glog.Println("occur abnormal event!")
					if !r.handleException(ctx, robot, message.AbnormalEvents) {
						// not been solved, then remove from registedList
						err := r.removeFromRegistedList(ctx, robot)
						if err != nil {
							glog.Println("Failed to remove from registedList", err)
							return
						}
						return
					}
				}
				glog.Println("finish process heartbeat message")

				// Exceptions have been solved or no exception
				// update the status and send command
				err = r.updateStatus(ctx, conn, robot, message)
				if err != nil {
					glog.Println("Failed to update status", err)
					return
				}
			} else if message.MsgType == "task-feedback" {
				glog.Println("into: ", message.MsgType)
				// Only update task and robot information. Robot changes will trigger a new round of robot controller
				if *robot.Status.TaskInfo.OrderID != 114514 {
					// You can add conditions to specifically judge the message returned.
					glog.Println("searching task with")
					glog.Print("OrderId: ", *robot.Status.TaskInfo.OrderID)
					glog.Println("TaskId: ", *robot.Status.TaskInfo.TaskID)
					task, err := r.getTaskByID(ctx, *robot.Status.TaskInfo.OrderID, *robot.Status.TaskInfo.TaskID)
					if err != nil {
						glog.Println("Failed to get task by Id", err)
						return
					}
					glog.Println("task: ", task)
					pointSequence := task.Spec.PointStateSequence
					glog.Println("length of pointSequence: ", len(pointSequence))
					// No other tasks point. The task can be deleted and the robot information can be modified at the same time.
					if len(pointSequence) == 1 {
						// Triggering task graceful deletion will trigger handleTaskDelete in the task controller to modify the robot task status.
						glog.Println("deleting task in handlewebsocket")
						r.Delete(ctx, task)
					} else if len(pointSequence) >= 1 {
						// Remove completed mission points
						task.Spec.PointStateSequence = task.Spec.PointStateSequence[1:]
						r.Update(ctx, task)
					} else {
						glog.Println("invalid status of task")
					}
				} else {
					robot.Status.UnderTask = false
					r.Status().Update(ctx, robot)
				}
			}
		}
	}
}

func (r *RobotReconciler) handleRobotDelete(ctx context.Context, robot *cloudrobotkubeedgev1beta1.Robot, node *corev1.Node) error {
	// 1. Untag the node
	if err := r.Get(ctx, types.NamespacedName{Namespace: "", Name: node.ObjectMeta.Name}, node); err != nil {
		return err
	}
	delete(node.Labels, RobotNodeTagName)
	delete(node.Labels, "robotId")
	glog.Println("UNTAG the labels and annotations")
	if err := r.Update(ctx, node); err != nil {
		glog.Println("Failed to untag node "+node.ObjectMeta.Name, err)
		return err
	}

	// 2. release the robotID
	robotID := robot.Spec.RobotID
	SharedIDPool.Release(robotID)
	glog.Println("release the robotID")

	// 3. remove robot from RegistedList
	err := r.removeFromRegistedList(ctx, robot)
	if err != nil {
		glog.Println("Failed to remove from registedList", err)
		return err
	}
	glog.Println("remove robot from RegistedList")

	// 4. delete
	if err := r.Delete(ctx, robot); err != nil {
		glog.Println("Failed to delete robot", err)
		return err
	}

	return nil
}

func (r *RobotReconciler) updateBasicInfo(ctx context.Context, robot *cloudrobotkubeedgev1beta1.Robot, message Message, robotNotFound bool) error {
	glog.Println("updating robot infos and lastheartbeat in robotsync ...")
	robot.Status.Sensors = message.Sensors
	robot.Status.BatteryStatus = message.BatteryStatus
	robot.Status.Position = message.Position
	robot.Status.ResourceStatus = message.ResourceStatus
	robot.Status.RunningStatus = message.RunningStatus
	robot.Status.AbnormalEvents = message.AbnormalEvents
	if err := r.Status().Update(ctx, robot); err != nil {
		glog.Println("Failed to update basic info of robot", err)
		return err
	}

	// update the heartbeat info
	robotSyncList := &cloudrobotkubeedgev1beta1.RobotSyncList{}
	if err := r.List(ctx, robotSyncList); err != nil {
		glog.Println("Failed to list the robotSync", err)
		return err
	}
	robotSync := robotSyncList.Items[0]
	if robotSync.Status.LastHeartbeat == nil {
		robotSync.Status.LastHeartbeat = map[string]metav1.Time{}
	}

	glog.Println("updateing lastheartbeat in robotsync, received lastheatbeat: ", message.Time)
	robotSync.Status.LastHeartbeat[robot.ObjectMeta.Name] = metav1.NewTime(message.Time)

	err := r.Status().Update(ctx, &robotSync)
	if err != nil {
		glog.Println("Failed to update heartbeat info", err)
		return err
	}

	return nil
}

func (r *RobotReconciler) handleException(ctx context.Context, robot *cloudrobotkubeedgev1beta1.Robot, abnormalEvents []cloudrobotkubeedgev1beta1.AbnormalEvents) bool {
	glog.Println("length of abnormalEvents", len(abnormalEvents))
	for _, event := range abnormalEvents {
		glog.Println("ExceptionLevel: ", event.ExceptionLevel)
		glog.Println("EventCode: ", event.EventCode)
		if event.ExceptionLevel == 0x00 { // 0x00: info
			glog.Println("abnormal event info: ", abnormalEvents[0].Description)
		} else if event.ExceptionLevel == 0x01 { // 0x01: warnin
			switch event.EventCode {
			case 0x5013:
				// Low battery abnormality
				robot.Status.TaskInfo = ChargingTask
				robot.Status.UnderTask = true
				glog.Println("target_x of robot: ", *robot.Status.TaskInfo.PointStateSequence[0].Position_x)
				glog.Println("target_y of robot:", *robot.Status.TaskInfo.PointStateSequence[0].Position_y)
				err := r.Status().Update(ctx, robot)
				if err != nil {
					glog.Println("Failed to update robot status", err)
					return false
				}
			default:
				glog.Println("invalid event code")
			}
		} else if event.ExceptionLevel == 0x02 { // 0x02: error
			switch event.EventCode {
			case 0x01:

			default:
				glog.Println("invalid event code")
			}
		}
	}
	return true
}

func (r *RobotReconciler) removeFromRegistedList(ctx context.Context, robot *cloudrobotkubeedgev1beta1.Robot) error {
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

	for i, registedRobotId := range robotSync.Status.RegistedRobots {
		if registedRobotId == robot.Spec.RobotID {
			// remove from registedList
			robotSync.Status.RegistedRobots = append(robotSync.Status.RegistedRobots[:i], robotSync.Status.RegistedRobots[i+1:]...)
			break
		}
	}

	r.Status().Update(ctx, &robotSync)

	return nil
}

func (r *RobotReconciler) updateStatus(ctx context.Context, conn *websocket.Conn, robot *cloudrobotkubeedgev1beta1.Robot, message Message) error {
	// send commands and update robot Task status
	// 1. compare the sensors
	if !reflect.DeepEqual(robot.Status.Sensors, message.Sensors) {
		if r.isRelatedToTask(robot) {
			err := r.setTaskUnassigned(ctx, robot)
			if err != nil {
				glog.Println("Failed to reset the task", err)
				return err
			}
			err = r.updateRobotTaskStatus(ctx, robot)
			if err != nil {
				glog.Println("Failed to reset the robot", err)
				return err
			}
			return nil
		}
	}

	// 2. send up-to-date task information to robot node
	var taskInfoJSON []byte
	var err error
	if robot.Status.UnderTask {
		glog.Println("sending up-to-date task information to robot node")
		taskInfoJSON, err = json.Marshal(robot.Status.TaskInfo)
		if err != nil {
			glog.Println("Failed to marshal the taskInfo", err)
			return err
		}

		// send command(taskInfo) to robot node
		if err := conn.WriteMessage(websocket.TextMessage, taskInfoJSON); err != nil {
			glog.Println("Failed to send message", err)
			return err
		}
		glog.Println("message sent successfully!")
	} else {
		// When there is no message to send, notify the agent to close the connection.
		closeMessage := json.RawMessage(`{"close": true}`)
		// send message to remind agent to close the connection
		if err := conn.WriteMessage(websocket.TextMessage, closeMessage); err != nil {
			glog.Println("Failed to send message", err)
			return err
		}
	}
	return nil
}

func (r *RobotReconciler) updateRobotTaskStatus(ctx context.Context, robot *cloudrobotkubeedgev1beta1.Robot) error {
	robot.Status.UnderTask = false
	r.Update(ctx, robot)
	return nil
}

func (r *RobotReconciler) isRelatedToTask(robot *cloudrobotkubeedgev1beta1.Robot) bool {
	requiredSensors := robot.Status.TaskInfo.RequiredSensors
	for _, sensor := range requiredSensors {
		// if sensor is in robot.Spec.Sensors
		if !r.sensorInRobot(robot, sensor) {
			return false
		}
	}
	return false
}

func (r *RobotReconciler) sensorInRobot(robot *cloudrobotkubeedgev1beta1.Robot, sensor cloudrobotkubeedgev1beta1.Sensor) bool {
	for _, s := range robot.Status.Sensors {
		if reflect.DeepEqual(s, sensor) {
			return true
		}
	}
	return false
}

func (r *RobotReconciler) setTaskUnassigned(ctx context.Context, robot *cloudrobotkubeedgev1beta1.Robot) error {
	orderID := robot.Status.TaskInfo.OrderID
	taskID := robot.Status.TaskInfo.TaskID
	// get Task CR List
	taskList := &cloudrobotkubeedgev1beta1.TaskList{}
	if err := r.List(ctx, taskList); err != nil {
		glog.Println("Failed to list the task", err)
		return err
	}

	for _, task := range taskList.Items {
		if task.Spec.OrderID == orderID && task.Spec.TaskID == taskID {
			task.Spec.Allocated = -1
			break
		}
	}
	return nil
}

func (r *RobotReconciler) getNodeByRobotId(ctx context.Context, robotId uint) (*corev1.Node, error) {
	nodeList := &corev1.NodeList{}
	if err := r.List(ctx, nodeList); err != nil {
		glog.Println("Failed to list the node", err)
		return nil, err
	}

	for _, node := range nodeList.Items {
		if node.Labels["robotId"] == strconv.Itoa(int(robotId)) {
			return &node, nil
		}
	}
	return nil, nil

}

func (r *RobotReconciler) getRobotByID(ctx context.Context, robotId uint) (*cloudrobotkubeedgev1beta1.Robot, error) {
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
	glog.Println("Failed to match robotID in robot controller")
	return nil, nil
}

func (r *RobotReconciler) getTaskByID(ctx context.Context, orderId uint, taskId uint) (*cloudrobotkubeedgev1beta1.Task, error) {
	taskList := &cloudrobotkubeedgev1beta1.TaskList{}
	if err := r.List(ctx, taskList); err != nil {
		glog.Println("Failed to list the task", err)
		return nil, err
	}
	task := cloudrobotkubeedgev1beta1.Task{}
	for _, task = range taskList.Items {
		if *task.Spec.OrderID == orderId && *task.Spec.TaskID == taskId {
			if err := r.Get(ctx, types.NamespacedName{Namespace: NamespacedName, Name: task.ObjectMeta.Name}, &task); err == nil {
				return &task, nil
			}
		}
	}
	glog.Println("Failed to match robotID in robot controller")
	return nil, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *RobotReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cloudrobotkubeedgev1beta1.Robot{}).
		Complete(r)
}
