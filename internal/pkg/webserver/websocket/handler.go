package websocket

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/chat"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/general"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/types"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/webserver/errors"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/webserver/websocket/chatHandler"
	"net/http"
	"strings"
)

var websocketUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	Error: func(w http.ResponseWriter, r *http.Request, status int, reason error) {
		general.Log.Error("General failure in websocket: ", reason)
	},
}

func HandleWebsocket(writer http.ResponseWriter, request *http.Request, user *types.User) {
	connection, err := websocketUpgrader.Upgrade(writer, request, nil)
	if err != nil {
		writer.WriteHeader(500)
		writer.Write([]byte("{\"error\": \"websocket_upgrade_failed\"}"))
		general.Log.Error("Failed to upgrade websocket: ", err)
	}

	chat.ConnectionManager.AddConnection(user.Id, connection)

	go func() {
		defer chat.ConnectionManager.RemoveConnection(user.Id)

	loop:
		for {
			messageType, message, err := connection.ReadMessage()
			if err != nil {
				general.Log.Error("Failed to read websocket message: ", err)
				break
			}

			if messageType == websocket.CloseMessage {
				break
			} else if messageType == websocket.PingMessage {
				err = connection.WriteMessage(websocket.PongMessage, message)
				if err != nil {
					general.Log.Error("Failed to send pong to websocket: ", err)
					break
				}
			}

			var messageMap map[string]interface{}
			err = json.Unmarshal(message, &messageMap)
			if err != nil {
				general.Log.Error("Failed to decode json: ", err)
				err = connection.WriteJSON(errors.InvalidRequest)
				if err != nil {
					general.Log.Error("Failed to send response to websocket: ", err)
					break
				}

				continue
			}

			switch strings.ToLower(messageMap["type"].(string)) {
			case strings.ToLower(types.ChatMessagePacketName):
				packet := types.CreateChatMessagePacket()
				jsonBytes, err := json.Marshal(messageMap)
				if err != nil {
					general.Log.Error("Failed encode map to json: ", err)
					err = connection.WriteJSON(errors.InternalError)
					if err != nil {
						general.Log.Error("Failed to send response to websocket: ", err)
						break loop
					}
					continue loop
				}

				err = json.Unmarshal(jsonBytes, packet)
				if err != nil {
					general.Log.Error("Failed encode json to ChatMessagePacket: ", err)
					err = connection.WriteJSON(errors.InternalError)
					if err != nil {
						general.Log.Error("Failed to send response to websocket: ", err)
						break loop
					}
					continue loop
				}

				err = chatHandler.HandleMessagePacket(connection, user, packet)
				if err != nil {
					general.Log.Error("Failed to send response to websocket: ", err)
					break loop
				}
			case strings.ToLower(types.ChatRoomAddUserPacketName):
				packet := types.CreateChatRoomAddUserPacket()
				jsonBytes, err := json.Marshal(messageMap)
				if err != nil {
					general.Log.Error("Failed encode map to json: ", err)
					err = connection.WriteJSON(errors.InternalError)
					if err != nil {
						general.Log.Error("Failed to send response to websocket: ", err)
						break loop
					}
					continue loop
				}

				err = json.Unmarshal(jsonBytes, packet)
				if err != nil {
					general.Log.Error("Failed encode json to ChatRoomAddUserPacket: ", err)
					err = connection.WriteJSON(errors.InternalError)
					if err != nil {
						general.Log.Error("Failed to send response to websocket: ", err)
						break loop
					}
					continue loop
				}

				err = chatHandler.HandleRoomAddUserPacket(connection, user, packet)
				if err != nil {
					general.Log.Error("Failed to send response to websocket: ", err)
					break loop
				}
			case strings.ToLower(types.ChatRoomLeaveUserPacketName):
				packet := types.CreateChatRoomLeaveUserPacket()
				jsonBytes, err := json.Marshal(messageMap)
				if err != nil {
					general.Log.Error("Failed encode map to json: ", err)
					err = connection.WriteJSON(errors.InternalError)
					if err != nil {
						general.Log.Error("Failed to send response to websocket: ", err)
						break loop
					}
					continue loop
				}

				err = json.Unmarshal(jsonBytes, packet)
				if err != nil {
					general.Log.Error("Failed encode json to ChatRoomLeaveUserPacket: ", err)
					err = connection.WriteJSON(errors.InternalError)
					if err != nil {
						general.Log.Error("Failed to send response to websocket: ", err)
						break loop
					}
					continue loop
				}

				err = chatHandler.HandleRoomLeaveUserPacket(connection, user, packet)
				if err != nil {
					general.Log.Error("Failed to send response to websocket: ", err)
					break loop
				}
			case strings.ToLower(types.ChatRoomCreatePacketName):
				packet := types.CreateChatRoomCreatePacket()
				jsonBytes, err := json.Marshal(messageMap)
				if err != nil {
					general.Log.Error("Failed encode map to json: ", err)
					err = connection.WriteJSON(errors.InternalError)
					if err != nil {
						general.Log.Error("Failed to send response to websocket: ", err)
						break loop
					}
					continue loop
				}

				err = json.Unmarshal(jsonBytes, packet)
				if err != nil {
					general.Log.Error("Failed encode json to ChatRoomCreatePacket: ", err)
					err = connection.WriteJSON(errors.InternalError)
					if err != nil {
						general.Log.Error("Failed to send response to websocket: ", err)
						break loop
					}
					continue loop
				}

				err = chatHandler.HandleCreateRoom(connection, user, packet)
				if err != nil {
					general.Log.Error("Failed to send response to websocket: ", err)
					break loop
				}
			case strings.ToLower(types.ChatUnreadMessagesPacketName):
				packet := types.CreateChatUnreadMessagesPacket()
				jsonBytes, err := json.Marshal(messageMap)
				if err != nil {
					general.Log.Error("Failed encode map to json: ", err)
					err = connection.WriteJSON(errors.InternalError)
					if err != nil {
						general.Log.Error("Failed to send response to websocket: ", err)
						break loop
					}
					continue loop
				}

				err = json.Unmarshal(jsonBytes, packet)
				if err != nil {
					general.Log.Error("Failed encode json to ChatUnreadMessagesPacket: ", err)
					err = connection.WriteJSON(errors.InternalError)
					if err != nil {
						general.Log.Error("Failed to send response to websocket: ", err)
						break loop
					}
					continue loop
				}

				err = chatHandler.HandleUnreadMessagesPacket(connection, user, packet)
				if err != nil {
					general.Log.Error("Failed to send response to websocket: ", err)
					break loop
				}
			case strings.ToLower(types.ChatRoomMessagesPacketName):
				packet := types.CreateChatRoomMessagesPacket()
				jsonBytes, err := json.Marshal(messageMap)
				if err != nil {
					general.Log.Error("Failed encode map to json: ", err)
					err = connection.WriteJSON(errors.InternalError)
					if err != nil {
						general.Log.Error("Failed to send response to websocket: ", err)
						break loop
					}
					continue loop
				}

				err = json.Unmarshal(jsonBytes, packet)
				if err != nil {
					general.Log.Error("Failed encode json to ChatRoomMessagesPacket", err)
					err = connection.WriteJSON(errors.InternalError)
					if err != nil {
						general.Log.Error("Failed to send response to websocket: ", err)
						break loop
					}
					continue loop
				}

				err = chatHandler.HandleRoomMessagesPacket(connection, user, packet)
				if err != nil {
					general.Log.Error("Failed to send response to websocket: ", err)
					break loop
				}
			case strings.ToLower(types.ChatRoomsPacketName):
				packet := types.CreateChatRoomsPacket()
				jsonBytes, err := json.Marshal(messageMap)
				if err != nil {
					general.Log.Error("Failed encode map to json: ", err)
					err = connection.WriteJSON(errors.InternalError)
					if err != nil {
						general.Log.Error("Failed to send response to websocket: ", err)
						break loop
					}
					continue loop
				}

				err = json.Unmarshal(jsonBytes, packet)
				if err != nil {
					general.Log.Error("Failed encode json to ChatRoomsPacket: ", err)
					err = connection.WriteJSON(errors.InternalError)
					if err != nil {
						general.Log.Error("Failed to send response to websocket: ", err)
						break loop
					}
					continue loop
				}

				err = chatHandler.HandleRoomsPacket(connection, user, packet)
				if err != nil {
					general.Log.Error("Failed to send response to websocket: ", err)
					break loop
				}
			default:
				err = connection.WriteJSON(errors.InvalidRequest)
				if err != nil {
					general.Log.Error("Failed to send response to websocket: ", err)
					break
				}
			}
		}
	}()
}
