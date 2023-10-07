package database

import (
	"app/config"
	"app/domain"
	"strconv"
)

type ItemDatabaseInterface interface {
	Find(id string) (*domain.Item, error)
	Create(title string, description string) (*domain.Item, error)
	Update(id string, title *string, description *string) (*domain.Item, error)
	Delete(id string) error
}

type itemDatabase struct {
	db *config.DB
}

func NewItemDatabase(db *config.DB) ItemDatabaseInterface {
	return &itemDatabase{db: db}
}

func (d *itemDatabase) Find(id string) (*domain.Item, error) {
	repo := d.db.Connect()
	var item domain.Item
	if err := repo.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (d *itemDatabase) Create(title string, description string) (*domain.Item, error) {
	repo := d.db.Connect()
	item := domain.Item{
		Title:       title,
		Description: description,
	}
	if err := repo.Create(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (d *itemDatabase) Update(id string, title *string, description *string) (*domain.Item, error) {
	repo := d.db.Connect()
	numId, _ := strconv.Atoi(id)
	var item domain.Item = domain.Item{
		ID: numId,
	}
	if title != nil {
		item.Title = *title
	}
	if description != nil {
		item.Description = *description
	}

	if err := repo.
		Model(&domain.Item{}).
		Where("id = ?", id).
		Updates(&item).Error; err != nil {
		return nil, err
	}

	return &item, nil
}

func (d *itemDatabase) Delete(id string) error {
	repo := d.db.Connect()
	if err := repo.Delete(&domain.Item{}, id).Error; err != nil {
		return err
	}
	return nil
}
