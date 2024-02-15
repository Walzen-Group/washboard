package portainer

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/kpango/glg"
)

type Endpoint struct {
	Id int `json:"Id"`
	Name string `json:"Name"`
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

func GetStackContainers(endpointId int, stackLabel string) ([]map[string]interface{}, error){
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

	return containers, nil
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
	
	//TODO: Check if body is empty ([]), meaning wrong containerId

	var container map[string]interface{}
	err = json.Unmarshal(body, &container)
	if err != nil {
		glg.Errorf("Failed to unmarshal JSON: %s", err)
		return "", err
	}

	return container["Status"].(string), nil
}