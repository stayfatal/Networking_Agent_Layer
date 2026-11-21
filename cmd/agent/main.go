package main

import (
	"github.com/stayfatal/Networking_Agent_Layer/internal/adapter/httpapi"
	"github.com/stayfatal/Networking_Agent_Layer/internal/adapter/ollama"
	"github.com/stayfatal/Networking_Agent_Layer/internal/adapter/transport"
	"github.com/stayfatal/Networking_Agent_Layer/internal/config"
	"github.com/stayfatal/Networking_Agent_Layer/internal/usecase"
	"log"
	"net/http"
)

func main() {
	cfg := config.Load()
	service := usecase.New(ollama.New(cfg.OllamaURL, cfg.OllamaModel), transport.New(cfg.TransportURL))
	log.Fatal(http.ListenAndServe(":"+cfg.Port, httpapi.New(service)))
}
