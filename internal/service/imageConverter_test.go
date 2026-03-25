package service

import (
	"encoding/base64"
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func TestPictureEncoding(t *testing.T) {
	tmpDir := t.TempDir()
	imagePath := createTempJPEG(t, tmpDir)

	base64String := ProcessImage(imagePath)
	if !IsBase64(base64String) {
		t.Errorf("Expected valid base64 string to pass validation: %s", base64String)
	}
}

func IsBase64(s string) bool {
	// Check if string length is valid (must be multiple of 4)
	if len(s)%4 != 0 {
		// Exception: Base64 strings with no padding can have length not divisible by 4
		// but removing padding chars should make them divisible
		s = strings.TrimRight(s, "=")
		if len(s)%4 != 0 {
			return false
		}
	}

	pattern := "^[A-Za-z0-9+/]*={0,2}$"
	match, err := regexp.MatchString(pattern, s)
	if err != nil || !match {
		return false
	}

	// Try to decode the string
	_, err = base64.StdEncoding.DecodeString(s)
	return err == nil
}

func TestReduceImageSize(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)
	imagePath := createTempJPEG(t, tmpDir)

	reducedImagePath := ReduceImageSize(imagePath)
	if reducedImagePath != "" && !strings.HasSuffix(reducedImagePath, ".jpg") {
		t.Errorf("Expected reduced image path to end with '.jpg': %s", reducedImagePath)
	}

	if _, err := os.Stat(reducedImagePath); err != nil {
		t.Fatalf("Expected reduced image file to exist: %v", err)
	}
}

func createTempJPEG(t *testing.T, dir string) string {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	img.Set(0, 0, color.RGBA{R: 255, A: 255})

	path := filepath.Join(dir, "test.jpg")
	file, err := os.Create(path)
	if err != nil {
		t.Fatalf("failed to create temp image: %v", err)
	}
	defer file.Close()

	if err := jpeg.Encode(file, img, nil); err != nil {
		t.Fatalf("failed to encode temp image: %v", err)
	}

	return path
}
