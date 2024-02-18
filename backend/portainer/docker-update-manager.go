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

	"washboard/helper"

	"github.com/kpango/glg"
	"github.com/patrickmn/go-cache"
)

type EndpointDto struct {
	Id         int            `json:"id"`
	Name       string         `json:"name"`
	Containers []ContainerDto `json:"containers"`
}

type ContainerDto struct {
	Id       string                 `json:"id"`
	Name     string                 `json:"name"`
	Image    string                 `json:"image"`
	UpToDate string                 `json:"upToDate"`
	Status   string                 `json:"status"`
	Ports    []int                  `json:"ports"`
	Labels   map[string]interface{} `json:"labels"`
}

type StackDto struct {
	Id         int            `json:"id"`
	Name       string         `json:"name"`
	Containers []ContainerDto `json:"containers"`
}

type StackUpdateStatus struct {
	EndpointId int    `json:"endpointId"`
	StackId    int    `json:"stackId"`
	Status     string `json:"status"`
	Details    string `json:"details"`
}

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

	var endpoints []EndpointDto
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

// GetStacks returns the stacks for the given endpoint
func GetStacks(endpointId int) ([]StackDto, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/stacks", appState.Config.PortainerUrl), nil)
	if err != nil {
		glg.Errorf("Failed to create request: %s", err)
		return nil, err
	}

	q := req.URL.Query()
	q.Add("filters", fmt.Sprintf("{\"EndpointId\":%d}", endpointId))
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
	var stacksDict = make(map[string]map[string]interface{})
	for _, stack := range stacks {
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

	stacksDto, err := buildStackDto(stacksDict, endpointId)

	return stacksDto, err
}

// GetContainers returns the containers for the given endpoint. If stackLabel is provided, only the containers of the stack with the given label are returned, otherwise all containers are returned
func GetContainers(endpointId int, stackLabel string) ([]ContainerDto, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/endpoints/%d/docker/containers/json", appState.Config.PortainerUrl, endpointId), nil)
	if err != nil {
		glg.Errorf("Failed to create request: %s", err)
		return nil, err
	}

	q := req.URL.Query()
	q.Add("all", "true")
	if stackLabel != "" {
		q.Add("filters", fmt.Sprintf("{\"label\":[\"com.docker.compose.project=%s\"]}", stackLabel))
	}
	//?filter="{"label":["com.docker.compose.project=stackName"]}"
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

	containersDto := buildStackContainerDto(containers, endpointId)

	return containersDto, nil
}

func buildStackDto(stacks map[string]map[string]interface{}, endpointId int) ([]StackDto, error) {
	var stacksDto = make(map[string]*StackDto)
	containers, err := GetContainers(endpointId, "")
	if err != nil {
		glg.Errorf("Failed to get stack containers: %s", err)
		return nil, err
	}

	for _, container := range containers {
		var label string

		if labelRaw, ok := container.Labels["com.docker.compose.project"]; !ok {
			continue
		} else if labelParsed, ok := labelRaw.(string); !ok {
			glg.Warnf("label %s is not a string", labelRaw)
			continue
		} else {
			label = labelParsed
		}
		if val, ok := stacksDto[label]; ok {
			val.Containers = append(val.Containers, container)
		} else {
			if val, ok := stacks[label]; ok {
				stacksDto[label] = &StackDto{
					Id:         int(val["Id"].(float64)),
					Name:       val["Name"].(string),
					Containers: []ContainerDto{container},
				}
			}
		}
	}

	stacksDtoList := make([]StackDto, 0, len(stacksDto))
	for _, stack := range stacksDto {
		stacksDtoList = append(stacksDtoList, *stack)
	}

	return stacksDtoList, nil
}

func buildStackContainerDto(containers []map[string]interface{}, endpointId int) []ContainerDto {
	var containersDto []ContainerDto
	for _, container := range containers {
		portsData := container["Ports"].([]interface{})
		// Get unique public ports
		uniquePorts := make(map[int]struct{})
		for _, portData := range portsData {
			portMap := portData.(map[string]interface{})
			if publicPort, ok := portMap["PublicPort"].(float64); ok {
				uniquePorts[int(publicPort)] = struct{}{}
			}
		}
		publicPorts := make([]int, 0, len(uniquePorts))
		for port := range uniquePorts {
			publicPorts = append(publicPorts, port)
		}
		name := container["Names"].([]interface{})[0].(string)
		name = helper.RemoveFirstIfMatch(name, "/")
		containersDto = append(containersDto, ContainerDto{
			Id:       container["Id"].(string),
			Name:     name,
			Image:    container["Image"].(string),
			UpToDate: "",
			Status:   container["State"].(string),
			Ports:    publicPorts,
			Labels:   container["Labels"].(map[string]interface{}),
		})
	}

	// Fetch UpToDate status for each container
	var wg sync.WaitGroup
	statusChan := make(chan struct {
		index    int
		upToDate string
		cached   bool
	}, len(containersDto))

	for i, container := range containersDto {
		wg.Add(1)
		go func(i int, container ContainerDto) {
			defer wg.Done()

			cachedStatus, found := portainerCache.Get(container.Id)
			var status string
			if found {
				status = cachedStatus.(string)
				// glg.Debugf("found cached status %s for container %s", status, container.Id)
			} else {
				liveStatus, err := GetImageStatus(endpointId, container.Id)
				if err != nil {
					glg.Errorf("Error fetching image status for container id %s", container.Id)
					portainerCache.Set(container.Id, liveStatus, time.Minute*5)
					return
				}
				portainerCache.Set(container.Id, liveStatus, cache.DefaultExpiration)
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
	glg.Logf("cached images: %d, uncached images: %d", cachedCount, uncachedCount)
	return containersDto
}

func GetImageStatus(endpointId int, containerId string) (string, error) {
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
// - endpointId: the id of the endpoint
//   where the stack is running
// - stackId: the id of the stack to update
// - prune: whether to prune the stack
// - pullImage: whether to pull the image´
//
// Creates a StackUpdateStatus object with the following values depending on the result of the operation:
// - Status: "queued", "done",
//   "error"
// - Details: the error message if the operation fails
//
func EnqueueUpdateStack(endpointId int, stackId int, prune bool, pullImage bool) (float64, error) {
	id := getUpdateOperationId(endpointId, stackId)
	if val, ok := appState.StackUpdateQueue.Get(id); ok {
		data := val.(StackUpdateStatus)
		if data.Status != "error" && data.Status != "done"{
			glg.Infof("stack update already queued: %s", val)
			retErr := errors.New("stack update already queued")
			return -1, retErr
		}
	}
	glg.Infof("enqueueing stack id: %d", stackId)

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
		updateStatus := StackUpdateStatus{
			EndpointId: endpointId,
			StackId:    stackId,
			Status:     "queued",
			Details:    "",
		}
		appState.StackUpdateQueue.Set(id, updateStatus, time.Minute*30)
		_, err := updateStack(endpointId, stackId, reqBodyByte)
		if err != nil {
			glg.Errorf("Failed to update stack: %s", err)
			updateStatus.Status = "error"
			updateStatus.Details = err.Error()
		} else {
			updateStatus.Status = "done"
		}
		appState.StackUpdateQueue.Set(id, updateStatus, cache.DefaultExpiration)
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
	stackFileContentRaw = strings.ReplaceAll(stackFileContentRaw, "\n", "\\n")
	stackFileContentRaw = strings.ReplaceAll(stackFileContentRaw, "\"", "\\\"")
	return stackFileContentRaw, nil
}
