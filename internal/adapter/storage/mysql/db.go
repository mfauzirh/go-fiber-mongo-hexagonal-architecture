package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/Masterminds/squirrel"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mfauzirh/go-fiber-mongo-hexarch/internal/adapter/config"
)

/*
 * This is wrapper for MySQL database connection,
 * It holds a reference to squirrel query builder and MySQL driver
 */
type DB struct {
	*sql.DB
	QueryBuilder *squirrel.StatementBuilderType
	url          string
}

/*
 * Create new database connection
 * using configuration from config
 */
func New(ctx context.Context, config *config.DB) (*DB, error) {
	url := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Name,
	)

	db, err := sql.Open("mysql", url)
	if err != nil {
		return nil, err
	}

	// Ping to check the connection
	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	mySQL := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question) // Use ? placeholder for MySQL

	return &DB{
		DB:           db,
		QueryBuilder: &mySQL,
		url:          url,
	}, nil
}

// ErrorCode returns the error code of the given error.
func (db *DB) ErrorCode(err error) string {
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		return fmt.Sprintf("%d", mysqlErr.Number)
	}
	return ""
}

// Close closes the database connection.
func (db *DB) Close() {
	if err := db.DB.Close(); err != nil {
		log.Println("Error closing the database connection:", err)
	}
}
