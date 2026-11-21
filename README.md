# Networking Agent Layer

Сервис генерирует названия товаров через Ollama и передаёт описание и результат транспортному сервису сегментами.

## Запуск

Требуется Go 1.23 или новее. Запустите Ollama с нужной моделью и транспортный сервис, затем:

```sh
go run ./cmd/agent
```

| Переменная | Значение по умолчанию |
| --- | --- |
| `PORT` | `8083` |
| `TRANSPORT_URL` | `http://localhost:8082` |
| `OLLAMA_URL` | `http://localhost:11434` |
| `OLLAMA_MODEL` | `qwen2.5:7b` |

`POST /api/v1/process` принимает JSON с полями `sender`, `sent_at` (RFC3339), `request_id`, `product_description` и возвращает 202. Обработка выполняется асинхронно. `GET /health` проверяет доступность сервиса.

## Структура

- `cmd/agent` — сборка зависимостей и запуск HTTP-сервера.
- `internal/domain` — модель и проверка запроса, контракты обмена.
- `internal/usecase` — сценарий обработки и интерфейсы внешних зависимостей.
- `internal/adapter` — HTTP API, клиенты Ollama и транспорта.
- `internal/config` — конфигурация из окружения.

Проверка: `go test ./...`.
