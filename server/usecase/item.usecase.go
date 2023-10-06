package usecase

import (
	"app/config"
	"app/database"
	"app/graph/model"
)

type ItemUsecaseInterface interface {
	CreateItem(title string, description string) (*model.Item, error)
	UpdateItem(id string, title *string, description *string) (*model.Item, error)
	DeleteItem(id string) error
	SearchItems() ([]*model.Item, error)
}

type itemUsecase struct {
	db           *config.DB
	itemDatabase database.ItemDatabaseInterface
}

func NewItemUsecase(db *config.DB, itemDatabase database.ItemDatabaseInterface) ItemUsecaseInterface {
	return &itemUsecase{db: db, itemDatabase: itemDatabase}
}

func (u *itemUsecase) CreateItem(title string, description string) (*model.Item, error) {
	item, err := u.itemDatabase.Create(title, description)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (u *itemUsecase) UpdateItem(id string, title *string, description *string) (*model.Item, error) {
	item, err := u.itemDatabase.Update(id, title, description)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (u *itemUsecase) DeleteItem(id string) error {
	err := u.itemDatabase.Delete(id)
	if err != nil {
		return err
	}

	return nil
}

func (u *itemUsecase) SearchItems() ([]*model.Item, error) {
	return nil, nil
}
