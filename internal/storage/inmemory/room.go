package inmemory

import (
	"context"
	"errors"
	"sync"

	"github.com/necutya/meeting-app/internal/models"

	"github.com/google/uuid"
)

type RoomRepo struct {
	Mutex sync.RWMutex
	Map   map[uuid.UUID][]models.RoomParticipant
}

func NewRoomRepo() *RoomRepo {
	rr := &RoomRepo{}
	rr.Map = make(map[uuid.UUID][]models.RoomParticipant)
	return rr
}

func (rr *RoomRepo) GetOne(ctx context.Context, roomID uuid.UUID) (*models.Room, error) {
	rr.Mutex.Lock()
	defer rr.Mutex.Unlock()
	participants, ok := rr.Map[roomID]
	if !ok {
		return nil, errors.New("no such a room")
	}
	return &models.Room{
		UUID:         roomID,
		Participants: participants,
	}, nil
}

func (rr *RoomRepo) Create(ctx context.Context) (*models.Room, error) {
	rr.Mutex.Lock()
	defer rr.Mutex.Unlock()
	roomID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	rr.Map[roomID] = []models.RoomParticipant{}
	return &models.Room{
		UUID:         roomID,
		Participants: rr.Map[roomID],
	}, nil
}

func (rr *RoomRepo) InsertParticipant(
	ctx context.Context,
	roomID uuid.UUID,
	newParticipant models.RoomParticipant,
) (*models.Room, error) {
	rr.Mutex.Lock()
	defer rr.Mutex.Unlock()

	for _, participant := range rr.Map[roomID] {
		if participant.Username == newParticipant.Username {
			return nil, errors.New("user with such username already exists")
		}
	}

	rr.Map[roomID] = append(rr.Map[roomID], newParticipant)
	return &models.Room{
		UUID:         roomID,
		Participants: rr.Map[roomID],
	}, nil
}

func (rr *RoomRepo) RemoveParticipant(
	ctx context.Context,
	roomID uuid.UUID,
	participantToDelete models.RoomParticipant,
) error {
	rr.Mutex.Lock()
	defer rr.Mutex.Unlock()

	for index, participant := range rr.Map[roomID] {
		if participant.Username == participantToDelete.Username {
			rr.Map[roomID] = append(rr.Map[roomID][:index], rr.Map[roomID][index+1:]...)
			return nil
		}
	}

	return errors.New("no user with such username")
}

func (rr *RoomRepo) DeleteRoom(ctx context.Context, roomID uuid.UUID) error {
	rr.Mutex.Lock()
	defer rr.Mutex.Unlock()

	_, ok := rr.Map[roomID]
	if !ok {
		return errors.New("no room with such uuid")
	}

	delete(rr.Map, roomID)

	return nil
}
