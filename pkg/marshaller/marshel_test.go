package marshaller

import (
	"descriptinator/pkg/file_supply"
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

var mockDemoArtNum = "njdsjnds"

var MockEntry = Entry{
	KundenNr: &mockDemoArtNum,
	Title:    &mockTitle,
	Subtitle: &mockSubtitle,
	Article: Article{
		GeneralInfo: &mockArticleGeneralInfo,
		Description: &mockArticleDescription,
		Fitting:     &mockArticleFitting,
		shipping:    &mockArticleShipping,
		Condition:   &mockArticleCondition,
	},
	shipping: &mockShipping,
	legal:    &mockLegal,
	auction:  &mockAuction,
	seller:   &mockSeller,
	dsgvo:    &mockDsgvo,
}

func (s *MarshallerTestSuite) TestMarshalHTMLPage() {
	rootPath, err := file_supply.GotoTmpl()
	require.NoError(s.T(), err)

	tmplPath, err := file_supply.GetTmplFile(rootPath)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), tmplPath)
	err = marshalOne(tmplPath, &MockEntry)

	fileData, ok := file_supply.LoadFile(getFileDestination(&MockEntry))
	require.True(s.T(), ok)
	require.NotNil(s.T(), fileData)
}
