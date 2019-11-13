package mail

import "github.com/fidelfly/gox/logx"

type MailData map[string]interface{}

type Receiver interface {
	GetAddress() []string
}

type Dispatcher interface {
	GetReceiver() []Receiver
}

type Producer interface {
	Produce(Receiver, interface{}) []MessageDecorator
}

type DecoratorAsProducer []MessageDecorator

func (dap DecoratorAsProducer) Produce(receiver Receiver, data interface{}) []MessageDecorator {
	return dap
}

func SendDynamicMail(dispatcher Dispatcher, data interface{}, producers ...Producer) {
	receivers := dispatcher.GetReceiver()
	if len(receivers) > 0 {
		for _, receiver := range receivers {
			if address := receiver.GetAddress(); len(address) > 0 {
				decorators := []MessageDecorator{To(address...)}
				for _, producer := range producers {
					decorators = append(decorators, producer.Produce(receiver, data)...)
				}
				logx.CaptureError(SendMail(CreateMessage(decorators...)))
			}
		}
	}
}
