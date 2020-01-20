package data

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
)

type PageResults struct {
	TotalPages  int
	RecordCount int
}

type PagedBuilder struct {
	builder sq.SelectBuilder
	columns []string
}

type PagedSelectBuilder struct {
	pageNumber int
	pageSize   int
	orderBys   []string
	builder    sq.SelectBuilder
}

type PagedCountBuilder struct {
	builder sq.SelectBuilder
}

func RecordColumns(tables ...string) []string {
	columns := make([]string, 0)
	for _, t := range tables {
		columns = append(columns, fmt.Sprintf(`%s.id "%s.id"`, t, t))
		columns = append(columns, fmt.Sprintf(`%s.organisationid "%s.organisationid"`, t, t))
		columns = append(columns, fmt.Sprintf(`%s.version "%s.version"`, t, t))
		columns = append(columns, fmt.Sprintf(`%s.record "%s.record"`, t, t))
	}

	return columns
}

func Paged(columns ...string) PagedBuilder {
	return PagedBuilder{
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar).Select(),
		columns: columns,
	}
}

func Select(columns ...string) sq.SelectBuilder {
	return sq.StatementBuilder.PlaceholderFormat(sq.Dollar).Select(columns...)
}

func Insert(into string) sq.InsertBuilder {
	return sq.StatementBuilder.PlaceholderFormat(sq.Dollar).Insert(into)
}

func Update(table string) sq.UpdateBuilder {
	return sq.StatementBuilder.PlaceholderFormat(sq.Dollar).Update(table)
}

func Delete(from string) sq.DeleteBuilder {
	return sq.StatementBuilder.PlaceholderFormat(sq.Dollar).Delete(from)
}

func (b PagedSelectBuilder) OrderBy(orderBys ...string) PagedSelectBuilder {
	b.orderBys = orderBys
	return b
}

func (b PagedSelectBuilder) ToSql() (string, []interface{}, error) {
	return b.builder.OrderBy(b.orderBys...).Offset(uint64(b.pageNumber * b.pageSize)).Limit(uint64(b.pageSize)).ToSql()
}

func (b PagedCountBuilder) ToSql() (string, []interface{}, error) {
	return b.builder.ToSql()
}

func (b PagedBuilder) ToSelect(pageNumber int, pageSize int) PagedSelectBuilder {
	return PagedSelectBuilder{
		builder:    b.builder.Columns(b.columns...),
		pageNumber: pageNumber,
		pageSize:   pageSize,
		orderBys:   []string{"paginationId"},
	}
}

func (b PagedBuilder) ToCount() PagedCountBuilder {
	return PagedCountBuilder{builder: b.builder.Columns("COUNT(*)")}
}

func (b PagedBuilder) Prefix(sql string, args ...interface{}) PagedBuilder {
	return PagedBuilder{builder: b.builder.Prefix(sql, args), columns: b.columns}
}

func (b PagedBuilder) Distinct() PagedBuilder {
	return PagedBuilder{builder: b.builder.Distinct(), columns: b.columns}
}

func (b PagedBuilder) Options(options ...string) PagedBuilder {
	return PagedBuilder{builder: b.builder.Options(options...), columns: b.columns}
}

func (b PagedBuilder) From(from string) PagedBuilder {
	return PagedBuilder{builder: b.builder.From(from), columns: b.columns}
}

func (b PagedBuilder) FromSelect(from sq.SelectBuilder, alias string) PagedBuilder {
	return PagedBuilder{builder: b.builder.FromSelect(from, alias), columns: b.columns}
}

func (b PagedBuilder) JoinClause(pred interface{}, args ...interface{}) PagedBuilder {
	return PagedBuilder{builder: b.builder.JoinClause(pred, args...), columns: b.columns}
}

func (b PagedBuilder) Join(join string, rest ...interface{}) PagedBuilder {
	return PagedBuilder{builder: b.builder.Join(join, rest...), columns: b.columns}
}

func (b PagedBuilder) LeftJoin(join string, rest ...interface{}) PagedBuilder {
	return PagedBuilder{builder: b.builder.LeftJoin(join, rest...), columns: b.columns}
}

func (b PagedBuilder) RightJoin(join string, rest ...interface{}) PagedBuilder {
	return PagedBuilder{builder: b.builder.RightJoin(join, rest...), columns: b.columns}
}

func (b PagedBuilder) Where(pred interface{}, args ...interface{}) PagedBuilder {
	return PagedBuilder{builder: b.builder.Where(pred, args...), columns: b.columns}
}

func (b PagedBuilder) GroupBy(groupBys ...string) PagedBuilder {
	return PagedBuilder{builder: b.builder.GroupBy(groupBys...), columns: b.columns}
}

func (b PagedBuilder) Having(pred interface{}, rest ...interface{}) PagedBuilder {
	return PagedBuilder{builder: b.builder.Having(pred, rest...), columns: b.columns}
}

func (b PagedBuilder) Suffix(sql string, args ...interface{}) PagedBuilder {
	return PagedBuilder{builder: b.builder.Suffix(sql, args...), columns: b.columns}
}
