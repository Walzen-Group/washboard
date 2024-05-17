package api

import (
	"errors"
	"fmt"
	"net/http"
	"time"
	"washboard/db"
	"washboard/portainer"
	"washboard/types"
	"washboard/werrors"

	"github.com/gin-gonic/gin"
	"github.com/kpango/glg"
)

func SyncWithPortainer(c *gin.Context) {
	if _, ok := appState.StateQueue.Get("sync"); ok {
		handleError(c, fmt.Errorf("Sync already in progress"), "Sync already in progress", http.StatusConflict)
		return
	}
	appState.StateQueue.Add("sync", "inprog", time.Minute*2)
	defer appState.StateQueue.Delete("sync")

	glg.Infof("starting setting sync portainer->washbdb")
	var syncOptions *types.SyncOptions = &types.SyncOptions{}
	if err := c.ShouldBindJSON(&syncOptions); err != nil {
		handleError(c, err, "Failed to bind json. Check the request body and ensure that the correct fields are present.", http.StatusBadRequest)
		return
	}

	err := portainer.PerformSync(syncOptions)
	if err != nil {
		handleError(c, err, err.Error(), http.StatusInternalServerError)
		return
	}

	glg.Infof("sync completed")
}

func CreateStackSettings(c *gin.Context) {
	var stackSettings *types.StackSettings = &types.StackSettings{}
	if err := c.ShouldBindJSON(&stackSettings); err != nil {
		errorMessage := "Failed to bind json. Check the request body and ensure that the correct fields are present."
		handleError(c, err, errorMessage, http.StatusBadRequest)
		return
	}
	glg.Infof("Creating stack settings: %+v", stackSettings)
	err := db.CreateStackSettings(stackSettings)
	if err != nil {
		target := &werrors.CannotInsertError{}
		if errors.As(err, &target) {
			glg.Errorf("%s %s", err.(*werrors.CannotInsertError), err)
			c.JSON(http.StatusNotAcceptable, gin.H{
				"message": "There was an issue with the insert operation",
				"error":   err,
			})
			return
		}
		handleError(c, err, "Failed to create stack settings", http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message":       "Stack settings created successfully.",
		"stackSettings": stackSettings,
	})
}

func GetStackSettings(c *gin.Context) {
	name := c.Param("name")

	if name == "" {
		stacks, err := db.GetAllStackSettings()
		if err != nil {
			handleError(c, err, "Failed to get stack settings", http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message":       "ok",
			"stackSettings": stacks,
		})
		return
	}
	stack, err := db.GetStackSettings(name)
	target := &werrors.DoesNotExistError{}
	if errors.As(err, &target) {
		handleError(c, err, "No result", http.StatusNotFound)
		return
	} else if err != nil {
		handleError(c, err, "Failed to get stack settings", http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":       "ok",
		"stackSettings": stack,
	})
}

func UpdateStackSettings(c *gin.Context) {
	name := c.Param("name")
	updatePriority := c.DefaultQuery("updatePrio", "false")



	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "stack name is required",
		})
		return
	}

	var stackSettings *types.StackSettings = &types.StackSettings{}
	if err := c.ShouldBindJSON(&stackSettings); err != nil {
		errorMessage := "Failed to bind json. Check the request body and ensure that the correct fields are present."
		handleError(c, err, errorMessage, http.StatusBadRequest)
		return
	}
	var err error
	if updatePriority == "true" {
		err = db.UpdateStackPriority(stackSettings)
	} else {
		err = db.UpdateStackSettings(stackSettings, name)
	}


	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to update stack settings",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":       "Stack settings updated successfully.",
		"stackSettings": stackSettings,
	})

}

func DeleteStackSettings(c *gin.Context) {
	name := c.Param("name")

	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "stack name is required",
		})
		return
	}

	err := db.DeleteStackSettings(name)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to delete stack settings",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Stack settings deleted successfully.",
	})
}
