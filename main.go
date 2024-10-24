package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path"

	"github.com/g00dv1n/sol-rpc-router/pkg/router"
)

type ProxyServerConfig struct {
	Port       int          `json:"port"`
	Host       string       `json:"host,omitempty"`
	RegularRpc router.Route `json:"regularRpc"`
	DasRpc     router.Route `json:"dasRpc"`
}

var logger *slog.Logger

func main() {
	logger = slog.New(slog.NewTextHandler(os.Stdout, nil))

	exDir, dirErr := os.Getwd()
	if dirErr != nil {
		log.Fatal(dirErr)
	}

	// Get Config Path
	var configPath string
	defaultConfigPath := path.Join(exDir, "proxy_config.json")

	flag.StringVar(&configPath, "c", defaultConfigPath, "config file path")
	flag.Parse()

	configFileRaw, fileErr := os.ReadFile(configPath)

	if fileErr != nil {
		log.Fatalf("Can't open config file: %s", fileErr)
	}

	var proxyServerConfig ProxyServerConfig
	jsonErr := json.Unmarshal(configFileRaw, &proxyServerConfig)

	if jsonErr != nil {
		log.Fatalf("Can't parse config: %s \n", jsonErr)
	}

	logger.Info("Config file loaded", "path", path.Base(configPath))

	addr := fmt.Sprintf("%s:%d", proxyServerConfig.Host, proxyServerConfig.Port)

	serverErr := NewProxyServer(addr, proxyServerConfig.RegularRpc, proxyServerConfig.DasRpc)
	if serverErr != nil {
		log.Fatal(serverErr)
	}
}

func NewProxyServer(addr string, regular router.Route, das router.Route) error {
	server := http.NewServeMux()
	router, err := router.NewRouter(regular, das)

	if err != nil {
		return err
	}

	server.Handle("/", router)

	logger.Info("Running proxy server", "addr", addr)

	return http.ListenAndServe(addr, server)
}
