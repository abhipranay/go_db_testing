package tidb

import (
	"context"
	"time"

	tidblite "github.com/WangXiangUSTC/tidb-lite"
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

	_, err = dbConn.ExecContext(ctx, "create database shopee")
	c.Assert(err, IsNil)
	_, err = dbConn.ExecContext(ctx, "create table shopee.orders(id int primary key, value float(10, 7), reserved_fee float(10, 7), status int)")
	c.Assert(err, IsNil)
	_, err = dbConn.ExecContext(ctx, "insert into shopee.orders values(1, 100.0, 10.0, 0),(2, 120.0, 10.0, 1),(3, 130.0, 10.0, 0)")
	c.Assert(err, IsNil)

	tidbServer.Close()
}
