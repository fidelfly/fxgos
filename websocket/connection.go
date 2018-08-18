package websocket

import (
	"log"

	"sync"

	"net/http"

	gws "github.com/gorilla/websocket"
)

const (
	STANDBY = iota
	OPENED
	CLOSED
)

const (
	TextMessage   = 1
	BinaryMessage = 2
)

type WsConnect struct {
	Code             string
	Decoder          WsDecoder
	Encoder          WsEncoder
	Status           uint
	Conn             *gws.Conn
	receivers        []WsReceiver
	closeHandlers    []WsCloseHandler
	receiveLock      *sync.Mutex
	closeHandlerLock *sync.Mutex
	writerChan       chan interface{}
}

type WsDecoder func([]byte) (interface{}, error)
type WsEncoder func(interface{}) (int, []byte, error)
type WsCloseHandler func(int, string) error

type WsWriter interface {
	EncodeWsMessage() (int, []byte, error)
}

type WsReceiver func(interface{})

func (wsc *WsConnect) AddReceiver(receiver WsReceiver) {
	wsc.receiveLock.Lock()
	defer wsc.receiveLock.Unlock()
	if wsc.receivers == nil {
		wsc.receivers = make([]WsReceiver, 1)
	}
	wsc.receivers = append(wsc.receivers, receiver)
}

func (wsc *WsConnect) AddCloseHandler(handler WsCloseHandler) {
	wsc.closeHandlerLock.Lock()
	defer wsc.closeHandlerLock.Unlock()
	if wsc.closeHandlers == nil {
		wsc.closeHandlers = make([]WsCloseHandler, 1)
	}
	wsc.closeHandlers = append(wsc.closeHandlers, handler)
}

func (wsc *WsConnect) SendMessage(message interface{}) {
	if wsc.Status == OPENED {
		wsc.writerChan <- message
	}
}

func (wsc *WsConnect) SetupConnection(ws *gws.Conn) {
	wsc.Conn = ws
	wsc.Status = OPENED
	wsc.writerChan = make(chan interface{}, 100)
	wsc.Conn.SetCloseHandler(wsc.onClose)
}

func (wsc *WsConnect) GetStatus() uint {
	return wsc.Status
}

func (wsc *WsConnect) IsOpen() bool {
	return wsc.Status == OPENED
}

func (wsc *WsConnect) ListenAndServe() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		wsc.startReader()
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		wsc.startWriter()
	}()

	wg.Wait()
}

func (wsc *WsConnect) onClose(code int, text string) error {
	wsc.Status = CLOSED
	if len(wsc.closeHandlers) > 0 {
		for _, handler := range wsc.closeHandlers {
			err := handler(code, text)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (wsc *WsConnect) notifyReceiver(message interface{}) {
	if len(wsc.receivers) > 0 {
		for _, receiver := range wsc.receivers {
			receiver(message)
		}
	}
}

func (wsc *WsConnect) startReader() {
	if wsc.Status == OPENED {
		for {
			if wsc.Status == OPENED {
				_, p, err := wsc.Conn.ReadMessage()
				if err != nil {
					if gws.IsUnexpectedCloseError(err, gws.CloseGoingAway, gws.CloseAbnormalClosure) {
						log.Printf("error: %v", err)
					}
					break
				}
				var message interface{}
				if wsc.Decoder != nil {
					message, err = wsc.Decoder(p)
					if err != nil {
						message = p
					}
				} else {
					message = p
				}

				wsc.notifyReceiver(message)
			} else {
				break
			}
		}
	}
}

func (wsc *WsConnect) startWriter() {
	for {
		select {
		case message := <-wsc.writerChan:
			if wsc.Encoder != nil {
				messageType, data, err := wsc.Encoder(message)
				if err != nil {
					wsc.Conn.WriteMessage(messageType, data)
				}
			} else {
				if encoder, ok := message.(WsWriter); ok {
					messageType, data, err := encoder.EncodeWsMessage()
					if err == nil {
						wsc.Conn.WriteMessage(messageType, data)
						break
					}
				}
				if text, ok := message.(string); ok {
					wsc.Conn.WriteMessage(TextMessage, []byte(text))
					break
				} else {
					wsc.Conn.WriteJSON(message)
				}
			}

			break
		default:
			if wsc.Status != OPENED {
				return
			}
		}
	}
}

func SetupWebsocket(wsc *WsConnect, w http.ResponseWriter, r *http.Request, headers ...map[string]string) (err error) {
	respHeader := http.Header{}
	if len(headers) > 0 {
		for _, header := range headers {
			for key, value := range header {
				respHeader.Set(key, value)
			}
		}
	}

	webSocket, err := upgrader.Upgrade(w, r, respHeader)
	if err != nil {
		return
	}

	wsc.SetupConnection(webSocket)
	return
}

var upgrader = gws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
