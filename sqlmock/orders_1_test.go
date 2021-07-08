package sqlmock

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

// NEWSQLMOCK OMIT
// will test that order with a different status, cannot be cancelled
func TestShouldNotCancelOrderWithNonPendingStatus(t *testing.T) {
	db, mock, err := sqlmock.New() // HL12
	if err != nil {
		t.Errorf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	// ENDMOCK OMIT
	columns := []string{"o_id", "o_status"}
	// SETEXPECTTATIONS OMIT
	mock.ExpectBegin() // HL12
	mock.ExpectQuery("SELECT (.+) FROM orders AS o INNER JOIN users AS u (.+) FOR UPDATE"). // HL12
		WithArgs(1). // HL12
		WillReturnRows(sqlmock.NewRows(columns).FromCSVString("1,1")) // HL12
	mock.ExpectRollback() // HL12
	// ENDEXPECTATIONS OMIT
	// RUNFUNC OMIT
	err = CancelOrder(1, db)
	if err != nil {
		t.Errorf("Expected no error, but got %s instead", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil { // HL12
		t.Errorf("there were unfulfilled expectations: %s", err) // HL12
	} // HL12
	// DONEFUNC OMIT
}

// END
