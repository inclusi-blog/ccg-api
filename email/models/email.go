package models

type Email struct {
	From                string
	To                  []string
	Subject             string
	Body                MessageBody
	Attachments         []Attachment
	IncludeBaseTemplate bool
}

type Attachment struct {
	FileName string
	Data     []byte
}

type MessageBody struct {
	MimeType string
	Content  string
}
