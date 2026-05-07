package api

import (
	"bufio"
	"bytes"
	"encoding/json"
	"net/http"
	"time"
	"washboard/portainer"
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
	glg.Debugf("ws connection request from %s", c.ClientIP())
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	glg.Debugf("Upgraded to websocket")
	if err != nil {
		glg.Errorf("error while upgrading to websocket: %s", err)
		return
	}
	glg.Infof("client %s connected to status websocket", c.ClientIP())

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
	// Force the first refresh-state push so reconnecting clients sync immediately.
	const firstPush = ^uint64(0)
	lastRefreshVersion := firstPush
	for {
		select {
		case msg := <-oniiChan:
			if msg == CMD_STOP {
				glg.Debugf("stopping websocket push")
				return
			}
		default:
		}

		// Always push the stack-update queue (cheap, and the frontend reconciles).
		items := appState.StackUpdateQueue.Items()
		groupedItems := make(map[string]map[string]types.StackUpdateStatus)
		for _, item := range items {
			stackUpdateStatus := item.Object.(types.StackUpdateStatus)
			status := stackUpdateStatus.Status
			if _, ok := groupedItems[status]; !ok {
				groupedItems[status] = make(map[string]types.StackUpdateStatus)
			}
			groupedItems[status][stackUpdateStatus.StackName] = stackUpdateStatus
		}
		if err := writeEnvelope(ws, types.WsMsgStackUpdateQueue, groupedItems); err != nil {
			glg.Errorf("error writing stack-update-queue envelope: %s", err)
			return
		}

		// Push image-refresh state only on first connection or when it has changed.
		currentVersion := portainer.Refresh.Version()
		if currentVersion != lastRefreshVersion {
			if err := writeEnvelope(ws, types.WsMsgImageRefreshState, portainer.Refresh.Snapshot()); err != nil {
				glg.Errorf("error writing image-refresh-state envelope: %s", err)
				return
			}
			lastRefreshVersion = currentVersion
		}

		time.Sleep(1 * time.Second)
	}
}

func writeEnvelope(ws *websocket.Conn, msgType string, data interface{}) error {
	out, err := encodeJson(types.WsEnvelope{Type: msgType, Data: data})
	if err != nil {
		glg.Warnf("error while marshaling envelope to json: %s", err)
		return err
	}
	return ws.WriteMessage(websocket.TextMessage, out)
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
