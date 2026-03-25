package service

import (
	"context"
	"demo-basic-ai-chat-bot/internal/dto"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"google.golang.org/genai"
)

func GoogleChatService(imagePath string) (*dto.ProductResponse, error) {
	ctx := context.Background()
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY environment variable is not set")
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	bytes, err := os.ReadFile(imagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read image: %w", err)
	}

	parts := []*genai.Part{
		genai.NewPartFromBytes(bytes, "image/jpeg"),
		/**
		You can specify json you want to return, and the model will try to follow that format as closely as possible. The model will still try to return the information even if it can't fill in all the fields, but it will be more likely to
		include the fields you specify if you provide a format.
		*/
		genai.NewPartFromText("Describe the key details in this image and return a concise response in JSON."),
	}

	contents := []*genai.Content{
		genai.NewContentFromParts(parts, genai.RoleUser),
	}
	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		contents,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %w", err)
	}

	raw := result.Text()

	clean := strings.TrimSpace(raw)
	clean = strings.TrimPrefix(clean, "```json")
	clean = strings.TrimSuffix(clean, "```")
	clean = strings.TrimSpace(clean)

	var info dto.ProductInfo
	if err := json.Unmarshal([]byte(clean), &info); err != nil {
		return nil, fmt.Errorf("failed to decode model output: %w", err)
	}

	var response dto.ProductResponse

	response.Name = info.Product.ProductName
	response.Brand = info.Product.Brands
	response.ExpirationDate = &info.Product.ExpirationDate
	response.Quantity = &info.Product.Quantity
	response.Category = info.Product.ProductType
	response.Unit = &info.Product.ProductQuantityUnit
	response.Status = dto.StatusNew
	response.Priority = dto.PriorityMedium

	return &response, nil
}
