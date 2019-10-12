package pub

import "github.com/fidelfly/fxgo/pubsubx"

var myPubSub = pubsubx.New(3)

func Subscribe(topic string, handler ...pubsubx.SubscriberHandler) pubsubx.Subscriber {
	return myPubSub.Subscribe(topic, handler...)
}

func UnSubscribe(sub pubsubx.Subscriber, topics ...string) {
	myPubSub.UnSubscribe(sub, topics...)
}

func Publish(msg interface{}, topics ...string) {
	myPubSub.Publish(msg, topics...)
}
