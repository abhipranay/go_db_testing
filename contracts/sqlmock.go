package contracts

import "github.com/DATA-DOG/go-sqlmock"

// START OMIT
type Sqlmock interface {
	// ExpectClose queues an expectation for this database action to be triggered.
	ExpectClose() *sqlmock.ExpectedClose // HL12
	// ExpectationsWereMet checks whether all queued expectations were met in order. If any of them was not met - an error is returned.
	ExpectationsWereMet() error // HL12
	// ExpectPrepare expects Prepare() to be called with sql query which match sqlRegexStr given regexp.
	ExpectPrepare(sqlRegexStr string) *sqlmock.ExpectedPrepare // HL12
	// ExpectQuery expects Query() or QueryRow() to be called with sql query which match sqlRegexStr given regexp.
	ExpectQuery(sqlRegexStr string) *sqlmock.ExpectedQuery // HL12
	// ExpectExec expects Exec() to be called with sql query which match sqlRegexStr given regexp.
	ExpectExec(sqlRegexStr string) *sqlmock.ExpectedExec // HL12
	// ExpectBegin expects *sql.DB.Begin to be called.
	ExpectBegin() *sqlmock.ExpectedBegin // HL12
	// ExpectCommit expects *sql.Tx.Commit to be called.
	ExpectCommit() *sqlmock.ExpectedCommit // HL12
	// ExpectRollback expects *sql.Tx.Rollback to be called.
	ExpectRollback() *sqlmock.ExpectedRollback // HL12
	// By default it is set to - true. But if you use goroutines
	MatchExpectationsInOrder(bool) // HL12
}

// END OMIT
