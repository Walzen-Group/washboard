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

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func wsHandler(c *gin.Context, in interface{}) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	glg.Infof("Client %s connected", c.ClientIP())

	oniiChan := make(chan string)
	go readData(ws, oniiChan)
	go pushData(ws, in, oniiChan)
}

func readData(ws *websocket.Conn, oniiChan chan string) {
	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			return
		}
		glg.Infof("Received message: %s", message)
	}
}

func pushData(ws *websocket.Conn, in interface{}, oniiChan chan string) {
	defer ws.Close()
	for {
		select {
		case msg := <-oniiChan:
			if msg == "STOP" {
				glg.Debugf("stopping websocket push")
				return
			}
		default:
		}
		out, err := encodeJson(in)
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
