package portainer

import (
	"time"
	"washboard/state"

	"github.com/patrickmn/go-cache"
)

var appState *state.Data = state.Instance()
var portainerCache = declarePortainerCache()
var fallbackCache = cache.New(cache.NoExpiration, cache.NoExpiration)

const FallbackCacheLastUpdatedKey = "__fallback_cache_last_updated__"

func declarePortainerCache() (*cache.Cache) {
	return cache.New(time.Duration(appState.Config.CacheDurationMinutes) * time.Minute, 10*time.Minute)
}
