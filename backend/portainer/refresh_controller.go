package portainer

import (
	"sync"
	"sync/atomic"
	"time"

	"washboard/types"

	"github.com/kpango/glg"
)

type ImageRefreshController struct {
	mu         sync.Mutex
	running    bool
	startedAt  int64
	finishedAt int64
	lastError  string
	version    atomic.Uint64
}

var Refresh = &ImageRefreshController{}

// TriggerRefresh kicks off runUpdateCheck in a goroutine if no refresh is in flight.
// Returns true if a new run was started, false if a run was already in progress.
func (c *ImageRefreshController) TriggerRefresh(endpointId int) bool {
	c.mu.Lock()
	if c.running {
		c.mu.Unlock()
		glg.Debugf("image refresh already running, skipping trigger")
		return false
	}
	c.running = true
	c.startedAt = time.Now().Unix()
	c.lastError = ""
	c.mu.Unlock()
	c.version.Add(1)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				glg.Errorf("panic during image refresh: %v", r)
				c.finish("panic during image refresh")
				return
			}
			c.finish("")
		}()
		runUpdateCheck(endpointId)
	}()
	return true
}

func (c *ImageRefreshController) finish(errMsg string) {
	c.mu.Lock()
	c.running = false
	c.finishedAt = time.Now().Unix()
	c.lastError = errMsg
	c.mu.Unlock()
	c.version.Add(1)
}

func (c *ImageRefreshController) Snapshot() types.ImageRefreshState {
	c.mu.Lock()
	defer c.mu.Unlock()
	return types.ImageRefreshState{
		Running:    c.running,
		StartedAt:  c.startedAt,
		FinishedAt: c.finishedAt,
		Error:      c.lastError,
	}
}

func (c *ImageRefreshController) Version() uint64 {
	return c.version.Load()
}
