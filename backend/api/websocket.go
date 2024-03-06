package api

import (
	"bufio"
	"bytes"
	"encoding/json"
	"net/http"
	"time"
	"washboard/types"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/kpango/glg"
)

type WsState struct {
	StackUpdateInProgressIds []int `json:"stackUpdateInProgressIds"`
}

const (
	CMD_STOP = "CMD_STOP"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WsHandler(c *gin.Context) {
	glg.Infof("WS connection from %s", c.ClientIP())
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	glg.Infof("Upgraded to websocket")
	if err != nil {
		glg.Errorf("error while upgrading to websocket: %s", err)
		return
	}
	glg.Infof("Client %s connected", c.ClientIP())

	oniiChan := make(chan string)
	go readData(ws, oniiChan)
	go pushData(ws, oniiChan)
}

func readData(ws *websocket.Conn, oniiChan chan string) {
	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseNoStatusReceived, websocket.CloseAbnormalClosure) {
				glg.Infof("client %s disconnected from websocket", ws.RemoteAddr())
			} else {
				glg.Infof("client %s disconnected from websocket: %s", ws.RemoteAddr(), err)
			}
			oniiChan <- CMD_STOP
			return
		}
		glg.Infof("Received message: %s", message)
	}
}

func pushData(ws *websocket.Conn, oniiChan chan string) {
	defer ws.Close()
	for {
		select {
		case msg := <-oniiChan:
			if msg == CMD_STOP {
				glg.Debugf("stopping websocket push")
				return
			}
		default:
		}

		items := appState.StackUpdateQueue.Items()

		// group items by status and create new map[string]map[string]cache.Item
		groupedItems := make(map[string]map[string]types.StackUpdateStatus)
		for _, item := range items {
			stackUpdateStatus := item.Object.(types.StackUpdateStatus)

			status := stackUpdateStatus.Status
			if _, ok := groupedItems[status]; !ok {
				groupedItems[status] = make(map[string]types.StackUpdateStatus)
			}
			groupedItems[status][stackUpdateStatus.StackName] = stackUpdateStatus
		}

		out, err := encodeJson(groupedItems)
		if err != nil {
			glg.Warnf("error while marshaling encoder to json: %s", err)
			break
		}
		err = ws.WriteMessage(websocket.TextMessage, out)
		if err != nil {
			glg.Errorf("error while writing to websocket: %s", err)
			break
		}
		time.Sleep(1 * time.Second)
	}
}

func encodeJson(in interface{}) ([]byte, error) {
	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	err := encoder.Encode(in)
	if err != nil {
		return nil, err
	}
	err = writer.Flush()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
