package service

import (
	"context"

	"server/internal/models"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

func (s *Service) CreateRoom(ctx context.Context) (*models.Room, error){

	room, err := s.roomRepo.Create(ctx)
	if err != nil {
		return nil, err
	}

	return room, err
}
var broadcast = make(chan broadcastMsg)


func (s *Service) JoinRoom(ctx context.Context, roomID uuid.UUID, participant models.RoomParticipant){
	room, err := s.roomRepo.InsertParticipant(ctx, roomID, participant)
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster(room.Participants)

	for {
		var msg broadcastMsg

		err = participant.Conn.ReadJSON(&msg.Message)
		if err != nil {
			log.Error("Read Error: ", err)
			participant.Conn.Close()
			break
		}

		msg.Client = participant.Conn
		msg.RoomID = roomID.String()

		log.Println(msg.Message)

		broadcast <- msg
	}
}

type broadcastMsg struct {
	Message map[string]interface{}
	RoomID  string
	Client    *websocket.Conn
}


func broadcaster(participants []models.RoomParticipant) {
	for {
		msg := <- broadcast

		for _, client := range participants {
			if client.Conn != msg.Client {
				err := client.Conn.WriteJSON(msg.Message)

				if err != nil {
					log.Error(err)
					client.Conn.Close()
				}
			}
		}
	}
}