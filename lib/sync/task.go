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
	taskPool      sync.Pool
	taskChan      chan *Task
	taskProcessor map[string]Processor
)

func init() {
	processor := runtime.NumCPU()
	if config.Config.Sync.Processor > 0 {
		processor = int(math.Min(float64(processor), float64(config.Config.Sync.Processor)))
	}
	taskPool = sync.Pool{New: func() interface{} {
		return &Task{}
	}}
	taskChan = make(chan *Task, processor)
	taskProcessor = make(map[string]Processor)
	InitTaskReceiver(processor)
}

type Processor func(param interface{})

func RegisterProcessor(key string, processor Processor) {
	taskProcessor[key] = processor
}
func UnregisterProcessor(key string) {
	delete(taskProcessor, key)
}

func DispatchProcessor(key string) Processor {
	return taskProcessor[key]
}

type Task struct {
	Key       string
	Param     interface{}
	UUID      string
	processor Processor
}

func NewTask(key string, param interface{}) error {
	uuid := random.Uuid(false)
	task := &Task{}
	task.UUID = uuid
	task.Key = key
	task.Param = param
	return AddTask(task)
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
		task.processor(task.Param)
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
				taskProcess.processor = DispatchProcessor(taskData.Key)
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
