package fileconfig

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type FileConfig struct {
	BasePath      string
	AllowedExtMap map[string][]string
	mu            sync.Mutex
}

var Config *FileConfig

func InitFileConfig() {
	basePath := os.Getenv("FILE_BASE_PATH")
	if basePath == "" {
		basePath = "./uploads"
	}

	Config = &FileConfig{
		BasePath:      basePath,
		AllowedExtMap: make(map[string][]string),
	}
}

func (fc *FileConfig) GetPath(contentType string) string {
	contentType = sanitize(contentType)
	return filepath.Join(fc.BasePath, contentType)
}

func (fc *FileConfig) GetAllowedExtensions(contentType string) []string {
	contentType = sanitize(contentType)
	fc.mu.Lock()
	defer fc.mu.Unlock()
	return fc.AllowedExtMap[contentType]
}

func (fc *FileConfig) SetAllowedExtensions(contentType string, extensions []string) {
	contentType = sanitize(contentType)
	fc.mu.Lock()
	defer fc.mu.Unlock()
	fc.AllowedExtMap[contentType] = extensions

	dir := fc.GetPath(contentType)
	if err := os.MkdirAll(dir, 0755); err != nil {
		panic("Klasör oluşturulamadı: " + dir + " | Hata: " + err.Error())
	}
}

func (fc *FileConfig) IsExtensionAllowed(contentType, ext string) bool {
	ext = strings.ToLower(strings.TrimPrefix(ext, "."))
	for _, allowed := range fc.GetAllowedExtensions(contentType) {
		if allowed == ext {
			return true
		}
	}
	return false
}

func sanitize(str string) string {
	str = strings.ToLower(str)
	str = strings.TrimSpace(str)
	str = strings.ReplaceAll(str, " ", "_")
	return str
}
