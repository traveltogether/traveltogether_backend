package chatHandler

import (
	"github.com/forestgiant/sliceutil"
	"github.com/gorilla/websocket"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/chat"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/general"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/types"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/users"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/webserver/errors"
)

func HandleCreateRoom(conn *websocket.Conn, user *types.User, packet *types.ChatRoomCreatePacket) error {
	finalParticipants := packet.Information.Participants
	if !sliceutil.Contains(packet.Information.Participants, int64(user.Id)) {
		finalParticipants = append(packet.Information.Participants, int64(user.Id))
	}

	id, err := chat.ConnectionManager.CreateChatRoom(&finalParticipants, packet.Information.Group)
	if err != nil {
		if err == users.UserNotFound {
			return conn.WriteJSON(errors.UserNotFound)
		} else if err == chat.PrivateChatCanOnlyContainTwoUsers {
			return conn.WriteJSON(errors.PrivateChatCanOnlyContainTwoUsers)
		}
		general.Log.Error("Failed to send response to websocket", err)
		return conn.WriteJSON(errors.InternalError)
	}

	packet.Information.Id = id

	return conn.WriteJSON(packet)
}
