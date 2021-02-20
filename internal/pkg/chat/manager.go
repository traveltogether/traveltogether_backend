package chat

import (
	"errors"
	"github.com/forestgiant/sliceutil"
	"github.com/gorilla/websocket"
	"github.com/lib/pq"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/types"
)

var (
	InvalidChatPacket = errors.New("invalid chat packet")
	NoPermission      = errors.New("no permission")

	ConnectionManager = &Manager{connections: map[int]*websocket.Conn{}}
)

type Manager struct {
	connections map[int]*websocket.Conn
}

func (manager *Manager) AddConnection(userId int, conn *websocket.Conn) {
	manager.connections[userId] = conn
}

func (manager *Manager) RemoveConnection(userId int) {
	delete(manager.connections, userId)
}

func (manager *Manager) SendMessage(message *types.ChatMessage, userId int, readBySender bool) (int, error) {
	if message.ChatId == nil {
		if message.ReceiverId == nil {
			return -1, InvalidChatPacket
		}

		err := manager.addChatIdOfPrivateChatToMessageObject(message, userId)

		if err != nil {
			return -1, err
		}
	}

	message.SenderId = &userId

	readBy := pq.Int64Array{}
	if readBySender {
		readBy = append(readBy, int64(userId))
	}

	id, err := saveMessageIntoDatabase(message.Message, *message.ChatId, *message.SenderId,
		&readBy, message.Time)
	if err != nil {
		return -1, err
	}

	var receivers *pq.Int64Array
	if message.ReceiverId != nil {
		receivers = &pq.Int64Array{int64(*message.ReceiverId)}
		message.ReceiverId = nil
	} else {
		receivers, err = getParticipantsOfChatRoomAsUser(*message.ChatId, userId)
	}

	if receivers == nil {
		return id, nil
	}

	for receiver := range *receivers {
		if receiver == *message.SenderId {
			continue
		}

		if connection, ok := manager.connections[receiver]; ok {
			message.Id = &id

			packet := types.CreateChatMessagePacket()
			packet.ChatMessage = *message
			err = connection.WriteJSON(packet)

			if err != nil {
				return id, err
			}

			err = setMessageReadInDatabase(id, receiver)
			if err != nil {
				return id, err
			}
		}
	}

	return id, nil
}

func (manager *Manager) CreateChatRoom(participants *pq.Int64Array, group bool) (int, error) {
	return createChatRoomInDatabase(participants, group)
}

func (manager *Manager) AddUserToChatRoom(information *types.ChatRoomAddUserInformation, userId int) error {
	participants, err := getParticipantsOfChatRoomAsUser(information.ChatId, userId)
	if err != nil {
		return err
	}

	if !sliceutil.Contains(participants, userId) {
		return NoPermission
	}

	return addUserToChatRoomInDatabase(information.ChatId, information.UserId)
}

func (manager *Manager) LeaveChatRoom(chatId int, userId int) error {
	return leaveChatRoomFromDatabase(chatId, userId)
}

func (manager *Manager) GetUnreadMessages(userId int) ([]*types.ChatMessage, error) {
	return getUnreadMessagesFromDatabase(userId)
}

func (manager *Manager) GetMessagesFromChatRoom(chatId int, userId int) ([]*types.ChatMessage, error) {
	return getMessagesOfChatRoomFromDatabase(chatId, userId)
}

func (manager *Manager) GetChatRoomsJoinedByUser(userId int) ([]*types.ChatRoom, error) {
	return getChatRoomsAsUser(userId)
}

func (manager *Manager) addChatIdOfPrivateChatToMessageObject(message *types.ChatMessage, userId int) error {
	chatId, err := getPrivateChatIdOfUsersFromDatabase(userId, *message.ReceiverId)
	if err != nil {
		return err
	}

	if chatId == -1 {
		chatId, err = manager.CreateChatRoom(&pq.Int64Array{int64(userId), int64(*message.ReceiverId)}, false)
		if err != nil {
			return err
		}
	}

	*message.ChatId = chatId
	return nil
}
