package service

import (
	"context"
	"sync"

	"server/config"
	"server/internal/models"

	"github.com/google/uuid"
)

var (
	service *Service
	once    sync.Once
)

type Service struct {
	cfg      *config.Config
	roomRepo RoomRepo
	// redis RedisCli
}

func New(
	cfg *config.Config,
	roomRepo RoomRepo,
	// rds RedisCli,
) *Service {
	once.Do(func() {
		service = &Service{
			// redis: rds,
			cfg: cfg,
			roomRepo: roomRepo,
		}
	})

	return service
}

func Get() *Service {
	return service
}

type RoomRepo interface {
	GetOne(ctx context.Context, roomID uuid.UUID) (*models.Room, error)
	Create(ctx context.Context) (*models.Room, error)
	InsertParticipant(ctx context.Context, roomID uuid.UUID, newParticipant models.RoomParticipant) (*models.Room, error)
	RemoveParticipant(ctx context.Context, roomID uuid.UUID, participantToDelete models.RoomParticipant) error
	DeleteRoom(ctx context.Context, roomID uuid.UUID) error
}
