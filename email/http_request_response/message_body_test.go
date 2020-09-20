package http_request_response

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type messageBodyTestSuite struct {
	suite.Suite
}

func TestMessageBodyTestSuite(t *testing.T) {
	suite.Run(t, new(messageBodyTestSuite))
}

func (suite *messageBodyTestSuite) TestGetDecodedData_ShouldReturnBase64DecodedData() {
	messageBody := MessageBody{
		Content: "QmFzZTY0IGVuY29kZWQgYm9keSE=",
	}

	decodedData, err := messageBody.GetDecodedContent()
	suite.Equal("Base64 encoded body!", decodedData)
	suite.Nil(err)
}

func (suite *messageBodyTestSuite) TestGetDecodedData_ShouldReturnErrorIfUnableToDecode() {
	messageBody := MessageBody{
		Content: "Message body not in base64 encoded format",
	}

	_, err := messageBody.GetDecodedContent()
	suite.NotNil(err)
}
