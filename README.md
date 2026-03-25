# AI Vision Chat API (Go)

Production-style demo API that accepts text + image uploads and routes them to:
- local multimodal inference via Ollama
- cloud vision models via OpenAI and Google Gemini

This project is designed as a portfolio backend sample focused on API design, multipart handling, provider abstraction, and secure environment-based configuration.

## Why this project

- Demonstrates practical Go backend skills with Gin routing and handler separation
- Shows multimodel AI integration in one service (`Ollama`, `OpenAI`, `Gemini`)
- Includes DTO/service layering to keep request/response logic organized
- Uses env vars for credentials (no hardcoded API keys)

## Tech Stack

- Go 1.20+
- Gin (`github.com/gin-gonic/gin`)
- OpenAI Chat Completions (vision request format)
- Google GenAI SDK (`google.golang.org/genai`)
- Ollama local API (`http://localhost:11434`)

## Project Structure

```text
cmd/
  main.go                 # app entrypoint and route registration
internal/
  dto/                    # request/response models
  service/                # OpenAI, Gemini, image utilities
web/
  app/
    chat_app.go           # HTTP handlers
```

## API Endpoints

Base URL: `http://localhost:8089/api/v1`

- `GET /responses`
  - Health-style sample chatbot response
- `GET /responses/version`
  - Proxies Ollama tag/version info
- `POST /responses`
  - Local Ollama multimodal chat
  - Multipart fields: `message`, `file`
- `POST /openai/vision`
  - OpenAI vision analysis
  - Multipart fields: `message`, `file`
- `POST /google/vision`
  - Gemini-based product extraction flow
  - Multipart fields: `message`, `data`

## Quick Start

1. Clone and enter the repo.
2. Copy environment template.
3. Run the server.

```bash
cp .env.example .env
go run ./cmd
```

Server listens on `:8089`.

## Configuration

Set values in `.env` (or export in your shell):

- `OPENAI_API_KEY` for `POST /api/v1/openai/vision`
- `GEMINI_API_KEY` for `POST /api/v1/google/vision`

Template file: `.env.example`

## Example Requests

### OpenAI vision

```bash
curl -X POST "http://localhost:8089/api/v1/openai/vision" \
  -F "message=Describe what you see" \
  -F "file=@/absolute/path/to/image.jpg"
```

### Gemini vision

```bash
curl -X POST "http://localhost:8089/api/v1/google/vision" \
  -F "message=Extract product details" \
  -F "data=@/absolute/path/to/image.jpg"
```

### Ollama local vision

```bash
curl -X POST "http://localhost:8089/api/v1/responses" \
  -F "message=What is in this image?" \
  -F "file=@/absolute/path/to/image.jpg"
```

## Security Notes

- Do not commit `.env`.
- Keep API keys in environment variables only.
- `.gitignore` is configured to ignore `.env` and `.env.*` (except `.env.example`).

## Portfolio Highlights

- Multipart upload handling and file persistence
- External API orchestration with structured error handling
- Model output normalization into typed DTO responses
- Clear path to evolve into provider interface abstraction and tests per provider

## License

MIT - see `LICENSE`.
