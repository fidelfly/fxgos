package api

import (
	"net/http"
	"sync"
	"time"

	"github.com/fidelfly/gox"
	"github.com/fidelfly/gox/cachex/mcache"
	"github.com/fidelfly/gox/httprxr"
	"github.com/fidelfly/gox/logx"
	"github.com/fidelfly/gox/pkg/randx"

	"github.com/fidelfly/fxgos/cmd/pkg/db"
	"github.com/fidelfly/fxgos/cmd/service/audit/res"
)

//Task Menu
const (
	TaskMenuRole  = "task.menu.role"
	TaskMenuUsers = "task.menu.users"
)

//Task Operation
const (
	TaskOperationCreate  = "task.operation.create"
	TaskOperationModify  = "task.operation.modify"
	TaskOperationDelete  = "task.operation.delete"
	TaskOperationDisable = "task.operation.disable"
	TaskOperationEnable  = "task.operation.enable"
)

var keyMux = sync.Mutex{}
var taskCache = mcache.NewCache(mcache.DefaultExpiration, 5*time.Minute)

type TaskStatus int

const (
	TaskStatusOk     TaskStatus = 1
	TaskStatusFailed TaskStatus = -1
)

type Task struct {
	*res.Systrail
	done bool
}

func registerTask(code string, operation string, r *http.Request, w http.ResponseWriter) (task *Task, status bool) {

	if len(code) == 0 {
		code = r.RequestURI
	}
	task, err := _registerTask(code, operation)
	if err != nil {
		if w != nil {
			httprxr.ResponseJSON(w, http.StatusInternalServerError, httprxr.ExceptionMessage(err))
		}
		status = false
		return
	}

	if r != nil {
		userInfo := GetUserInfo(r)
		if userInfo != nil {
			task.ExecUser = userInfo.Id
		}
		requestId := GetRequestId(r)
		if len(requestId) > 0 {
			task.RequestId = requestId
		}
	}
	status = true
	return
}

func GetAttachTask(r *http.Request) *Task {
	if r != nil {
		keyObj := httprxr.ContextGet(r, "TaskKey")
		if keyObj != nil {
			if key, ok := keyObj.(string); ok {
				return GetTask(key)
			}
		}
	}
	return nil
}

func GetTask(key string) *Task {
	if cacheObj, find := taskCache.Get(key); !find {
		return nil
	} else if task, ok := cacheObj.(*Task); ok {
		return task
	}
	return nil
}

func _registerTask(code string, operation string) (task *Task, err error) {
	keyMux.Lock()
	key := randx.GetUUID(code)
	task = &Task{
		Systrail: &res.Systrail{
			Key:       key,
			Code:      code,
			Operation: operation,
		},
	}
	taskCache.Set(key, task)
	keyMux.Unlock()
	return
}

func (task *Task) StartTrail(args ...int64) {
	if len(args) > 0 {
		task.ExecUser = args[0]
	}
	task.StartTime = time.Now()
}

func (task *Task) TrailDone(status TaskStatus, message string) {
	task.EndTime = time.Now()
	task.Duration = int64(task.EndTime.Sub(task.StartTime) / time.Millisecond)
	task.Status = int64(status)
	task.PutStatusDescription(message)
}

func (task *Task) Done(args ...interface{}) {
	task.done = true
	task.Destroy()
}

func (task *Task) IsDone() bool {
	return task.done
}

func (task *Task) LogTrailDone(arg interface{}) {
	if arg == nil {
		task.TrailDone(TaskStatusOk, "")
	} else {
		if err, ok := arg.(error); ok {
			task.TrailDone(TaskStatusFailed, err.Error())
		} else if msg, ok := arg.(string); ok {
			task.TrailDone(TaskStatusOk, msg)
		} else {
			task.TrailDone(TaskStatusOk, "")
		}
	}
	task.LogTask()
}

func (task *Task) LogTask() {
	go func() {
		defer gox.CapturePanicAndRecover("Log Task")
		_, err := db.Create(task.Systrail)
		if err != nil {
			logx.Error(err)
		}
	}()
}

func (task *Task) Destroy() {
	taskCache.Remove(task.Key)
}

func (task *Task) PutStatusDescription(desc string) {
	if len(desc) > 0 {
		if len(task.StatusDescription) > 0 {
			task.StatusDescription = task.StatusDescription + "\n" + desc
		} else {
			task.StatusDescription = desc
		}
	}
}

func (task *Task) SetField(field string, value string) {
	if task.Info == nil {
		task.Info = make(map[string]string)
	}
	task.Info[field] = value
}

func (task *Task) GetField(field string) string {
	if task.Info == nil {
		return ""
	}
	if value, found := task.Info[field]; found {
		return value
	}
	return ""
}

func (task *Task) IsFieldExist(field string) bool {
	_, found := task.Info[field]
	return found
}
