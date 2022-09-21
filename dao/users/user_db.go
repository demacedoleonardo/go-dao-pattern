package users

import (
	"database/sql"
	"go-dao-pattern/domain"
	"go-dao-pattern/pkg/context"
	"go-dao-pattern/pkg/storage/mysql"
	"go-dao-pattern/pkg/storage/mysql/db"
)

// Ensure type implements interface.
var _ DataAccess = (*userStorage)(nil)

const (
	users db.Table = "users"

	id   db.Column = "id"
	name db.Column = "name"
	age  db.Column = "age"
)

type userStorage struct {
	storage *mysql.StorageClient
}

func NewUserStorage(options mysql.ConnectionOptions) *userStorage {
	return &userStorage{
		storage: mysql.InitConnection(options),
	}
}

func (us *userStorage) beginTx(ctx *context.Context) (*sql.Tx, error) {
	return us.storage.BeginTx(ctx.Context(), nil)
}

func (us *userStorage) Create(ctx *context.Context, u domain.User) error {
	panic("implement me")
}

func (us *userStorage) CreateTx(ctx *context.Context, tx *sql.Tx, u User) error {
	panic("implement me")
}

func (us *userStorage) Search(ctx *context.Context, f Filters) (domain.UserPages, error) {
	wheres, args := f.projections()

	sql := db.Select(f.Fields...).From(users)
	sql.Limit(f.Offset, f.Limit)

	for _, ko := range wheres {
		sql.Where(ko.key, ko.Op)
	}

	query, _ := sql.Build()

	var up domain.UserPages
	rows, err := db.ExecQuery(ctx.Context(), us.storage, string(users), query, args...)

	if err != nil {
		return up, err
	}
	defer rows.Close()

	users := make(domain.Users, 0)
	for rows.Next() {
		user := new(domain.User)
		if err := rows.Scan(user.Cols(f.Fields)...); err != nil {
			return up, err
		}
		users = append(users, *user)
	}

	up.Offset = f.Offset
	up.Limit = f.Limit
	up.Total = 1
	up.Users = users
	return up, nil
}

func (u *User) args() []interface{} {
	args := make([]interface{}, 0)

	if u.ID > 0 {
		args = append(args, u.ID)
	}

	if len(u.Name) > 0 {
		args = append(args, u.Name)
	}

	if u.Age > 0 {
		args = append(args, u.Age)
	}

	return args
}

func (f Filters) projections() ([]KeyOperator, []interface{}) {
	args := make([]interface{}, 0)
	ko := make([]KeyOperator, 0)

	if f.Id.HasValue() {
		f.Id.key = id
		ko = append(ko, f.Id)
		args = append(args, f.Id.Value)
	}

	if f.Name.HasValue() {
		f.Name.key = name
		ko = append(ko, f.Name)
		args = append(args, f.Name.Value)
	}

	if f.Age.HasValue() {
		f.Age.key = age
		ko = append(ko, f.Age)
		args = append(args, f.Age.Value)
	}

	return ko, args
}

func (ko KeyOperator) HasValue() bool {
	return ko.Value != nil && len(ko.Op) > 0
}
