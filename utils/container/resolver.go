package container

import (
	"database/sql"
	"no-q-solution/externals/adapters"
	"no-q-solution/externals/repositories"
	"no-q-solution/utils/config"
)

func Resolve(config config.Config) (Containers, error) {
	adaptrs, err := resolveAdapters(config)
	if err != nil {
		return Containers{}, err
	}

	repos, err := resolveRepostories(adaptrs.Db)
	if err != nil {
		return Containers{}, err
	}

	cont := Containers{
		Adapters:     adaptrs,
		Repositories: repos,
	}

	return cont, nil
}

func resolveAdapters(config config.Config) (Adapters, error) {

	mysql, err := adapters.NewDB(config.Database)
	if err != nil {
		return Adapters{}, err
	}

	adapters := Adapters{
		Db: mysql,
	}

	return adapters, nil
}

func resolveRepostories(db *sql.DB) (Repositories, error) {
	merchantRepo := repositories.NewMerchantRepository(db)
	queueRepo := repositories.NewQueueRepository(db)

	repos := Repositories{
		Merchant: merchantRepo,
		Queue:    queueRepo,
	}

	return repos, nil
}
