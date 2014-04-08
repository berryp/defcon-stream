package defcon

import (
	"fmt"
	zmq "github.com/alecthomas/gozmq"
)

const url = "tcp://vdev-tabo.paix.yougov.local:8888"

func Zero() {
	context, _ := zmq.NewContext()
	defer context.close()

	socket, _ := context.NewSocket(zmq.SUB)
	defer socket.close()

	socket.Connect(url)
	fmt.Println("Zero listening on "+url)

	socket.SetSubscribe("")

	for {
		msg, _ := socket.Recv(0)
		fmt.Println(string(msg))
	}
}
