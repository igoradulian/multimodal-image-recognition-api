package app

import (
	"bytes"
	"demo-basic-ai-chat-bot/internal/dto"
	"demo-basic-ai-chat-bot/internal/service"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const SlmVersionUrl = "http://localhost:11434/api/tags"
const SlmChatUrl = "http://localhost:11434/api/chat"

func GetBotResponse(c *gin.Context) {

	chatResp := dto.ChatResponse{ID: "1", TEXT: "Hello, I'm chatbot"}
	c.IndentedJSON(http.StatusOK, chatResp)
}

func PostQuestion(c *gin.Context) {
	var newMessage dto.UseRequest

	newMessage.MESSAGE = c.PostForm("message")

	// single file
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "image file is required", "details": err.Error()})
		return
	}
	log.Println(file.Filename)

	filename := filepath.Base(file.Filename)
	dst := fmt.Sprintf("./files/%s", filename)

	fmt.Println(dst)

	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to save the file",
			"details": err.Error(),
		})
		return
	}

	// Log the received data for debugging
	log.Printf("Received response: ID=%s, TEXT=%s, FILENAME=%s",
		newMessage.ID, newMessage.MESSAGE, newMessage.FILENAME)

	// Read the file content
	imageReqDto := dto.ImageRequest{}
	imageReqDto.Model = "llava" //"smollm2:135m"
	imageReqDto.Stream = false
	messageDto := dto.Messages{}

	messageDto.Role = "user"
	messageDto.Content = newMessage.MESSAGE
	messageDto.Images = append(messageDto.Images, service.ProcessImage(dst))
	imageReqDto.Messages = append(imageReqDto.Messages, messageDto)

	fmt.Printf("%+v\n", imageReqDto)
	b, err := json.Marshal(imageReqDto)

	if err != nil {
		log.Printf("Error marshalling JSON: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp, err := http.Post(SlmChatUrl, "application/json", bytes.NewBuffer(b))

	if err != nil {
		log.Printf("Error making POST request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Error response from server: %s", resp.Status)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get response from server"})
		return
	}

	var botResposne map[string]interface{}

	if err := json.NewDecoder(resp.Body).Decode(&botResposne); err != nil {
		log.Printf("Error decoding JSON response: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, botResposne)
}

func CheckBotVersion(c *gin.Context) {

	resp, err := http.Get(SlmVersionUrl)

	if err != nil {
		log.Printf("Error fetching version: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch version"})
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Printf("Error fetching version: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch version"})
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to read response body",
			"details": err.Error(),
		})
		return
	}
	var jsonData map[string]interface{}

	if err := json.Unmarshal(body, &jsonData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to parse JSON response",
			"details": err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusOK, jsonData)
}

// PostOpenAIVision handles image+prompt requests using OpenAI vision capable model.
// Form fields: message (text prompt), file (image upload)
func PostOpenAIVision(c *gin.Context) {
	prompt := c.PostForm("message")

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "image file is required", "details": err.Error()})
		return
	}

	filename := filepath.Base(file.Filename)
	dst := fmt.Sprintf("./files/%s", filename)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file", "details": err.Error()})
		return
	}

	// Call OpenAI vision service.
	analysis, err := service.AnalyzeImageWithOpenAI(c.Request.Context(), dst, prompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "openai vision failed", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"filename": filename,
		"prompt":   prompt,
		"result":   analysis,
	})
}

func GoogleChatService(c *gin.Context) {
	prompt := c.PostForm("message")
	_ = prompt

	file, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "image file is required", "details": err.Error()})
		return
	}

	directory := uuid.NewString()
	homeDir, err := os.UserHomeDir()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to determine user home directory", "details": err.Error()})
		return
	}

	dirPath := filepath.Join(homeDir, "img", "files", directory)
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create directory", "details": err.Error()})
	}

	filename := filepath.Base(file.Filename)
	imagePath := filepath.Join(dirPath, filename)
	if err := c.SaveUploadedFile(file, imagePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file", "details": err.Error()})
		return
	}

	reducePath := service.ReduceImageSize(imagePath)
	//_ = reducePath
	result, err := service.GoogleChatService(reducePath)

	//result := DummyProduct() // Replace with actual call to GoogleChatService

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file", "details": err.Error()})
		return
	}

	payload, err := json.Marshal(result)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to marshal response", "details": err.Error()})
		return
	}

	c.Data(http.StatusOK, "application/json", payload)
}

func DummyProduct() dto.ProductResponse {
	var quantity *string
	q := "4"
	quantity = &q

	var unit *string
	u := "fl oz"
	unit = &u

	var expirationDate *string
	exp := ""
	expirationDate = &exp

	return dto.ProductResponse{
		Name:           "Children's Motrin Oral Suspension (NSAID) Pain Reliever/Fever Reducer",
		Category:       "MEDICATION",
		Brand:          "Motrin",
		ExpirationDate: expirationDate,
		Quantity:       quantity,
		Unit:           unit,
		Status:         "NEW",
		Priority:       "MEDIUM",
	}
}
