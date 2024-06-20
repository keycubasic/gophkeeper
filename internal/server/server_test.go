package server

import (
	"context"
	"testing"
	"time"

	"github.com/keycubasic/gophkeeper/internal/app"
	"github.com/keycubasic/gophkeeper/internal/config"
	"github.com/keycubasic/gophkeeper/internal/server/testhelpers"
	"github.com/keycubasic/gophkeeper/internal/storage/postgres"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/metadata"

	apb "github.com/keycubasic/gophkeeper/internal/auth/proto"
	gpb "github.com/keycubasic/gophkeeper/internal/grpc/proto"
)

type GophKeeperTestSuite struct {
	suite.Suite
	ctx         context.Context
	cfg         *config.Config
	pgContainer *testhelpers.PostgresContainer
}

func (suite *GophKeeperTestSuite) SetupSuite() {
	t := suite.T()
	suite.ctx = context.Background()
	pgContainer, err := testhelpers.CreatePostgresContainer(suite.ctx)
	if err != nil {
		t.Fatal(err)
	}
	suite.pgContainer = pgContainer
	cfg := new(config.Config)
	cfg.ParseConfig()
	cfg.DatabaseDSN = pgContainer.ConnectionString
	suite.cfg = cfg

	db, err := postgres.NewStorage(suite.ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}

	err = db.Migrate(suite.ctx)
	if err != nil {
		t.Fatal(err)
	}

	app := app.NewApp(suite.ctx, cfg, db)
	gsrv := RegisterGrpcServer(suite.ctx, cfg, app)
	go RunServer(suite.ctx, cfg, gsrv)
}

func (suite *GophKeeperTestSuite) TearDownSuite() {
	t := suite.T()
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		t.Fatalf("terminating postgres container: %s", err)
	}
}

func (suite *GophKeeperTestSuite) TestRegisterUserOK() {
	t := suite.T()
	authClient, err := testhelpers.NewAuthClient(suite.ctx, suite.cfg.ServerAddress)
	if err != nil {
		t.Fatal(err)
	}

	user := &apb.User{UserName: "user1", Password: "password1"}
	_, err = authClient.Client.RegisterUser(suite.ctx, &apb.RegisterUserRequest{User: user})
	require.NoError(t, err)

	if err := authClient.Conn.Close(); err != nil {
		t.Fatalf("auth client closing connection: %s", err)
	}
}

func (suite *GophKeeperTestSuite) TestRegisterUserKO() {
	t := suite.T()
	authClient, err := testhelpers.NewAuthClient(suite.ctx, suite.cfg.ServerAddress)
	if err != nil {
		t.Fatal(err)
	}

	user := &apb.User{UserName: "user2", Password: "password2"}
	_, err = authClient.Client.RegisterUser(suite.ctx, &apb.RegisterUserRequest{User: user})
	require.NoError(t, err)

	_, err = authClient.Client.RegisterUser(suite.ctx, &apb.RegisterUserRequest{User: user})
	require.Error(t, err)

	if err := authClient.Conn.Close(); err != nil {
		t.Fatalf("auth client closing connection: %s", err)
	}
}

func (suite *GophKeeperTestSuite) TestLoginUserOK() {
	t := suite.T()
	authClient, err := testhelpers.NewAuthClient(suite.ctx, suite.cfg.ServerAddress)
	if err != nil {
		t.Fatal(err)
	}

	user := &apb.User{UserName: "user3", Password: "password3"}
	_, err = authClient.Client.RegisterUser(suite.ctx, &apb.RegisterUserRequest{User: user})
	require.NoError(t, err)

	_, err = authClient.Client.LoginUser(suite.ctx, &apb.LoginUserRequest{User: user})
	require.NoError(t, err)

	if err := authClient.Conn.Close(); err != nil {
		t.Fatalf("auth client closing connection: %s", err)
	}
}

func (suite *GophKeeperTestSuite) TestLoginUserKO() {
	t := suite.T()
	authClient, err := testhelpers.NewAuthClient(suite.ctx, suite.cfg.ServerAddress)
	if err != nil {
		t.Fatal(err)
	}

	user := &apb.User{UserName: "userNotFound", Password: "password"}
	_, err = authClient.Client.LoginUser(suite.ctx, &apb.LoginUserRequest{User: user})
	require.Error(t, err)

	if err := authClient.Conn.Close(); err != nil {
		t.Fatalf("auth client closing connection: %s", err)
	}
}

func (suite *GophKeeperTestSuite) TestStoreDataOK() {
	t := suite.T()
	authClient, err := testhelpers.NewAuthClient(suite.ctx, suite.cfg.ServerAddress)
	if err != nil {
		t.Fatal(err)
	}

	user := &apb.User{UserName: "user4", Password: "password4"}
	resAuth, err := authClient.Client.RegisterUser(suite.ctx, &apb.RegisterUserRequest{User: user})
	require.NoError(t, err)

	if err := authClient.Conn.Close(); err != nil {
		t.Fatalf("auth client closing connection: %s", err)
	}

	gkClient, err := testhelpers.NewGophKeeperClient(suite.ctx, suite.cfg.ServerAddress)
	if err != nil {
		t.Fatal(err)
	}

	newCtx := metadata.AppendToOutgoingContext(suite.ctx, "authorization", "Bearer "+resAuth.AccessToken)

	creds := &gpb.Credentials{
		Login:    "login1",
		Password: "123456",
	}
	res, err := gkClient.Client.StoreData(newCtx, &gpb.StoreDataRequest{Site: "somesite.com", Credentials: creds, Metadata: "somemeta"})
	require.NoError(t, err)
	require.Equal(t, "user data will be stored", res.Info)

	if err := gkClient.Conn.Close(); err != nil {
		t.Fatalf("data client closing connection: %s", err)
	}
}

func (suite *GophKeeperTestSuite) TestStoreDataKO() {
	t := suite.T()
	authClient, err := testhelpers.NewAuthClient(suite.ctx, suite.cfg.ServerAddress)
	if err != nil {
		t.Fatal(err)
	}

	user := &apb.User{UserName: "user5", Password: "password5"}
	_, err = authClient.Client.RegisterUser(suite.ctx, &apb.RegisterUserRequest{User: user})
	require.NoError(t, err)

	if err := authClient.Conn.Close(); err != nil {
		t.Fatalf("auth client closing connection: %s", err)
	}

	gkClient, err := testhelpers.NewGophKeeperClient(suite.ctx, suite.cfg.ServerAddress)
	if err != nil {
		t.Fatal(err)
	}

	creds := &gpb.Credentials{
		Login:    "login1",
		Password: "123456",
	}
	_, err = gkClient.Client.StoreData(suite.ctx, &gpb.StoreDataRequest{Site: "somesite.com", Credentials: creds, Metadata: "somemeta"})
	require.Error(t, err)

	if err := gkClient.Conn.Close(); err != nil {
		t.Fatalf("data client closing connection: %s", err)
	}
}

func (suite *GophKeeperTestSuite) TestGetDataOK() {
	t := suite.T()
	authClient, err := testhelpers.NewAuthClient(suite.ctx, suite.cfg.ServerAddress)
	if err != nil {
		t.Fatal(err)
	}

	user := &apb.User{UserName: "user6", Password: "password6"}
	resAuth, err := authClient.Client.RegisterUser(suite.ctx, &apb.RegisterUserRequest{User: user})
	require.NoError(t, err)

	if err := authClient.Conn.Close(); err != nil {
		t.Fatalf("auth client closing connection: %s", err)
	}

	gkClient, err := testhelpers.NewGophKeeperClient(suite.ctx, suite.cfg.ServerAddress)
	if err != nil {
		t.Fatal(err)
	}

	newCtx := metadata.AppendToOutgoingContext(suite.ctx, "authorization", "Bearer "+resAuth.AccessToken)

	creds := &gpb.Credentials{
		Login:    "login1",
		Password: "123456",
	}
	resStore, err := gkClient.Client.StoreData(newCtx, &gpb.StoreDataRequest{Site: "somesite.com", Credentials: creds, Metadata: "somemeta"})
	require.NoError(t, err)
	require.Equal(t, "user data will be stored", resStore.Info)

	time.Sleep(2 * time.Second)

	resGet, err := gkClient.Client.GetData(newCtx, &gpb.GetDataRequest{Site: "somesite.com"})
	require.NoError(t, err)
	require.Equal(t, "123456", resGet.Credentials.Password)

	if err := gkClient.Conn.Close(); err != nil {
		t.Fatalf("data client closing connection: %s", err)
	}
}

func (suite *GophKeeperTestSuite) TestGetDataKO() {
	t := suite.T()
	authClient, err := testhelpers.NewAuthClient(suite.ctx, suite.cfg.ServerAddress)
	if err != nil {
		t.Fatal(err)
	}

	user := &apb.User{UserName: "user7", Password: "password7"}
	resAuth, err := authClient.Client.RegisterUser(suite.ctx, &apb.RegisterUserRequest{User: user})
	require.NoError(t, err)

	if err := authClient.Conn.Close(); err != nil {
		t.Fatalf("auth client closing connection: %s", err)
	}

	gkClient, err := testhelpers.NewGophKeeperClient(suite.ctx, suite.cfg.ServerAddress)
	if err != nil {
		t.Fatal(err)
	}

	newCtx := metadata.AppendToOutgoingContext(suite.ctx, "authorization", "Bearer "+resAuth.AccessToken)

	_, err = gkClient.Client.GetData(newCtx, &gpb.GetDataRequest{Site: "wrongsite.com"})
	require.Error(t, err)

	if err := gkClient.Conn.Close(); err != nil {
		t.Fatalf("data client closing connection: %s", err)
	}
}

func TestGophKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(GophKeeperTestSuite))
}
