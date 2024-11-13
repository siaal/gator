package cli

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
	"github.com/siaal/gator/internal/config"
	"github.com/siaal/gator/internal/database"
)

type State struct {
	Config *config.Config
	DB     *database.Queries
}

func NewState() (State, error) {
	state := State{}
	conf, err := config.ReadConfig(config.DefaultPaths)
	if err != nil {
		return state, err
	}
	state.Config = &conf

	db, err := openDatabase(conf.DBUrl)
	if err != nil {
		return state, err
	}
	state.DB = db
	return state, err
}

func openDatabase(dbUrl string) (*database.Queries, error) {
	var dbType string
	switch {
	case strings.HasPrefix(dbUrl, "postgres"):
		dbType = "postgres"
	default:
		return nil, fmt.Errorf("Unable to determine database type using dburl: %s", dbUrl)
	}
	db, err := sql.Open(dbType, dbUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return database.New(db), nil
}
