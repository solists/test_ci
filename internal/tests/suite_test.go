package tests

import (
	"context"
	"fmt"
	"mymod/internal/repository"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/pkg/errors"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
	dockerPool *dockertest.Pool
	db         *sqlx.DB
	pgdb       *sqlx.DB
	resource   *dockertest.Resource
	pgresource *dockertest.Resource
	repo       repository.IRepository
	pgrepo     repository.IRepository
}

func (s *Suite) TearDownSuite() {
	if s.db != nil {
		err := s.db.Close()
		require.Nil(s.T(), err)
	}

	if s.dockerPool != nil && s.resource != nil {
		err := s.dockerPool.Purge(s.resource)
		require.Nil(s.T(), err)
	}
}

func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}
func (s *Suite) SetupSuite() {
	var err error
	s.dockerPool, err = dockertest.NewPool("")
	require.Nil(s.T(), err)

	s.initPG()
}
func (s *Suite) TearDownTest() {
	truncateTablesQuery := `truncate
			items_log;`

	_, err := s.db.ExecContext(context.TODO(), truncateTablesQuery)
	require.Nil(s.T(), err)
	_, err = s.pgdb.ExecContext(context.TODO(), truncateTablesQuery)
	require.Nil(s.T(), err)

}

func (s *Suite) initPG() {
	const migrationsPGPath = "../../db/migrations"
	if dsn := os.Getenv("DATABASE_DSN"); dsn != "" {
		var err error
		s.pgdb, err = sqlx.Connect("postgres", dsn)
		require.Nil(s.T(), err)

		err = goose.Up(s.pgdb.DB, migrationsPGPath, goose.WithAllowMissing())
		require.Nil(s.T(), err)

		s.pgrepo = repository.NewRepository(s.pgdb)
		return
	}

	var err error
	s.pgresource, err = s.dockerPool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "12",
		Env: []string{
			"POSTGRES_PASSWORD=pass",
			"POSTGRES_USER=user",
			"POSTGRES_DB=test",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})

	require.Nil(s.T(), err)

	err = s.dockerPool.Retry(func() error {
		port := s.pgresource.GetPort("5432/tcp")
		dsn := fmt.Sprintf("postgres://user:pass@127.0.0.1:%s/test?sslmode=disable&binary_parameters=yes", port)

		if s.pgdb, err = sqlx.Connect("postgres", dsn); err != nil {
			return errors.Wrap(err, "db.Connect")
		}

		return s.pgdb.Ping()
	})
	require.Nil(s.T(), err)

	require.Nil(s.T(), goose.SetDialect("postgres"))
	err = goose.Up(s.pgdb.DB, migrationsPGPath, goose.WithAllowMissing())
	require.Nil(s.T(), err)

	s.pgrepo = repository.NewRepository(s.pgdb)
}
