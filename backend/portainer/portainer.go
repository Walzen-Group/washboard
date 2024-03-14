package portainer

import (
	"time"
	"washboard/state"

	"github.com/patrickmn/go-cache"
)

var appState *state.Data = state.Instance()
var portainerCache = cache.New(time.Duration(appState.Config.CacheDurationMinutes), 10*time.Minute)
