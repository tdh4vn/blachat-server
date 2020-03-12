package services

import (
	"blachat-server/config"
	"blachat-server/entities"
	"encoding/json"
	"github.com/valyala/fasthttp"
	"time"
)

const NEW_MESSAGE = "new_message"
const NEW_CHANNEL = "new_channel"
const UPDATE_CHANNEL = "update_channel"
const UPDATE_USER = "update_user"
const PRESENT_USER = "present_user"
const MARK_SEEN = "mark_seen"
const MARK_RECEIVE = "mark_receive"
const TYPING_EVENT = "typing_event"
const USER_IS_ONLINE = "user_is_online"
const USER_IS_OFFLINE = "user_is_offline"

type BasePublishMessage struct {
	Method string `json:"method"`
	Params interface{} `json:"params"`
}

type Event struct {
	Type string `json:"type"`
	Payload interface{} `json:"payload"`
}

type PublishToChannelMessage struct {
	Channel string `json:"channel"`
	Data Event `json:"data"`
}

type BroadcastParams struct {
	Channels []string `json:"channels"`
	Data Event `json:"payload"`
}

func SendTypingEvent(receiveID string, cID string, userTyping string, isTyping bool) {
	body := BasePublishMessage{
		Method: "publish",
		Params: PublishToChannelMessage{
			Channel: "chat#" + receiveID,
			Data: Event {
				Type: TYPING_EVENT,
				Payload: map[string]interface{} {
					"user_id": userTyping,
					"channel_id": cID,
					"is_typing": isTyping,
					"time": time.Now().Unix(),
				},
			},
		},
	}

	sendDataToCentrifugo(body)
}

func SendReceiveMessageEvent(messageID string, channelID string, receiveID string, actionActorID string) {
	body := BasePublishMessage{
		Method: "publish",
		Params: PublishToChannelMessage{
			Channel: "chat#" + receiveID,
			Data: Event {
				Type: MARK_RECEIVE,
				Payload: map[string]interface{} {
					"message_id": messageID,
					"channel_id": channelID,
					"actor_id": actionActorID,
					"time": time.Now(),
				},
			},
		},
	}

	sendDataToCentrifugo(body)
}

func SendSeenMessageEvent(messageID string, channelID string, receiveID string, actionActorID string) {
	body := BasePublishMessage{
		Method: "publish",
		Params: PublishToChannelMessage{
			Channel: "chat#" + receiveID,
			Data: Event {
				Type: MARK_SEEN,
				Payload: map[string]interface{} {
					"message_id": messageID,
					"channel_id": channelID,
					"actor_id": actionActorID,
					"time": time.Now(),
				},
			},
		},
	}

	sendDataToCentrifugo(body)
}


func SendMessageViaCentrigufo(message *entities.Message, toUser string) {

	body := BasePublishMessage{
		Method: "publish",
		Params: PublishToChannelMessage{
			Channel: "chat#" + toUser,
			Data: Event {
				Type: NEW_MESSAGE,
				Payload: message,
			},
		},
	}

	sendDataToCentrifugo(body)
}


func SendNewChannel(channel *entities.Channel, toUser string){
	body := BasePublishMessage{
		Method: "publish",
		Params: PublishToChannelMessage{
			Channel: "chat#" + toUser,
			Data: Event {
				Type: NEW_CHANNEL,
				Payload: channel,
			},
		},
	}
	sendDataToCentrifugo(body)
}

func SendUserOnline(idOfUserOnline string, toUsers []string) {
	body := BasePublishMessage{
		Method: "broadcast",
		Params: BroadcastParams {
			Channels: toUsers,
			Data: Event{
				Type: USER_IS_ONLINE,
				Payload: idOfUserOnline,
			},
		},
	}

	sendDataToCentrifugo(body)
}

func SendUserOffline(idOfUserOnline string, toUsers []string) {
	body := BasePublishMessage{
		Method: "broadcast",
		Params: BroadcastParams {
			Channels: toUsers,
			Data: Event{
				Type: USER_IS_OFFLINE,
				Payload: idOfUserOnline,
			},
		},
	}

	sendDataToCentrifugo(body)
}

func sendDataToCentrifugo(body BasePublishMessage){
	url := config.GetConfig().GetString("centrifugo_url")
	key := config.GetConfig().GetString("centrifugo_key")

	bodyString, _ := json.Marshal(body)

	var strPost = []byte("POST")
	var strRequestURI = []byte(url + "/api")

	req := fasthttp.AcquireRequest()
	req.SetBody(bodyString)
	req.Header.SetMethodBytes(strPost)
	req.SetRequestURIBytes(strRequestURI)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", "apikey " + key)

	res := fasthttp.AcquireResponse()

	if err := fasthttp.Do(req, res); err != nil {
		panic("handle error")
	}
}