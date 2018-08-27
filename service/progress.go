package service

import (
	"github.com/lyismydg/fxgos/system"
	"github.com/lyismydg/fxgos/websocket"

	"time"

	"encoding/json"

	"sync"

	"github.com/sirupsen/logrus"
)

const (
	PROGRESS_ACTIVE    = "active"
	PROGRESS_EXCEPTION = "exception"
	PROGRESS_SUCCESS   = "success"
)

type ProgressGetter interface {
	GetPercent() int
	GetStatus() string
	GetMessage() interface{}
}

type ProgressSetter interface {
	ProgressGetter
	Set(percent int, status string, message ...interface{})
	update(percent int, status string, message ...interface{})
}

type ProgressSubscriber interface {
	ProgressSet(percent int, status string, messages ...interface{})
}

type ProgressSuperior interface {
	ProgressChanged(subProgress *SubProgress)
}

type SubProgress struct {
	superior    ProgressSuperior
	Proporition int
	Code        string
	Message     interface{}
	Percent     int
	Status      string
	Propagation bool
	done        bool
}

func (sp *SubProgress) ProgressSet(percent int, status string, message ...interface{}) {
	sp.Percent = percent
	sp.Status = status
	if len(message) > 0 {
		sp.Message = message[0]
	}

	sp.superior.ProgressChanged(sp)

	sp.done = percent >= 100
}

func (sp *SubProgress) IsDone() bool {
	return sp.done
}

func NewProgressDispatcher(code string, subscriber ...ProgressSubscriber) *ProgressDispatcher {
	return &ProgressDispatcher{Code: code, Subscribers: subscriber}
}

// Progress Dispatcher
type ProgressDispatcher struct {
	Code        string
	Subscribers []ProgressSubscriber
	message     interface{}
	percent     int
	status      string
	auto        *AutoProgress
	sub         []*SubProgress
	mux         sync.Mutex
	data        map[string]interface{}
	notifyLock  sync.Mutex
}

func (pd *ProgressDispatcher) GetPercent() int {
	return pd.percent
}
func (pd *ProgressDispatcher) GetStatus() string {
	return pd.status
}
func (pd *ProgressDispatcher) GetMessage() interface{} {
	return pd.message
}

func (pd *ProgressDispatcher) SetStatus(status string, message ...interface{}) {
	pd.Set(pd.percent, status, message...)
}

func (pd *ProgressDispatcher) Exception(percent int, message ...interface{}) {
	pd.Set(percent, PROGRESS_EXCEPTION, message...)
}

func (pd *ProgressDispatcher) Active(percent int, message ...interface{}) {
	pd.Set(percent, PROGRESS_ACTIVE, message...)
}

func (pd *ProgressDispatcher) Success(message ...interface{}) {
	if len(message) == 0 {
		pd.Set(100, PROGRESS_SUCCESS, "")
	} else {
		pd.Set(100, PROGRESS_SUCCESS, message...)
	}
}

func (pd *ProgressDispatcher) Done(status string, message ...interface{}) {
	if len(status) == 0 {
		if pd.status == PROGRESS_ACTIVE {
			status = PROGRESS_SUCCESS
		} else {
			status = pd.status
		}
	}
	if len(message) == 0 {
		pd.Set(100, status, "")
	} else {
		pd.Set(100, status, message...)
	}
}

func (pd *ProgressDispatcher) notifySubscriber() {
	pd.notify(pd.percent, pd.status, pd.message)
}

func (pd *ProgressDispatcher) updateData(percent int, status string, message interface{}) bool {
	dataChange := false
	if pd.data == nil {
		pd.data = make(map[string]interface{})
	}

	if pd.data["percent"] != percent {
		dataChange = true
		pd.data["percent"] = percent
	}

	if pd.data["status"] != status {
		pd.data["status"] = status
	}

	if pd.data["message"] != message {
		pd.data["message"] = message
	}

	return dataChange
}
func (pd *ProgressDispatcher) notify(percent int, status string, message interface{}) {
	pd.notifyLock.Lock()
	defer pd.notifyLock.Unlock()
	if percent < 0 {
		percent = pd.percent
	}
	if len(status) == 0 {
		status = pd.status
	}
	if message == nil {
		message = pd.message
	}

	if !pd.updateData(percent, status, message) {
		return
	}
	if len(pd.Subscribers) > 0 {
		for _, subscriber := range pd.Subscribers {
			subscriber.ProgressSet(percent, status, message)
		}
	} else {
		msg := ""
		if message != nil {
			if msgText, ok := message.(string); ok {
				msg = msgText
			} else {
				if msgData, err := json.Marshal(message); err == nil {
					msg = string(msgData)
				}
			}
		}
		logrus.Infof("Progress(%s) : percent = %d%%, status = %s, message = %s", pd.Code, percent, status, msg)
	}
}
func (pd *ProgressDispatcher) update(percent int, status string, message ...interface{}) {
	pd.percent = percent
	pd.status = status
	if len(message) > 0 {
		pd.message = message[0]
	}
	pd.notifySubscriber()
}
func (pd *ProgressDispatcher) Set(percent int, status string, message ...interface{}) {
	if pd.auto != nil {
		pd.auto.Stop()
		pd.auto = nil
	}
	pd.update(percent, status, message...)
}

func (pd *ProgressDispatcher) AutoProgress(stepValue int, duration time.Duration, maxValue int, message ...interface{}) {
	if len(message) > 0 {
		pd.message = message[0]
		pd.notifySubscriber()
	}

	if pd.auto != nil {
		pd.auto.Stop()
		pd.auto = nil
	}
	pd.auto = newAutoProgress(pd, stepValue, duration, maxValue)
	pd.auto.Start()
}

func (pd *ProgressDispatcher) Step(stepValue int, message ...interface{}) {
	pd.Set(pd.percent+stepValue, PROGRESS_ACTIVE, message...)
}

func (pd *ProgressDispatcher) NewSubProgress(proporition int) *SubProgress {
	pd.mux.Lock()
	defer pd.mux.Unlock()
	sp := &SubProgress{superior: pd, Proporition: proporition}
	pd.sub = append(pd.sub, sp)
	return sp
}

func (pd *ProgressDispatcher) ProgressChanged(subProgress *SubProgress) {
	pd.mux.Lock()
	defer pd.mux.Unlock()
	if pd.sub != nil && len(pd.sub) > 0 {
		subValue := int(0)
		index := -1
		for i, sp := range pd.sub {
			if sp == subProgress && sp.IsDone() {
				index = i
				pd.percent += sp.Percent
			} else {
				value := 0
				if sp.IsDone() {
					value = sp.Proporition
				} else {
					value = sp.Proporition * sp.Percent / 100
				}
				subValue += value
			}
		}

		if index >= 0 {
			newSub := make([]*SubProgress, 0)
			if index > 0 {
				newSub = append(newSub, pd.sub[:index]...)
			}
			if index < len(pd.sub)-1 {
				newSub = append(newSub, pd.sub[index+1:]...)
			}
			pd.sub = newSub
		}

		percent := pd.percent + subValue
		msg := pd.message
		if subProgress.Message != nil {
			msg = subProgress.Message
		}

		if subProgress.Propagation && subProgress.Status == PROGRESS_EXCEPTION {
			pd.status = subProgress.Status
		}
		pd.notify(percent, pd.status, msg)
	}

}

//Struct AutoProgress
type AutoProgress struct {
	progress  ProgressSetter
	stepValue int
	maxValue  int
	duration  time.Duration
	ticker    *time.Ticker
}

func (ap *AutoProgress) Start() {
	ap.ticker = time.NewTicker(ap.duration)
	go func() {
		defer ap.ticker.Stop()
		for _ = range ap.ticker.C {
			percent := ap.progress.GetPercent() + ap.stepValue
			if percent > ap.maxValue {
				percent = ap.maxValue
			}
			ap.progress.update(percent, ap.progress.GetStatus(), ap.progress.GetMessage())
			if percent >= ap.maxValue {
				return
			}
		}
	}()
}

func (ap *AutoProgress) Stop() {
	if ap.ticker != nil {
		ap.ticker.Stop()
	}
	if ap.progress.GetPercent() < ap.maxValue {
		ap.progress.update(ap.maxValue, ap.progress.GetStatus(), ap.progress.GetMessage())
	}
}

func newAutoProgress(progress ProgressSetter, stepValue int, duration time.Duration, maxValue int) *AutoProgress {
	return &AutoProgress{progress: progress, stepValue: stepValue, duration: duration, maxValue: maxValue}
}

//Core Struct : WsProgress
type WsProgress struct {
	ws           *websocket.WsConnect `json:"_"`
	Code         string               `json:"code"`
	Message      interface{}          `json:"message"`
	Percent      int                  `json:"percent"`
	Status       string               `json:"status"`
	auto         *AutoProgress
	sub          []*SubProgress
	mux          sync.Mutex
	data         map[string]interface{}
	delayed      bool
	delayMessage bool
	senderLock   sync.Mutex
}

func GetProgress(key string, code string) *WsProgress {
	if conn, ok := system.SocketCache.Get(key); ok {
		return &WsProgress{ws: conn.(*websocket.WsConnect), Code: code}
	}
	return &WsProgress{Code: code}
}

func (wsp *WsProgress) GetPercent() int {
	return wsp.Percent
}

func (wsp *WsProgress) GetStatus() string {
	return wsp.Status
}
func (wsp *WsProgress) GetMessage() interface{} {
	return wsp.Message
}

func (wsp *WsProgress) SetStatus(status string, message ...interface{}) {
	wsp.Set(wsp.Percent, status, message...)
}

func (wsp *WsProgress) Exception(percent int, message ...interface{}) {
	wsp.Set(percent, PROGRESS_EXCEPTION, message...)
}

func (wsp *WsProgress) Active(percent int, message ...interface{}) {
	wsp.Set(percent, PROGRESS_ACTIVE, message...)
}

func (wsp *WsProgress) Success(message ...interface{}) {
	wsp.Set(100, PROGRESS_SUCCESS, message...)
}

func (wsp *WsProgress) Done(status string, message ...interface{}) {
	wsp.Set(100, status, message...)
}

func (wsp *WsProgress) update(percent int, status string, message ...interface{}) {
	wsp.Percent = percent
	wsp.Status = status
	if len(message) > 0 {
		wsp.Message = message[0]
	}

	wsp.SendMsg()
}

func (wsp *WsProgress) Set(percent int, status string, message ...interface{}) {
	if wsp.auto != nil {
		wsp.auto.Stop()
		wsp.auto = nil
	}
	wsp.update(percent, status, message...)
}

func (wsp *WsProgress) updateData(percent int, status string, message interface{}) bool {
	dataChange := false
	if wsp.data == nil {
		wsp.data = make(map[string]interface{})
	}

	if wsp.data["percent"] != percent {
		dataChange = true
		wsp.data["percent"] = percent
	}

	if wsp.data["status"] != status {
		wsp.data["status"] = status
	}

	if wsp.data["message"] != message {
		wsp.data["message"] = message
	}

	wsp.data["code"] = wsp.Code

	return dataChange
}

func (wsp *WsProgress) delaySend(timer *time.Timer) {
	wsp.delayed = true
	go func() {
		defer timer.Stop()
		for {
			select {
			case <-timer.C:
				wsp.senderLock.Lock()
				if wsp.delayMessage {
					if wsp.ws != nil && wsp.ws.IsOpen() {
						wsp.ws.SendMessage(wsp.data)
						timer.Reset(100 * time.Millisecond)
						wsp.delayMessage = false
						wsp.senderLock.Unlock()
					} else {
						wsp.delayed = false
						wsp.delayMessage = false
						wsp.senderLock.Unlock()
						return
					}
				} else {
					wsp.delayed = false
					wsp.delayMessage = false
					wsp.senderLock.Unlock()
					return
				}
			default:
			}
		}

	}()
}

func (wsp *WsProgress) Send(percent int, status string, message interface{}) {
	wsp.senderLock.Lock()
	defer wsp.senderLock.Unlock()
	if !wsp.updateData(percent, status, message) {
		return
	}
	if wsp.ws != nil && wsp.ws.IsOpen() {
		if !wsp.delayed {
			wsp.ws.SendMessage(wsp.data)
			wsp.delaySend(time.NewTimer(100 * time.Millisecond))
		} else {
			wsp.delayMessage = true
		}
	} else {
		msg := ""
		if message != nil {
			if msgText, ok := message.(string); ok {
				msg = msgText
			} else {
				if msgData, err := json.Marshal(message); err == nil {
					msg = string(msgData)
				}
			}
		}
		logrus.Infof("Progress(%s) : percent = %d%%, status = %s, message = %s", wsp.Code, percent, status, msg)
	}
}

func (wsp *WsProgress) SendMsg() {
	wsp.Send(wsp.Percent, wsp.Status, wsp.Message)
	/*if wsp.ws != nil && wsp.ws.IsOpen() {
		wsp.ws.SendMessage(wsp)
	} else {
		msg := ""
		if wsp.Message != nil {
			if msgText, ok := wsp.Message.(string); ok {
				msg = msgText
			} else {
				if msgData, err := json.Marshal(wsp.Message); err == nil {
					msg = string(msgData)
				}
			}
		}
		logrus.Infof("Progress(%s) : percent = %d%%, status = %s, message = %s", wsp.Code, wsp.Percent, wsp.Status, msg)
	}*/
}

func (wsp *WsProgress) AutoProgress(stepValue int, duration time.Duration, maxValue int, message ...interface{}) {
	if len(message) > 0 {
		wsp.Message = message[0]
		wsp.SendMsg()
	}

	if wsp.auto != nil {
		wsp.auto.Stop()
		wsp.auto = nil
	}

	wsp.auto = newAutoProgress(wsp, stepValue, duration, maxValue)
	wsp.auto.Start()
}

func (wsp *WsProgress) Step(stepValue int, message ...interface{}) {
	wsp.Set(wsp.Percent+stepValue, PROGRESS_ACTIVE, message...)
}

func (wsp *WsProgress) NewSubProgress(proporition int) *SubProgress {
	wsp.mux.Lock()
	defer wsp.mux.Unlock()
	sp := &SubProgress{superior: wsp, Proporition: proporition}
	if wsp.sub == nil {
		wsp.sub = make([]*SubProgress, 0)
	}
	wsp.sub = append(wsp.sub, sp)
	return sp
}

func (wsp *WsProgress) ProgressChanged(subProgress *SubProgress) {
	wsp.mux.Lock()
	defer wsp.mux.Unlock()
	if wsp.sub != nil && len(wsp.sub) > 0 {
		subValue := int(0)
		index := -1
		for i, sp := range wsp.sub {
			if sp == subProgress && sp.IsDone() {
				index = i
				wsp.Percent += sp.Percent
			} else {
				value := 0
				if sp.IsDone() {
					value = sp.Proporition
				} else {
					value = sp.Proporition * sp.Percent / 100
				}
				subValue += value
			}
		}

		if index >= 0 {
			newSub := make([]*SubProgress, 0)
			if index > 0 {
				newSub = append(newSub, wsp.sub[:index]...)
			}
			if index < len(wsp.sub)-1 {
				newSub = append(newSub, wsp.sub[index+1:]...)
			}
			wsp.sub = newSub
		}

		percent := wsp.Percent + subValue
		msg := wsp.Message
		if subProgress.Message != nil {
			msg = subProgress.Message
		}

		if subProgress.Propagation && subProgress.Status == PROGRESS_EXCEPTION {
			wsp.Status = subProgress.Status
		}
		wsp.Send(percent, wsp.Status, msg)
	}

}
