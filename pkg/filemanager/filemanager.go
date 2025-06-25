package filemanager

import (
	"crypto/rand"
	"davet.link/configs/fileconfig"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var (
	ErrFileNotProvided = errors.New("dosya sağlanmadı")
	ErrInvalidFileType = errors.New("geçersiz dosya türü veya uzantısı")
	ErrFileTooLarge    = errors.New("dosya boyutu çok büyük")
)

const (
	DefaultMaxFileSize = 2 * 1024 * 1024
)

func DeleteFile(contentType, fileName string) {
	if fileName == "" || contentType == "" {
		return
	}

	go func() {
		const maxRetries = 5
		const retryDelay = 1 * time.Second

		absolutePath, err := filepath.Abs(filepath.Join(fileconfig.Config.GetPath(contentType), fileName))
		if err != nil {
			return
		}

		for i := 0; i < maxRetries; i++ {
			err = os.Remove(absolutePath)
			if err == nil || os.IsNotExist(err) {
				return
			}
			time.Sleep(retryDelay)
		}
	}()
}

func UploadFile(c *fiber.Ctx, formFieldName, contentType string) (string, error) {
	file, err := c.FormFile(formFieldName)
	if err != nil {
		if err == http.ErrMissingFile {
			return "", ErrFileNotProvided
		}
		return "", err
	}
	if err := validateFile(file, contentType); err != nil {
		return "", err
	}
	newFileName, err := generateUniqueFileName(file.Filename)
	if err != nil {
		return "", err
	}
	destination := filepath.Join(fileconfig.Config.GetPath(contentType), newFileName)
	if err := c.SaveFile(file, destination); err != nil {
		return "", err
	}
	return newFileName, nil
}

func validateFile(file *multipart.FileHeader, contentType string) error {
	if file.Size > DefaultMaxFileSize { return ErrFileTooLarge }
	ext := filepath.Ext(file.Filename)
	if !fileconfig.Config.IsExtensionAllowed(contentType, ext) { return ErrInvalidFileType }
	return nil
}

func generateUniqueFileName(originalName string) (string, error) {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil { return "", err }
	randomStr := fmt.Sprintf("%x", b)
	ext := filepath.Ext(originalName)
	safeBaseName := regexp.MustCompile(`[^a-zA-Z0-9_-]+`).ReplaceAllString(strings.TrimSuffix(originalName, ext), "")
	if safeBaseName == "" { safeBaseName = "file" }
	return fmt.Sprintf("%s-%s%s", randomStr, safeBaseName, ext), nil
}