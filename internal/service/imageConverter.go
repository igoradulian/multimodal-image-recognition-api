package service

import (
	"encoding/base64"
	"image"
	"image/jpeg"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/nfnt/resize"
)

func ProcessImage(imagePath string) string {
	// Read the image file
	file, err := os.Open(imagePath)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer file.Close()

	// Read the file content
	imageData, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("failed to read file: %v", err)
	}

	// Encode the image data to Base64
	base64String := base64.StdEncoding.EncodeToString(imageData)

	// Print the Base64 string
	return base64String
}

func ReduceImageSize(imagePath string) string {
	file, err := os.Open(imagePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	newImage := resize.Resize(800, 600, img, resize.Lanczos3)

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("failed to determine user home directory: %v", err)
	}

	reducedDir := filepath.Join(homeDir, "img", "reduced")
	if err := os.MkdirAll(reducedDir, 0o755); err != nil {
		log.Fatalf("failed to ensure reduced directory: %v", err)
	}

	outName := uuid.NewString() + ".jpg"
	outPath := filepath.Join(reducedDir, outName)
	out, err := os.Create(outPath)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	if err := jpeg.Encode(out, newImage, nil); err != nil {
		log.Fatal(err)
	}

	return outPath // return the reduced image relative path
}
