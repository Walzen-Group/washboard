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

	var containers []*types.ContainerDto

	for _, endpoint := range syncOptions.EndpointIds {
		tmp, err := portainer.GetContainers(endpoint, "")

		if err != nil {
			handleError(c, err, fmt.Sprintf("Failed to get containers for endpoint %d", endpoint), http.StatusBadRequest)
			return
		}
		containers = append(containers, tmp...)
	}

	collectedStacks := make(map[string]*types.StackSettings)
	targetError := &werrors.DoesNotExistError{}

	// add missing grups and stakcs (das bleibt so)
	for _, container := range containers {
		if stackName, ok := container.Labels[types.StackLabel]; ok {
			stack, err := db.GetStackSettings(stackName.(string))
			if errors.As(err, &targetError) {
				stack = &types.StackSettings{StackName: stackName.(string)}
				glg.Infof("adding missing stack %s", stackName.(string))
				db.CreateStackSettings(stack)
			} else if err != nil {
				handleError(c, err, fmt.Sprintf("Failed to sync settings with db for: %s", stackName), http.StatusInternalServerError)
				return
			}
			collectedStacks[stackName.(string)] = stack
		}
	}

	allStackSettings, err2 := db.GetAllStackSettings()
	if err2 != nil {
		handleError(c, err2, "Failed to get all stack settings", http.StatusInternalServerError)
		return
	}

	stacksToRemove := make([]string, 0)

	for _, stackSettings := range allStackSettings {
		if _, ok := collectedStacks[stackSettings.StackName]; !ok {
			stacksToRemove = append(stacksToRemove, stackSettings.StackName)
		}
	}

	for _, stack := range stacksToRemove {
		glg.Infof("removing orphaned stack %s", stack)
		err := db.DeleteStackSettings(stack)
		if err != nil {
			glg.Errorf("Failed to delete orphaned stack %s", stack)
		}
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

	err := db.UpdateStackSettings(stackSettings, name)

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

func CreateGroupSettings(c *gin.Context) {
	var groupSettings *types.GroupSettings = &types.GroupSettings{}
	if err := c.ShouldBindJSON(&groupSettings); err != nil {
		errorMessage := "Failed to bind json. Check the request body and ensure that the correct fields are present."
		glg.Errorf("%s %s", errorMessage, err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": errorMessage,
			"error":   err,
		})
		return
	}

	if groupSettings.GroupName == "" {
		errorMessage := "Name is required"
		glg.Errorf("%s", errorMessage)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "notok",
			"error":   errorMessage,
		})
		return
	}

	glg.Infof("Creating group settings: %+v", groupSettings)
	err := db.CreateGroupSettings(groupSettings)
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
		errorMessage := "Failed to create group settings."
		glg.Errorf("%s %s", errorMessage, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": errorMessage,
			"error":   err,
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message":       "Group settings created successfully.",
		"groupSettings": groupSettings,
	})
}

func GetGroupSettings(c *gin.Context) {
	name := c.Param("name")

	if name == "" {
		stacks, err := db.GetAllGroupSettings()
		if err != nil {
			handleError(c, err, "Failed to get group settings", http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message":       "ok",
			"stackSettings": stacks,
		})
		return
	}

	group, err := db.GetGroupSettings(name)
	target := &werrors.DoesNotExistError{}
	if errors.As(err, &target) {
		handleError(c, err, "No result", http.StatusNotFound)
		return
	} else if err != nil {
		handleError(c, err, "Failed to get group settings", http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":       "ok",
		"groupSettings": group,
	})
}

func UpdateGroupSettings(c *gin.Context) {
	name := c.Param("name")

	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "name is required",
		})
		return
	}

	var groupSettings *types.GroupSettings = &types.GroupSettings{}
	if err := c.ShouldBindJSON(&groupSettings); err != nil {
		errorMessage := "Failed to bind json. Check the request body and ensure that the correct fields are present."
		glg.Errorf("%s %s", errorMessage, err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": errorMessage,
			"error":   err,
		})
		return
	}

	err := db.UpdateGroupSettings(groupSettings, name)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to update group settings",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":       "Group settings updated successfully.",
		"groupSettings": groupSettings,
	})
}

func DeleteGroupSettings(c *gin.Context) {
	name := c.Param("name")

	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "name is required",
		})
		return
	}

	err := db.DeleteGroupSettings(name)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to delete group settings",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Group settings deleted successfully.",
	})
}
