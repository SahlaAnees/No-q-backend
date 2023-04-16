package repositories

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"no-q-solution/domain/entities"
	"no-q-solution/domain/interfaces"
	"strings"
	"time"
)

type QueueRepository struct {
	db *sql.DB
}

func NewQueueRepository(db *sql.DB) interfaces.QueueRepository {
	repo := &QueueRepository{
		db: db,
	}

	return repo
}

func (repo QueueRepository) GetByMerchant(ctx context.Context, merchantID int64) ([]entities.Queue, error) {

	query := `
		SELECT q.id, q.name, q.merchant_id, q.intervals, q.start_time, q.end_time, q.is_available, GROUP_CONCAT(ua.date) as unavailable_dates, q.created_at 
		FROM queue q LEFT JOIN unavailable ua on q.id = ua.queue_id WHERE merchant_id = ? 
		GROUP BY q.id;`

	stmt, err := repo.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, merchantID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	queues := make([]entities.Queue, 0)

	for rows.Next() {
		queue := entities.Queue{}

		var unAvailableDates sql.NullString

		err := rows.Scan(
			&queue.ID,
			&queue.Name,
			&queue.MerchantID,
			&queue.Interval,
			&queue.StartTime,
			&queue.EndTime,
			&queue.IsAvailable,
			&unAvailableDates,
			&queue.CreatedAt,
		)

		if err != nil {
			log.Println(err)
			continue
		}

		var dates []time.Time

		if unAvailableDates.Valid {

			for _, unAvailableDate := range strings.Split(unAvailableDates.String, ",") {
				date, err := time.Parse("2006-01-02 15:04:05", unAvailableDate) // use the format of your date/time string
				if err != nil {
					continue
				}

				dates = append(dates, date)
			}

			queue.UnavailableDates = dates
		}

		queues = append(queues, queue)
	}

	return queues, nil
}

func (repo QueueRepository) GetSlotsByDate(ctx context.Context, queueID int64, date time.Time) (entities.Queue, error) {

	query := `SELECT EXISTS(SELECT 1 FROM queue WHERE id = ?);`

	stmt, err := repo.db.PrepareContext(ctx, query)

	if err != nil {
		return entities.Queue{}, err
	}

	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, queueID)

	var exists bool

	err = row.Scan(&exists)

	if err != nil {
		return entities.Queue{}, err
	}

	if !exists {
		return entities.Queue{}, errors.New("there are no such queue exists")
	}

	queue := entities.Queue{}

	queue.ID = queueID

	query = `
		SELECT rs.token_no, rs.queue_id, rs.start_time, rs.end_time, rs.created_at, u.id, u.name, u.phone, u.email  
		FROM reserved_slots rs INNER JOIN user u on rs.reserved_by = u.id
		WHERE rs.queue_id = ? AND DATE(rs.start_time) = DATE(?);`

	stmt, err = repo.db.PrepareContext(ctx, query)
	if err != nil {
		return entities.Queue{}, err
	}

	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, queueID, date)
	if err != nil {
		return entities.Queue{}, err
	}

	defer rows.Close()

	reservedSlots := make([]entities.ReservedSlots, 0)

	for rows.Next() {

		reservedSlot := entities.ReservedSlots{}
		user := entities.User{}

		err := rows.Scan(
			&reservedSlot.TokenNo,
			&reservedSlot.QueueID,
			&reservedSlot.StartTime,
			&reservedSlot.EndTime,
			&reservedSlot.CreatedAt,
			&user.ID,
			&user.Name,
			&user.Phone,
			&user.Email,
		)

		if err != nil {
			log.Println(err)
			continue
		}

		reservedSlot.ReservedBy = user

		reservedSlots = append(reservedSlots, reservedSlot)
	}

	queue.ReservedSlots = reservedSlots

	return queue, nil
}

func (repo QueueRepository) MakeItAvailable(ctx context.Context, merchantID int64, queueID int64) (bool, error) {

	query := `UPDATE queue SET is_available = 1 WHERE id = ? AND merchant_id = ?;`

	stmt, err := repo.db.PrepareContext(ctx, query)
	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(queueID, merchantID)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (repo QueueRepository) MakeItUnAvailable(ctx context.Context, merchantID int64, queueID int64) (bool, error) {

	query := `UPDATE queue SET is_available = 0 WHERE id = ? AND merchant_id = ?;`

	stmt, err := repo.db.PrepareContext(ctx, query)
	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(queueID, merchantID)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (repo QueueRepository) MakeDatesUnAvailable(ctx context.Context, queueID int64, dates []time.Time) (bool, error) {

	query := `INSERT INTO unavailable (queue_id, date) VALUES (?, ?)`

	stmt, err := repo.db.PrepareContext(ctx, query)
	if err != nil {
		return false, err
	}

	defer stmt.Close()

	for _, date := range dates {

		result, err := stmt.ExecContext(ctx, queueID, date)
		if err != nil {
			continue
		}

		_, err = result.LastInsertId()
		if err != nil {
			continue
		}

	}

	return true, nil
}

func (repo QueueRepository) MakeDatesAvailable(ctx context.Context, queueID int64, dates []time.Time) (bool, error) {

	query := `DELETE FROM unavailable WHERE queue_id = ? AND date = ?;`

	stmt, err := repo.db.PrepareContext(ctx, query)

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	for _, date := range dates {

		_, err = stmt.Exec(queueID, date)
		if err != nil {
			return false, err
		}
	}

	return true, nil
}

func (repo QueueRepository) Create(ctx context.Context, queue entities.Queue) (entities.Queue, error) {

	query := `INSERT INTO queue (merchant_id, name, intervals, start_time, end_time) VALUES (?, ?, ?, ?, ?);`

	stmt, err := repo.db.PrepareContext(ctx, query)
	if err != nil {
		return entities.Queue{}, err
	}

	defer stmt.Close()

	result, err := stmt.ExecContext(
		ctx,
		queue.MerchantID,
		queue.Name,
		queue.Interval,
		queue.StartTime,
		queue.EndTime,
	)
	if err != nil {
		return entities.Queue{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return entities.Queue{}, err
	}

	queue.ID = id

	return queue, nil
}

func (repo QueueRepository) ReserveSlot(ctx context.Context, reserve entities.ReservedSlots) (entities.ReservedSlots, error) {

	query := `SELECT EXISTS(SELECT 1 FROM queue WHERE id = ?);`

	stmt, err := repo.db.PrepareContext(ctx, query)

	if err != nil {
		return entities.ReservedSlots{}, err
	}

	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, reserve.QueueID)

	var exists bool

	err = row.Scan(&exists)

	if err != nil {
		return entities.ReservedSlots{}, err
	}

	if !exists {
		return entities.ReservedSlots{}, errors.New("there are no such queue exists")
	}

	query = `
        SELECT COUNT(*) 
        FROM reserved_slots 
        WHERE (? <= end_time) AND (? >= start_time) AND queue_id = ?
    `

	stmt, err = repo.db.PrepareContext(ctx, query)
	if err != nil {
		return entities.ReservedSlots{}, err
	}

	defer stmt.Close()

	var count int

	err = stmt.QueryRowContext(ctx, reserve.StartTime, reserve.EndTime, reserve.QueueID).Scan(&count)
	if err != nil {
		return entities.ReservedSlots{}, err
	}
	if count > 0 {
		return entities.ReservedSlots{}, errors.New("the slot already reserved")
	}

	query = `INSERT INTO user (name, phone, email) VALUES (?, ?, ?);`

	stmt, err = repo.db.PrepareContext(ctx, query)
	if err != nil {
		return entities.ReservedSlots{}, err
	}

	defer stmt.Close()

	result, err := stmt.ExecContext(
		ctx,
		reserve.ReservedBy.Name,
		reserve.ReservedBy.Phone,
		reserve.ReservedBy.Email,
	)
	if err != nil {
		return entities.ReservedSlots{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return entities.ReservedSlots{}, err
	}

	reserve.ReservedBy.ID = id

	query = `INSERT INTO reserved_slots (queue_id, start_time, end_time, reserved_by) VALUES (?, ?, ?, ?);`

	stmt, err = repo.db.PrepareContext(ctx, query)
	if err != nil {
		return entities.ReservedSlots{}, err
	}

	defer stmt.Close()

	result, err = stmt.ExecContext(
		ctx,
		reserve.QueueID,
		reserve.StartTime,
		reserve.EndTime,
		id,
	)
	if err != nil {
		return entities.ReservedSlots{}, err
	}

	id, err = result.LastInsertId()
	if err != nil {
		return entities.ReservedSlots{}, err
	}

	reserve.TokenNo = id

	return reserve, nil
}

func (repo QueueRepository) UnReserveSlot(ctx context.Context, tokenNo int64) (bool, error) {

	query := `SELECT EXISTS(SELECT 1 FROM reserved_slots WHERE token_no = ?);`

	stmt, err := repo.db.PrepareContext(ctx, query)

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, tokenNo)

	var exists bool

	err = row.Scan(&exists)

	if err != nil {
		return false, err
	}

	if !exists {
		return false, errors.New("there are no such token no")
	}

	query = `DELETE FROM reserved_slots WHERE token_no = ?;`

	stmt, err = repo.db.PrepareContext(ctx, query)

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(tokenNo)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (repo QueueRepository) Delete(ctx context.Context, merchantID int64, queueID int64) (bool, error) {

	query := `DELETE FROM queue WHERE id = ? AND merchant_id = ?;`

	stmt, err := repo.db.PrepareContext(ctx, query)

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(queueID, merchantID)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (repo QueueRepository) IsQueueBelongsToMerchant(ctx context.Context, merchantID int64, queueID int64) (bool, error) {

	query := `SELECT EXISTS(SELECT 1 FROM queue WHERE id = ? AND merchant_id =?);`

	stmt, err := repo.db.PrepareContext(ctx, query)

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, queueID, merchantID)

	var exists bool

	err = row.Scan(&exists)

	if err != nil {
		return false, err
	}

	if !exists {
		return false, errors.New("queue is not blongs to merchant")
	}

	return exists, nil
}
