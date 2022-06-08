package marshaller

import (
	"descriptinator/pkg/server"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestMarshallerTestSuite(t *testing.T) {
	suite.Run(t, new(MarshallerTestSuite))
}

type MarshallerTestSuite struct {
	suite.Suite
}

func (s *MarshallerTestSuite) SetupSuite() {
	log.SetLevel(log.DebugLevel)
}

var MockEntry = Entry{
	Title:    mockTitle,
	Subtitle: mockSubtitle,
	Article: Article{
		GeneralInfo: mockArticleGeneralInfo,
		Description: mockArticleDescription,
		Fitting:     mockArticleFitting,
		Shipping:    mockArticleShipping,
		Condition:   mockArticleCondition,
	},
	Shipping: mockShipping,
	Legal:    mockLegal,
	Auction:  mockAuction,
	Seller:   mockSeller,
	Dsgvo:    mockDsgvo,
}

func (s *MarshallerTestSuite) TestMarshalHTMLPage() {
	fileNames, err := server.GetHtmlFiles()
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), fileNames)

	log.Debugf("%v", *fileNames)

	err = Marshal(*fileNames, MockEntry)
	require.NoError(s.T(), err)
}
