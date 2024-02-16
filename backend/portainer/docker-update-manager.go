package portainer

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/kpango/glg"
)

type Endpoint struct {
	Id int `json:"Id"`
	Name string `json:"Name"`
}

type Container struct {
	Id string `json:"Id"`
	Name string `json:"Name"`
	Image string `json:"Image"`
	UpToDate bool `json:"UpToDate"`
	Status string `json:"Status"`
	Ports []int `json:"Ports"`
	Labels map[string]interface{} `json:"Labels"`
}

type Stack struct {
	Id int `json:"Id"`
	Name string `json:"Name"`
	Containers []Container `json:"Containers"`
}

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

func GetStacks(endpointId int) ([]map[string]interface{}, error) {
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

	return stacks, nil
}

func GetStackContainers(endpointId int, stackLabel string) ([]Container, error){
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
		containersDto = append(containersDto, Container{
			Id: container["Id"].(string),
			Name: container["Names"].([]interface{})[0].(string),
			Image: container["Image"].(string),
			UpToDate: false,
			Status: container["State"].(string),
			Ports: publicPorts,
			Labels: container["Labels"].(map[string]interface{}),
		})
	}

	// Fetch UpToDate status for each container
	var wg sync.WaitGroup
	statusChan := make(chan struct {
		index int
		upToDate bool
	}, len(containersDto))

	for i, container := range containersDto {
		wg.Add(1)
		go func(i int, container Container) {
			defer wg.Done()
			status, err := GetImageStatus(endpointId, container.Id)
            if err != nil {
                glg.Errorf("Error fetching UpToDate status for container %s: %v\n", container.Id, err)
                return
            }
			statusChan <- struct {
				index int
				upToDate bool
			}{i, status == "updated"}
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