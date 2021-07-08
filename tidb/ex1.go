package tidb

import (
	"database/sql"
	"fmt"
)

import (
	"context"

	"github.com/pingcap/errors"
	"github.com/pingcap/log"
	"go.uber.org/zap"
)

// GetRowCount returns row count of the table.
// if not specify where condition, return total row count of the table.
func GetRowCount(ctx context.Context, db *sql.DB, schemaName string, tableName string, where string) (int64, error) {
	/*
		select count example result:
		mysql> SELECT count(1) cnt from `test`.`itest` where id > 0;
		+------+
		| cnt  |
		+------+
		|  100 |
		+------+
	*/

	query := fmt.Sprintf("SELECT COUNT(1) cnt FROM `%s`.`%s`", schemaName, tableName)
	if len(where) > 0 {
		query += fmt.Sprintf(" WHERE %s", where)
	}
	log.Debug("get row count", zap.String("sql", query))

	var cnt sql.NullInt64
	err := db.QueryRowContext(ctx, query).Scan(&cnt)
	if err != nil {
		return 0, errors.Trace(err)
	}
	if !cnt.Valid {
		return 0, errors.NotFoundf("table `%s`.`%s`", schemaName, tableName)
	}

	return cnt.Int64, nil
}
