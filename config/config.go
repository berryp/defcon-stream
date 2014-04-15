package config

import (
	"gopkg.in/yaml.v1"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

const defaultStaticRoot = "src/defcon/app"
const defaultHttpPort = 8080
const defaultZeroMqUrl = "http://localhost:8888"

type DefconConfig struct {
	StaticRoot string `yaml:"static_root"`
	HttpPort   int    `yaml:"port"`
	ZeroMqUrl  string `yaml:"zeromq_url"`
}

func FromYaml(filename string) DefconConfig {
	dc := DefconConfig{}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	err = yaml.Unmarshal([]byte(data), &dc)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	return dc
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
