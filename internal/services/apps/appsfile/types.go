package appsfile

import (
	"sync"
)

type AppsService struct {
	File string
	mu   sync.Mutex
}
