package types

import (
	"github.com/lib/pq"
	"reflect"
)

const (
	ChatMessagePacketName        = "ChatMessagePacket"
	ChatRoomAddUserPacketName    = "ChatRoomAddUserPacket"
	ChatRoomLeaveUserPacketName  = "ChatRoomLeaveUserPacket"
	ChatRoomCreatePacketName     = "ChatRoomCreatePacket"
	ChatUnreadMessagesPacketName = "ChatUnreadMessagesPacket"
	ChatRoomMessagesPacketName   = "ChatRoomMessagesPacket"
	ChatRoomsPacketName          = "ChatRoomsPacket"
)

var (
	GroupInformationType        = reflect.TypeOf(&GroupInformation{})
	ChatMessageType             = reflect.TypeOf(&ChatMessage{})
	ChatRoomType                = reflect.TypeOf(&ChatRoom{})
	ParticipantsInformationType = reflect.TypeOf(&ParticipantsInformation{})
)

type GroupInformation struct {
	Group bool `db:"group_chat"`
}

type ParticipantsInformation struct {
	Participants pq.Int64Array `db:"participants"`
}

type ChatMessage struct {
	Id         *int   `json:"id" db:"id"`
	ChatId     *int   `json:"chat_id,omitempty" db:"chat_id"`
	SenderId   *int   `json:"sender_id,omitempty" db:"sender_id"`
	ReceiverId *int   `json:"user_id,omitempty" db:"user_id"`
	Message    string `json:"message" db:"message"`
	Time       int    `json:"time" db:"time"`
}

type ChatRoom struct {
	Id           int           `json:"id" db:"id"`
	Participants pq.Int64Array `json:"participants" db:"participants"`
	Group        bool          `json:"group" db:"group_chat"`
}

type ChatMessagePacket struct {
	Type        string      `json:"type"`
	ChatMessage ChatMessage `json:"chat_message"`
}

type ChatRoomAddUserPacket struct {
	Type        string                     `json:"type"`
	Information ChatRoomAddUserInformation `json:"information"`
}

type ChatRoomLeaveUserPacket struct {
	Type        string                       `json:"type"`
	Information ChatRoomLeaveUserInformation `json:"information"`
}

type ChatRoomCreatePacket struct {
	Type        string   `json:"type"`
	Information ChatRoom `json:"information"`
}

type ChatRoomAddUserInformation struct {
	ChatId int `json:"chat_id"`
	UserId int `json:"user_id"`
}

type ChatRoomLeaveUserInformation struct {
	ChatId int `json:"chat_id"`
}

type ChatUnreadMessagesPacket struct {
	Type         string         `json:"type"`
	ChatMessages []*ChatMessage `json:"chat_messages"`
}

type ChatRoomMessagesPacket struct {
	Type         string         `json:"type"`
	ChatId       int            `json:"chat_id"`
	ChatMessages []*ChatMessage `json:"chat_messages"`
}

type ChatRoomsPacket struct {
	Type      string      `json:"type"`
	ChatRooms []*ChatRoom `json:"chat_rooms"`
}

func CreateChatMessagePacket() *ChatMessagePacket {
	return &ChatMessagePacket{
		Type: ChatMessagePacketName,
	}
}

func CreateChatRoomAddUserPacket() *ChatRoomAddUserPacket {
	return &ChatRoomAddUserPacket{
		Type: ChatRoomAddUserPacketName,
	}
}

func CreateChatRoomLeaveUserPacket() *ChatRoomLeaveUserPacket {
	return &ChatRoomLeaveUserPacket{
		Type: ChatRoomLeaveUserPacketName,
	}
}

func CreateChatRoomCreatePacket() *ChatRoomCreatePacket {
	return &ChatRoomCreatePacket{
		Type: ChatRoomCreatePacketName,
	}
}

func CreateChatUnreadMessagesPacket() *ChatUnreadMessagesPacket {
	return &ChatUnreadMessagesPacket{
		Type: ChatUnreadMessagesPacketName,
	}
}

func CreateChatRoomMessagesPacket() *ChatRoomMessagesPacket {
	return &ChatRoomMessagesPacket{
		Type: ChatRoomMessagesPacketName,
	}
}

func CreateChatRoomsPacket() *ChatRoomsPacket {
	return &ChatRoomsPacket{
		Type: ChatRoomsPacketName,
	}
}
