package api

import (
	"washboard/state"

	"github.com/gin-gonic/gin"
	"github.com/kpango/glg"
)

var appState *state.Data = state.Instance()

func handleError(c *gin.Context, err error, context string, statusCode int) {
	glg.Errorf("%s %s", context, err)
	c.JSON(statusCode, gin.H{
		"message": context,
		"error":   err.Error(),
	})
}
