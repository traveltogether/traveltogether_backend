package chatHandler

import (
	"github.com/gorilla/websocket"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/chat"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/general"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/types"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/users"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/webserver/errors"
)

func HandleRoomAddUserPacket(conn *websocket.Conn, user *types.User, packet *types.ChatRoomAddUserPacket) error {
	err := chat.ConnectionManager.AddUserToChatRoom(&packet.Information, user.Id)
	if err != nil {
		if err == chat.RoomIsNotAGroup {
			return conn.WriteJSON(errors.RoomIsNotAGroup)
		} else if err == chat.RoomDoesNotExistOrUserAlreadyInRoom {
			return conn.WriteJSON(errors.RoomDoesNotExistOrUserIsAlreadyInRoom)
		} else if err == chat.NoPermission {
			return conn.WriteJSON(errors.Forbidden)
		} else if err == users.UserNotFound {
			return conn.WriteJSON(errors.UserNotFound)
		} else {
			general.Log.Error("Failed to send response to websocket: ", err)
			return conn.WriteJSON(errors.InternalError)
		}
	}
	return conn.WriteJSON(packet)
}

func HandleRoomLeaveUserPacket(conn *websocket.Conn, user *types.User, packet *types.ChatRoomLeaveUserPacket) error {
	err := chat.ConnectionManager.LeaveChatRoom(packet.Information.ChatId, user.Id)
	if err != nil {
		if err == chat.RoomIsNotAGroup {
			return conn.WriteJSON(errors.RoomIsNotAGroup)
		} else if err == chat.NotInRoom {
			return conn.WriteJSON(errors.NotInRoom)
		} else {
			general.Log.Error("Failed to send response to websocket: ", err)
			return conn.WriteJSON(errors.InternalError)
		}
	}
	return conn.WriteJSON(packet)
}
