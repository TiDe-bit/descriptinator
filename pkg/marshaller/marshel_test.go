package marshaller

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type MarshallerTestSuite struct {
	suite.Suite
}

func NewMarshallerTestSuite() MarshallerTestSuite {
	return MarshallerTestSuite{}
}

func TestMarshaller(t *testing.T) {
	NewMarshallerTestSuite()
}

func (s *MarshallerTestSuite) TestMarshalHTMLPage() {

}
