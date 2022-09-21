package db

import (
	"errors"
	"fmt"
	"strings"
)

var (
	SqlBuilderJoinMismatchLenErr    = errors.New("join statement should be create by equal number of tables")
	SqlBuilderFromClauseErr         = errors.New("from clause should provide a valida table name")
	SqlBuilderMissingActionErr      = errors.New("action should be select, update")
	SqlBuilderMissingOrderFieldsErr = errors.New("order by should provide a valid fields")
)

const (
	And               Operator  = " AND "
	Or                Operator  = " OR "
	Equal             Operator  = " = "
	GreaterThan       Operator  = " > "
	LessThan          Operator  = " < "
	GreaterEqualsThan Operator  = " >= "
	LessEqualsThan    Operator  = " <= "
	Asc               OrderType = " ASC "
	Desc              OrderType = " DESC "

	space           = " "
	defaultMaxPages = 10
)

type (
	Table  string
	Column string

	Operator  string
	OrderType string

	condition struct {
		union string
		key   string
		op    Operator
	}

	sort struct {
		operator string
		values   []string
	}

	pagination struct {
		limit  int
		offset int
	}

	conditions struct {
		c     []condition
		union string
	}

	WhereOptions func(w *conditions)
	SetOptions   func(w *conditions)

	tableInfo struct {
		name string
		key  string
	}

	query struct {
		action      string
		columns     string
		sql         string
		table       tableInfo
		joins       []tableInfo
		withjoins   []tableInfo
		sets        []condition
		wheres      []condition
		args        []interface{}
		withcounter bool
		sort        sort
		pagination  pagination
	}

	Filters struct {
		Fields     []string
		SetValues  []SetOptions
		AndFilters []WhereOptions
		OrFilters  []WhereOptions
		Limit      int
		Offset     int
	}

	beforeSelect struct {
		q *query
	}

	beforeCounter struct {
		q *query
	}

	beforeFrom struct {
		q *query
	}

	beforeTable struct {
		q *query
	}

	beforeWhere struct {
		q *query
	}

	beforeLimit struct {
		q *query
	}

	beforeConditional struct {
		q *query
	}

	beforeUpdate struct {
		q *query
	}

	beforeSet struct {
		q *query
	}

	finish struct {
		q *query
	}
)

func Select(fields ...Column) *beforeSelect {
	columns := make([]string, len(fields))
	for i, value := range fields {
		columns[i] = string(value)
	}

	cols := "*"
	if len(columns) > 0 {
		cols = strings.Join(columns, ", ")
	}

	q := &query{
		sql:     "SELECT " + cols + space,
		action:  "select",
		columns: cols,
		table:   tableInfo{},
	}

	return &beforeSelect{q: q}
}

func (q *beforeSelect) From(t Table) *beforeFrom {
	q.q.table.name = string(t)
	return &beforeFrom{q: q.q}
}

func (q *beforeCounter) From(t Table) *beforeFrom {
	q.q.table.name = string(t)
	return &beforeFrom{q: q.q}
}

func (q *beforeSelect) WithCounter() *beforeCounter {
	q.q.withcounter = true
	return &beforeCounter{q: q.q}
}

func (q *beforeFrom) Join(t Table, k Column) *beforeTable {
	if q.q.joins == nil {
		q.q.joins = make([]tableInfo, 0)
	}

	q.q.joins = append(q.q.joins, tableInfo{
		name: string(t),
		key:  string(k),
	})

	return &beforeTable{q: q.q}
}

func (q *beforeTable) Table(t Table, k Column) *beforeFrom {
	if q.q.withjoins == nil {
		q.q.withjoins = make([]tableInfo, 0)
	}

	q.q.withjoins = append(q.q.withjoins, tableInfo{
		name: string(t),
		key:  string(k),
	})

	return &beforeFrom{q: q.q}
}

func (q *beforeFrom) Where(c Column, o Operator) *beforeWhere {
	if q.q.wheres == nil {
		q.q.wheres = make([]condition, 0)
	}

	q.q.wheres = append(q.q.wheres, condition{
		key: string(c),
		op:  o,
	})

	return &beforeWhere{q: q.q}
}

func (q *beforeSet) Where(c Column, o Operator) *beforeWhere {
	if q.q.wheres == nil {
		q.q.wheres = make([]condition, 0)
	}

	q.q.wheres = append(q.q.wheres, condition{
		key: string(c),
		op:  o,
	})

	return &beforeWhere{q: q.q}
}

func (q *beforeConditional) Where(c Column, o Operator) *beforeWhere {
	if q.q.wheres == nil {
		q.q.wheres = make([]condition, 0)
	}

	q.q.wheres = append(q.q.wheres, condition{
		key: string(c),
		op:  o,
	})

	return &beforeWhere{q: q.q}
}

func (q *beforeWhere) OrderBy(sort OrderType, columns ...Column) *beforeLimit {
	fields := make([]string, len(columns))
	for i, v := range columns {
		fields[i] = string(v)
	}

	q.q.sort.operator = string(sort)
	q.q.sort.values = fields
	return &beforeLimit{q: q.q}
}

func (q *beforeFrom) Limit(offset, limit int) *finish {
	if limit == 0 {
		limit = defaultMaxPages
	}

	q.q.pagination.limit = limit
	q.q.pagination.offset = offset
	return &finish{q: q.q}
}

func (q *beforeWhere) Limit(offset, limit int) *finish {
	q.q.pagination.limit = limit
	q.q.pagination.offset = offset
	return &finish{q: q.q}
}

func (q *beforeLimit) Limit(offset, limit int) *finish {
	q.q.pagination.limit = limit
	q.q.pagination.offset = offset
	return &finish{q: q.q}
}

func (q *beforeWhere) And() *beforeConditional {
	q.q.wheres[len(q.q.wheres)-1].union = string(And)
	return &beforeConditional{q: q.q}
}

func (q *beforeWhere) Or() *beforeConditional {
	q.q.wheres[len(q.q.wheres)-1].union = string(Or)
	return &beforeConditional{q: q.q}
}

func Update(table Table) *beforeUpdate {
	q := &query{}
	q.action = "update"
	q.sql = "UPDATE " + string(table) + space
	q.table.name = string(table)
	return &beforeUpdate{q: q}
}

func (q *beforeUpdate) Set(column Column, operator Operator) *beforeSet {
	if q.q.sets == nil {
		q.q.sets = make([]condition, 0)
	}

	if len(q.q.sets) > 0 {
		q.q.sets[len(q.q.sets)-1].union = " ,"
	}

	q.q.sets = append(q.q.sets, condition{key: string(column), op: operator})
	return &beforeSet{q: q.q}
}

func (q *beforeSet) Set(column Column, operator Operator) *beforeSet {
	if q.q.sets == nil {
		q.q.sets = make([]condition, 0)
	}

	if len(q.q.sets) > 0 {
		q.q.sets[len(q.q.sets)-1].union = ", "
	}

	q.q.sets = append(q.q.sets, condition{key: string(column), op: operator})
	return &beforeSet{q: q.q}
}

func (q *beforeFrom) Build() (string, error) {
	statement, err := selectStmt(q.q)
	return statement, err
}

func (q *finish) Build() (string, error) {
	statement, err := selectStmt(q.q)
	return statement, err
}

func (q *beforeWhere) Build() (string, error) {
	switch q.q.action {
	case "select":
		statement, err := selectStmt(q.q)
		return statement, err
	case "update":
		statement, err := updateStmt(q.q)
		return statement, err
	default:
		return "", SqlBuilderMissingActionErr
	}
}

func selectStmt(q *query) (string, error) {
	if len(q.joins) != len(q.withjoins) {
		return "", SqlBuilderJoinMismatchLenErr
	}

	if len(q.table.name) == 0 {
		return "", SqlBuilderFromClauseErr
	}

	if len(q.sort.operator) > 0 && len(q.sort.values) == 0 {
		return "", SqlBuilderMissingOrderFieldsErr
	}

	var sb strings.Builder
	sb.WriteString(sel(q))
	sb.WriteString(from(q.table.name))
	sb.WriteString(join(q.joins, q.withjoins))
	sb.WriteString(wheres(q.wheres))
	sb.WriteString(orderBy(q.sort))
	sb.WriteString(pages(q.pagination))
	sb.WriteString(";")
	return sb.String(), nil
}

func sel(q *query) string {
	if !q.withcounter {
		return "SELECT " + q.columns
	}

	var sb strings.Builder
	sb.WriteString("SELECT count(*)")
	sb.WriteString(from(q.table.name))
	sb.WriteString(join(q.joins, q.withjoins))
	sb.WriteString(wheres(q.wheres))

	return fmt.Sprintf("SELECT %s, (%s) as total", q.columns, sb.String())
}

func pages(p pagination) string {
	if p.offset >= 0 && p.limit >= 0 && (p.offset > 0 || p.limit > 0) {
		return fmt.Sprintf(" LIMIT %d, %d", p.offset, p.limit)
	}
	return ""
}

func orderBy(s sort) string {
	if len(s.operator) > 0 {
		return " ORDER BY " + strings.Join(s.values, ", ")
	}
	return ""
}

func from(name string) string {
	return " FROM " + name
}

func wheres(w []condition) string {
	out := ""
	if len(w) > 0 {
		for _, where := range w {
			out += where.key + string(where.op) + "?" + where.union
		}
		out = " WHERE " + out
	}
	return out
}

func join(joins []tableInfo, withjoins []tableInfo) string {
	out := ""
	if len(joins) > 0 {
		for i, join := range joins {
			withTable := withjoins[i]
			out += withTable.name + " ON " + withTable.name + "." + withTable.key + " = " + join.name + "." + join.key

			if len(joins)-i > 1 {
				out += " JOIN "
			}
		}
		out = " JOIN " + out
	}
	return out
}

func updateStmt(q *query) (string, error) {
	if len(q.sets) > 0 {
		sets := ""
		for _, set := range q.sets {
			sets += set.key + string(set.op) + "?" + set.union
		}
		q.sql += "SET " + sets
	}

	if len(q.wheres) > 0 {
		wheres := ""
		for _, where := range q.wheres {
			wheres += where.key + string(where.op) + "?" + where.union
		}
		q.sql += " WHERE " + wheres
	}

	q.sql = strings.TrimSuffix(q.sql, space) + ";"
	return q.sql, nil
}
