package usecases

import (
	"context"
	"errors"
	"no-q-solution/domain/entities"
	"no-q-solution/domain/interfaces"
	"time"
)

type QueuetUsecase struct {
	repo interfaces.QueueRepository
}

func NewQueuetUsecase(repo interfaces.QueueRepository) QueuetUsecase {
	usecase := QueuetUsecase{
		repo: repo,
	}

	return usecase
}

func (usecase QueuetUsecase) GetByMerchant(ctx context.Context, merchantID int64) ([]entities.Queue, error) {

	return usecase.repo.GetByMerchant(ctx, merchantID)
}

func (usecase QueuetUsecase) GetSlotsByDate(ctx context.Context, queueID int64, date time.Time) (entities.Queue, error) {

	return usecase.repo.GetSlotsByDate(ctx, queueID, date)
}

func (usecase QueuetUsecase) MakeItAvailable(ctx context.Context, merchantID int64, queueID int64) (bool, error) {

	_, err := usecase.repo.IsQueueBelongsToMerchant(ctx, merchantID, queueID)
	if err != nil {
		return false, err
	}

	return usecase.repo.MakeItAvailable(ctx, merchantID, queueID)
}

func (usecase QueuetUsecase) MakeItUnAvailable(ctx context.Context, merchantID int64, queueID int64) (bool, error) {

	_, err := usecase.repo.IsQueueBelongsToMerchant(ctx, merchantID, queueID)
	if err != nil {
		return false, err
	}

	return usecase.repo.MakeItUnAvailable(ctx, merchantID, queueID)
}

func (usecase QueuetUsecase) MakeDatesAvailable(ctx context.Context, merchantID int64, queueID int64, dates []time.Time) (bool, error) {

	_, err := usecase.repo.IsQueueBelongsToMerchant(ctx, merchantID, queueID)
	if err != nil {
		return false, err
	}

	return usecase.repo.MakeDatesAvailable(ctx, queueID, dates)
}

func (usecase QueuetUsecase) MakeDatesUnAvailable(ctx context.Context, merchantID int64, queueID int64, dates []time.Time) (bool, error) {

	_, err := usecase.repo.IsQueueBelongsToMerchant(ctx, merchantID, queueID)
	if err != nil {
		return false, err
	}

	return usecase.repo.MakeDatesUnAvailable(ctx, queueID, dates)
}

func (usecase QueuetUsecase) Create(ctx context.Context, queue entities.Queue) (entities.Queue, error) {

	if len(queue.Name) == 0 {
		return entities.Queue{}, errors.New("name cannot be emtpy")
	}

	if queue.StartTime.After(queue.EndTime) {
		return entities.Queue{}, errors.New("given time range is wrong")
	}

	return usecase.repo.Create(ctx, queue)
}

func (usecase QueuetUsecase) ReserveSlot(ctx context.Context, reserve entities.ReservedSlots) (entities.ReservedSlots, error) {

	return usecase.repo.ReserveSlot(ctx, reserve)
}

func (usecase QueuetUsecase) UnReserveSlot(ctx context.Context, tokenNo int64) (bool, error) {

	return usecase.repo.UnReserveSlot(ctx, tokenNo)
}

func (usecase QueuetUsecase) Delete(ctx context.Context, merchantID int64, queueID int64) (bool, error) {

	_, err := usecase.repo.IsQueueBelongsToMerchant(ctx, merchantID, queueID)
	if err != nil {
		return false, err
	}

	return usecase.repo.Delete(ctx, merchantID, queueID)
}
