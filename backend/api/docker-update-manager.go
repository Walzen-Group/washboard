package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"washboard/portainer"

	"github.com/gin-gonic/gin"
	"github.com/kpango/glg"
)

// PortainerGetEndpoint handles an HTTP GET request to retrieve the identifier of a Portainer endpoint
// based on the provided endpoint name. It attempts to fetch the endpoint ID using the 'endpointName' query
// parameter. If the 'endpointName' parameter is not provided in the request, a default value of "Quasar" is used.
// The function queries the Portainer API to find the endpoint ID corresponding to the given endpoint name.
// Upon successful retrieval, the function returns the endpoint ID. If the query fails due to an error
// or if the endpoint cannot be found, it returns an error message and an appropriate HTTP status code.
//
// Query Parameters:
// - endpoint (optional, default "Quasar"): The name of the Portainer endpoint for which the ID is requested.
//   If this parameter is omitted from the query, the function defaults to "Quasar" as the endpoint name.
//
// Responses:
// - 200 OK: Successfully found the endpoint ID. The response includes the endpoint ID in the format
//   {"endpoint": "endpoint_id_here"}.
// - 500 Internal Server Error: Encountered if there is an error in querying the Portainer API for the
//   endpoint ID, or if the endpoint name provided does not match any existing endpoints. The response
//   contains an error message indicating the failure to retrieve endpoint information.
//
// Note: This function utilizes the Gin web framework for routing and handling HTTP requests and responses,
// facilitating the extraction of query parameters and the delivery of structured JSON responses. Logging
// is employed to record any errors encountered during the operation. This function ensures efficient
// communication with the Portainer API and provides clear feedback to the client regarding the operation's
// outcome.
func PortainerGetEndpoint(c *gin.Context) {
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

// PortainerGetContainers handles an HTTP GET request to retrieve a list of containers filtered by a specific
// stack label within a specified Portainer endpoint. It extracts 'endpointId' and 'stackName' from the query
// parameters, validates and converts 'endpointId' to an integer, and then queries the Portainer API to get
// containers associated with the given stack label and endpoint. The function returns a list of containers
// if successful or appropriate error messages and HTTP status codes in case of errors or if no containers are
// found matching the criteria.
//
// Query Parameters:
// - endpointId (optional, default "1"): The unique identifier of the Portainer endpoint from which to retrieve
//   containers. If not provided, a default value of "1" is assumed. It should be a valid integer.
// - stackName (required): The label of the stack used to filter containers. Containers associated with this
//   stack label within the specified endpoint are returned.
//
// Responses:
// - 200 OK: Successfully retrieved the list of containers. The response body includes the containers' data.
// - 404 Not Found: No containers found matching the specified stack label in the given endpoint. Returns a
//   message indicating that the specified stack is not found.
// - 500 Internal Server Error: Encountered if 'endpointId' cannot be converted to an integer or if there is
//   an error querying the Portainer API for containers. The response includes an error message detailing the
//   issue.
//
// Note: This function leverages the Gin web framework for routing and handling HTTP requests and responses.
// It uses logging to record errors encountered during operation. Careful handling of query parameters and
// clear communication of errors through HTTP status codes and messages are emphasized to ensure a smooth
// user experience and effective troubleshooting.
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
	c.JSON(http.StatusOK, res)
}

func PortainerGetContainers(c *gin.Context) {
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

	stackName := c.Query("stackName")
	if stackName == "" {
		glg.Error("stackName is empty")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "stackName is empty",
		})
		return
	}

	res, err := portainer.GetContainers(endpointId, stackName)
	if err != nil {
		glg.Error("failed to get stacks")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to get stacks",
		})
		return
	}

	// TODO: Need another check to check if the stack exists
	// if len(res) == 0 {
	// 	c.JSON(http.StatusNotFound, gin.H{
	// 		"message": fmt.Sprintf("Stack \"%s\" not found in environment %d", stackName, endpointId),
	// 	})
	// 	return
	// }
	c.JSON(http.StatusOK, res)
}

// PortainerGetImageStatus handles an HTTP GET request to retrieve the status of a specific container image
// within a given Portainer endpoint. It extracts the 'endpointId' and 'containerId' from the query parameters,
// performs input validation, and then proceeds to query the Portainer API for the image status using these
// identifiers. If successful, it returns the status of the image. In case of errors during input validation
// or while fetching the image status, appropriate HTTP status codes and error messages are returned.
//
// Query Parameters:
// - endpointId (optional, default "1"): The unique identifier of the Portainer endpoint. If not provided,
//   a default value of "1" is used. This parameter is expected to be a valid integer.
// - containerId (required): The unique identifier of the container whose image status is being queried.
//
// Responses:
// - 200 OK: Successfully retrieved the image status. The response body includes the status in the format
//   {"status": "image_status_here"}.
// - 500 Internal Server Error: Encountered if 'endpointId' cannot be converted to an integer or if 'containerId'
//   is not provided in the query parameters. The response includes an error message detailing the issue.
// - Additionally, if the function fails to get the image status due to issues with the Portainer API, it logs
//   the error but does not explicitly return a response to the client. Implementers may want to handle this case
//   by returning a specific error response to the client.
//
// Note: This function uses the Gin web framework for routing and handling HTTP requests and responses. It employs
// logging to record errors encountered during the operation. It is recommended to enhance error handling to ensure
// that all error states are communicated back to the client with appropriate HTTP status codes and error messages.
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
		"cacheTimestamp": time.Now(),
	})
}

// PortainerUpdateContainer handles the HTTP request to update a container configuration in Portainer.
// It reads from a JSON request body to extract necessary parameters such as the container's identifier,
// the endpoint ID where the container is located, and a boolean flag indicating whether to pull the
// image again. The function validates the presence and types of these fields, returning appropriate
// error messages and HTTP status codes for any discrepancies found.
//
// Parameters:
// - c *gin.Context: The Gin context for HTTP request and response handling, facilitating the extraction
//   of JSON request body data and the delivery of HTTP responses.
//
// JSON Request Body Fields:
// - pullImage (required, boolean): Specifies whether the container's image should be pulled again
//   as part of the update process.
// - endpointId (required, int): The unique identifier of the Portainer endpoint where the container
//   is deployed. This is expected to be a numeric value but comes as a float64 from JSON and is converted
//   to an int.
// - containerId (required, string): The unique identifier of the container to be updated.
//
// Responses:
// - 200 OK: The container update was successful. The response includes the updated container's identifier.
// - 400 Bad Request: Missing required fields or incorrect data types in the request body. An error message
//   detailing the specific issue is returned.
// - 404 Not Found: The update operation failed due to an error with the Portainer API call, with an error
//   message provided.
//
// Note: This function employs logging to record warnings for missing or incorrectly typed fields and errors
// for unsuccessful update attempts. It uses the Gin web framework to parse JSON request bodies and to send
// JSON responses, ensuring a structured and consistent API interface.
func PortainerUpdateContainer(c *gin.Context) {
	// Set pullImage
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


	var pullImageBool bool
	var endpointId int
	var containerId string

	// we will forget how how the error handling works here in the future
	// val ok ok val
	if pullImage, ok := reqBody["pullImage"]; !ok {
		glg.Warn("pullImage field is missing")
		c.JSON(http.StatusBadRequest, gin.H{"message": "pullImage field is missing"})
		return
	} else if pullImageBool, ok = pullImage.(bool); !ok {
		glg.Warn("pullImage field is not a boolean")
		c.JSON(http.StatusBadRequest, gin.H{"message": "pullImage field is not a boolean"})
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

	if containerRaw, ok := reqBody["containerId"]; !ok {
		glg.Warn("containerId field is missing")
		c.JSON(http.StatusBadRequest, gin.H{"message": "containerId field is missing"})
		return
	} else if containerId, ok = containerRaw.(string); !ok {
		glg.Warn("containerId field is not a string")
		c.JSON(http.StatusBadRequest, gin.H{"message": "containerId field is not a boolean"})
		return
	}

	res, err := portainer.UpdateContainer(endpointId, containerId, pullImageBool)
	if err != nil {
		glg.Errorf("failed to update container: %s", err)
		c.JSON(http.StatusNotFound, gin.H{
			"message": "failed to update container",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": res,
	})
}


// PortainerUpdateStack handles the request to update a Portainer stack. It extracts necessary
// information such as endpointId, stackId, prune, and pullImage options from the JSON request body.
// If any required field is missing or if there is a type mismatch, it responds with an appropriate
// error message and status code. Upon successful validation of request data, it calls the Portainer
// API to update the stack with the provided details. If the stack update is successful, it responds
// with the stack ID, otherwise, it returns an error message and status code indicating the failure.
//
// Parameters:
// - c *gin.Context: The Gin context which holds request and response information.
//
// JSON Request Body Fields:
// - endpointId (required, int): The unique identifier for the Portainer endpoint where the stack is deployed.
// - stackId (required, int): The unique identifier for the stack to be updated.
// - prune (required, boolean): Indicates whether unused services should be removed after the update.
// - pullImage (required, boolean): Determines whether images should be pulled fresh before deployment.
//
// Responses:
// This function can respond with several HTTP status codes depending on the outcome of the request processing:
// - 200 OK: Successfully updated the stack and returns the stack ID.
// - 400 Bad Request: The request body is missing required fields or has type mismatches.
// - 404 Not Found: Failed to update the stack due to an issue with the Portainer API call.
//
// Note: This function logs warnings for missing fields or type mismatches and errors for issues encountered
// during stack update attempts. It uses the Gin framework for handling HTTP requests and responses.
func PortainerUpdateStack(c *gin.Context) {
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
	var stackId int
	var prune bool
	var pullImage bool

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

	if stackRaw, ok := reqBody["stackId"]; !ok {
		glg.Warn("stackId field is missing")
		c.JSON(http.StatusBadRequest, gin.H{"message": "stackId field is missing"})
		return
	} else if stackIdFloat, ok := stackRaw.(float64); !ok {
		glg.Warn("stackId field is not a string")
		c.JSON(http.StatusBadRequest, gin.H{"message": "stackId field is not a boolean"})
		return
	} else {
		stackId = int(stackIdFloat)
	}

	if pullImageRaw, ok := reqBody["pullImage"]; !ok {
		glg.Warn("pullImage field is missing")
		c.JSON(http.StatusBadRequest, gin.H{"message": "pullImage field is missing"})
		return
	} else if pullImage, ok = pullImageRaw.(bool); !ok {
		glg.Warn("pullImage field is not a boolean")
		c.JSON(http.StatusBadRequest, gin.H{"message": "pullImage field is not a boolean"})
		return
	}

	if pruneRaw, ok := reqBody["prune"]; !ok {
		glg.Warn("prune field is missing")
		c.JSON(http.StatusBadRequest, gin.H{"message": "prune field is missing"})
		return
	} else if prune, ok = pruneRaw.(bool); !ok {
		glg.Warn("prune field is not a boolean")
		c.JSON(http.StatusBadRequest, gin.H{"message": "prune field is not a boolean"})
		return
	}


	res, err := portainer.EnqueueUpdateStack(endpointId, stackId, prune, pullImage)

	if res == -2 {
		glg.Errorf("Endpoint: %d, Stack: %d, %s", endpointId, stackId, err)
		c.JSON(http.StatusAccepted, gin.H{
			"message": "Failed to update stack",
			"error":  err.Error(),
		})
		return
	}
	if err != nil {
		glg.Errorf("Failed to update stack: %s", err)
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Failed to update stack",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": res,
	})
}

