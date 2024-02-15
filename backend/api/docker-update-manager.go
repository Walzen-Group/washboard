package api

import (
	"net/http"
	"strconv"

	"github.com/Walzen-Group/washboard/backend/portainer"
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
	endpoint := c.DefaultQuery("endpointId", "1")
	endpointId, err := strconv.Atoi(endpoint)
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
