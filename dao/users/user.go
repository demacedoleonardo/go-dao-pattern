package users

import (
	"go-dao-pattern/domain"
	"go-dao-pattern/pkg/context"
	"go-dao-pattern/pkg/storage"
	"go-dao-pattern/pkg/storage/mysql/db"
)

var (
	c DataAccess
)

const (
	MySql StorageType = iota + 1
	Memory
)

type (
	StorageType int

	User      domain.User
	Users     []User
	UserPages domain.UserPages

	KeyOperator struct {
		key   db.Column
		Op    db.Operator
		Value interface{}
	}

	Filters struct {
		Fields []db.Column
		Id     KeyOperator
		Name   KeyOperator
		Age    KeyOperator
		Offset int
		Limit  int
	}

	DataAccess interface {
		Search(*context.Context, Filters) (domain.UserPages, error)
		Create(*context.Context, domain.User) error
	}
)

func Search(ctx *context.Context, f Filters) (domain.UserPages, error) {
	return c.Search(ctx, f)
}

func Create(ctx *context.Context, u domain.User) error {
	return c.Create(ctx, u)
}

func InitDataAccess(st StorageType, cfg *storage.Config) {
	switch st {
	case MySql:
		c = NewUserStorage(cfg.Db)
	case Memory:
		c = NewUserMemoryStorage()
	default:
		c = nil
	}
}
