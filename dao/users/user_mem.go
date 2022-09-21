package users

import (
	"fmt"
	"go-dao-pattern/domain"
	"go-dao-pattern/pkg/context"
	"go-dao-pattern/pkg/storage/memory"
)

// Ensure type implements interface.
var _ DataAccess = (*userMemory)(nil)

type userMemory struct {
	storage *memory.StorageClient
}

func (u *userMemory) Search(context *context.Context, filters Filters) (domain.UserPages, error) {
	key := fmt.Sprintf("%s-%s", fmt.Sprint(filters.Id.Value), filters.Name.Value)
	data, err := u.storage.Get(context, key)
	if err != nil {
		return domain.UserPages{}, err
	}

	users := make(domain.Users, 1)
	users[0] = data.(domain.User)
	up := domain.UserPages{
		Limit:  0,
		Offset: 0,
		Total:  0,
		Users:  users,
	}
	return up, err
}

func (u *userMemory) Create(context *context.Context, user domain.User) error {
	key := fmt.Sprintf("%s-%s", fmt.Sprint(user.ID), user.Name)
	return u.storage.Save(context, key, user)
}

func NewUserMemoryStorage() *userMemory {
	return &userMemory{
		storage: memory.InitConnection(),
	}
}
