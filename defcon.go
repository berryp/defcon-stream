package main

import (
	"defcon/zero"
	"fmt"
	"github.com/googollee/go-socket.io"
	"log"
	"net/http"
	"os"
	"strconv"
)

const defaultStaticRoot = "src/defcon/app"
const defaultHttpPort = 8080
const defaultZeroMqUrl = "http://localhost:8888"

type DefconConfig struct {
	StaticRoot string
	HttpPort int
	ZeroMqUrl string
}

func NewConfig() DefconConfig {
	var staticRoot = os.Getenv("DEFCON_STATIC_ROOT")
	var httpPortStr = os.Getenv("DEFCON_HTTP_PORT")
	var zeroMqUrl = os.Getenv("DEFCON_ZEROMQ_URL")

	if staticRoot == "" {
		staticRoot = defaultStaticRoot
	}

	var httpPort, err = strconv.Atoi(httpPortStr)
	if err != nil {
		httpPort = defaultHttpPort
	}

	if zeroMqUrl == "" {
		zeroMqUrl = defaultZeroMqUrl
	}

	return DefconConfig{
		StaticRoot: staticRoot,
		HttpPort: httpPort,
		ZeroMqUrl: zeroMqUrl,
	}
}

var config = NewConfig()

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

	http.ServeFile(w, r, config.StaticRoot+path)

	return nil
}

func onConnect(ns *socketio.NameSpace) {
	fmt.Println("[Socket] Connected:", ns.Id(), "in channel", ns.Endpoint())
}

func main() {
	sock_config := &socketio.Config{}
	sock_config.HeartbeatTimeout = 2
	sock_config.ClosingTimeout = 4

	fmt.Println("*** STARTING DEFCON ***")

	sio := socketio.NewSocketIOServer(sock_config)

	sio.On("connect", onConnect)

	sio.Handle("/", appHandler(mainHandler))

	consumer := zero.NewZeroConsumer(config.ZeroMqUrl)
	go consumer.Consume(sio.Broadcast)

	fmt.Printf("[Web] Listening on 0.0.0.0:%v\n", config.HttpPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", config.HttpPort), sio))
}
