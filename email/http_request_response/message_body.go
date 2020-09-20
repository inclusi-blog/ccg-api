package http_request_response

import "encoding/base64"

type MessageBody struct {
	MimeType string `json:"mime_type" example:"text/html"` //Use text/plain as default
	Content  string `json:"base64_encoded_content" validate:"required,base64,notblankbase64" example:"base64 encoded value"`
}

func (message MessageBody) GetDecodedContent() (string, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(message.Content)
	if err != nil {
		return "", err
	}
	return string(decodedBytes), nil
}
