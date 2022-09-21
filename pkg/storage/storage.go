package storage

import "go-dao-pattern/pkg/storage/mysql"

type Config struct {
	Db mysql.ConnectionOptions
}
