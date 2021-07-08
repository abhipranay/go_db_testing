package sqlmock

import (
	"database/sql"
	"fmt"

	"github.com/kisielk/sqlstruct"
)

const ORDER_PENDING = 0
const ORDER_CANCELLED = 1

type User struct {
	Id       int     `sql:"id"`
	Username string  `sql:"username"`
	Balance  float64 `sql:"balance"`
}

type Order struct {
	Id          int     `sql:"id"`
	Value       float64 `sql:"value"`
	ReservedFee float64 `sql:"reserved_fee"`
	Status      int     `sql:"status"`
}

func CancelOrder(id int, db *sql.DB) (err error) {
	var order Order
	var user User

	// SELECTQ OMIT
	tx, err := db.Begin() // <-- tx start // HL12
	if err != nil { // OMIT
		return // OMIT
	} // OMIT
	sql := fmt.Sprintf(`
SELECT %s, %s
FROM shopee.orders AS o
INNER JOIN shopee.users AS u ON o.buyer_id = u.id
WHERE o.id = ?
FOR UPDATE`,
		sqlstruct.ColumnsAliased(order, "o"),
		sqlstruct.ColumnsAliased(user, "u"))
	// ENDSELECTQ OMIT
	// fetch order to cancel
	rows, err := tx.Query(sql, id)
	if err != nil {
		tx.Rollback()
		return
	}

	defer rows.Close()
	// no rows, nothing to do
	if !rows.Next() {
		tx.Rollback()
		return
	}

	// read order
	err = sqlstruct.ScanAliased(&order, rows, "o")
	if err != nil {
		tx.Rollback()
		return
	}

	// ensure order status
	if order.Status != ORDER_PENDING {
		tx.Rollback()
		return
	}

	// read user
	err = sqlstruct.ScanAliased(&user, rows, "u")
	if err != nil {
		tx.Rollback()
		return
	}
	rows.Close() // manually close before other prepared statements

	// refund order value
	// UPDATEQ OMIT
	sql = "UPDATE shopee.users SET balance = balance + ? WHERE id = ?"
	refundStmt, err := tx.Prepare(sql)
	if err != nil { // OMIT
		tx.Rollback() // OMIT
		return // OMIT
	} // OMIT
	defer refundStmt.Close()
	_, err = refundStmt.Exec(order.Value+order.ReservedFee, user.Id)
	// ENDUPDATE OMIT
	if err != nil {
		tx.Rollback()
		return
	}

	// update order status
	// CANCELLEDQ OMIT
	order.Status = ORDER_CANCELLED
	sql = "UPDATE shopee.orders SET status = ?, updated = NOW() WHERE id = ?"
	orderUpdStmt, err := tx.Prepare(sql)
	// ENDCANCELLED OMIT
	if err != nil {
		tx.Rollback()
		return
	}
	defer orderUpdStmt.Close()
	_, err = orderUpdStmt.Exec(order.Status, order.Id)
	if err != nil {
		tx.Rollback()
		return
	}
	return tx.Commit()
}

func GetBalance(id int, db *sql.DB) (float64, error) {
	var user User
	sql := fmt.Sprintf(`
SELECT %s
FROM shopee.users
WHERE id = ?
`,
		sqlstruct.Columns(User{}))
	// fetch order to cancel
	rows, err := db.Query(sql, id)
	if err != nil {
		return 0, err
	}
	rows.Next()
	err = sqlstruct.Scan(&user, rows)
	if err != nil {
		return 0, err
	}

	return user.Balance, nil
}
