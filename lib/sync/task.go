package sync

import (
	"pmsGo/lib/config"
	"pmsGo/lib/security/random"
	"runtime"
	"sync"
)

const (
	StateWaiting   = "waiting"
	StateCompleted = "completed"
	StateError     = "failed"
	StateNone      = "none"
	StateOverdue   = "overdue"
)

var (
	taskStateMutex sync.Mutex
	taskPool       sync.Pool
	taskChan       chan *Task
	taskState      map[string]string
)

func init() {
	processor := runtime.NumCPU()
	if config.Config.Sync.Processor > 0 {
		processor = config.Config.Sync.Processor
	}
	taskChan = make(chan *Task, processor)
	taskState = make(map[string]string)
	InitTaskReceiver(processor)
}

type Task struct {
	Param   interface{}
	Handler Handler
	UUID    string
}

type Handler func(string, interface{}) (string, error)

func NewTask(param interface{}, handler Handler) *Task {
	uuid := random.Uuid(false)
	task := taskPool.Get()
	if task == nil {
		return &Task{
			param,
			handler,
			uuid,
		}
	} else {
		task := task.(*Task)
		task.Param = param
		task.Handler = handler
		task.UUID = uuid
		return task
	}
}

func AddTask(task *Task) string {
	go func() {
		taskChan <- task
	}()
	uuid := task.UUID
	UpdateTaskState(uuid, StateWaiting)
	return uuid
}

func UpdateTaskState(uuid, state string) {
	taskStateMutex.Lock()
	defer taskStateMutex.Unlock()

	taskState[uuid] = state
}

func GetTaskState(uuid string) (state string) {
	taskStateMutex.Lock()
	defer taskStateMutex.Unlock()

	resultState, exists := taskState[uuid]
	if !exists {
		state = StateNone
	} else {
		state = resultState
	}
	return
}

func taskReceiver() {
	var taskUUID string
	var err error
	for {
		task := <-taskChan
		taskUUID, err = task.Handler(task.UUID, task.Param)
		if err != nil {
			UpdateTaskState(taskUUID, StateError)
			taskPool.Put(task)
		} else {
			UpdateTaskState(taskUUID, StateCompleted)
			taskPool.Put(task)
		}
	}
}

func InitTaskReceiver(num int) {
	for i := 0; i < num; i++ {
		go taskReceiver()
	}
}
