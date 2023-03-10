package server

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestServerTestSuite(t *testing.T) {
	suite.Run(t, new(ServerTestSuite))
}

type ServerTestSuite struct {
	suite.Suite
}

func (s *ServerTestSuite) SetupSuite() {
	log.SetLevel(log.DebugLevel)
}

func (s *ServerTestSuite) TestExtractQueryParams() {
	mockParams := gin.Params{
		gin.Param{Key: "amount", Value: "3"},
		gin.Param{Key: "legal", Value: "4"},
	}
	extracted := extractQueryParams(mockParams)
	require.NotEmpty(s.T(), extracted)
	require.Equal(s.T(), 2, len(extracted))
}
