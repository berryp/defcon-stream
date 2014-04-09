package main

import (
	"defcon/zero"
	"fmt"
	"github.com/googollee/go-socket.io"
	"log"
	"net/http"
	"os"
)

const staticRoot = "src/defcon/app"

type appHandler func(http.ResponseWriter, *http.Request) error

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if err := fn(w, r); err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func mainHandler(w http.ResponseWriter, r *http.Request) error {
	var path = r.URL.Path

	if path == "/" {
		path = "/index.html"
	}

	http.ServeFile(w, r, staticRoot+path)

	return nil
}

func onConnect(ns *socketio.NameSpace) {
	fmt.Println("connected:", ns.Id(), "in channel", ns.Endpoint())
}

func main() {
	url := os.Getenv("DEFCON_ZEROMQ_URL")

	sock_config := &socketio.Config{}
	sock_config.HeartbeatTimeout = 2
	sock_config.ClosingTimeout = 4

	sio := socketio.NewSocketIOServer(sock_config)

	sio.On("connect", onConnect)

	sio.Handle("/", appHandler(mainHandler))

	consumer := zero.NewZeroConsumer(url)
	go consumer.Consume(sio.Broadcast)

	port := os.Getenv("PORT")

	fmt.Printf("Listening on 0.0.0.0:%v\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), sio))
}
