package file_supply

import (
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
	KundenNr: "Demo",
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
	rootPath, err := gotoTmpl()
	require.NoError(s.T(), err)

	tmplPath, err := getTmplFile(rootPath)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), tmplPath)
	err = marshalOne(tmplPath, MockEntry)

	fileData, ok := LoadFile(getFileDestination(MockEntry))
	require.True(s.T(), ok)
	require.NotNil(s.T(), fileData)
}
