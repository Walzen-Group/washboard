package control

import (
	"errors"
	"sort"
	"time"
	"washboard/db"
	"washboard/portainer"
	"washboard/types"
	"washboard/werrors"

	"github.com/patrickmn/go-cache"

	"github.com/kpango/glg"
)

var controlCache = declareControlCache()

func declareControlCache() (*cache.Cache) {
	return cache.New(time.Duration(30 * time.Second), time.Duration(1 * time.Minute))
}


func StopAllStacks(endpointId int) error {
	if _, found := controlCache.Get("stopAllStacks"); found {
		glg.Infof("stopAllStacks already in progress")
		return werrors.NewAlreadyInProgressError(errors.New("Operation can only be performed once"), "A stop operation is already in progress.")
	}
	controlCache.Set("stopAllStacks", true, cache.DefaultExpiration)

	portainer.PerformSync(&types.SyncOptions{EndpointIds: []int{endpointId}})
	settings, err := db.GetAllStackSettings()
	if err != nil {
		controlCache.Delete("stopAllStacks")
		return err
	}

	stacks, err := portainer.GetStacks(endpointId, true)
	if err != nil {
		controlCache.Delete("stopAllStacks")
		return err
	}

	stackMap := make(map[string]types.StackDto)
	for _, stack := range stacks {
		stackMap[stack.Name] = stack
	}

	sort.Slice(settings, func(i, j int) bool {
		return settings[i].Priority > settings[j].Priority
	})

	for _, setting := range settings {
		if stack, ok := stackMap[setting.StackName]; ok {
			if types.CheckWashbImage(stack) {
				glg.Infof("not modifying stack containing a washboard image")
				continue
			}
			//portainer.StartOrStopStack(endpointId, setting.StackId, "stop")
			glg.Infof("stopped %s", setting.StackName)
		}
	}
	controlCache.Delete("stopAllStacks")
	return nil
}

func SyncAutoStartState(endpointId int) error {
	if _, found := controlCache.Get("syncAutoStartState"); found {
		glg.Infof("syncAutoStartState already in progress")
		return werrors.NewAlreadyInProgressError(errors.New("Operation can only be performed once"), "A container state sync operation is already in progress.")
	}
	controlCache.Set("syncAutoStartState", true, cache.DefaultExpiration)

	portainer.PerformSync(&types.SyncOptions{EndpointIds: []int{endpointId}})
	settings, err := db.GetAllStackSettings()
	if err != nil {
		controlCache.Delete("syncAutoStartState")
		return err
	}

	stacks, err := portainer.GetStacks(endpointId, true)
	if err != nil {
		controlCache.Delete("syncAutoStartState")
		return err
	}

	stackMap := make(map[string]types.StackDto)
	for _, stack := range stacks {
		stackMap[stack.Name] = stack
	}

	sort.Slice(settings, func(i, j int) bool {
		return settings[i].Priority < settings[j].Priority
	})

	for _, setting := range settings {
		if setting.AutoStart {
			if stack, ok := stackMap[setting.StackName]; ok {

				if types.CheckWashbImage(stack) {
					glg.Infof("not modifying stack containing a washboard image")
					continue
				}

				if len(stack.Containers) > 0 {
					// check if all containers are stopped. This is an indicator that it probably got stopped by the user
					// It's unlikely that all containers crashed at once
					// and thus warrants a restart attempt
					allStopped := true
					for _, container := range stack.Containers {
						if container.Status == types.ContainerRunning {
							allStopped = false
							break
						}
					}
					if allStopped {
						portainer.StartOrStopStack(endpointId, setting.StackId, "stop")
						portainer.StartOrStopStack(endpointId, setting.StackId, "start")
						glg.Infof("synced autostart state for %s (restore state)", setting.StackName)
					}
				} else {
					// starts those that are not running but have the autostart flag
					portainer.StartOrStopStack(endpointId, setting.StackId, "start")
					glg.Infof("synced autostart state for %s (start)", setting.StackName)
				}
			} else {
				glg.Infof("synced autostart state for %s (no action)", setting.StackName)
			}
		}
	}
	controlCache.Delete("syncAutoStartState")
	return nil
}


