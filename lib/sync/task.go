package sync

import (
	"errors"
	"math"
	"pmsGo/lib/cache"
	"pmsGo/lib/config"
	"pmsGo/lib/log"
	"pmsGo/lib/security/random"
	"runtime"
	"sync"
	"time"
)

const (
	TaskQueueKey = "task:queue"
)

var (
	poolNum    int
	taskPool   sync.Pool
	taskChan   chan *Task
	taskWorker map[string]Worker
)

func init() {
	poolNum = runtime.NumCPU()
	if config.Config.Sync.Processor > 0 {
		poolNum = int(math.Min(float64(poolNum), float64(config.Config.Sync.Processor)))
	}
	taskPool = sync.Pool{New: func() interface{} {
		return &Task{}
	}}
	taskChan = make(chan *Task, poolNum)
	taskWorker = make(map[string]Worker)
}
func Run() {
	InitTaskReceiver(poolNum)
}

func RegisterWorker(key string, processor Worker) error {
	if taskWorker[key] != nil {
		return errors.New("worker already existed")
	}
	taskWorker[key] = processor
	return nil
}

func UnregisterWorker(key string) {
	delete(taskWorker, key)
}

func DispatchWorker(key string) Worker {
	return taskWorker[key]
}

type Task struct {
	Key    string
	Param  interface{}
	UUID   string
	Worker Worker
}

func NewTask(key string, param interface{}) (*Task, error) {
	uuid := random.Uuid(false)
	task := &Task{}
	task.UUID = uuid
	task.Key = key
	task.Param = param
	err := AddTask(task)
	return task, err
}

func AddTask(task *Task) error {
	err := cache.Push(TaskQueueKey, task)
	if err != nil {
		return err
	}
	return nil
}

func GetTask() (*Task, error) {
	value, err := cache.Pop(TaskQueueKey)
	if err != nil {
		return nil, err
	}
	if value != nil {
		decode := value.(map[string]interface{})
		task := &Task{}
		task.UUID = decode["UUID"].(string)
		task.Key = decode["Key"].(string)
		task.Param = decode["Param"]
		return task, nil
	}
	return nil, errors.New("task data is empty")
}

func taskReceiver() {
	for {
		task := <-taskChan
		result, err := task.Worker.Process(task.UUID, task.Param)
		if err != nil {
			task.Worker.Fallback(task.UUID, task.Param, err)
		} else {
			task.Worker.Callback(task.UUID, task.Param, result)
		}
		taskPool.Put(task)
	}
}

func taskDispatcher() {
	for {
		taskData, err := GetTask()
		if err != nil {
			time.Sleep(time.Second * 2)
			continue
		}
		for {
			taskWorker := taskPool.Get()
			if taskWorker != nil {
				taskProcess := taskWorker.(*Task)
				taskProcess.UUID = taskData.UUID
				taskProcess.Param = taskData.Param
				taskProcess.Worker = DispatchWorker(taskData.Key)
				taskChan <- taskProcess
				break
			} else {
				time.Sleep(time.Second * 2)
			}
		}
	}
}

func InitTaskReceiver(num int) {
	log.Debugf("Start sync workers: %v \n", num)
	for i := 0; i < num; i++ {
		go taskReceiver()
	}
	log.Debug("Start sync dispatcher \n")
	go taskDispatcher()
}
