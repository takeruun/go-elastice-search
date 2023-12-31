package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.38

import (
	"app/graph/model"
	"context"
	"fmt"
)

// CreateItem is the resolver for the createItem field.
func (r *mutationResolver) CreateItem(ctx context.Context, title string, description string) (*model.Item, error) {
	item, err := r.ItemUsecase.CreateItem(title, description)
	if err != nil {
		return nil, err
	}

	return item, nil
}

// UpdateItem is the resolver for the updateItem field.
func (r *mutationResolver) UpdateItem(ctx context.Context, id string, title *string, description *string) (*model.Item, error) {
	item, err := r.ItemUsecase.UpdateItem(id, title, description)
	if err != nil {
		return nil, err
	}

	return item, nil
}

// DeleteItem is the resolver for the deleteItem field.
func (r *mutationResolver) DeleteItem(ctx context.Context, id string) (string, error) {
	err := r.ItemUsecase.DeleteItem(id)
	if err != nil {
		return "", err
	}

	return "success", nil
}

// SearchItems is the resolver for the searchItems field.
func (r *queryResolver) SearchItems(ctx context.Context, where *model.ItemWhere) ([]*model.Item, error) {
	items, _, err := r.ItemUsecase.SearchItems(ctx, where)
	if err != nil {
		return nil, err
	}

	return items, nil
}

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//   - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//     it when you're done.
//   - You have helper methods in this file. Move them out to keep these resolver files clean.
func (r *queryResolver) Items(ctx context.Context, where *model.ItemWhere) ([]*model.Item, error) {
	panic(fmt.Errorf("not implemented: Items - items"))
}
