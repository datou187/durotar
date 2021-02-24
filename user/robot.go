package user

import (
	"container/list"
	"sync"
)

//机器人容器，并发安全
//TODO：目前直接一次性将所有机器人加载到内存，在机器人数量过多(>10000)时，考虑进行分段加载，以减少内存占用
type RobotQueue struct {
	robots *list.List
	sync.Mutex
}

func NewRobotQueue(robots []Account) *RobotQueue {
	pool := &RobotQueue{
		robots: list.New(),
	}
	for i := 0; i < len(robots); i++ {
		pool.Enqueue(&robots[i])
	}
	return pool
}

func (rc *RobotQueue) Dequeue() *Account {
	rc.Lock()
	defer rc.Unlock()
	elem := rc.robots.Front()
	if elem == nil {
		return nil
	}
	inst := elem.Value.(*Account)
	rc.robots.Remove(elem)
	return inst
}

func (rc *RobotQueue) Enqueue(rbs ...*Account) {
	rc.Lock()
	defer rc.Unlock()
	for _, rb := range rbs {
		rc.robots.PushBack(rb)
	}
}

func (rc *RobotQueue) Size() int {
	rc.Lock()
	defer rc.Unlock()
	return rc.robots.Len()
}
