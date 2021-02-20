package chat

import (
	"errors"
	"github.com/lib/pq"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/database"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/types"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/users"
)

var (
	NotInRoom                           = errors.New("not in room")
	RoomDoesNotExistOrUserAlreadyInRoom = errors.New("room does not exist or user is already in room")
	RoomDoesNotExist                    = errors.New("room does not exist")
	RoomIsNotAGroup                     = errors.New("room is not a group")
	PrivateChatCanOnlyContainTwoUsers   = errors.New("private chat can only contain two users")
)

func saveMessageIntoDatabase(message string, chatId int, senderId int, readBy *pq.Int64Array, time int) (int, error) {
	roomExists, err := doesRoomExist(chatId)
	if err != nil {
		return -1, err
	}

	if !roomExists {
		return -1, RoomDoesNotExist
	}

	participants, err := getParticipantsOfChatRoomAsUser(chatId, senderId)
	if err != nil {
		return -1, err
	}

	if len(*participants) == 0 {
		return -1, NotInRoom
	}

	slice, err := database.QueryAsync(database.DefaultTimeout, types.IdInformationType,
		"INSERT INTO "+
			"chat_messages(chat_id, message, sender, read_by, time) "+
			"VALUES($1, $2, $3, $4, $5) RETURNING id",
		chatId, message, senderId, readBy, time)

	if err != nil {
		return -1, err
	}

	idSlice := slice.([]*types.IdInformation)

	return idSlice[0].Id, nil
}

func setMessageReadInDatabase(messageId int, userId int) error {
	return database.PrepareAsync(database.DefaultTimeout, "UPDATE chat_messages SET read_by = "+
		"array_append(read_by, $1) WHERE id = $2 AND (read_by IS NULL OR NOT $3 = ANY(read_by))",
		userId, messageId, userId)
}

func createChatRoomInDatabase(participants *pq.Int64Array, group bool) (int, error) {
	if !group && len(*participants) > 2 {
		return -1, PrivateChatCanOnlyContainTwoUsers
	}

	for userId := range *participants {
		_, err := users.GetUserById(userId)
		if err != nil {
			return -1, err
		}
	}

	slice, err := database.QueryAsync(database.DefaultTimeout, types.IdInformationType,
		"INSERT INTO "+
			"chat_rooms(participants, group_chat) "+
			"VALUES($1, $2) RETURNING id",
		participants, group)

	if err != nil {
		return -1, err
	}

	idSlice := slice.([]*types.IdInformation)

	return idSlice[0].Id, nil
}

func isRoomGroup(chatId int) (bool, error) {
	slice, err := database.QueryAsync(database.DefaultTimeout, types.GroupInformationType,
		"SELECT group_chat FROM chat_rooms WHERE id = $1",
		chatId)

	if err != nil {
		return false, err
	}

	groupInformation := slice.([]*types.GroupInformation)
	if len(groupInformation) == 0 {
		return false, RoomDoesNotExist
	}

	return groupInformation[0].Group, nil
}

func doesRoomExist(chatId int) (bool, error) {
	slice, err := database.QueryAsync(database.DefaultTimeout, types.IdInformationType,
		"SELECT id FROM chat_rooms WHERE id = $1",
		chatId)

	if err != nil {
		return false, err
	}

	if len(slice.([]*types.IdInformation)) == 0 {
		return false, nil
	}

	return true, nil
}

func leaveChatRoomFromDatabase(chatId int, userId int) error {
	isGroup, err := isRoomGroup(chatId)
	if err != nil {
		return err
	}

	if !isGroup {
		return RoomIsNotAGroup
	}

	slice, err := database.QueryAsync(database.DefaultTimeout, types.IdInformationType,
		"UPDATE chat_rooms "+
			"SET participants = array_remove(participants, $1) "+
			"WHERE id = $2 "+
			"AND participants IS NOT NULL AND $3 = ANY(participants) "+
			"RETURNING id",
		userId, chatId, userId)

	if err != nil {
		return err
	}

	if len(slice.([]*types.IdInformation)) == 0 {
		return NotInRoom
	}

	slice, err = database.QueryAsync(database.DefaultTimeout, types.Int64ArrayType,
		"SELECT participants FROM chat_rooms WHERE id = $1",
		chatId)
	if err != nil {
		return err
	}

	if len(slice.(pq.Int64Array)) == 0 {
		return database.PrepareAsync(database.DefaultTimeout, "DELETE FROM chat_rooms WHERE id = $1", chatId)
	}

	return nil
}

func addUserToChatRoomInDatabase(chatId int, userId int) error {
	isGroup, err := isRoomGroup(chatId)
	if err != nil {
		return err
	}

	if !isGroup {
		return RoomIsNotAGroup
	}

	_, err = users.GetUserById(userId)
	if err != nil {
		return err
	}

	slice, err := database.QueryAsync(database.DefaultTimeout, types.IdInformationType,
		"UPDATE chat_rooms "+
			"SET participants = array_append(participants, $1) WHERE id = $2 RETURNING id", userId, chatId)

	if err != nil {
		return err
	}

	if len(slice.([]*types.IdInformation)) == 0 {
		return RoomDoesNotExistOrUserAlreadyInRoom
	}

	return nil
}

func getParticipantsOfChatRoom(chatId int) (*pq.Int64Array, error) {
	slice, err := database.QueryAsync(database.DefaultTimeout, types.Int64ArrayType,
		"SELECT participants FROM chat_rooms WHERE id = $1 AND participants IS NOT NULL",
		chatId)

	if err != nil {
		return nil, err
	}

	return slice.(*pq.Int64Array), nil
}

func getParticipantsOfChatRoomAsUser(chatId int, userId int) (*pq.Int64Array, error) {
	slice, err := database.QueryAsync(database.DefaultTimeout, types.Int64ArrayType,
		"SELECT participants FROM chat_rooms WHERE id = $1 AND participants IS NOT NULL AND $2 = ANY(participants)",
		chatId, userId)

	if err != nil {
		return nil, err
	}

	return slice.(*pq.Int64Array), nil
}

func getUnreadMessagesFromDatabase(userId int) ([]*types.ChatMessage, error) {
	slice, err := database.QueryAsync(database.DefaultTimeout, types.ChatMessageType,
		"UPDATE chat_messages "+
			"SET read_by = array_append(read_by, $1) "+
			"FROM chat_rooms "+
			"WHERE chat_messages.chat_id = chat_rooms.id "+
			"AND chat_rooms.participants IS NOT NULL "+
			"AND $2 = ANY(chat_rooms.participants) "+
			"AND (chat_messages.read_by IS NULL OR NOT $3 = ANY(chat_messages.read_by)) "+
			"RETURNING chat_messages.id AS id, chat_messages.chat_id AS chat_id, "+
			"chat_messages.message AS message, chat.messages_sender_id AS sender_id, chat_messages.time AS time",
		userId, userId)

	if err != nil {
		return nil, err
	}

	return slice.([]*types.ChatMessage), nil
}

func getPrivateChatIdOfUsersFromDatabase(user1Id int, user2Id int) (int, error) {
	slice, err := database.QueryAsync(database.DefaultTimeout, types.IdInformationType,
		"SELECT id FROM chat_rooms "+
			"WHERE participants IS NOT NULL "+
			"AND $1 = ANY(participants) "+
			"AND $2 = ANY(participants) "+
			"AND group_chat = false",
		user1Id, user2Id)

	if err != nil {
		return -1, err
	}

	idSlice := slice.([]*types.IdInformation)
	if len(idSlice) == 0 {
		return -1, nil
	}

	return idSlice[0].Id, nil
}

func getMessagesOfChatRoomFromDatabase(chatId int, userId int) ([]*types.ChatMessage, error) {
	slice, err := database.QueryAsync(database.DefaultTimeout, types.ChatMessageType,
		"SELECT chat_messages.id, chat_messages.chat_id, chat_messages.sender_id, "+
			"chat_messages.message, chat_messages.time "+
			"FROM chat_messages, chat_rooms "+
			"WHERE chat_messages.chat_id = chat_rooms.id "+
			"AND chat_messages.chat_id = $1 "+
			"AND chat_rooms.participants IS NOT NULL "+
			"AND $2 = ANY(chat_rooms.participants)",
		chatId, userId)

	if err != nil {
		return nil, err
	}

	messages := slice.([]*types.ChatMessage)

	if len(messages) != 0 {
		err = database.PrepareAsync(database.DefaultTimeout,
			"UPDATE chat_messages "+
				"SET read_by = array_append(read_by, $1) "+
				"FROM chat_rooms "+
				"WHERE chat_messages.chat_id = chat_rooms.id "+
				"AND chat_messages.chat_id = $2 "+
				"AND chat_rooms.participants IS NOT NULL "+
				"AND $3 = ANY(chat_rooms.participants)",
			userId, userId)

		if err != nil {
			return nil, err
		}
	}

	return messages, nil
}

func getChatRoomsAsUser(userId int) ([]*types.ChatRoom, error) {
	slice, err := database.QueryAsync(database.DefaultTimeout, types.ChatRoomType,
		"SELECT * FROM chat_rooms WHERE participants IS NOT NULL AND $1 = ANY(participants)",
		userId)

	if err != nil {
		return nil, err
	}

	return slice.([]*types.ChatRoom), err
}
