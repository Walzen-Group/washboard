package api

import (
	"errors"
	"fmt"
	"net/http"
	"washboard/db"
	"washboard/portainer"
	"washboard/types"

	"github.com/gin-gonic/gin"
	"github.com/kpango/glg"
)

func SyncWithPortainer(c *gin.Context) {
	var syncOptions *types.SyncOptions = &types.SyncOptions{}
	if err := c.ShouldBindJSON(&syncOptions); err != nil {
		errorMessage := "Failed to bind json. Check the request body and ensure that the correct fields are present."
		glg.Errorf("%s %s", errorMessage, err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": errorMessage,
			"error":   err,
		})
		return
	}

	var containers []*types.ContainerDto

	for _, endpoint := range syncOptions.EndpointIds {
		tmp, err := portainer.GetContainers(endpoint, "")

		if err != nil {
			errorMessage := fmt.Sprintf("Failed to get containers for endpoint %d", endpoint)
			glg.Errorf("%s %s", errorMessage, err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": errorMessage,
				"error":   err,
			})
			return
		}
		containers = append(containers, tmp...)
	}

	collectedGroups := make(map[string]*types.GroupSettings)
	collectedStacks := make(map[string]*types.StackSettings)
	targetError := &db.DoesNotExistWrappedError{}

	// add missing grups and stakcs (das bleibt so)
	for _, container := range containers {
		if groupName, ok := container.Labels[types.StackGroupLabel]; ok {
			group, err := db.GetGroupSettings(groupName.(string))
			if errors.As(err, &targetError) {
				group = &types.GroupSettings{GroupName: groupName.(string)}
				glg.Infof("adding missing group %s", groupName.(string))
				db.CreateGroupSettings(group)
			} else if err != nil {
				errorMessage := fmt.Sprintf("Failed to sync settings with db %s", groupName)
				glg.Errorf("%s %s", errorMessage, err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": errorMessage,
					"error":   err,
				})
				return
			}
			collectedGroups[groupName.(string)] = group
		} else if stackName, ok := container.Labels[types.StackLabel]; ok {
			stack, err := db.GetStackSettings(stackName.(string))
			if errors.As(err, &targetError) {
				stack = &types.StackSettings{StackName: stackName.(string)}
				glg.Infof("adding missing stack %s", stackName.(string))
				db.CreateStackSettings(stack)
			} else if err != nil {
				errorMessage := fmt.Sprintf("Failed to sync settings with db %s", groupName)
				glg.Errorf("%s %s", errorMessage, err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": errorMessage,
					"error":   err,
				})
				return
			}
			collectedStacks[stackName.(string)] = stack
		}
	}

	allGroupSettings, err := db.GetAllGroupSettings()
	allStackSettings, err2 := db.GetAllStackSettings()
	if err != nil {
		errorMessage := "Failed to get all group settings"
		glg.Errorf("%s %s", errorMessage, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": errorMessage,
			"error":   err,
		})
		return
	}
	if err2 != nil {
		errorMessage := "Failed to get all stack settings"
		glg.Errorf("%s %s", errorMessage, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": errorMessage,
			"error":   err,
		})
		return
	}

	stacksToRemove := make([]string, 0)
	groupsToRemove := make([]string, 0)

	for _, stackSettings := range allStackSettings {
		if _, ok := collectedStacks[stackSettings.StackName]; !ok {
			stacksToRemove = append(stacksToRemove, stackSettings.StackName)
		}
	}

	for _, groupSettings := range allGroupSettings {
		if _, ok := collectedGroups[groupSettings.GroupName]; !ok {
			groupsToRemove = append(groupsToRemove, groupSettings.GroupName)
		}
	}

	for _, stack := range stacksToRemove {
		glg.Infof("removing orphaned stack %s", stack)
		db.DeleteStackSettings(stack)
	}

	for _, group := range groupsToRemove {
		glg.Infof("removing orphaned group %s", group)
		db.DeleteGroupSettings(group)
	}
}

func CreateStackSettings(c *gin.Context) {
	var stackSettings *types.StackSettings = &types.StackSettings{}
	if err := c.ShouldBindJSON(&stackSettings); err != nil {
		errorMessage := "Failed to bind json. Check the request body and ensure that the correct fields are present."
		glg.Errorf("%s %s", errorMessage, err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": errorMessage,
			"error":   err,
		})
		return
	}
	glg.Infof("Creating stack settings: %+v", stackSettings)
	err := db.CreateStackSettings(stackSettings)
	if err != nil {
		target := &db.CannotInsertWrappedError{}
		if errors.As(err, &target) {
			glg.Errorf("%s %s", err.(*db.CannotInsertWrappedError), err)
			c.JSON(http.StatusNotAcceptable, gin.H{
				"message": "There was an issue with the insert operation",
				"error":   err,
			})
			return
		}
		errorMessage := "Failed to create stack settings."
		glg.Errorf("%s %s", errorMessage, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": errorMessage,
			"error":   err,
		})
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
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to get stack settings",
				"error":   err,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message":       "ok",
			"stackSettings": stacks,
		})
		return
	}
	stack, err := db.GetStackSettings(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get stack settings",
			"error":   fmt.Sprintf("%s", err),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":       "ok",
		"stackSettings": stack,
	})
	return
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
		glg.Errorf("%s %s", errorMessage, err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": errorMessage,
			"error":   err,
		})
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
		target := &db.CannotInsertWrappedError{}
		if errors.As(err, &target) {
			glg.Errorf("%s %s", err.(*db.CannotInsertWrappedError), err)
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
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "id is required",
		})
		return
	}

	group, err := db.GetGroupSettings(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get group settings",
			"error":   fmt.Sprintf("%s", err),
		})
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
