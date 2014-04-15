package main

import (
	"defcon/config"
	"defcon/zero"
	"fmt"
	"github.com/googollee/go-socket.io"
	"log"
	"os"
	"net/http"
)

var dcConfig config.DefconConfig

type appHandler func(http.ResponseWriter, *http.Request) error

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("[Web] %v %v 200 OK\n", r.Method, r.URL.Path)
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

	http.ServeFile(w, r, dcConfig.StaticRoot+path)

	return nil
}

func onConnect(ns *socketio.NameSpace) {
	fmt.Println("[Socket] Connected:", ns.Id(), "in channel", ns.Endpoint())
}

func main() {
	if _, err := os.Stat("settings.yaml"); err == nil {
		dcConfig = config.FromYaml("settings.yaml")
	} else {
		dcConfig = config.FromEnv()
	}

	sock_config := &socketio.Config{}
	sock_config.HeartbeatTimeout = 2
	sock_config.ClosingTimeout = 4

	fmt.Println("*** STARTING DEFCON ***")

	sio := socketio.NewSocketIOServer(sock_config)

	sio.On("connect", onConnect)

	sio.Handle("/", appHandler(mainHandler))

	consumer := zero.NewZeroConsumer(dcConfig.ZeroMqUrl)
	go consumer.Consume(sio.Broadcast)

	fmt.Printf("[Web] Listening on 0.0.0.0:%v\n", dcConfig.HttpPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", dcConfig.HttpPort), sio))
}
