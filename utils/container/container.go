package container

import (
	"database/sql"
	"no-q-solution/domain/interfaces"
)

type Containers struct {
	Adapters     Adapters
	Repositories Repositories
}

type Adapters struct {
	Db *sql.DB
}

type Repositories struct {
	Merchant interfaces.MerchantRepository
	Queue    interfaces.QueueRepository
}
