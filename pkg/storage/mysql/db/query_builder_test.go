package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuery_BuildSelectAll(t *testing.T) {
	var (
		users Table = "users"
	)

	q, err := Select().
		From(users).
		Where("name", Equal).
		And().
		Where("age", Equal).
		Build()

	expected := `SELECT * FROM users WHERE name = ? AND age = ?;`

	assert.Nil(t, err)
	assert.Equal(t, expected, q)
}

func TestQuery_BuildSelectWithCounter(t *testing.T) {
	var (
		users Table = "users"
	)

	q, err := Select().
		WithCounter().
		From(users).
		Where("name", Equal).
		And().
		Where("age", Equal).
		Build()

	expected := `SELECT *, (SELECT count(*) FROM users WHERE name = ? AND age = ?) as total FROM users WHERE name = ? AND age = ?;`

	assert.Nil(t, err)
	assert.Equal(t, expected, q)
}

func TestQuery_BuildSelect(t *testing.T) {
	var (
		columns       = []Column{"id", "name", "age"}
		users   Table = "users"
	)

	q, err := Select(columns...).
		From(users).
		Where("name", Equal).
		And().
		Where("age", Equal).
		Build()

	expected := `SELECT id, name, age FROM users WHERE name = ? AND age = ?;`

	assert.Nil(t, err)
	assert.Equal(t, expected, q)
}

func TestQuery_BuildJoin(t *testing.T) {
	var (
		columns           = []Column{"id", "name", "age"}
		users       Table = "users"
		credentials Table = "credentials"
		history     Table = "history"
	)

	q, err := Select(columns...).
		From(users).
		Join(users, "id").
		Table(credentials, "users_id").
		Join(users, "id").
		Table(history, "users_id").
		Where("id", Equal).
		And().
		Where("name", Equal).
		Build()

	expected := `SELECT id, name, age FROM users JOIN credentials ON credentials.users_id = users.id JOIN history ON history.users_id = users.id WHERE id = ? AND name = ?;`

	assert.Nil(t, err)
	assert.Equal(t, expected, q)
}

func TestQuery_BuildUpdate(t *testing.T) {
	var (
		users Table  = "users"
		id    Column = "id"
		name  Column = "name"
		age   Column = "age"
	)

	q, err := Update(users).
		Set(name, Equal).
		Set(age, Equal).
		Where(id, Equal).Or().
		Where(age, GreaterThan).Build()

	expected := `UPDATE users SET name = ?, age = ? WHERE id = ? OR age > ?;`

	assert.Nil(t, err)
	assert.Equal(t, expected, q)
}

func TestQuery_BuildSelectWithOrderAndLimit(t *testing.T) {
	var (
		sort          = []Column{"id", "name"}
		columns       = []Column{"id", "name", "age"}
		users   Table = "users"
	)

	q, err := Select(columns...).
		From(users).
		Where("name", Equal).
		And().
		Where("age", Equal).
		OrderBy(Asc, sort...).
		Limit(0, 10).
		Build()

	expected := `SELECT id, name, age FROM users WHERE name = ? AND age = ? ORDER BY id, name LIMIT 0, 10;`

	assert.Nil(t, err)
	assert.Equal(t, expected, q)
}
