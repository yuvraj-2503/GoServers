package postgresdb

import (
	"bufio"
	"context"
	"log"
	"os"
	"strings"
)

type SchemaExecutor struct {
	connection *DBConnection
}

func NewSchemaExecutor(conn *DBConnection) *SchemaExecutor {
	return &SchemaExecutor{connection: conn}
}

// CreateTables reads queries from file and executes the queries and creates tables
func (schema *SchemaExecutor) CreateTables(ctx *context.Context, sqlFilePath string) error {
	queries, err := schema.readQueriesFromFile(sqlFilePath)
	if err != nil {
		return err
	}

	err = schema.createTables(ctx, queries)
	if err != nil {
		return err
	}

	return nil
}

// readQueriesFromFile reads and returns SQL queries from a specified file.
func (schema *SchemaExecutor) readQueriesFromFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var queries []string
	scanner := bufio.NewScanner(file)
	var queryBuilder strings.Builder

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "--") { // Skip empty lines and comments
			continue
		}

		queryBuilder.WriteString(line)
		if strings.HasSuffix(line, ";") { // End of a query
			queries = append(queries, queryBuilder.String())
			queryBuilder.Reset()
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return queries, nil
}

// createTables executes a list of CREATE TABLE queries.
func (schema *SchemaExecutor) createTables(ctx *context.Context, queries []string) error {
	for _, query := range queries {
		_, err := schema.connection.GetPool().Exec(*ctx, query)
		if err != nil {
			return err
		}
	}
	log.Printf("CREATE TABLE : Transaction Succeeded")
	return nil
}
