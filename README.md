# Demo Basic AI Chat Bot

Minimal Go API that accepts text + image uploads and forwards them to local LLMs (Ollama) or cloud vision models (OpenAI, Gemini).

## Features

- Local inference via Ollama for lightweight models
- Cloud vision endpoints using OpenAI and Gemini
- Multipart form uploads for image + prompt
- Simple REST interface for easy integration

## Prerequisites

- Go 1.20+ (based on `go.mod`)
- Ollama (optional, only for local `/responses` endpoint)
- OpenAI and/or Gemini API keys for cloud endpoints

## Setup

1) Copy environment variables

```bash
cp .env.example .env
```

2) Fill in API keys in `.env`

- `OPENAI_API_KEY` for `/api/v1/openai/vision`
- `GEMINI_API_KEY` for `/api/v1/google/vision`

## Run

```bash
go run ./cmd
```

Server starts on `:8089`.

## Endpoints

- `POST /api/v1/openai/vision` (form-data: `message`, `file`)
- `POST /api/v1/google/vision` (form-data: `message`, `data`)
- `POST /api/v1/responses` (form-data: `message`, `file`)
  - Forwards to local Ollama at `http://localhost:11434`

## Notes

- Uploaded images are written to `./files/` or to `$HOME/img/files/<uuid>` for the Google endpoint.
- Keep `.env` out of version control; use `.env.example` as a template.
- If you do not use a cloud provider, you can omit its API key.
