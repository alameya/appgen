package integration

import (
	"math/rand"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"

	"github.com/your-project/proto"
	"github.com/your-project/repository"
)

type IntegrationTestSuite struct {
	suite.Suite
	db         *sqlx.DB
	pool       *dockertest.Pool
	resource   *dockertest.Resource
	grpcConn   *grpc.ClientConn
	courier    proto.CourierServiceClient
	location   proto.LocationServiceClient
	repository *repository.Repository
	httpPort   string
	r          *rand.Rand
}

func (s *IntegrationTestSuite) SetupSuite() {
	// ... existing setup ...
	s.httpPort = os.Getenv("PORT")
	if s.httpPort == "" {
		s.httpPort = "8080"
	}
	s.r = rand.New(rand.NewSource(time.Now().UnixNano()))
}
