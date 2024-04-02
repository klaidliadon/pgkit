package pgkit_test

import (
	"strings"
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/goware/pgkit/v2"
	"github.com/stretchr/testify/require"
)

type T struct{}

func TestPagination(t *testing.T) {
	const (
		DefaultSize = 2
		MaxSize     = 5
		Sort        = "ID"
	)
	paginator := pgkit.NewPaginator[T](
		pgkit.WithColumnFunc(strings.ToLower),
		pgkit.WithDefaultSize(DefaultSize),
		pgkit.WithMaxSize(MaxSize),
		pgkit.WithSort(Sort),
	)
	page := pgkit.NewPage(0, 0)
	result, query := paginator.PrepareQuery(sq.Select("*").From("t"), page)
	require.Len(t, result, 0)
	require.Equal(t, &pgkit.Page{Page: 1, Size: MaxSize}, page)

	sql, args, err := query.ToSql()
	require.NoError(t, err)
	require.Equal(t, "SELECT * FROM t ORDER BY id ASC LIMIT 6 OFFSET 0", sql)
	require.Empty(t, args)

	result = paginator.PrepareResult(make([]T, 0), page)
	require.Len(t, result, 0)
	require.Equal(t, &pgkit.Page{Page: 1, Size: MaxSize}, page)

	result = paginator.PrepareResult(make([]T, MaxSize), page)
	require.Len(t, result, MaxSize)
	require.Equal(t, &pgkit.Page{Page: 1, Size: MaxSize}, page)

	result = paginator.PrepareResult(make([]T, MaxSize+2), page)
	require.Len(t, result, MaxSize)
	require.Equal(t, &pgkit.Page{Page: 1, Size: MaxSize, More: true}, page)
}
