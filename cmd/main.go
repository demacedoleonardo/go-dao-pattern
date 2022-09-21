package main

import (
	"fmt"
	"go-dao-pattern/cfg"
	"go-dao-pattern/dao/users"
	"go-dao-pattern/domain"
	"go-dao-pattern/pkg/context"
	"go-dao-pattern/pkg/storage"
	"go-dao-pattern/pkg/storage/mysql/db"
)

func main() {
	FindUserDataBase()
	//FindUserMemory()
}

func FindUserDataBase() {
	config := &storage.Config{
		Db: cfg.MysqlConfig,
	}

	users.InitDataAccess(users.MySql, config)

	f := users.Filters{
		Fields: []db.Column{"age", "name"},
		Id: users.KeyOperator{
			Op:    db.Equal,
			Value: 1,
		},
	}

	ctx := context.NewBackgroundContext()
	data, err := users.Search(ctx, f)

	if err != nil {
		println("BOOM!!!")
		return
	}

	print(data)
}

func FindUserMemory() {
	users.InitDataAccess(users.Memory, nil)
	ctx := context.NewBackgroundContext()

	// Given User
	user := domain.User{
		ID:   1,
		Age:  38,
		Name: "Leonardo",
	}

	err := users.Create(ctx, user)

	//Search User
	f := users.Filters{
		Id: users.KeyOperator{
			Op:    db.Equal,
			Value: 1,
		},
		Name: users.KeyOperator{
			Op:    db.Equal,
			Value: "Leonardo",
		},
	}

	data, err := users.Search(ctx, f)

	if err != nil {
		println("BOOM!!!")
		return
	}

	print(data)
}

func print(data domain.UserPages) {
	println("Users Len: ", data.Total)

	for _, d := range data.Users {
		println(fmt.Sprintf("Users Data: [id: %d] [name: %s] [age: %d]", d.ID, d.Name, d.Age))
	}
}
