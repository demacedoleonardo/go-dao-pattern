package domain

import "go-dao-pattern/pkg/storage/mysql/db"

type (
	User struct {
		ID   int    `json:"id"`
		Age  int    `json:"age"`
		Name string `json:"name"`
	}

	Users []User

	UserPages struct {
		Limit  int   `json:"limit"`
		Offset int   `json:"offset"`
		Total  int   `json:"total"`
		Users  Users `json:"users"`
	}
)

func (u *User) Cols(fields []db.Column) []interface{} {
	if len(fields) == 0 {
		return []interface{}{
			&u.ID, &u.Name, &u.Age,
		}
	}

	c := map[string]interface{}{
		"id":   &u.ID,
		"age":  &u.Age,
		"name": &u.Name,
	}

	cols := make([]interface{}, len(fields))
	for i, f := range fields {
		if col, found := c[string(f)]; found {
			cols[i] = col
		}
	}
	return cols
}
