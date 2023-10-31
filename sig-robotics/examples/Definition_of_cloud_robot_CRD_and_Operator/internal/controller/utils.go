package controller

import (
	"strconv"
	"sync"
	"time"

	kerrors "errors"
	glog "log"

	"github.com/gorilla/websocket"
	cloudrobotkubeedgev1beta1 "github.com/ospp2023/api/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

type IDPool struct {
	mutex sync.Mutex
	ids   []uint
}

var SharedIDPool *IDPool

const (
	RobotNodeTagName = "edge.cloudrobot.kubeedge.io/robot"
	MaxRobotNum      = 100
)

func init() {
	SharedIDPool = NewIDPool(MaxRobotNum)
}

func NewIDPool(maxID uint) *IDPool {
	ids := make([]uint, 0, maxID)
	for i := uint(1); i < maxID; i++ {
		ids = append(ids, i)
	}
	return &IDPool{ids: ids}
}

func (p *IDPool) Allocate() (uint, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if len(p.ids) == 0 {
		glog.Println("Failed to allocate id ")
		return 0, kerrors.New("no available id to allocate")
	}

	id := p.ids[0]
	p.ids = p.ids[1:]

	return id, nil
}

func (p *IDPool) Release(id uint) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	

	if !ContainsElement(p.ids, id) {
		p.ids = append(p.ids, id)
	}
}


// nodeLabelFilter filters nodes that have the desired label
func nodeLabelFilter() predicate.Predicate {
	return predicate.NewPredicateFuncs(func(object client.Object) bool {
		node, ok := object.(*corev1.Node)
		if !ok {
			return false
		}
		_, exists := node.Labels[RobotNodeTagName]
		return exists
	})
}

func getBase() (*cloudrobotkubeedgev1beta1.Robot, error) {
	// TODO: implement
	robotID, err := SharedIDPool.Allocate()
	if err != nil {
		glog.Println("Failed to allocate id ", err)
		return nil, err
	}
	resourceName := "robot" + strconv.Itoa(int(robotID))

	robot := &cloudrobotkubeedgev1beta1.Robot{
		ObjectMeta: metav1.ObjectMeta{
			Name:      resourceName,
			Namespace: NamespacedName,
			Labels: map[string]string{
				"robotId": strconv.Itoa(int(robotID)),
			},
		},
		Spec: cloudrobotkubeedgev1beta1.RobotSpec{
			RobotID: robotID,
		},
	}

	return robot, nil
}

func closeWebSocketConnection(conn *websocket.Conn) error {
	closeMessagge := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")
	var err error

	// retry max 3 times
	for i := 0; i < 3; i++ {
		err = conn.WriteControl(websocket.CloseMessage, closeMessagge, time.Now().Add(5*time.Second))
		if err == nil {
			break
		}
		glog.Println("Failed to send close message, retrying", err)
		time.Sleep(2 * time.Second)
	}

	// force close
	if err != nil {
		glog.Println("Force close websocket connection", err)
		conn.Close()
		return err
	}
	return nil
}

func isSubset(subSensors []cloudrobotkubeedgev1beta1.Sensor, sensors []cloudrobotkubeedgev1beta1.Sensor) bool {
	sensorMap := make(map[cloudrobotkubeedgev1beta1.Sensor]bool)
	for _, s := range sensors {
		sensorMap[s] = true
	}
	for _, s := range subSensors {
		if !sensorMap[s] {
			return false
		}
	}
	return true
}

func ContainsElement(arr []uint, target uint) bool {
	for _, item := range arr {
		if item == target {
			return true
		}
	}
	return false
}

func float32Ptr(f float32) *float32 {
    return &f
}

func uintPtr(i uint) *uint {
    return &i
}
