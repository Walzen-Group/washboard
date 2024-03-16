package portainer

import (
	"time"
	"washboard/state"

	"github.com/patrickmn/go-cache"
)

var appState *state.Data = state.Instance()
var portainerCache = declarePortainerCache()

func declarePortainerCache() (*cache.Cache) {
	return cache.New(time.Duration(appState.Config.CacheDurationMinutes) * time.Minute, 10*time.Minute)
}
