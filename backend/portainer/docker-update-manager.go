package portainer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	"washboard/helper"

	"github.com/kpango/glg"
	"github.com/patrickmn/go-cache"
)

type Endpoint struct {
	Id         int         `json:"id"`
	Name       string      `json:"name"`
	Containers []Container `json:"containers"`
}

type Container struct {
	Id       string                 `json:"id"`
	Name     string                 `json:"name"`
	Image    string                 `json:"image"`
	UpToDate string                 `json:"upToDate"`
	Status   string                 `json:"status"`
	Ports    []int                  `json:"ports"`
	Labels   map[string]interface{} `json:"labels"`
}

type Stack struct {
	Id         int         `json:"id"`
	Name       string      `json:"name"`
	Containers []Container `json:"containers"`
}


// GetEndpointId returns the id of the endpoint with the given name, which is also the environment in Portainer
func GetEndpointId(endpointName string) (int, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/endpoints", appState.PortainerUrl), nil)
	if err != nil {
		glg.Errorf("Failed to create request: %s", err)
		return -1, err
	}

	req.Header.Add("X-API-Key", appState.PortainerSecret)
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

	var endpoints []Endpoint
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
func GetStacks(endpointId int) ([]Stack, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/stacks", appState.PortainerUrl), nil)
	if err != nil {
		glg.Errorf("Failed to create request: %s", err)
		return nil, err
	}

	q := req.URL.Query()
	q.Add("filters", fmt.Sprintf("{\"EndpointId\":%d}", endpointId))
	req.URL.RawQuery = q.Encode()
	req.Header.Add("X-API-Key", appState.PortainerSecret)

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
func GetContainers(endpointId int, stackLabel string) ([]Container, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/endpoints/%d/docker/containers/json", appState.PortainerUrl, endpointId), nil)
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
	req.Header.Add("X-API-Key", appState.PortainerSecret)

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

func buildStackDto(stacks map[string]map[string]interface{}, endpointId int) ([]Stack, error) {
	var stacksDto = make(map[string]*Stack)
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
				stacksDto[label] = &Stack{
					Id:         int(val["Id"].(float64)),
					Name:       val["Name"].(string),
					Containers: []Container{container},
				}
			}
		}
	}

	stacksDtoList := make([]Stack, 0, len(stacksDto))
	for _, stack := range stacksDto {
		stacksDtoList = append(stacksDtoList, *stack)
	}

	return stacksDtoList, nil
}

func buildStackContainerDto(containers []map[string]interface{}, endpointId int) []Container {
	var containersDto []Container
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
		containersDto = append(containersDto, Container{
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
	}, len(containersDto))

	for i, container := range containersDto {
		wg.Add(1)
		go func(i int, container Container) {
			defer wg.Done()


			cachedStatus, found := portainerCache.Get(container.Id)
			var status string
			if found {
				status = cachedStatus.(string)
				// glg.Debugf("found cached status %s for container %s", status, container.Id)
			} else {
				liveStatus, err := GetImageStatus(endpointId, container.Id)
				if err != nil {
					glg.Warnf("Error fetching UpToDate status for container %s: %v\n", container.Id, err)
					return
				}
				portainerCache.Set(container.Id, liveStatus, cache.DefaultExpiration)
				status = liveStatus
			}
			statusChan <- struct {
				index    int
				upToDate string
			}{i, status}
		}(i, container)
	}

	go func() {
		wg.Wait()
		close(statusChan)
	}()

	for status := range statusChan {
		containersDto[status.index].UpToDate = status.upToDate
	}
	return containersDto
}

func GetImageStatus(endpointId int, containerId string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/docker/%d/containers/%s/image_status", appState.PortainerUrl, endpointId, containerId), nil)
	if err != nil {
		glg.Errorf("Failed to create request: %s", err)
		return "", err
	}

	req.Header.Add("X-API-Key", appState.PortainerSecret)
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
	return container["Status"].(string), nil
}

func UpdateContainer(endpointId int, containerId string, pullImage bool) (string, error) {
	client := &http.Client{}
	reqBody := []byte(fmt.Sprintf("{\"PullImage\":%t}", pullImage))

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/docker/%d/containers/%s/recreate", appState.PortainerUrl, endpointId, containerId), bytes.NewBuffer(reqBody))
	if err != nil {
		glg.Errorf("Failed to create request: %s", err)
		return "", err
	}

	req.Header.Add("X-API-Key", appState.PortainerSecret)
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

func UpdateStack(endpointId int, stackId int, prune bool, pullImage bool) (float64, error) {
	client := &http.Client{}
	stackFileContent, err := getStackFile(stackId)
	if err != nil {
		glg.Errorf("Failed to get stack file: %s", err)
		return -1, err
	}
	type RequestBody struct {
		Env             []string `json:"Env"`
		Id              int      `json:"id"`
		Prune           bool     `json:"Prune"`
		PullImage       bool     `json:"PullImage"`
		StackFileContent string   `json:"StackFileContent"`
		Webhook         string   `json:"Webhook"`
	}

	reqBodyRaw := fmt.Sprintf(`{"Env":[],"id":%d,"Prune":%t,"PullImage":%t,"StackFileContent":"%s","Webhook":null}`, stackId, prune, pullImage, stackFileContent)
	reqBodyByte := []byte(reqBodyRaw)

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/stacks/%d", appState.PortainerUrl, stackId), bytes.NewBuffer(reqBodyByte))
	if err != nil {
		glg.Errorf("Failed to create request: %s", err)
		return -1, err
	}

	q := req.URL.Query()
	q.Add("endpointId", fmt.Sprintf("%d", endpointId))
	req.URL.RawQuery = q.Encode()

	req.Header.Add("X-API-Key", appState.PortainerSecret)
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

func getStackFile(stackId int) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/stacks/%d/file", appState.PortainerUrl, stackId), nil)
	if err != nil {
		glg.Errorf("Failed to create request: %s", err)
		return "", err
	}

	req.Header.Add("X-API-Key", appState.PortainerSecret)
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
