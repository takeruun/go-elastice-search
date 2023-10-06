package usecase

import (
	"app/config"
	"app/database"
	"app/graph/model"
	"bytes"
	"encoding/json"
	"strconv"
)

const (
	ES_INDEX_KEY = "go-elastic-search_item"
)

type ItemUsecaseInterface interface {
	CreateItem(title string, description string) (*model.Item, error)
	UpdateItem(id string, title *string, description *string) (*model.Item, error)
	DeleteItem(id string) error
	SearchItems() ([]*model.Item, error)
}

type itemUsecase struct {
	itemDatabase database.ItemDatabaseInterface
	es           *config.ElasticSearch
}

func NewItemUsecase(itemDatabase database.ItemDatabaseInterface, es *config.ElasticSearch) ItemUsecaseInterface {
	return &itemUsecase{itemDatabase: itemDatabase, es: es}
}

func (u *itemUsecase) CreateItem(title string, description string) (*model.Item, error) {
	item, err := u.itemDatabase.Create(title, description)
	if err != nil {
		return nil, err
	}

	document := map[string]interface{}{
		"title":       item.Title,
		"description": item.Description,
	}
	data, err := json.Marshal(document)
	if err != nil {
		return nil, err
	}

	_, err = u.es.Index(ES_INDEX_KEY, bytes.NewReader(data), u.es.Index.WithDocumentID(strconv.Itoa(item.ID)))
	if err != nil {
		u.itemDatabase.Delete(strconv.Itoa(item.ID))
		return nil, err
	}

	return &model.Item{
		ID:          strconv.Itoa(item.ID),
		Title:       item.Title,
		Description: item.Description,
	}, nil
}

func (u *itemUsecase) UpdateItem(id string, title *string, description *string) (*model.Item, error) {
	item, err := u.itemDatabase.Update(id, title, description)
	if err != nil {
		return nil, err
	}

	return &model.Item{
		ID:          strconv.Itoa(item.ID),
		Title:       item.Title,
		Description: item.Description,
	}, nil
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
