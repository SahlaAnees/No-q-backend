package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"no-q-solution/domain/entities"
	"no-q-solution/domain/interfaces"
)

type MerchantRepository struct {
	db *sql.DB
}

func NewMerchantRepository(db *sql.DB) interfaces.MerchantRepository {
	repo := &MerchantRepository{
		db: db,
	}

	return repo
}

func (repo MerchantRepository) GetAll(ctx context.Context, paginator entities.Paginator) ([]entities.Merchant, error) {
	offset := (paginator.Page - 1) * paginator.Size

	query := `
		SELECT id, name, category, facebook, instagram, website, created_at, updated_at
		FROM merchant ORDER BY id ASC LIMIT ? OFFSET ?;`

	stmt, err := repo.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, paginator.Size, offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	merchantes := make([]entities.Merchant, 0)

	for rows.Next() {
		merchant := entities.Merchant{}

		err := rows.Scan(
			&merchant.ID,
			&merchant.Name,
			&merchant.Category,
			&merchant.Facebook,
			&merchant.Instagram,
			&merchant.Website,
			&merchant.CreatedAt,
			&merchant.UpdatedAt,
		)

		if err != nil {
			log.Println(err)
			continue
		}

		merchantes = append(merchantes, merchant)
	}

	return merchantes, nil
}

func (repo MerchantRepository) GetCategories(ctx context.Context) ([]entities.Category, error) {

	query := `
		SELECT category, created_at, updated_at
		FROM category ORDER BY category ASC;`

	stmt, err := repo.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	categories := make([]entities.Category, 0)

	for rows.Next() {
		category := entities.Category{}

		err := rows.Scan(
			&category.Name,
			&category.CreatedAt,
			&category.UpdatedAt,
		)

		if err != nil {
			log.Println(err)
			continue
		}

		categories = append(categories, category)
	}

	return categories, nil
}

func (repo MerchantRepository) GetByCategory(ctx context.Context, category string) ([]entities.Merchant, error) {

	query := `
		SELECT id, name, category, facebook, instagram, website, created_at, updated_at
		FROM merchant WHERE category = ? ORDER BY id ASC;
	`

	stmt, err := repo.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, category)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	merchantes := make([]entities.Merchant, 0)

	for rows.Next() {
		merchant := entities.Merchant{}

		err := rows.Scan(
			&merchant.ID,
			&merchant.Name,
			&merchant.Category,
			&merchant.Facebook,
			&merchant.Instagram,
			&merchant.Website,
			&merchant.CreatedAt,
			&merchant.UpdatedAt,
		)

		if err != nil {
			log.Println(err)
			continue
		}

		merchantes = append(merchantes, merchant)
	}

	return merchantes, nil
}

func (repo MerchantRepository) GetSingle(ctx context.Context, id int64) (*entities.Merchant, error) {

	query := `SELECT EXISTS(SELECT 1 FROM merchant WHERE id = ?);`

	stmt, err := repo.db.PrepareContext(ctx, query)

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, id)

	var exists bool

	err = row.Scan(&exists)

	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, errors.New("there are no such merchant exists")
	}

	merchant := entities.Merchant{}

	query = `
	SELECT id, name, category, facebook, instagram, website, created_at, updated_at
	FROM merchant WHERE id = ?;
	`

	stmt, err = repo.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	row = stmt.QueryRowContext(ctx, id)

	err = row.Scan(
		&merchant.ID,
		&merchant.Name,
		&merchant.Category,
		&merchant.Facebook,
		&merchant.Instagram,
		&merchant.Website,
		&merchant.CreatedAt,
		&merchant.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &merchant, nil
}

func (repo MerchantRepository) Search(ctx context.Context, input string) ([]entities.Merchant, error) {
	query := fmt.Sprintf(`
		SELECT id, name, category, facebook, instagram, website, created_at, updated_at
		FROM merchant WHERE name LIKE '%%%s%%' ORDER BY id ASC;
	`, input)

	stmt, err := repo.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	merchantes := make([]entities.Merchant, 0)

	for rows.Next() {
		merchant := entities.Merchant{}

		err := rows.Scan(
			&merchant.ID,
			&merchant.Name,
			&merchant.Category,
			&merchant.Facebook,
			&merchant.Instagram,
			&merchant.Website,
			&merchant.CreatedAt,
			&merchant.UpdatedAt,
		)

		if err != nil {
			log.Println(err)
			continue
		}

		merchantes = append(merchantes, merchant)
	}

	return merchantes, nil
}

func (repo MerchantRepository) Create(ctx context.Context, merchant entities.Merchant) (int64, error) {

	query := `INSERT INTO merchant (category, name, email, password, facebook, instagram, website) VALUES (?, ?, ?, ?, ?, ?, ?);`

	stmt, err := repo.db.PrepareContext(ctx, query)
	if err != nil {
		return 0, err
	}

	defer stmt.Close()

	result, err := stmt.ExecContext(
		ctx,
		merchant.Category,
		merchant.Name,
		merchant.Email,
		merchant.Password,
		merchant.Facebook,
		merchant.Instagram,
		merchant.Website,
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (repo MerchantRepository) Login(ctx context.Context, login entities.Login) (*entities.Merchant, error) {
	merchant := entities.Merchant{}

	query := `
	SELECT id, name, category, facebook, instagram, website, created_at, updated_at
	FROM merchant WHERE email = ? and password = ?;
	`

	stmt, err := repo.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, login.Email, login.Password)

	err = row.Scan(
		&merchant.ID,
		&merchant.Name,
		&merchant.Category,
		&merchant.Facebook,
		&merchant.Instagram,
		&merchant.Website,
		&merchant.CreatedAt,
		&merchant.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &merchant, nil
}

func (repo MerchantRepository) CreateToken(ctx context.Context, merchantID int64, token string) (string, error) {

	query := `INSERT INTO token (merchant_id, auth_token) VALUES (?, ?);`

	stmt, err := repo.db.PrepareContext(ctx, query)
	if err != nil {
		return "", err
	}

	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, merchantID, token)
	if err != nil {
		return "", err
	}

	_, err = result.LastInsertId()
	if err != nil {
		return "", err
	}

	return token, nil
}

func (repo MerchantRepository) ValidateToken(ctx context.Context, token string) (int64, error) {

	query := `SELECT merchant_id FROM token WHERE auth_token = ?;`

	stmt, err := repo.db.PrepareContext(ctx, query)
	if err != nil {
		return 0, err
	}

	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, token)

	var merchantID int64

	err = row.Scan(&merchantID)

	if err == sql.ErrNoRows {
		return 0, errors.New("token is not valid")
	}

	if err != nil {
		return 0, err
	}

	return merchantID, nil
}

func (repo MerchantRepository) Logout(ctx context.Context, merchantID int64) (bool, error) {

	query := `DELETE FROM token WHERE merchant_id = ?;`

	stmt, err := repo.db.PrepareContext(ctx, query)

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(merchantID)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (repo MerchantRepository) Delete(ctx context.Context, id int64) (bool, error) {

	query := `SELECT EXISTS(SELECT 1 FROM merchant WHERE id = ?);`

	stmt, err := repo.db.PrepareContext(ctx, query)

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, id)

	var exists bool

	err = row.Scan(&exists)

	if err != nil {
		return false, err
	}

	if !exists {
		return false, errors.New("there are no such merchant exists")
	}

	query = `DELETE FROM merchant WHERE id = ?;`

	stmt, err = repo.db.PrepareContext(ctx, query)

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return false, err
	}

	return true, nil
}
