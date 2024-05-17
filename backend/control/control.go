package control

import (
	"sort"
	"washboard/db"
	"washboard/portainer"
	"washboard/types"

	"github.com/kpango/glg"
)

func StopAllStacks(endpointId int) error {
	settings, err := db.GetAllStackSettings()
	if err != nil {
		return err
	}

	stacks, err := portainer.GetStacks(endpointId, true)
	if err != nil {
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
		if _, ok := stackMap[setting.StackName]; ok {
			portainer.StartOrStopStack(endpointId, setting.StackId, "stop")
			glg.Infof("stopped %s", setting.StackName)
		}
	}

	return nil
}

func SyncAutoStartState(endpointId int) error {
	settings, err := db.GetAllStackSettings()
	if err != nil {
		return err
	}

	stacks, err := portainer.GetStacks(endpointId, true)
	if err != nil {
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
	return nil
}
