package config

import "os"

type Config struct{ Port, TransportURL, OllamaURL, OllamaModel string }

func Load() Config {
	return Config{env("PORT", "8083"), env("TRANSPORT_URL", "http://localhost:8082"), env("OLLAMA_URL", "http://localhost:11434"), env("OLLAMA_MODEL", "qwen2.5:7b")}
}
func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
