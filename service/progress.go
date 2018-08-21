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

type ProgressSetter interface {
	GetPercent() int
	GetStatus() string
	GetMessage() interface{}
	Set(percent int, status string, message ...interface{})
}

type ProgressSubscriber interface {
	ProgressSet(percent int, status string, messages ...interface{})
}

type ProgressSuperior interface {
	ProgressChanged(status string, message ...interface{})
}

type SubProgress struct {
	Superior    ProgressSuperior
	Proporition int
	Code        string
	Message     interface{}
	Percent     int
	Status      string
}

func (sp *SubProgress) ProgressSet(percent int, status string, message ...interface{}) {
	sp.Percent = percent
	sp.Status = status
	if len(message) > 0 {
		sp.Message = message[0]
	}

	sp.Superior.ProgressChanged(sp.Status, message...)
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

func (pd *ProgressDispatcher) Done(message ...interface{}) {
	pd.Set(100, PROGRESS_SUCCESS, message...)
}
func (pd *ProgressDispatcher) notifySubscriber() {
	pd.notify(pd.percent, pd.status, pd.message)
}
func (pd *ProgressDispatcher) notify(percent int, status string, message interface{}) {
	if percent < 0 {
		percent = pd.percent
	}
	if len(status) == 0 {
		status = pd.status
	}
	if message == nil {
		message = pd.message
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

func (pd *ProgressDispatcher) Set(percent int, status string, message ...interface{}) {
	if pd.auto != nil {
		pd.auto.Stop()
		pd.auto = nil
	}
	pd.percent = percent
	pd.status = status
	if len(message) > 0 {
		pd.message = message[0]
	}
	pd.notifySubscriber()
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
	sp := &SubProgress{Superior: pd, Proporition: proporition}
	pd.sub = append(pd.sub, sp)
	return sp
}

func (pd *ProgressDispatcher) ProgressChanged(status string, message ...interface{}) {
	pd.mux.Lock()
	defer pd.mux.Unlock()
	msg := pd.message
	if len(message) > 0 {
		msg = message[0]
	}
	percent := pd.percent

	if pd.sub != nil && len(pd.sub) > 0 {
		subok := int(0)
		okindex := make([]int, len(pd.sub))
		for index, sp := range pd.sub {
			value := sp.Proporition * sp.Percent / 100
			if sp.Status != PROGRESS_ACTIVE {
				subok += value
				okindex = append(okindex, index)
			}
			percent += value
		}

		if subok > 0 {
			newSub := make([]*SubProgress, 0)
			pd.percent += subok
			index := int(0)
			for _, i := range okindex {
				if i > index {
					newSub = append(newSub, pd.sub[index:i]...)

				}
				index++
			}
			if index < len(pd.sub) {
				newSub = append(newSub, pd.sub[index:]...)
			}
		}
	}

	pd.notify(percent, status, msg)
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
		for _ = range ap.ticker.C {
			percent := ap.progress.GetPercent() + ap.stepValue
			if percent > ap.maxValue {
				percent = ap.maxValue
			}
			ap.progress.Set(percent, ap.progress.GetStatus(), ap.progress.GetMessage())
			if percent >= ap.maxValue {
				break
			}
		}
	}()
}

func (ap *AutoProgress) Stop() {
	ap.ticker.Stop()
}

func newAutoProgress(progress ProgressSetter, stepValue int, duration time.Duration, maxValue int) *AutoProgress {
	return &AutoProgress{progress: progress, stepValue: stepValue, duration: duration, maxValue: maxValue}
}

//Core Struct : WsProgress
type WsProgress struct {
	ws      *websocket.WsConnect `json:"_"`
	Code    string               `json:"code"`
	Message interface{}          `json:"message"`
	Percent int                  `json:"percent"`
	Status  string               `json:"status"`
	auto    *AutoProgress
	sub     []*SubProgress
	mux     sync.Mutex
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

func (wsp *WsProgress) Done(message ...interface{}) {
	wsp.Set(100, PROGRESS_SUCCESS, message...)
}

func (wsp *WsProgress) Set(percent int, status string, message ...interface{}) {
	if wsp.auto != nil {
		wsp.auto.Stop()
		wsp.auto = nil
	}
	wsp.Percent = percent
	wsp.Status = status
	if len(message) > 0 {
		wsp.Message = message[0]
	}

	wsp.SendMsg()
}

func (wsp *WsProgress) Send(percent int, status string, message interface{}) {
	if wsp.ws != nil && wsp.ws.IsOpen() {
		wsp.ws.SendMessage(map[string]interface{}{
			"code":    wsp.Code,
			"percent": percent,
			"status":  status,
			"message": message,
		})
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
	if wsp.ws != nil && wsp.ws.IsOpen() {
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
	}
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
	sp := &SubProgress{Superior: wsp, Proporition: proporition}
	if wsp.sub == nil {
		wsp.sub = make([]*SubProgress, 0)
	}
	wsp.sub = append(wsp.sub, sp)
	return sp
}

func (wsp *WsProgress) ProgressChanged(status string, message ...interface{}) {
	wsp.mux.Lock()
	defer wsp.mux.Unlock()
	msg := wsp.Message
	if len(message) > 0 {
		msg = message[0]
	}
	percent := wsp.Percent

	if wsp.sub != nil && len(wsp.sub) > 0 {
		subok := int(0)
		okindex := make([]int, len(wsp.sub))
		for index, sp := range wsp.sub {
			value := sp.Proporition * sp.Percent / 100
			if sp.Status != PROGRESS_ACTIVE {
				subok += value
				okindex = append(okindex, index)
			}
			percent += value
		}

		if subok > 0 {
			newSub := make([]*SubProgress, 0)
			wsp.Percent += subok
			index := int(0)
			for _, i := range okindex {
				if i > index {
					newSub = append(newSub, wsp.sub[index:i]...)

				}
				index++
			}
			if index < len(wsp.sub) {
				newSub = append(newSub, wsp.sub[index:]...)
			}
		}
	}

	wsp.Send(percent, status, msg)
}
