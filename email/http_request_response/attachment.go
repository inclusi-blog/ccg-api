package http_request_response

import "encoding/base64"

type Attachment struct {
	FileName string `json:"file_name" validate:"required,validFileExtension" example:"fileName.pdf"`
	Data     string `json:"base64_encoded_data" validate:"required,base64" example:"base64 encoded value"`
}

func (attachment Attachment) GetDecodedData() ([]byte, error) {
	return base64.StdEncoding.DecodeString(attachment.Data)
}
