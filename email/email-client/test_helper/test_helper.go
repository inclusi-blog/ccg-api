package test_helper

import (
	"bufio"
	"bytes"
	"gopkg.in/gomail.v2"
	"regexp"
	"strings"
)

func SanitizeMessageContent(message string) string {
	//Remove date component
	index := strings.Index(message, "Content-Type:")
	message = message[index:]

	//Remove message boundary identifier
	index = strings.Index(message, "boundary=") + len("boundary=")
	messageIdentifier := message[index : index+len("5de10fdd5e00e6fadb59e85a1eb0e114db94cd6b1aa36c370687be1c1e57")]
	message = strings.ReplaceAll(message, messageIdentifier, "")
	message = removeContentTypeForAttachments(message)
	//Replace \r\n with |
	re := regexp.MustCompile(`\r?\n`)
	return re.ReplaceAllString(message, "|")
}

func removeContentTypeForAttachments(message string) string {
	lines := strings.Split(message, "\n")
	var sb strings.Builder
	for _, line := range lines {
		if strings.HasPrefix(string(line), "Content-Type") {
			contentData := strings.Split(string(line), ";")
			nameData := contentData[len(contentData)-1]
			if len(contentData) >= 2 && strings.Contains(nameData, "name=") {
				sb.WriteString(nameData)
			} else {
				sb.WriteString(string(line))
			}
		} else {
			sb.WriteString(string(line))
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func SanitizeMessage(message *gomail.Message) string {
	return SanitizeMessageContent(toString(message))
}

func toString(gomailMessage *gomail.Message) string {
	var gomailMessageBuffer bytes.Buffer
	gomailMessageWriter := bufio.NewWriter(&gomailMessageBuffer)
	contentLength, _ := gomailMessage.WriteTo(gomailMessageWriter)
	_ = gomailMessageWriter.Flush()
	gomailMessageContent := gomailMessageBuffer.Next(int(contentLength))
	return string(gomailMessageContent)
}
