package http_request_response

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type attachmentTestSuite struct {
	suite.Suite
}

func TestAttachmentTestSuite(t *testing.T) {
	suite.Run(t, new(attachmentTestSuite))
}

func (suite *attachmentTestSuite) TestGetDecodedData_ShouldReturnBase64DecodedData() {
	attachment := Attachment{
		FileName: "attachment.pdf",
		Data:     "QXR0YWNobWVudCB3aXRoIHNvbWUgZGF0YSE=",
	}

	decodedData, err := attachment.GetDecodedData()
	suite.Equal([]byte("Attachment with some data!"), decodedData)
	suite.Nil(err)
}

func (suite *attachmentTestSuite) TestGetDecodedData_ShouldReturnErrorIfUnableToDecode() {
	attachment := Attachment{
		FileName: "attachment.pdf",
		Data:     "Not base64 encoded data",
	}

	_, err := attachment.GetDecodedData()
	suite.NotNil(err)
}
