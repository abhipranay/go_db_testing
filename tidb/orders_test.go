package tidb

import (
	"context"
	"time"

	tidblite "github.com/WangXiangUSTC/tidb-lite"
	"github.com/abhipranay/dbtesting/sqlmock"
	. "github.com/pingcap/check"
)

var _ = Suite(&testOrderCancelSuite{})

type testOrderCancelSuite struct{}

func (t *testExampleSuite) Test_OrderCancelled(c *C) {
	tidbServer, err := tidblite.NewTiDBServer(tidblite.NewOptions(c.MkDir())) // HL12
	c.Assert(err, IsNil)

	dbConn, err := tidbServer.CreateConn() // HL12
	c.Assert(err, IsNil)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// START OMIT
	_, err = dbConn.ExecContext(ctx, "create database shopee") // HL12
	c.Assert(err, IsNil)
	_, err = dbConn.ExecContext(ctx, "create table shopee.orders(id int primary key, buyer_id int, value float(10, 7), reserved_fee float(10, 7), status int, updated datetime default CURRENT_TIMESTAMP)") // HL12
	c.Assert(err, IsNil)
	_, err = dbConn.ExecContext(ctx, "create table shopee.users(id int primary key, username varchar(32), balance float(10, 7))") // HL12
	c.Assert(err, IsNil)
	_, err = dbConn.ExecContext(ctx, "insert into shopee.orders values(1, 10, 100.0, 10.0, 0, '0000-00-00 00:00:00'),(2, 20, 120.0, 10.0, 3, '0000-00-00 00:00:00'),(3, 10, 130.0, 10.0, 0, '0000-00-00 00:00:00')") // HL12
	c.Assert(err, IsNil)
	_, err = dbConn.ExecContext(ctx, "insert into shopee.users values(10, 'test1', 100.0),(20, 'test2', 200.0)") // HL12
	c.Assert(err, IsNil)

	err = sqlmock.CancelOrder(1, dbConn)
	c.Assert(err, IsNil)
	balance, err := sqlmock.GetBalance(10, dbConn)
	c.Assert(err, IsNil)
	c.Assert(balance == 210.0, IsTrue)
	// END OMIT

	// START2 OMIT
	err = sqlmock.CancelOrder(2, dbConn)
	c.Assert(err, IsNil)
	balance, err = sqlmock.GetBalance(20, dbConn)
	c.Assert(err, IsNil)
	c.Assert(balance == 200.0, IsTrue)
	// END2 OMIT
	tidbServer.Close()
}
