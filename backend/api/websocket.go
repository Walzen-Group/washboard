package api

import (
	"bufio"
	"bytes"
	"encoding/json"
	"net/http"
	"time"

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

	wsState := &WsState{}

	oniiChan := make(chan string)
	go readData(ws, oniiChan, wsState)
	go pushData(ws, oniiChan, wsState)
}

func readData(ws *websocket.Conn, oniiChan chan string, wsState *WsState) {
	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseNoStatusReceived, websocket.CloseAbnormalClosure) {
				glg.Infof("client %s disconnected from websocket", ws.RemoteAddr())
			} else {
				glg.Warnf("error while writing to websocket: %s", err)
			}
			oniiChan <- CMD_STOP
			return
		}
		glg.Infof("Received message: %s", message)
	}
}

func pushData(ws *websocket.Conn, oniiChan chan string, wsState *WsState) {
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
		out, err := encodeJson(items)
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
