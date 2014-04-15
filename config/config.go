package config

import (
	"os"
	"strconv"
)

const defaultStaticRoot = "src/defcon/app"
const defaultHttpPort = 8080
const defaultZeroMqUrl = "http://localhost:8888"

type DefconConfig struct {
	StaticRoot string
	HttpPort   int
	ZeroMqUrl  string
}

func FromEnv() DefconConfig {
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
		HttpPort:   httpPort,
		ZeroMqUrl:  zeroMqUrl,
	}
}
