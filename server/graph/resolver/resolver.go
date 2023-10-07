package resolver

import (
	"app/config"
	"app/database"
	"app/usecase"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	ItemUsecase usecase.ItemUsecaseInterface
}

func NewResolver(db *config.DB, es *config.ElasticSearch) *Resolver {
	itemDatabase := database.NewItemDatabase(db)
	return &Resolver{
		ItemUsecase: usecase.NewItemUsecase(itemDatabase, es),
	}
}
