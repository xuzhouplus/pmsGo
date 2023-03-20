package worker

import (
	"context"
	"pmsGo/lib/cache"
	"pmsGo/lib/sync"
	"time"
)

const TaskProcessCachePrefix = "task_process:"

func init() {
	sync.RegisterWorker(CarouselWorkerName, &CarouselWorker{})
	sync.RegisterWorker(ImageWorkerName, &ImageWorker{})
	sync.RegisterWorker(VideoWorkerName, &VideoWorker{})
}

func SetTaskProcessStatus(taskId string, stepName string, status map[string]interface{}) {
	if taskId == "" {
		return
	}
	cache.Redis.HSet(context.TODO(), cache.Key(TaskProcessCachePrefix+taskId), stepName, status)
}

func ClearTaskProcessStatus(taskId string) {
	if taskId == "" {
		return
	}
	cache.Redis.Del(context.TODO(), cache.Key(TaskProcessCachePrefix+taskId))
}

func SetStatusExpire(taskId string) {
	if taskId == "" {
		return
	}
	cache.Redis.Expire(context.TODO(), cache.Key(TaskProcessCachePrefix+taskId), time.Hour)
}

func TaskSteps(steps interface{}) []string {
	taskStep := make([]string, 0)
	if steps == nil {
		return taskStep
	}
	switch steps.(type) {
	case []interface{}:
		actionSteps := steps.([]interface{})
		for _, step := range actionSteps {
			actionStep := step.(string)
			taskStep = append(taskStep, actionStep)
		}
		return taskStep
	case []string:
		return steps.([]string)
	default:
		panic("steps类型错误")
	}
}
