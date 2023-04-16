package usecases

import (
	"context"
	"errors"
	"no-q-solution/domain/entities"
	"no-q-solution/domain/interfaces"
	"regexp"

	"github.com/google/uuid"
)

type MerchantUsecase struct {
	repo interfaces.MerchantRepository
}

func NewMerchantUsecase(repo interfaces.MerchantRepository) MerchantUsecase {
	usecase := MerchantUsecase{
		repo: repo,
	}

	return usecase
}

func (usecase MerchantUsecase) GetAll(ctx context.Context, paginator entities.Paginator) ([]entities.Merchant, error) {

	return usecase.repo.GetAll(ctx, paginator)
}

func (usecase MerchantUsecase) GetCategories(ctx context.Context) ([]entities.Category, error) {
	return usecase.repo.GetCategories(ctx)
}

func (usecase MerchantUsecase) GetByCategory(ctx context.Context, category string) ([]entities.Merchant, error) {

	return usecase.repo.GetByCategory(ctx, category)
}

func (usecase MerchantUsecase) GetSingle(ctx context.Context, id int64) (entities.Merchant, error) {

	merchant, err := usecase.repo.GetSingle(ctx, id)
	if err != nil {
		return entities.Merchant{}, nil
	}

	if merchant == nil {
		return entities.Merchant{}, errors.New("not found")
	}

	return *merchant, nil
}

func (usecase MerchantUsecase) Search(ctx context.Context, input string) ([]entities.Merchant, error) {

	return usecase.repo.Search(ctx, input)
}

func (usecase MerchantUsecase) Create(ctx context.Context, merchant entities.Merchant) (int64, error) {

	if len(merchant.Name) == 0 {
		return 0, errors.New("name not found")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	if !emailRegex.MatchString(merchant.Email) {
		return 0, errors.New("invalid email")
	}

	if len(merchant.Password) == 0 {
		return 0, errors.New("password not found")
	}

	return usecase.repo.Create(ctx, merchant)
}

func (usecase MerchantUsecase) Login(ctx context.Context, login entities.Login) (string, error) {

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	if !emailRegex.MatchString(login.Email) {
		return "", errors.New("invalid email")
	}

	if len(login.Password) == 0 {
		return "", errors.New("password not found")
	}

	merchant, err := usecase.repo.Login(ctx, login)
	if err != nil {
		return "", err
	}

	if merchant == nil {
		return "", errors.New("invalid email or password")
	}

	token := uuid.New().String()

	token, err = usecase.repo.CreateToken(ctx, merchant.ID, token)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (usecase MerchantUsecase) Logout(ctx context.Context, id int64) (bool, error) {

	return usecase.repo.Logout(ctx, id)
}

func (usecase MerchantUsecase) Delete(ctx context.Context, id int64) (bool, error) {

	return usecase.repo.Delete(ctx, id)
}
