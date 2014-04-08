package defcon

import (
	"fmt"
	zmq "github.com/alecthomas/gozmq"
)

type ZeroConsumer struct {
	url string
}

type ZeroConsumerCallback func(string, ...interface{})

func (z *ZeroConsumer) Consume(callback ZeroConsumerCallback) {
	context, _ := zmq.NewContext()
	defer context.Close()

	socket, _ := context.NewSocket(zmq.SUB)
	defer socket.Close()

	socket.Connect(z.url)
	fmt.Println("Zero listening on " + z.url)

	socket.SetSubscribe("")

	for {
		msg, _ := socket.Recv(0)
		fmt.Println(string(msg))
		callback("complete", string(msg))
	}
}

func NewZeroConsumer(url string) *ZeroConsumer {
	z := &ZeroConsumer{url: url}
	return z
}
