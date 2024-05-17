package portainer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"washboard/db"
	"washboard/types"

	"github.com/kpango/glg"
)

func StartOrStopStack(endpointId int, stackId int, starOrStop string) (string, int, error) {
	client := &http.Client{}
	reqBody := []byte(fmt.Sprintf(`{"endpointId":%d,"id":"%d"}`, endpointId, stackId))

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/stacks/%d/%s", appState.Config.PortainerUrl, stackId, starOrStop), bytes.NewBuffer(reqBody))
	if err != nil {
		glg.Errorf("Failed to create request: %s", err)
		return "", 500, err
	}

	q := req.URL.Query()
	q.Add("endpointId", fmt.Sprintf("%d", endpointId))
	req.URL.RawQuery = q.Encode()

	req.Header.Add("X-API-Key", appState.Config.PortainerSecret)
	resp, err := client.Do(req)
	if err != nil {
		glg.Errorf("Failed to send request: %s", err)
		return "", 500, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		glg.Errorf("Failed to read response: %s", err)
		return "", 500, err
	}

	var responseStack map[string]interface{}
	err = json.Unmarshal(body, &responseStack)
	if err != nil {
		glg.Errorf("Failed to unmarshal JSON: %s", err)
		return "", 500, err
	}


	switch resp.StatusCode {
	case 200:
		if responseName, ok := responseStack["Name"].(string); ok {
			return responseName, resp.StatusCode, nil
		}
		return "", 500, fmt.Errorf("response id from portainer is not a number")
	case 409:
		return "", resp.StatusCode, fmt.Errorf("%s: %d", responseStack["message"], stackId)
	default:
		errorMessage := fmt.Sprintf("%s: %d. %s", responseStack["message"], stackId, responseStack["details"])
		glg.Error(errorMessage)
		return "", resp.StatusCode, fmt.Errorf(errorMessage)
	}
}


func ManageContainer(endpointId int, containerId string, action types.ContainerAction) (string, error) {
	client := &http.Client{}
    reqBody := bytes.NewBuffer([]byte("{}"))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/endpoints/%d/docker/containers/%s/%s", appState.Config.PortainerUrl, endpointId, containerId, action), reqBody)
	if err != nil {
		glg.Errorf("Failed to create request: %s", err)
		return "", err
	}

	req.Header.Add("X-API-Key", appState.Config.PortainerSecret)

	// print req body
	glg.Infof("Request body: %+v", req.Body)
	resp, err := client.Do(req)
	if err != nil {
		glg.Errorf("Failed to send request: %s", err)
		return "", err
	}

	defer resp.Body.Close()

	// get status code
	if resp.StatusCode != 204 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			glg.Errorf("Failed to read response: %s", err)
			return "", err
		}
		return "", fmt.Errorf("Failed to %s container: %s", action, body)
	}
	return "success", nil
}

func PerformSync(syncOptions *types.SyncOptions) error {
	var stacks []types.StackDto

	for _, endpoint := range syncOptions.EndpointIds {
		tmp, err := GetStacks(endpoint, true)
		if err != nil {
			return fmt.Errorf("failed to get containers for endpoint %d: %w", endpoint, err)
		}
		stacks = append(stacks, tmp...)
	}

	collectedStackMap := make(map[string]*types.StackSettings)

	allStackSettings, err := db.GetAllStackSettings()
	if err != nil {
		return fmt.Errorf("failed to get all stack settings: %w", err)
	}

	for _, stack := range allStackSettings {
		collectedStackMap[stack.StackName] = &stack
	}

	newStackCount := 0
	stackSettingsToAdd := make([]*types.StackSettings, 0)

	for _, stack := range stacks {
		if stackSetting, ok := collectedStackMap[stack.Name]; !ok {
			autoStart := false
			if len(stack.Containers) > 0 {
				autoStart = true
			}
			stackSetting = &types.StackSettings{
				StackName: stack.Name,
				AutoStart: autoStart,
				Priority:  -1,
				StackId:   stack.Id,
			}
			stackSettingsToAdd = append(stackSettingsToAdd, stackSetting)
			newStackCount++
			collectedStackMap[stack.Name] = stackSetting
		} else {
			collectedStackMap[stack.Name] = stackSetting
		}
	}

	sort.Slice(stackSettingsToAdd, func(i, j int) bool {
		return stackSettingsToAdd[i].StackName < stackSettingsToAdd[j].StackName
	})

	for _, stackSetting := range stackSettingsToAdd {
		glg.Infof("adding missing stack %s", stackSetting.StackName)
		db.CreateStackSettings(stackSetting)
	}

	allStackSettings, err = db.GetAllStackSettings()
	if err != nil {
		return fmt.Errorf("failed to get all stack settings: %w", err)
	}

	stacksToRemove := make([]string, 0)

	for _, stackSettings := range allStackSettings {
		if _, ok := collectedStackMap[stackSettings.StackName]; !ok {
			stacksToRemove = append(stacksToRemove, stackSettings.StackName)
		}
	}

	for _, stack := range stacksToRemove {
		glg.Infof("removing orphaned stack %s", stack)
		err := db.DeleteStackSettings(stack)
		if err != nil {
			glg.Errorf("failed to delete orphaned stack %s", stack)
		}
	}

	allStackSettings, err = db.GetAllStackSettings()
	if err != nil {
		return fmt.Errorf("failed to get all stack settings: %w", err)
	}

	newIndex := 0
	if newStackCount > 0 {
		for _, settings := range allStackSettings {
			if settings.Priority == -1 {
				settings.Priority = newIndex
				newIndex++
				glg.Debugf("setting position of new stack %v", settings)
			} else {
				settings.Priority += newStackCount
			}
			db.UpdateStackSettings(&settings, settings.StackName)
		}
	}

	return nil
}

