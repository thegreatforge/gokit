package sql

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

// Driver is the type of sql driver to use
type Driver string

const (
	Postgres Driver = "postgres"
	Mysql    Driver = "mysql"
)

// Config is the configuration for the sql connection
type Config struct {
	Name      string
	sqldriver Driver
	Host      string
	Port      int
	User      string
	password  string
	Dbname    string
	logger    *zap.Logger
}

// NewConfig creates a new config for the sql connection
func NewConfig(name string, sqldriver Driver,
	host string, port int, user string, password string, dbname string,
	logger *zap.Logger) Config {

	return Config{
		Name:      name,
		sqldriver: sqldriver,
		Host:      host,
		Port:      port,
		User:      user,
		password:  password,
		Dbname:    dbname,
		logger:    logger,
	}
}

// Connection is a wrapper around the sql.DB object
type Connection struct {
	Db     *sql.DB
	Config Config
}

// Databases is a map of all the connections
var Databases map[string]*Connection

func init() {
	Databases = make(map[string]*Connection)
}

// driverConnectString returns the connection string for the sql driver
func driverConnectString(config Config) string {
	switch config.sqldriver {
	case Postgres:
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			config.Host, config.Port, config.User, config.password, config.Dbname)
	case Mysql:
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
			config.User, config.password, config.Host, config.Port, config.Dbname)
	default:
		return ""
	}
}

// NewConnection creates a new connection to the sql database
func NewConnection(config Config) error {
	db, err := sql.Open(string(config.sqldriver), driverConnectString(config))

	if err != nil {
		config.logger.Error("Error connecting to database", zap.Error(err))
		return err
	}

	// Check the connection
	err = db.Ping()
	if err != nil {
		config.logger.Error("Error connecting to database", zap.Error(err))
		return err
	}

	conn := &Connection{
		Db:     db,
		Config: config,
	}

	Databases[config.Name] = conn
	config.logger.Sugar().Infof("successfully connected to database instance - %s", config.Name)
	return nil
}

// Close closes the connection to the database
func (c *Connection) Close() error {
	c.Config.logger.Sugar().Infof("closing connection to database instance - %s", c.Config.Name)
	return c.Db.Close()
}

// Query executes a SELECT query on the database
func (c *Connection) Query(ctx context.Context, query string, args ...interface{}) ([]map[string]interface{}, error) {
	rows, err := c.Db.QueryContext(ctx, query, args...)
	if err != nil {
		c.Config.logger.Error("Error executing query", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	// Get the column names
	columns, err := rows.Columns()
	if err != nil {
		c.Config.logger.Error("Error getting columns", zap.Error(err))
		return nil, err
	}

	// Make a slice of interface{}'s to hold the values of each row
	values := make([]interface{}, len(columns))
	for i := range values {
		values[i] = new(interface{})
	}

	// Make a slice of maps to hold the results
	results := []map[string]interface{}{}

	for rows.Next() {

		err := rows.Scan(values...)
		if err != nil {
			c.Config.logger.Error("Error scanning row", zap.Error(err))
			return nil, err
		}

		// Make a map to hold the values of this row
		row := map[string]interface{}{}

		// Loop through each column
		for i, column := range columns {
			// Get the value of the interface{} and add it to the map
			val := *(values[i].(*interface{}))
			row[column] = val
		}

		// Add the map to the results slice
		results = append(results, row)
	}

	// Check for any errors in the rows loop
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return results, nil
}

// Exec executes a query on the database
// Mostly used for inserts, updates, and deletes
// Returns the number of rows affected
func (c *Connection) Exec(ctx context.Context, query string, args ...interface{}) (int64, error) {
	result, err := c.Db.ExecContext(ctx, query, args...)
	if err != nil {
		c.Config.logger.Error("Error executing query", zap.Error(err))
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.Config.logger.Sugar().Warnf("Error getting rows affected - %s", err.Error())
		return 0, nil
	}
	return rowsAffected, nil
}
