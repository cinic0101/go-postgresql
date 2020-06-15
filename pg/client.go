package pg

import (
	"context"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Client struct {
	MasterDbUrl string
	SlaveDbUrl  string
	masterPool  *pgxpool.Pool
	slavePool   *pgxpool.Pool
}

func (c *Client) InitConnPools() error {
	var err error
	if c.masterPool, err = pgxpool.Connect(context.Background(), c.MasterDbUrl); err != nil {
		return err
	}
	if c.slavePool, err = pgxpool.Connect(context.Background(), c.SlaveDbUrl); err != nil {
		return err
	}

	return nil
}

func (c *Client) Query(sql string, args ...interface{}) (pgx.Rows, error) {
	conn, err := c.slavePool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	return conn.Query(context.Background(), sql, args...)
}

func (c *Client) Exec(sql string, args ...interface{}) (pgconn.CommandTag, error) {
	conn, err := c.masterPool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	return conn.Exec(context.Background(), sql, args...)
}
