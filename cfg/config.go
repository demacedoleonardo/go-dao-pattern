package cfg

import "go-dao-pattern/pkg/storage/mysql"

var (
	MysqlConfig = mysql.ConnectionOptions{
		Host:            "127.0.0.1",
		User:            "root",
		Password:        "root",
		Schema:          "company",
		RegisterName:    "mysql",
		Port:            3306,
		ConnMaxOpen:     5,
		ConnMaxIdle:     2,
		ConnMaxLifetime: 200,
	}
)

func init() {

}
