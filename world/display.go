package world

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/sapphire-ai-dev/sapphire_display/server"
	"net/url"
)

type DisplayClient struct {
	conn                 *websocket.Conn
	displayMsgProcessors map[int]func(msg *server.WorldMsg)
}

func (c *DisplayClient) Send(data []byte) {
	PrintErr(c.conn.WriteMessage(websocket.TextMessage, data))
}

func (c *DisplayClient) SendState(state []byte) {
	resp := server.WorldResp{
		Method: server.MsgMethodUpdateState,
		State:  string(state),
	}
	data, err := json.Marshal(resp)
	PrintErr(err)
	c.Send(data)
}

func (c *DisplayClient) Read() {
	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			panic(err)
			return
		}

		var msg *server.WorldMsg
		err = json.Unmarshal(data, &msg)
		if err != nil {
			panic(err)
			return
		}

		c.processDisplayMsg(msg)
	}
}

func (c *DisplayClient) processDisplayMsg(msg *server.WorldMsg) {
	if processor, seen := c.displayMsgProcessors[msg.Method]; seen {
		processor(msg)
	}
}

func (c *DisplayClient) processDisplayMsgCreateViewerActors(msg *server.WorldMsg) {
	result := &server.WorldResp{
		Method:   msg.Method,
		ActorIds: []int{},
	}
	for i := 0; i < msg.ActorCount; i++ {
		newActorId, _ := NewActor()
		result.ActorIds = append(result.ActorIds, newActorId)
	}

	data, err := json.Marshal(result)
	PrintErr(err)
	c.Send(data)
}

func (c *DisplayClient) processDisplayMsgViewerSpeech(msg *server.WorldMsg) {
	speakerId, listenerId := msg.SpeechActorId, GetAgentSingleton()
	Speak(&speakerId, &listenerId, msg.Speech)
	Tick()
}

func NewDisplayClient(w World) *DisplayClient {
	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/world", RawQuery: fmt.Sprintf("name=%s", w.Name())}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil
	}

	result := &DisplayClient{conn: c}
	result.displayMsgProcessors = map[int]func(msg *server.WorldMsg){
		server.MsgMethodCreateViewerActors: result.processDisplayMsgCreateViewerActors,
		server.MsgMethodViewerSpeech:       result.processDisplayMsgViewerSpeech,
	}
	go result.Read()
	return result
}
