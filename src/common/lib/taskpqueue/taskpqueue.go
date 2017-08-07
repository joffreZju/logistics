package taskpqueue

import (
	"fmt"
	"sync"

	"github.com/Jeffail/tunny"
	lane "gopkg.in/oleiade/lane.v1"
)

type TaskQueue struct {
	QueueSize        int
	PQueue           *lane.PQueue
	PoolSize         int
	Pool             *tunny.WorkPool
	BaseTaskData     map[string]interface{}
	BaseTaskDataLock sync.RWMutex
	JobTimeOut       int
}

//NewTaskQueue(30, 10, 1000*60*3)
func NewTaskQueue(queuesize int, poolsize int, jobtimeout int, job func(interface{}) interface{}) (taskpqueue TaskQueue) {
	pool, _ := tunny.CreatePool(poolsize, job).Open()
	taskpqueue = TaskQueue{
		QueueSize:    queuesize,
		PQueue:       lane.NewPQueue(lane.MINPQ),
		PoolSize:     poolsize,
		Pool:         pool,
		BaseTaskData: map[string]interface{}{},
		JobTimeOut:   jobtimeout,
	}
	return
}

func (tpq *TaskQueue) Start() {
	go func() {
		for {
			tpq.BaseTaskDataLock.RLock()
			taskkey, pri := tpq.PQueue.Pop()
			if taskkey == nil && pri == 0 {
				tpq.BaseTaskDataLock.RUnlock()
				continue
			}
			fmt.Println("key:", taskkey, pri)
			key := taskkey.(string)
			if data, ok := tpq.BaseTaskData[key]; ok {
				//timeduration, _ := time.ParseDuration(fmt.Sprintf("%vms", tpq.JobTimeOut))
				//			fmt.Println("type:", reflect.TypeOf(data))
				go tpq.Pool.SendWork(data)
				fmt.Println("end :", data, tpq.Pool.NumPendingAsyncJobs())
			}
			tpq.BaseTaskDataLock.RUnlock()
		}
	}()
}

func (tpq *TaskQueue) AddTask(jobdata interface{}, jobkey string, joblevel int) (err error) {
	tpq.BaseTaskDataLock.Lock()
	defer tpq.BaseTaskDataLock.Unlock()
	tpq.BaseTaskData[jobkey] = jobdata
	if tpq.PQueue.Size() >= tpq.QueueSize {
		return fmt.Errorf("size is fill")
	}
	tpq.PQueue.Push(jobkey, joblevel)
	return
}
