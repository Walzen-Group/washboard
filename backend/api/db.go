package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"washboard/db"
	"washboard/types"

	"github.com/gin-gonic/gin"
	"github.com/kpango/glg"
)

func SyncWithPortainer(c *gin.Context) {

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
	id := c.Param("id")

	if id == "" {
		stacks, err := db.GetAllStackSettings()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to get stack settings",
				"error":   err,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
			"stackSettings": stacks,
		})
		return
	} else {
		intId, _ := strconv.Atoi(id)
		glg.Debugf("Getting stack settings for id: %d", intId)
		stack, err := db.GetStackSettings(intId)
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
}

func UpdateStackSettings(c *gin.Context) {
	idStr := c.Param("id")

	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "id is required",
		})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "id must be an integer",
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

	err = db.UpdateStackSettings(stackSettings, id)

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
	idStr := c.Param("id")

	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "id is required",
		})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "id must be an integer",
		})
		return
	}

	err = db.DeleteStackSettings(id)

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
			"error":  errorMessage,
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
