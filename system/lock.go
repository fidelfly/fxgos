package system

import (
	"sync"
	"strings"
	"time"
)

type SharedLock struct {
	Module string
}

type DataLock struct {
	Backup string
}

type ResourceLock struct {
	sync.Mutex
	tryLock sync.Mutex
	locked bool
	action *LockAction
}

func (rl *ResourceLock) TryLock(action *LockAction) (bool, *LockAction)  {
	if rl.locked {
		return  false, rl.action
	}
	rl.tryLock.Lock()
	defer rl.tryLock.Unlock()
	if rl.locked {
		return false, rl.action
	}
	rl.Mutex.Lock()
	rl.locked = true
	rl.action = action
	return true, action
}

func (rl *ResourceLock) Unlock() {
	rl.Mutex.Unlock()
	rl.locked = false
	rl.action = nil
}

func (rl *ResourceLock) IsLocked() bool {
	return rl.locked
}

type LockAction struct {
	Code string
	UserId int64
	Message string
}

type SystemLock struct {
	module *ResourceLock
}

func (sl *SystemLock) init() {
	sl.module = &ResourceLock{}
}

func (sl *SystemLock) GetResourceLock(resourceCode string, key... interface{}) *ResourceLock{
	switch resourceCode {
	case SharedLocker.Module:
		return sl.module
	default:
		return nil
	}
}

type systemResourceManager struct {
	initLock      sync.RWMutex
	systemLockers map[string]*SystemLock
	resourceLockers *MemCache
}

func newLockerManager() *systemResourceManager {
	lm := systemResourceManager{}
	lm.systemLockers = make(map[string]*SystemLock)
	lm.resourceLockers = CreateEnsureCache(120*time.Minute, 10*time.Minute, lockResolver)
	return &lm
}

func lockResolver(key string) interface{} {
	return &ResourceLock{}
}

func (lm *systemResourceManager) GetSystemLock(code string) *SystemLock {
	lm.initLock.RLock()
	lock := lm.systemLockers[code]
	if lock != nil {
		lm.initLock.RUnlock()
		return lock
	}
	lm.initLock.RUnlock()

	lm.initLock.Lock()
	lock = lm.systemLockers[code]
	if lock == nil {
		lock = &SystemLock{}
		lock.init()
		lm.systemLockers[code] = lock
	}
	lm.initLock.Unlock()
	return  lock
}

func (lm *systemResourceManager) resolveResourceKey(resourceCode string, key... string) string {
	resourceKey := strings.Builder{}
	resourceKey.WriteString(resourceCode)
	for _, keyValue := range key {
		resourceKey.WriteString("#")
		resourceKey.WriteString(keyValue)
	}
	return 	resourceKey.String()
}

func (lm *systemResourceManager) GetResourdeLock(resourceCode string, key... string) *ResourceLock {
	resourceKey := lm.resolveResourceKey(resourceCode, key...)
	lock, _ := lm.resourceLockers.Get(resourceKey)
	if lock != nil {
		return lock.(*ResourceLock)
	}
	return nil
}


