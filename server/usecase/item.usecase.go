package usecase

import (
	"app/config"
	"app/database"
	"app/graph/model"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

const (
	ES_INDEX_KEY = "go-elastic-search_item"
)

type ItemUsecaseInterface interface {
	CreateItem(title string, description string) (*model.Item, error)
	UpdateItem(id string, title *string, description *string) (*model.Item, error)
	DeleteItem(id string) error
	SearchItems(ctx context.Context, where *model.ItemWhere) ([]*model.Item, int, error)
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

	_, err = u.es.Delete(ES_INDEX_KEY, id)
	if err != nil {
		return err
	}

	return nil
}

func (u *itemUsecase) SearchItems(ctx context.Context, where *model.ItemWhere) ([]*model.Item, int, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"title": where.Title,
			},
		},
	}

	data, err := json.Marshal(query)
	if err != nil {
		return nil, 0, err
	}

	res, err := u.es.Search(
		u.es.Search.WithContext(ctx),
		u.es.Search.WithIndex(ES_INDEX_KEY),
		u.es.Search.WithBody(bytes.NewReader(data)),
	)
	if err != nil {
		return nil, 0, err
	}
	defer res.Body.Close()
	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, 0, err
		}
		return nil, 0, fmt.Errorf(
			"ElasticSearch error, status: %s, type: %s, reason: %s",
			res.Status(),
			e["error"].(map[string]interface{})["type"],
			e["error"].(map[string]interface{})["reason"],
		)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, 0, err
	}

	var ids []int
	var totalCount int
	hits, ok := result["hits"]
	if !ok || hits == nil {
		return nil, totalCount, errors.New("hits not exists")
	}

	hitsMap, ok := hits.(map[string]any)
	if !ok {
		return nil, totalCount, errors.New("hits cannot cast map[string]any")
	}

	total, ok := hitsMap["total"]
	if !ok || total == nil {
		return nil, totalCount, errors.New("total not exists")
	}
	totalMap, ok := total.(map[string]any)
	if !ok {
		return nil, totalCount, errors.New("total cannot cast map[string]any")
	}
	v, ok := totalMap["value"]
	if !ok || v == nil {
		return nil, totalCount, errors.New("total.value not exists")
	}
	totalCountF, ok := v.(float64)
	if !ok {
		return nil, totalCount, errors.New("total.value cannot cast float64")
	}
	totalCount = int(totalCountF)

	hitList, ok := hitsMap["hits"]
	if !ok || hitList == nil {
		return nil, totalCount, errors.New("hitList not exists")
	}

	hitAnyList, ok := hitList.([]any)
	if !ok {
		return nil, totalCount, errors.New("hitList cannot cast []any")
	}

	for _, hit := range hitAnyList {
		hitMap, ok := hit.(map[string]any)
		if !ok {
			return nil, totalCount, errors.New("hit cannot cast map[string]any")
		}

		i, ok := hitMap["_id"]
		if !ok || i == nil {
			return nil, totalCount, errors.New("_id not exists")
		}
		idStr, ok := i.(string)
		if !ok {
			return nil, totalCount, errors.New("_id cannot cast string")
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return nil, totalCount, err
		}
		source, ok := hitMap["_source"]
		if !ok || source == nil {
			return nil, totalCount, errors.New("_source not exists")
		}
		sourceMap, ok := source.(map[string]any)
		if !ok {
			return nil, totalCount, errors.New("_source cannot cast map[string]any")
		}
		fmt.Println(sourceMap)
		ids = append(ids, id)
	}

	items, err := u.itemDatabase.Search(ids)

	var itemsModel []*model.Item
	for _, item := range items {
		itemsModel = append(itemsModel, &model.Item{
			ID:          strconv.Itoa(item.ID),
			Title:       item.Title,
			Description: item.Description,
		})
	}

	return itemsModel, totalCount, nil
}
