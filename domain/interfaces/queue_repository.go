package interfaces

import (
	"context"
	"no-q-solution/domain/entities"
	"time"
)

type QueueRepository interface {
	GetByMerchant(ctx context.Context, merchantID int64) ([]entities.Queue, error)
	GetSlotsByDate(ctx context.Context, queueID int64, date time.Time) (entities.Queue, error)
	MakeItAvailable(ctx context.Context, merchantID int64, queueID int64) (bool, error)
	MakeItUnAvailable(ctx context.Context, merchantID int64, queueID int64) (bool, error)
	MakeDatesAvailable(ctx context.Context, queueID int64, dates []time.Time) (bool, error)
	MakeDatesUnAvailable(ctx context.Context, queueID int64, dates []time.Time) (bool, error)
	Create(ctx context.Context, queue entities.Queue) (entities.Queue, error)
	ReserveSlot(ctx context.Context, reserve entities.ReservedSlots) (entities.ReservedSlots, error)
	UnReserveSlot(ctx context.Context, tokenNo int64) (bool, error)
	Delete(ctx context.Context, merchantID int64, queueID int64) (bool, error)
	IsQueueBelongsToMerchant(ctx context.Context, merchantID int64, queueID int64) (bool, error)
}
