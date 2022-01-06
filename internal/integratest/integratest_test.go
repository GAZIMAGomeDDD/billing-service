//go:build integration
// +build integration

package integratest_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/GAZIMAGomeDDD/billing-service/internal/handler"
	"github.com/GAZIMAGomeDDD/billing-service/internal/server"
	"github.com/GAZIMAGomeDDD/billing-service/internal/storage/postgres"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

const (
	connString   = "postgres://user:password@localhost:54321/db?sslmode=disable"
	userID       = "b29f95a2-499a-4079-97f5-ff55c3854fcb"
	userBalance  = 99.97
	serveAddress = "localhost:9000"
)

type integraTestSuite struct {
	suite.Suite
	srv        *server.Server
	store      *postgres.Store
	dockerpool *dockertest.Pool
	resource   *dockertest.Resource
}

func TestIntegraTestSuite(t *testing.T) {
	suite.Run(t, &integraTestSuite{})
}

func (t *integraTestSuite) SetupSuite() {
	p, err := dockertest.NewPool("")
	if err != nil {
		logrus.Fatalf("Could not connect to docker: %s", err)
	}

	t.dockerpool = p

	opts := dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "14",
		Env: []string{
			"POSTGRES_USER=user",
			"POSTGRES_PASSWORD=password",
			"POSTGRES_DB=db",
		},
		ExposedPorts: []string{"5432"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432": {
				{HostIP: "0.0.0.0", HostPort: "54321"},
			},
		},
	}

	r, err := t.dockerpool.RunWithOptions(&opts)
	if err != nil {
		logrus.Fatalf("Could not start resource: %s", err.Error())
	}

	t.resource = r
	r.Expire(120)

	if err := t.dockerpool.Retry(func() error {
		db, err := sqlx.Connect("pgx", connString)
		if err != nil {
			return err
		}

		store, err := postgres.New(db)
		if err != nil {
			return err
		}

		t.store = store
		if _, err := db.Exec(
			"INSERT INTO users (id, balance) VALUES ($1, $2);",
			userID, userBalance); err != nil {
			return err
		}

		return nil

	}); err != nil {
		t.dockerpool.Purge(t.resource)

		logrus.Fatalf("Could not connect to docker: %s", err.Error())
	}

	h := handler.New(t.store, nil)
	srv := server.NewServer(h.Init(), serveAddress)
	t.srv = srv

	go t.srv.Run()
}

func (t *integraTestSuite) Test_getBalance() {
	user, err := t.store.GetBalance(userID)
	t.NoError(err)

	t.Equal(user.ID, userID)
	t.Equal(user.Balance, userBalance)

	body, err := json.Marshal(map[string]string{"id": userID})
	t.NoError(err)

	res, err := http.Post("http://"+serveAddress+"/getBalance", "application/json", bytes.NewReader(body))
	t.NoError(err)
	t.Equal(http.StatusOK, res.StatusCode)

	defer res.Body.Close()

	t.NoError(json.NewDecoder(res.Body).Decode(&user))

	t.Equal(user.ID, userID)
	t.Equal(user.Balance, userBalance)
}

func (t *integraTestSuite) TearDownSuite() {
	t.store.Close()
	t.srv.Stop()
	t.dockerpool.Purge(t.resource)
}
