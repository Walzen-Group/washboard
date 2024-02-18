package portainer

import (
	"time"
	"washboard/state"

	"github.com/patrickmn/go-cache"
)

var appState *state.Data = state.Instance()
var portainerCache = cache.New(1*time.Minute, 10*time.Minute)

