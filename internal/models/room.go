package models

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type RoomParticipant struct {
	IsHost bool
	Username string
	Conn   *websocket.Conn
}

type Room struct {
	UUID uuid.UUID
	Participants []RoomParticipant
}

type CreateRoomHTTPResponse struct {
	RoomID uuid.UUID `json:"room_id"`
}

type JoinRoomHTPPRequest struct {
	Username string `json:"username"`
}