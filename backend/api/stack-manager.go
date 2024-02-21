package api

import (
	"fmt"
	"net/http"

	"washboard/portainer"

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
