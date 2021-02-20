package chatHandler

import (
	"github.com/gorilla/websocket"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/chat"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/general"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/types"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/webserver/errors"
)

func HandleRoomsPacket(conn *websocket.Conn, user *types.User, packet *types.ChatRoomsPacket) error {
	rooms, err := chat.ConnectionManager.GetChatRoomsJoinedByUser(user.Id)
	if err != nil {
		general.Log.Error("Failed to send response to websocket", err)
		return conn.WriteJSON(errors.InternalError)
	}

	packet.ChatRooms = rooms
	return conn.WriteJSON(packet)
}
