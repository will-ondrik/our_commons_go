package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Pool *pgxpool.Pool
}

type Config struct {
	Host string
	Port string
	User string
	Password string
	DbName string
	PoolMaxConns int
}

func NewDb() *DB {

	dbConfig := &Config{
		Host: "localhost", 
		Port: "5432",
		User: "willondrik",
		Password: "Kamloops_1",
		DbName: "politicallycharged",
		PoolMaxConns: 10,
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.DbName)
	poolConf, err := pgxpool.ParseConfig(connStr)

	if err != nil {
		panic("failed to parse db config")
	}

	poolConf.MaxConns = int32(dbConfig.PoolMaxConns)

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConf)
	if err != nil {
		panic("db init failed")
	}

	fmt.Println("db connection successful")
	return &DB{
		Pool: pool,
	}
}