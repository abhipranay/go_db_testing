package sqlmock

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

// START OMIT
func TestShouldRollbackOnError(t *testing.T) {
	db, mock, err := sqlmock.New() // HL12
	if err != nil {
		t.Errorf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	mock.ExpectBegin() // HL12

	mock.ExpectQuery("SELECT (.+) FROM shopee.orders AS o INNER JOIN shopee.users AS u (.+) FOR UPDATE"). // HL12
		WithArgs(1). // HL12
		WillReturnError(fmt.Errorf("Some error")) // HL12

	mock.ExpectRollback() // HL12

	err = CancelOrder(1, db)
	if err == nil {
		t.Error("Expected error, but got none")
	}
	if err := mock.ExpectationsWereMet(); err != nil { // HL12
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
// END OMIT
