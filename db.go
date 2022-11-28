package jokit

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// DBConfig contains all configs required to connect to the DB, QueryParams are optional, defaultQueryParams are used.
type DBConfig struct {
	Driver                string
	User                  string
	Password              string
	Host                  string
	Port                  string
	Schema                string
	MaxOpenConnections    int
	MaxOpenConnectionsTTL int
	MaxIdleConnections    int
	MaxIdleConnectionsTTL int
	QueryParams           map[string]interface{} // Optional
}

const defaultQueryParams = "?charset=latin1&parseTime=True&loc=Local"

// DbConnect Returns a connection to db or error
// it takes in the db credentials and other configurations e.g max connections.
// if this function returns an error, it is advisable to stop execution of the application.
func DbConnect(dbConfig DBConfig) (*sql.DB, error) {
	// parse query params
	queryParams := parseQueryParams(dbConfig.QueryParams)
	dbUri := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s%s", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Schema, queryParams)

	if dbConfig.Driver != "mysql" {
		return nil, fmt.Errorf("%s only mysql driver is configured", LogPrefix)
	}

	Log("%s Connecting to DB...", LogPrefix)
	db, err := sql.Open(dbConfig.Driver, dbUri)
	if err != nil {
		return nil, err
	}

	Log("%s Testing connection...", LogPrefix)
	if err = db.Ping(); err != nil {
		return nil, err
	}
	Log("%s Connected to DB successfully...", LogPrefix)

	Log("%s Setting MAX_OPEN_CONNECTIONS...", LogPrefix)
	db.SetMaxOpenConns(dbConfig.MaxOpenConnections)

	Log("%s Setting MAX_IDLE_CONNECTIONS...", LogPrefix)
	db.SetMaxIdleConns(dbConfig.MaxIdleConnections)

	Log("%s Setting max open connection lifetime to %d", LogPrefix, dbConfig.MaxOpenConnectionsTTL)
	db.SetConnMaxLifetime(time.Duration(dbConfig.MaxOpenConnectionsTTL) * time.Second)

	Log("%s Setting max idle connection lifetime to %d", LogPrefix, dbConfig.MaxIdleConnectionsTTL)
	db.SetConnMaxIdleTime(time.Duration(dbConfig.MaxIdleConnectionsTTL) * time.Second)
	return db, nil
}

func DbxConnect(dbConfig DBConfig) (*sqlx.DB, error) {
	// parse query params
	queryParams := parseQueryParams(dbConfig.QueryParams)
	dbUri := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s%s", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Schema, queryParams)

	if dbConfig.Driver != "mysql" {
		return nil, fmt.Errorf("%s only mysql driver is configured", LogPrefix)
	}

	Log("%s Connecting to DB...", LogPrefix)
	db, err := sqlx.Open(dbConfig.Driver, dbUri)
	if err != nil {
		return nil, err
	}

	Log("%s Testing connection...", LogPrefix)
	if err = db.Ping(); err != nil {
		return nil, err
	}
	Log("%s Connected to DB successfully...", LogPrefix)

	Log("%s Setting MAX_OPEN_CONNECTIONS...", LogPrefix)
	db.SetMaxOpenConns(dbConfig.MaxOpenConnections)

	Log("%s Setting MAX_IDLE_CONNECTIONS...", LogPrefix)
	db.SetMaxIdleConns(dbConfig.MaxIdleConnections)

	Log("%s Setting max open connection lifetime to %d", LogPrefix, dbConfig.MaxOpenConnectionsTTL)
	db.SetConnMaxLifetime(time.Duration(dbConfig.MaxOpenConnectionsTTL) * time.Second)

	Log("%s Setting max idle connection lifetime to %d", LogPrefix, dbConfig.MaxIdleConnectionsTTL)
	db.SetConnMaxIdleTime(time.Duration(dbConfig.MaxIdleConnectionsTTL) * time.Second)
	return db, nil
}

func parseQueryParams(params map[string]interface{}) string {
	var queryParams string
	if len(params) > 0 {
		for key, value := range params {
			if len(queryParams) == 0 {
				queryParams = fmt.Sprintf("?%s=%v", key, value)
			} else {
				queryParams = fmt.Sprintf("%s&%s=%v", queryParams, key, value)
			}
		}
	} else {
		queryParams = defaultQueryParams
	}
	return queryParams
}
