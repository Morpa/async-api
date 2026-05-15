package fixtures

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/Morpa/async-api/config"
	"github.com/Morpa/async-api/store"
	"github.com/golang-migrate/migrate/v4"
	"github.com/stretchr/testify/require"
)

type TestEnv struct {
	Config *config.Config
	Db     *sql.DB
}

func NewTestEnv(t *testing.T) *TestEnv {
	os.Setenv("ENV", string(config.Env_Test))
	conf, err := config.New()
	require.NoError(t, err)

	db, err := store.NewPostgresDb(conf)
	require.NoError(t, err)

	return &TestEnv{
		Config: conf,
		Db:     db,
	}
}

func (te *TestEnv) SetupdDb(t *testing.T) func(t *testing.T) {
	m, err := migrate.New(
		fmt.Sprintf("file:///%s/migrations", te.Config.ProjectRoot),
		te.Config.DatabaseUrl(),
	)
	require.NoError(t, err)

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		require.NoError(t, err)
	}

	return te.TeardownDb
}

func (te *TestEnv) TeardownDb(t *testing.T) {
	_, err := te.Db.Exec(fmt.Sprintf("TRUNCATE TABLE %s;", strings.Join([]string{"users", "refresh_tokens", "reports"}, ", ")))
	require.NoError(t, err)
}
