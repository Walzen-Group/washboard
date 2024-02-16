package api

import (
	"fmt"
	"net/http"
	"strconv"

	"washboard/portainer"

	"github.com/gin-gonic/gin"
	"github.com/kpango/glg"
)


func PortainerGetEndpoints(c *gin.Context) {
	endpointName := c.DefaultQuery("endpoint", "Quasar")
	res, err := portainer.GetEndpointId(endpointName)
	if err != nil {
		glg.Error("failed to get endpoints")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to get endpoints",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"endpoint": res,
	})
}

func PortainerGetStacks(c *gin.Context) {
	// Set endpointId
	endpoint := c.DefaultQuery("endpointId", "1")
	endpointId, err := strconv.Atoi(endpoint)
	if err != nil {
		glg.Errorf("failed to convert endpointId to int: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("failed to convert endpointId \"%s\" to int", endpoint),
			"error": err,
		})
		return
	}
	
	res, err := portainer.GetStacks(endpointId)
	if err != nil {
		glg.Error("failed to get stacks")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to get stacks",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"stacks": res,
	})
}

func PortainerGetStackContainers(c *gin.Context) {
	// Set endpointId
	endpoint := c.DefaultQuery("endpointId", "1")
	endpointId, err := strconv.Atoi(endpoint)
	if err != nil {
		glg.Errorf("failed to convert endpointId to int: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("failed to convert endpointId \"%s\" to int", endpoint),
			"error": err,
		})
		return
	}
	
	// Set stackLabel
	stackLabel := c.Query("stackLabel")
	res, err := portainer.GetStackContainers(endpointId, stackLabel)
	if err != nil {
		glg.Error("failed to get stacks")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to get stacks",
		})
		return
	}

	// Check if stack is empty
	if len(res) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": fmt.Sprintf("Stack \"%s\" not found in environment %d", stackLabel, endpointId),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"stacks": res,
	})
}

func PortainerGetImageStatus(c *gin.Context) {
	// Set endpointId
	endpoint := c.DefaultQuery("endpointId", "1")
	endpointId, err := strconv.Atoi(endpoint)
	if err != nil {
		glg.Errorf("failed to convert endpointId to int: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to convert endpointId to int",
			"error": err,
		})
		return
	}

	// Set containerId
	containerId := c.Query("containerId")
	if containerId == "" {
		glg.Error("containerId is empty")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "containerId is empty",
		})
		return
	}

	res, err := portainer.GetImageStatus(endpointId, containerId)
	if err != nil {
		glg.Errorf("failed to get image status: %s", err)
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status": res,
	})
}