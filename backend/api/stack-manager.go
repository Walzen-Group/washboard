package api

import (
	"fmt"
	"net/http"
	"strconv"

	"washboard/control"
	"washboard/portainer"
	"washboard/types"

	"github.com/gin-gonic/gin"
	"github.com/kpango/glg"
)

func PortainerStartStack(c *gin.Context) {
	portainerStartOrStopStack(c, "start")
}

func PortainerStopStack(c *gin.Context) {
	portainerStartOrStopStack(c, "stop")
}

func portainerStartOrStopStack(c *gin.Context, startOrStop string) {
	stackIdStr := c.Param("id")

	if stackIdStr == "" {
		glg.Warn("stackId in path is missing")
		c.JSON(http.StatusBadRequest, gin.H{"message": "stackId in path is missing"})
		return
	}

	stackId, err := strconv.Atoi(stackIdStr)
	if err != nil {
		glg.Warn("stackId in path is not an int")
		c.JSON(http.StatusBadRequest, gin.H{"message": "stackId in path is not an int"})
		return
	}

	var reqBody map[string]interface{}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		errorMessage := "Failed to bind json. Check the request body and ensure that the pullImage field is present."
		glg.Errorf("%s %s", errorMessage, err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": errorMessage,
			"error": err,
		})
		return
	}


	var endpointId int

	if endpointRaw, ok := reqBody["endpointId"]; !ok {
		glg.Warn("endpointId field is missing")
		c.JSON(http.StatusBadRequest, gin.H{"message": "endpointId field is missing"})
		return
	} else if endpointIdFloat, ok := endpointRaw.(float64); !ok {
		glg.Warn("endpointId field is not an int")
		c.JSON(http.StatusBadRequest, gin.H{"message": "endpointId field is not an int"})
		return
	} else {
		endpointId = int(endpointIdFloat)
	}


	stackName, status, err := portainer.StartOrStopStack(endpointId, stackId, startOrStop)
	if err != nil {
		glg.Errorf("Failed to %s stack: %s", startOrStop, err)
		c.JSON(status, gin.H{
			"message":  fmt.Sprintf("Failed to %s stack", startOrStop),
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"name": stackName,
	})
}

func PortainerContainerAction(c *gin.Context) {
	containerId := c.Param("containerId")
	action := types.ContainerAction(c.Param("action"))

	if containerId == "" {
		glg.Warn("containerId field is missing")
		c.JSON(http.StatusBadRequest, gin.H{"message": "containerId field is missing"})
		return
	}

	if action == "" {
		glg.Warn("action field is missing")
		c.JSON(http.StatusBadRequest, gin.H{"message": "action field is missing"})
		return
	}

	var reqBody map[string]interface{}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		errorMessage := "Failed to bind json. Check the request body."
		glg.Errorf("%s %s", errorMessage, err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": errorMessage,
			"error": err,
		})
		return
	}


	var endpointId int
	if endpointRaw, ok := reqBody["endpointId"]; !ok {
		glg.Warn("endpointId field is missing")
		c.JSON(http.StatusBadRequest, gin.H{"message": "endpointId field is missing"})
		return
	} else if endpointIdFloat, ok := endpointRaw.(float64); !ok {
		glg.Warn("endpointId field is not an int")
		c.JSON(http.StatusBadRequest, gin.H{"message": "endpointId field is not an int"})
		return
	} else {
		endpointId = int(endpointIdFloat)
	}

	res, err := portainer.ManageContainer(endpointId, containerId, action)

	if err != nil {
		glg.Errorf("Failed to manage container: %s", err)
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Failed to manage container",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": res,
	})

}

func SyncAutoStartState(c *gin.Context) {
	var endpointId int

	var reqBody map[string]interface{}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		errorMessage := "Failed to bind json. Check the request body."
		glg.Errorf("%s %s", errorMessage, err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": errorMessage,
			"error": err,
		})
		return
	}

	if endpointRaw, ok := reqBody["endpointId"]; !ok {
		glg.Warn("endpointId field is missing")
		c.JSON(http.StatusBadRequest, gin.H{"message": "endpointId field is missing"})
		return
	} else if endpointIdFloat, ok := endpointRaw.(float64); !ok {
		glg.Warn("endpointId field is not an int")
		c.JSON(http.StatusBadRequest, gin.H{"message": "endpointId field is not an int"})
		return
	} else {
		endpointId = int(endpointIdFloat)
	}

	if err := control.SyncAutoStartState(endpointId); err != nil {
		glg.Errorf("Failed to sync auto start state: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to sync auto start state",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Auto start state synced successfully",
	})
}

func StopAllStacks(c *gin.Context) {
	var endpointId int

	var reqBody map[string]interface{}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		errorMessage := "Failed to bind json. Check the request body."
		glg.Errorf("%s %s", errorMessage, err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": errorMessage,
			"error": err,
		})
		return
	}

	if endpointRaw, ok := reqBody["endpointId"]; !ok {
		glg.Warn("endpointId field is missing")
		c.JSON(http.StatusBadRequest, gin.H{"message": "endpointId field is missing"})
		return
	} else if endpointIdFloat, ok := endpointRaw.(float64); !ok {
		glg.Warn("endpointId field is not an int")
		c.JSON(http.StatusBadRequest, gin.H{"message": "endpointId field is not an int"})
		return
	} else {
		endpointId = int(endpointIdFloat)
	}

	if err := control.StopAllStacks(endpointId); err != nil {
		glg.Errorf("Failed to sync auto start state: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to sync auto start state",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Auto start state synced successfully",
	})
}
