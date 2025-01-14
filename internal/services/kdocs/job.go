package kdocs

import (
	"KingExporter/internal/services/api"
)

type DownloadJob struct {
	Url      string
	FullPath string
}

type PreloadJob struct {
	File       api.File
	GroupID    int
	FullPath   string
	RetryCount int
	MaxRetries int
}
