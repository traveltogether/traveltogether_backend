package chatHandler

import (
	"github.com/gorilla/websocket"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/chat"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/general"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/types"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/users"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/webserver/errors"
)

func HandleMessagePacket(conn *websocket.Conn, user *types.User, packet *types.ChatMessagePacket) error {
	id, err := chat.ConnectionManager.SendMessage(&packet.ChatMessage, user.Id, true)
	if err != nil {
		if err == chat.RoomDoesNotExist {
			return conn.WriteJSON(errors.RoomDoesNotExist)
		} else if err == chat.NotInRoom {
			return conn.WriteJSON(errors.NotInRoom)
		} else if err == users.UserNotFound {
			return conn.WriteJSON(errors.UserNotFound)
		} else {
			general.Log.Error("Failed to send response to websocket", err)
			if id == -1 {
				return conn.WriteJSON(errors.InternalError)
			}
		}
	}
	return conn.WriteJSON(packet)
}
