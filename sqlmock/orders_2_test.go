package sqlmock

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

// will test order cancellation
func TestShouldRefundUserWhenOrderIsCancelled(t *testing.T) {
	// open database stub
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// START OMIT
	columns := []string{"o_id", "o_status", "o_value", "o_reserved_fee", "u_id", "u_balance"}
	mock.ExpectBegin() // HL12
	mock.ExpectQuery("SELECT (.+) FROM shopee.orders AS o INNER JOIN shopee.users AS u (.+) FOR UPDATE"). // HL12
		WithArgs(1). // HL12
		WillReturnRows(sqlmock.NewRows(columns).AddRow(1, 0, 25.75, 3.25, 2, 10.00)) // HL12

	mock.ExpectPrepare("UPDATE shopee.users SET balance").ExpectExec(). // HL12
		WithArgs(25.75+3.25, 2). // refund amount, user id // HL12
		WillReturnResult(sqlmock.NewResult(0, 1)) // no insert id, 1 affected row // HL12

	mock.ExpectPrepare("UPDATE shopee.orders SET status").ExpectExec(). // HL12
		WithArgs(ORDER_CANCELLED, 1). // status, id // HL12
		WillReturnResult(sqlmock.NewResult(0, 1)) // no insert id, 1 affected row // HL12

	mock.ExpectCommit() // HL12
	// END OMIT
	err = CancelOrder(1, db)
	if err != nil {
		t.Errorf("Expected no error, but got %s instead", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
