package chatHandler

import (
	"github.com/gorilla/websocket"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/chat"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/general"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/types"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/webserver/errors"
)

func HandleRoomMessagesPacket(conn *websocket.Conn, user *types.User, packet *types.ChatRoomMessagesPacket) error {
	messages, err := chat.ConnectionManager.GetMessagesFromChatRoom(packet.ChatId, user.Id)
	if err != nil {
		general.Log.Error("Failed to send response to websocket: ", err)
		return conn.WriteJSON(errors.InternalError)
	}

	packet.ChatMessages = messages
	return conn.WriteJSON(packet)
}
