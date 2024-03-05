package api

import (
	"fmt"
	"net/http"

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
	var stackId int

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

	if stackIdRaw, ok := reqBody["stackId"]; !ok {
		glg.Warn("id field is missing")
		c.JSON(http.StatusBadRequest, gin.H{"message": "id field is missing"})
		return
	} else if idFloat, ok := stackIdRaw.(float64); !ok {
		glg.Warn("id field is not a string")
		c.JSON(http.StatusBadRequest, gin.H{"message": "id field is not a number"})
		return
	} else {
		stackId = int(idFloat)
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


func PortainerStartContainer(c *gin.Context) {
	PortainerManageContainer(c, types.Start)
}

func PortainerStopContainer(c *gin.Context) {
	PortainerManageContainer(c, types.Stop)
}

func PortainerRestartContainer(c *gin.Context) {
	PortainerManageContainer(c, types.Restart)
}

func PortainerKillContainer(c *gin.Context) {
	PortainerManageContainer(c, types.Kill)
}

func PortainerPauseContainer(c *gin.Context) {
	PortainerManageContainer(c, types.Pause)
}

func PortainerResumeContainer(c *gin.Context) {
	PortainerManageContainer(c, types.Resume)
}


func PortainerManageContainer(c *gin.Context, action types.ActionType) {
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
	var containerId string
	// var action portainer.ActionType


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

	if containerRaw, ok := reqBody["containerId"]; !ok {
		glg.Warn("containerId field is missing")
		c.JSON(http.StatusBadRequest, gin.H{"message": "containerId field is missing"})
		return
	} else if containerId, ok = containerRaw.(string); !ok {
		glg.Warn("containerId field is not a string")
		c.JSON(http.StatusBadRequest, gin.H{"message": "containerId field is not a string"})
		return
	}

	// if actionRaw, ok := reqBody["action"]; !ok {
	// 	glg.Warn("action field is missing")
	// 	c.JSON(http.StatusBadRequest, gin.H{"message": "action field is missing"})
	// 	return
	// } else if actionStr, ok := actionRaw.(string); !ok {
	// 	glg.Warn("action field is not a string")
	// 	c.JSON(http.StatusBadRequest, gin.H{"message": "action field is not a string"})
	// 	return
	// } else {
	// 	action = portainer.ActionType(actionStr)
	// }

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
