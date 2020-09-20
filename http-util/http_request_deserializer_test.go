package http_util

import (
	"bytes"
	"ccg-api/util"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type httpDeserializerTestSuite struct {
	suite.Suite
	mockCtrl                *gomock.Controller
	recorder                *httptest.ResponseRecorder
	context                 *gin.Context
	httpRequestDeserializer HttpRequestDeserializer
}

type TestRequest struct {
	MandatoryStringField string   `json:"mandatory_string_field" binding:"required"`
	OptionalStringField  string   `json:"optional_string_field"`
	NonEmptyStringArray  []string `json:"non_empty_string_array" binding:"required" validate:"ne=0"`
	FixedWidthString     string   `json:"fixed_width_string" validate:"omitempty,len=5"`
}

func TestHttpDeserializerTestSuite(t *testing.T) {
	suite.Run(t, new(httpDeserializerTestSuite))
}

func (suite *httpDeserializerTestSuite) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.recorder = httptest.NewRecorder()
	suite.context, _ = gin.CreateTestContext(suite.recorder)
	suite.httpRequestDeserializer = NewHttpRequestDeserializer(validator.New())
}

func (suite httpDeserializerTestSuite) TestShouldBindJsonIfValid_ShouldDeserializeValidRequest() {

	expectedRequest := TestRequest{
		MandatoryStringField: "Some Value",
		NonEmptyStringArray:  []string{"First Element"},
		FixedWidthString:     "12345",
	}
	requestBody, _ := util.Encode(expectedRequest)
	suite.context.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString(requestBody))

	var actualRequest TestRequest
	bindError := suite.httpRequestDeserializer.ShouldBindJsonBodyIfValid(&actualRequest, suite.context)

	suite.Equal(expectedRequest, actualRequest)
	suite.Nil(bindError)
}

func (suite httpDeserializerTestSuite) TestShouldBindJsonIfValid_ShouldReturnErrorForIncompleteRequest() {
	erroneousRequest := TestRequest{
		NonEmptyStringArray: []string{"First Element"},
	}

	requestBody, _ := util.Encode(erroneousRequest)
	suite.context.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString(requestBody))

	var actualRequest TestRequest
	deserializationError := suite.httpRequestDeserializer.ShouldBindJsonBodyIfValid(&actualRequest, suite.context)

	suite.NotNil(deserializationError)
}

func (suite httpDeserializerTestSuite) TestShouldBindJsonIfValid_ShouldReturnErrorForInvalidRequestHavingEmptyArrayField() {
	erroneousRequest := TestRequest{
		MandatoryStringField: "Some Value",
		NonEmptyStringArray:  []string{},
	}

	requestBody, _ := util.Encode(erroneousRequest)
	suite.context.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString(requestBody))

	var actualRequest TestRequest
	deserializationError := suite.httpRequestDeserializer.ShouldBindJsonBodyIfValid(&actualRequest, suite.context)

	suite.NotNil(deserializationError)
}

func (suite httpDeserializerTestSuite) TestShouldBindJsonIfValid_ShouldReturnErrorForInvalidRequestHavingInvalidFixedWidthField() {
	erroneousRequest := TestRequest{
		MandatoryStringField: "Some Value",
		NonEmptyStringArray:  []string{"First Element", "Second Element"},
		FixedWidthString:     "ABCDE12345",
	}

	requestBody, _ := util.Encode(erroneousRequest)
	suite.context.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString(requestBody))

	var actualRequest TestRequest
	deserializationError := suite.httpRequestDeserializer.ShouldBindJsonBodyIfValid(&actualRequest, suite.context)

	suite.NotNil(deserializationError)
}
