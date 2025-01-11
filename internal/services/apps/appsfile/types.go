package appsfile

import (
	"sync"
)

type AppsFileService struct {
	File string
	mu   sync.Mutex
}
