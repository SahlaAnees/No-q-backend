package interfaces

import (
	"context"
	"no-q-solution/domain/entities"
)

type MerchantRepository interface {
	GetAll(ctx context.Context, paginator entities.Paginator) ([]entities.Merchant, error)
	GetCategories(ctx context.Context) ([]entities.Category, error)
	GetByCategory(ctx context.Context, category string) ([]entities.Merchant, error)
	GetSingle(ctx context.Context, id int64) (*entities.Merchant, error)
	Search(ctx context.Context, input string) ([]entities.Merchant, error)
	Create(ctx context.Context, merchant entities.Merchant) (int64, error)
	Login(ctx context.Context, login entities.Login) (*entities.Merchant, error)
	CreateToken(ctx context.Context, merchantID int64, token string) (string, error)
	ValidateToken(ctx context.Context, token string) (int64, error)
	Logout(ctx context.Context, merchantID int64) (bool, error)
	Delete(ctx context.Context, id int64) (bool, error)
}
