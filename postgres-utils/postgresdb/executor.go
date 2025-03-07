package postgresdb

import (
	"context"
	"log"
)

type QueryExecutor struct {
	connection *DBConnection
}

func NewQueryExecutor(conn *DBConnection) *QueryExecutor {
	return &QueryExecutor{connection: conn}
}

func (executor *QueryExecutor) Execute(ctx *context.Context, query string, tuple *Tuple) error {
	_, err := executor.connection.GetPool().Exec(*ctx, query, tuple.params...)
	if err != nil {
		log.Printf("Query execution failed: %v\n", err)
		return err
	}
	return nil
}

func (executor *QueryExecutor) Query(ctx *context.Context,
	query string, tuple *Tuple) (map[string]interface{}, error) {
	row := executor.connection.GetPool().QueryRow(*ctx, query, tuple.params...)
	result := make(map[string]interface{})
	err := row.Scan(&result)
	if err != nil {
		log.Printf("Failed to fetch row: %v\n", err)
		return nil, err
	}
	return result, nil
}
