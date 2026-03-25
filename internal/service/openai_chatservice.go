package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// OpenAI API endpoints and defaults.
const (
	openAIBaseURL      = "https://api.openai.com/v1"
	openAIChatPath     = "/chat/completions"
	defaultVisionModel = "gpt-4o-mini" // lightweight multimodal model
)

// ChatCompletionRequest represents the minimal structure needed for an image+text request.
type ChatCompletionRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
	Temperature float64       `json:"temperature,omitempty"`
}

type ChatMessage struct {
	Role    string               `json:"role"`
	Content []MessageContentPart `json:"content"`
}

type MessageContentPart struct {
	Type     string    `json:"type"`
	Text     string    `json:"text,omitempty"`
	ImageURL *ImageURL `json:"image_url,omitempty"`
}

type ImageURL struct {
	URL string `json:"url"`
}

// ChatCompletionResponse is a partial representation of OpenAI's response we care about.
type ChatCompletionResponse struct {
	ID      string `json:"id"`
	Model   string `json:"model"`
	Choices []struct {
		Index        int    `json:"index"`
		FinishReason string `json:"finish_reason"`
		Message      struct {
			Role    string               `json:"role"`
			Content []MessageContentPart `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Usage *struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage,omitempty"`
}

// AnalyzeImageWithOpenAI sends an image along with a user prompt to OpenAI's multimodal chat completion API.
// Returns the concatenated text content from the first choice.
func AnalyzeImageWithOpenAI(ctx context.Context, imagePath, prompt string) (string, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return "", errors.New("OPENAI_API_KEY environment variable is not set")
	}

	if prompt == "" {
		prompt = "Describe the key details in this image." // fallback prompt
	}

	// Read and base64-encode the image.
	b64, mimeType, err := encodeImageToDataURL(imagePath)
	if err != nil {
		return "", fmt.Errorf("failed to prepare image: %w", err)
	}

	// Build request payload.
	reqBody := ChatCompletionRequest{
		Model: defaultVisionModel,
		Messages: []ChatMessage{
			{
				Role: "user",
				Content: []MessageContentPart{
					{Type: "text", Text: prompt},
					{Type: "image_url", ImageURL: &ImageURL{URL: fmt.Sprintf("data:%s;base64,%s", mimeType, b64)}},
				},
			},
		},
		MaxTokens:   300,
		Temperature: 0.2,
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	endpoint := openAIBaseURL + openAIChatPath
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, io.NopCloser(bytesReader(payload)))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("request error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("openai api returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var completion ChatCompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&completion); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(completion.Choices) == 0 || len(completion.Choices[0].Message.Content) == 0 {
		return "", errors.New("no content returned by model")
	}

	// Aggregate any text parts in the first choice.
	var output string
	for _, part := range completion.Choices[0].Message.Content {
		if part.Type == "text" && part.Text != "" {
			output += part.Text + "\n"
		}
	}

	return output, nil
}

// encodeImageToDataURL reads an image file and returns base64 content and detected mime type.
func encodeImageToDataURL(path string) (string, string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", "", fmt.Errorf("open file: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return "", "", fmt.Errorf("read file: %w", err)
	}

	// Basic mime type inference from extension (fallback image/png)
	ext := filepath.Ext(path)
	mimeType := "image/png"
	switch ext {
	case ".jpg", ".jpeg":
		mimeType = "image/jpeg"
	case ".png":
		mimeType = "image/png"
	case ".gif":
		mimeType = "image/gif"
	case ".webp":
		mimeType = "image/webp"
	}

	return base64.StdEncoding.EncodeToString(data), mimeType, nil
}

// bytesReader returns an io.ReadCloser from bytes without importing bytes just for this small use.
// (We could import bytes; kept minimal here.)
func bytesReader(b []byte) io.ReadCloser {
	return io.NopCloser(&simpleByteReader{b: b})
}

type simpleByteReader struct {
	b []byte
	i int
}

func (r *simpleByteReader) Read(p []byte) (int, error) {
	if r.i >= len(r.b) {
		return 0, io.EOF
	}
	n := copy(p, r.b[r.i:])
	r.i += n
	return n, nil
}

func (r *simpleByteReader) Close() error { return nil }
