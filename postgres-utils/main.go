package main

import (
	"context"
	"log"
	"postgres-utils/postgresdb"
)

func main() {
	ctx := context.Background()
	dsn := "postgres://postgres:admin@localhost:5432/users"

	conn := postgresdb.NewDBConnection(&ctx, dsn)

	err := postgresdb.NewSchemaExecutor(conn).CreateTables(&ctx,
		"/Users/yuvrajsingh/GolandProjects/GoServers/postgres-utils/tables.sql")
	if err != nil {
		panic(err)
	}

	queryExecutor := postgresdb.NewQueryExecutor(conn)

	query := "INSERT INTO users (name, age) VALUES ($1, $2);"
	tuple := postgresdb.NewTuple()
	tuple.AddString("yuvraj")
	tuple.AddInt(24)
	err = queryExecutor.Execute(&ctx, query, tuple)
	if err != nil {
		log.Println("Error executing query: ", err)
		return
	}

}
