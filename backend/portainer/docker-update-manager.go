package portainer

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"washboard/db"
	"washboard/helper"
	"washboard/types"

	"github.com/kpango/glg"
	"github.com/patrickmn/go-cache"
)

// GetEndpointId returns the id of the endpoint with the given name, which is also the environment in Portainer
func GetEndpointId(endpointName string) (int, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/endpoints", appState.Config.PortainerUrl), nil)
	if err != nil {
		glg.Errorf("Failed to create request: %s", err)
		return -1, err
	}

	req.Header.Add("X-API-Key", appState.Config.PortainerSecret)
	resp, err := client.Do(req)
	if err != nil {
		glg.Errorf("Failed to send request: %s", err)
		return -1, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		glg.Errorf("Failed to read response: %s", err)
		return -1, err
	}

	var endpoints []types.EndpointDto
	err = json.Unmarshal(body, &endpoints)
	if err != nil {
		glg.Errorf("Failed to unmarshal JSON: %s", err)
		return -1, err
	}

	for _, endpoint := range endpoints {
		if endpoint.Name == endpointName {
			return endpoint.Id, nil
		}
	}

	glg.Infof("Endpoint %s not found", endpointName)
	return -1, nil
}

// StartBackgroundUpdateCheck starts a background job that checks for updates every 24 hours
func StartBackgroundUpdateCheck(endpointId int) {
	go func() {
		glg.Info("Starting background update check...")
		runUpdateCheck(endpointId)
		ticker := time.NewTicker(24 * time.Hour)
		for range ticker.C {
			if val, found := fallbackCache.Get(FallbackCacheLastUpdatedKey); found {
				if lastUpdated, ok := val.(time.Time); ok {
					if time.Since(lastUpdated) < 24*time.Hour {
						glg.Info("Skipping background update check because cache is fresh enough")
						continue
					}
				}
			}

			runUpdateCheck(endpointId)
		}
	}()
}

func runUpdateCheck(endpointId int) {
	glg.Info("Running background update check")
	stacks, err := GetStacks(endpointId, true)
	if err != nil {
		glg.Errorf("Failed to get stacks for background update check: %s", err)
		return
	}

	for _, stack := range stacks {
		// Optimization 1: Check if any container in the stack is already known to be outdated in fallback cache
		stackHasOutdated := false
		for _, container := range stack.Containers {
			if val, found := fallbackCache.Get(container.Id); found {
				if val.(string) == types.Outdated {
					stackHasOutdated = true
					break
				}
			}
		}

		if !stackHasOutdated {
			// Stack seems clean, try bulk check
			status, err := GetStackImagesStatus(stack.Id)
			if err != nil {
				glg.Errorf("Failed to get stack images status for stack %s: %s", stack.Name, err)
			} else if status == types.Updated {
				// Mark all containers as updated
				for _, container := range stack.Containers {
					fallbackCache.Set(container.Id, types.Updated, cache.NoExpiration)
				}
				continue
			}
		}

		// Check individual containers
		for _, container := range stack.Containers {
			// Skip if already outdated in fallback cache
			if val, found := fallbackCache.Get(container.Id); found && val.(string) == types.Outdated {
				continue
			}

			status, err := GetImageStatus(endpointId, container.Id)
			if err != nil {
				glg.Warnf("Error fetching image status for container id %s: %s", container.Id, err)
				fallbackCache.Delete(container.Id)
			} else {
				fallbackCache.Set(container.Id, status, cache.NoExpiration)
			}
		}
	}
	fallbackCache.Set(FallbackCacheLastUpdatedKey, time.Now(), cache.NoExpiration)
	glg.Info("Background update check finished")
}

// GetStacks returns the stacks for the given endpoint
func GetStacks(endpointId int, skeletonOnly bool) ([]types.StackDto, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/stacks", appState.Config.PortainerUrl), nil)
	if err != nil {
		glg.Errorf("Failed to create request: %s", err)
		return nil, err
	}

	q := req.URL.Query()
	q.Add("filters", fmt.Sprintf(`{"EndpointId":%d}`, endpointId))
	req.URL.RawQuery = q.Encode()
	req.Header.Add("X-API-Key", appState.Config.PortainerSecret)

	resp, err := client.Do(req)
	if err != nil {
		glg.Errorf("Failed to send request: %s", err)
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		glg.Errorf("Failed to read response: %s", err)
		return nil, err
	}

	var stacks []map[string]interface{}
	err = json.Unmarshal(body, &stacks)
	if err != nil {
		glg.Errorf("Failed to unmarshal JSON: %s", err)
		return nil, err
	}

	countCached := 0
	countUncached := 0
	var stacksDict = make(map[string]map[string]interface{})
	for _, stack := range stacks {
		// TODO: add key to stack with all images status with caching. If the status here is updated, we can skip individual container checks
		stackId := -1
		if id, ok := stack["Id"]; !ok {
			glg.Warnf("stack does not have id key")
			continue
		} else if val, ok := id.(float64); ok {
			stackId = int(val)
		}
		if stackId == -1 {
			glg.Warnf("stack id is not a number")
			continue
		}
		var allImagesStatus string
		// skip retrieval if only the stack skeleton is requested
		if skeletonOnly {
			stack["allImagesStatus"] = types.NotRequested
		} else {
			if val, ok := portainerCache.Get(fmt.Sprintf("stack-%d-images-status", stackId)); ok {
				allImagesStatus = val.(string)
				countCached++
			} else {
				allImagesStatus, err = GetStackImagesStatus(stackId)
				if err != nil {
					glg.Errorf("Failed to get stack images status: %s", err)
					allImagesStatus = "error"
				}
				portainerCache.Set(fmt.Sprintf("stack-%d-images-status", stackId), allImagesStatus, cache.DefaultExpiration)
				countUncached++
			}
			stack["allImagesStatus"] = allImagesStatus
		}

		if stackName, ok := stack["Name"]; !ok {
			glg.Warnf("stack does not have name key")
			continue
		} else if stackNameString, ok := stackName.(string); !ok {
			glg.Warnf("stack name is not a string")
			continue
		} else {
			stacksDict[stackNameString] = stack
		}
	}

	if skeletonOnly {
		glg.Infof("no image status requested")
	} else {
		glg.Infof("cached stack images status: %d, uncached stack images status: %d", countCached, countUncached)
	}

	stacksDto, err := buildStacksDto(stacksDict, endpointId)

	return stacksDto, err
}

// GetContainers returns the containers for the given endpoint. If stackName is provided, only the containers of the stack with the given label are returned, otherwise all containers are returned
func GetContainers(endpointId int, stackName string) ([]*types.ContainerDto, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/endpoints/%d/docker/containers/json", appState.Config.PortainerUrl, endpointId), nil)
	if err != nil {
		glg.Errorf("Failed to create request: %s", err)
		return nil, err
	}

	q := req.URL.Query()
	q.Add("all", "true")
	if stackName != "" {
		q.Add("filters", fmt.Sprintf(`{"label":["%s=%s"]}`, types.StackLabel, stackName))
	}
	req.URL.RawQuery = q.Encode()
	req.Header.Add("X-API-Key", appState.Config.PortainerSecret)

	resp, err := client.Do(req)
	if err != nil {
		glg.Errorf("Failed to send request: %s", err)
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		glg.Errorf("Failed to read response: %s", err)
		return nil, err
	}

	var containers []map[string]interface{}
	err = json.Unmarshal(body, &containers)
	if err != nil {
		glg.Errorf("Failed to unmarshal JSON: %s", err)
		return nil, err
	}

	containersDto := buildContainerDto(containers)

	return containersDto, nil
}

func buildStacksDto(stacks map[string]map[string]interface{}, endpointId int) ([]types.StackDto, error) {
	var stacksDto = make(map[string]*types.StackDto)
	containers, err := GetContainers(endpointId, "")
	if err != nil {
		glg.Errorf("Failed to get stack containers: %s", err)
		return nil, err
	}

	queryImageStatusContainers := make([]*types.ContainerDto, 0, len(containers))

	for _, container := range containers {
		var stackName string

		if labelRaw, ok := container.Labels[types.StackLabel]; !ok {
			continue
		} else if labelParsed, ok := labelRaw.(string); !ok {
			glg.Warnf("label %s is not a string", labelRaw)
			continue
		} else {
			stackName = labelParsed
		}
		if val, ok := stacks[stackName]["allImagesStatus"]; ok {
			if val.(string) == types.NotRequested {
				// Check fallback cache
				if cachedStatus, found := fallbackCache.Get(container.Id); found {
					container.UpToDate = cachedStatus.(string)
				} else {
					container.UpToDate = types.NotRequested
				}
			} else if val.(string) != types.Updated {
				queryImageStatusContainers = append(queryImageStatusContainers, container)
			} else {
				container.UpToDate = types.Updated
				// Update fallback cache
				fallbackCache.Set(container.Id, types.Updated, cache.NoExpiration)
			}
		}
		if val, ok := stacksDto[stackName]; ok {
			val.Containers = append(val.Containers, container)
		} else if val, ok := stacks[stackName]; ok {
			stacksDto[stackName] = &types.StackDto{
				Id:         int(val["Id"].(float64)),
				Name:       val["Name"].(string),
				Containers: []*types.ContainerDto{container},
			}
		}
	}

	for key, value := range stacks {
		if _, ok := stacksDto[key]; !ok {
			stacksDto[key] = &types.StackDto{
				Id:         int(value["Id"].(float64)),
				Name:       value["Name"].(string),
				Containers: make([]*types.ContainerDto, 0),
			}
		}
	}

	if len(queryImageStatusContainers) > 0 {
		queryContainerImageStatus(endpointId, queryImageStatusContainers)
	} else if !skeletonOnly {
		// if we didn't query any containers (all up to date), update timestamp
		fallbackCache.Set(FallbackCacheLastUpdatedKey, time.Now(), cache.NoExpiration)
	}

	stackSettings, err := db.GetAllStackSettings()
	if err == nil {
		for _, stackSetting := range stackSettings {
			if val, ok := stacksDto[stackSetting.StackName]; ok {
				val.Priority = stackSetting.Priority
				val.AutoStart = stackSetting.AutoStart
			}
		}
	}

	stacksDtoList := make([]types.StackDto, 0, len(stacksDto))
	for _, stack := range stacksDto {
		stacksDtoList = append(stacksDtoList, *stack)
	}

	return stacksDtoList, nil
}

func buildContainerDto(containers []map[string]interface{}) []*types.ContainerDto {
	var containersDto []*types.ContainerDto
	for _, container := range containers {
		portsData := container["Ports"].([]interface{})
		// Get unique public ports
		uniquePorts := make(map[int]int)
		for _, portData := range portsData {
			portMap := portData.(map[string]interface{})
			if publicPort, ok := portMap["PublicPort"].(float64); ok {
				if privatePort, ok := portMap["PrivatePort"].(float64); ok {
					uniquePorts[int(publicPort)] = int(privatePort)
				}
			}
		}
		outPorts := make([]string, 0, len(uniquePorts))
		for public, private := range uniquePorts {
			outPorts = append(outPorts, fmt.Sprintf("%d:%d", public, private))
		}

		networksData := container["NetworkSettings"].(map[string]interface{})["Networks"].(map[string]interface{})
		networkNames := make([]string, 0, len(networksData))
		for networkName := range networksData {
			networkNames = append(networkNames, networkName)
		}

		name := container["Names"].([]interface{})[0].(string)
		name = helper.RemoveFirstIfMatch(name, "/")
		containersDto = append(containersDto, &types.ContainerDto{
			Id:       container["Id"].(string),
			Name:     name,
			Image:    container["Image"].(string),
			UpToDate: "",
			Status:   container["State"].(string),
			Ports:    outPorts,
			Networks: networkNames,
			Labels:   container["Labels"].(map[string]interface{}),
		})
	}

	return containersDto
}

func queryContainerImageStatus(endpointId int, containersDto []*types.ContainerDto) {
	// Fetch UpToDate status for each container
	var wg sync.WaitGroup
	statusChan := make(chan struct {
		index    int
		upToDate string
		cached   bool
	}, len(containersDto))

	for i, container := range containersDto {
		wg.Add(1)
		go func(i int, container *types.ContainerDto) {
			defer wg.Done()

			cachedStatus, found := portainerCache.Get(container.Id)
			var status string
			if found {
				status = cachedStatus.(string)
				fallbackCache.Set(container.Id, status, cache.NoExpiration)
				// glg.Debugf("found cached status %s for container %s", status, container.Id)
			} else {
				liveStatus, err := GetImageStatus(endpointId, container.Id)
				if err != nil {
					glg.Warnf("Error fetching image status for container id %s", container.Id)
					portainerCache.Set(container.Id, liveStatus, time.Minute*5)
					fallbackCache.Delete(container.Id)
				} else {
					portainerCache.Set(container.Id, liveStatus, cache.DefaultExpiration)
					fallbackCache.Set(container.Id, liveStatus, cache.NoExpiration)
				}
				status = liveStatus
			}
			statusChan <- struct {
				index    int
				upToDate string
				cached   bool
			}{i, status, found}
		}(i, container)
	}

	go func() {
		wg.Wait()
		close(statusChan)
	}()

	cachedCount := 0
	uncachedCount := 0
	for status := range statusChan {
		containersDto[status.index].UpToDate = status.upToDate
		if status.cached {
			cachedCount++
		} else {
			uncachedCount++
		}
	}
	glg.Infof("cached container image statuses: %d, uncached container images statuses: %d", cachedCount, uncachedCount)
	fallbackCache.Set(FallbackCacheLastUpdatedKey, time.Now(), cache.NoExpiration)
}

func GetStackImagesStatus(stackId int) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/stacks/%d/images_status", appState.Config.PortainerUrl, stackId), nil)
	glg.Debugf("fetching images for stack %d", stackId)
	if err != nil {
		glg.Errorf("Failed to create request: %s", err)
		return "", err
	}

	req.Header.Add("X-API-Key", appState.Config.PortainerSecret)
	resp, err := client.Do(req)
	if err != nil {
		glg.Errorf("Failed to send request: %s", err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		glg.Errorf("Failed to read response: %s", err)
		return "", err
	}

	var imagesStatus map[string]interface{}
	err = json.Unmarshal(body, &imagesStatus)
	if err != nil {
		glg.Errorf("Failed to unmarshal JSON: %s", err)
		return "", err
	}
	if _, ok := imagesStatus["message"]; ok {
		errorMessage := fmt.Sprintf("%s: %s", imagesStatus["message"], imagesStatus["details"])
		return "", fmt.Errorf(errorMessage)
	}
	return imagesStatus["Status"].(string), nil
}

func GetImageStatus(endpointId int, containerId string) (string, error) {
	glg.Debugf("fetching images for container %s in endpoint %d", containerId, endpointId)
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/docker/%d/containers/%s/image_status", appState.Config.PortainerUrl, endpointId, containerId), nil)
	if err != nil {
		glg.Errorf("Failed to create request: %s", err)
		return "", err
	}

	req.Header.Add("X-API-Key", appState.Config.PortainerSecret)
	resp, err := client.Do(req)
	if err != nil {
		glg.Errorf("Failed to send request: %s", err)
		return "", err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		glg.Errorf("Failed to read response: %s", err)
		return "", err
	}

	var container map[string]interface{}
	err = json.Unmarshal(body, &container)
	if err != nil {
		glg.Errorf("Failed to unmarshal JSON: %s", err)
		return "", err
	}
	if _, ok := container["message"]; ok {
		errorMessage := fmt.Sprintf("%s: %s. %s", container["message"], containerId, container["details"])
		return "", fmt.Errorf(errorMessage)
	}
	return container["Status"].(string), nil
}

func UpdateContainer(endpointId int, containerId string, pullImage bool) (string, error) {
	client := &http.Client{}
	reqBody := []byte(fmt.Sprintf("{\"PullImage\":%t}", pullImage))

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/docker/%d/containers/%s/recreate", appState.Config.PortainerUrl, endpointId, containerId), bytes.NewBuffer(reqBody))
	if err != nil {
		glg.Errorf("Failed to create request: %s", err)
		return "", err
	}

	req.Header.Add("X-API-Key", appState.Config.PortainerSecret)
	resp, err := client.Do(req)
	if err != nil {
		glg.Errorf("Failed to send request: %s", err)
		return "", err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		glg.Errorf("Failed to read response: %s", err)
		return "", err
	}

	var container map[string]interface{}
	err = json.Unmarshal(body, &container)
	if err != nil {
		glg.Errorf("Failed to unmarshal JSON: %s", err)
		return "", err
	}
	if _, ok := container["message"]; ok {
		errorMessage := fmt.Sprintf("%s: %s. %s", container["message"], containerId, container["details"])
		glg.Error(errorMessage)
		return "", fmt.Errorf(errorMessage)
	}
	return container["Id"].(string), nil
}

func getUpdateOperationId(endpointId int, stackId int) string {
	return fmt.Sprintf("update-stack-%d-%d", endpointId, stackId)
}

// EnqueueUpdateStack enqueues a stack update operation. If the operation is already queued, it is not enqueued again
// Parameters:
//
// Query Parameters:
//   - endpointId: the id of the endpoint
//     where the stack is running
//   - stackId: the id of the stack to update
//   - prune: whether to prune the stack
//   - pullImage: whether to pull the image´
//
// Creates a StackUpdateStatus object with the following values depending on the result of the operation:
//   - Status: "queued", "done",
//     "error"
//   - Details: the error message if the operation fails
func EnqueueUpdateStack(endpointId int, stackId int, prune bool, pullImage bool) (float64, error) {
	id := getUpdateOperationId(endpointId, stackId)
	if val, ok := appState.StackUpdateQueue.Get(id); ok {
		data := val.(types.StackUpdateStatus)
		if data.Status != types.Error && data.Status != types.Done {
			glg.Infof("stack update already queued: %s", val)
			retErr := errors.New("stack update already queued")
			return -2, retErr
		}
	}
	glg.Infof("enqueueing stack id: %d, prune: %t", stackId, prune)

	stackData, err := getStackRaw(stackId)
	if err != nil {
		glg.Errorf("Failed to get stack data: %s", err)
		return -1, err
	}
	if val, ok := stackData["EndpointId"]; !ok {
		glg.Errorf("stack does not have endpoint id")
		return -1, fmt.Errorf("stack does not have endpoint id")
	} else if valInt, ok := val.(float64); !ok {
		glg.Errorf("stack endpoint id is not a number")
		return -1, fmt.Errorf("stack endpoint id is not a number")
	} else if int(valInt) != endpointId {
		glg.Errorf("stack endpoint id does not match")
		return -1, fmt.Errorf("stack endpoint id does not match")
	}

	var stackNameString string
	if stackName, ok := stackData["Name"]; !ok {
		glg.Errorf("stack does not have name data")
		return -1, fmt.Errorf("stack does not have name data")
	} else if stackNameString, ok = stackName.(string); !ok {
		glg.Errorf("stack name is not a string")
		return -1, fmt.Errorf("stack name is not a string")
	}

	stackFileContent, err := getStackFile(stackId)
	if err != nil {
		glg.Errorf("Failed to get stack file: %s", err)
		return -1, err
	}

	//
	//type RequestBody struct {
	//	Env              interface{} `json:"Env"`
	//	Id               int         `json:"id"`
	//	Prune            bool        `json:"Prune"`
	//	PullImage        bool        `json:"PullImage"`
	//	StackFileContent string      `json:"StackFileContent"`
	//	Webhook          string      `json:"Webhook"`
	//}

	envData, ok := stackData["Env"]
	if !ok {
		glg.Errorf("stack does not have env data")
		return -1, fmt.Errorf("stack does not have env data")
	}
	webhook, ok := stackData["Webhook"]
	if !ok {
		glg.Errorf("stack does not have webhook data")
		return -1, fmt.Errorf("stack does not have webhook data")
	}
	envDataByte, err := json.Marshal(envData)
	if err != nil {
		glg.Errorf("Failed to marshal env data: %s", err)
		return -1, err
	}

	envDataString := string(envDataByte)
	reqBodyRaw := fmt.Sprintf(`{"Env":%s,"id":%d,"Prune":%t,"PullImage":%t,"StackFileContent":"%s","Webhook":"%s"}`,
		envDataString, stackId, prune, pullImage, stackFileContent, webhook)
	//glg.Logf("%+v", reqBodyRaw)
	reqBodyByte := []byte(reqBodyRaw)

	go func() {
		updateStatus := types.StackUpdateStatus{
			EndpointId: endpointId,
			StackId:    stackId,
			StackName:  stackNameString,
			Status:     types.Queued,
			Timestamp:  int64(time.Now().Unix()),
			Details:    "",
		}
		appState.StackUpdateQueue.Set(id, updateStatus, time.Minute*30)
		_, err := updateStack(endpointId, stackId, reqBodyByte)
		if err != nil {
			glg.Errorf("No operation performed: %s", err)
			updateStatus.Status = types.Error
			updateStatus.Details = err.Error()
		} else {
			updateStatus.Status = types.Done
		}
		updateStatus.Timestamp = int64(time.Now().Unix())
		appState.StackUpdateQueue.Set(id, updateStatus, time.Hour*24*7)
	}()

	return float64(stackId), nil
}

func updateStack(endpointId int, stackId int, reqBodyByte []byte) (float64, error) {
	client := &http.Client{}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/stacks/%d", appState.Config.PortainerUrl, stackId), bytes.NewBuffer(reqBodyByte))
	if err != nil {
		glg.Errorf("Failed to create request: %s", err)
		return -1, err
	}

	q := req.URL.Query()
	q.Add("endpointId", fmt.Sprintf("%d", endpointId))
	req.URL.RawQuery = q.Encode()

	req.Header.Add("X-API-Key", appState.Config.PortainerSecret)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		glg.Errorf("Failed to send request: %s", err)
		return -1, err
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		glg.Errorf("Failed to read response: %s", err)
		return -1, err
	}

	var stack map[string]interface{}
	err = json.Unmarshal(respBody, &stack)
	if err != nil {
		glg.Errorf("Failed to unmarshal JSON: %s", err)
		return -1, err
	}
	if _, ok := stack["message"]; ok {
		errorMessage := fmt.Sprintf("%s: %d. %s", stack["message"], stackId, stack["details"])
		glg.Error(errorMessage)
		return -1, fmt.Errorf(errorMessage)
	}
	glg.Infof("Stack %s updated", stack["Name"])
	// remove cached images status when an update was performed
	portainerCache.Delete(fmt.Sprintf("stack-%d-images-status", stackId))
	return stack["Id"].(float64), nil
}

func getStackRaw(stackId int) (map[string]interface{}, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/stacks/%d", appState.Config.PortainerUrl, stackId), nil)
	if err != nil {
		glg.Errorf("Failed to create request: %s", err)
		return nil, err
	}

	req.Header.Add("X-API-Key", appState.Config.PortainerSecret)
	resp, err := client.Do(req)
	if err != nil {
		glg.Errorf("Failed to send request: %s", err)
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		glg.Errorf("Failed to read response: %s", err)
		return nil, err
	}

	var stack map[string]interface{}
	err = json.Unmarshal(body, &stack)
	if err != nil {
		glg.Errorf("Failed to unmarshal JSON: %s", err)
		return nil, err
	}
	return stack, nil
}

func getStackFile(stackId int) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/stacks/%d/file", appState.Config.PortainerUrl, stackId), nil)
	if err != nil {
		glg.Errorf("Failed to create request: %s", err)
		return "", err
	}

	req.Header.Add("X-API-Key", appState.Config.PortainerSecret)
	resp, err := client.Do(req)
	if err != nil {
		glg.Errorf("Failed to send request: %s", err)
		return "", err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		glg.Errorf("Failed to read response: %s", err)
		return "", err
	}

	var stackFileContent map[string]string
	err = json.Unmarshal(body, &stackFileContent)
	if err != nil {
		glg.Errorf("Failed to unmarshal JSON: %s", err)
		return "", err
	}

	stackFileContentRaw := stackFileContent["StackFileContent"]
	stackFileContentRaw = strings.ReplaceAll(stackFileContentRaw, "\\", "\\\\")
	stackFileContentRaw = strings.ReplaceAll(stackFileContentRaw, "\n", "\\n")
	stackFileContentRaw = strings.ReplaceAll(stackFileContentRaw, "\"", "\\\"")

	return stackFileContentRaw, nil
}
