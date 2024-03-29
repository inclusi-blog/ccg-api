// Code generated by MockGen. DO NOT EDIT.
// Source: email/configuration/email_client_config.go

// Package mocks is a generated GoMock package.
package mocks

import (
	configuration "ccg-api/configuration"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockEmailClientConfig is a mock of EmailClientConfig interface
type MockEmailClientConfig struct {
	ctrl     *gomock.Controller
	recorder *MockEmailClientConfigMockRecorder
}

// MockEmailClientConfigMockRecorder is the mock recorder for MockEmailClientConfig
type MockEmailClientConfigMockRecorder struct {
	mock *MockEmailClientConfig
}

// NewMockEmailClientConfig creates a new mock instance
func NewMockEmailClientConfig(ctrl *gomock.Controller) *MockEmailClientConfig {
	mock := &MockEmailClientConfig{ctrl: ctrl}
	mock.recorder = &MockEmailClientConfigMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockEmailClientConfig) EXPECT() *MockEmailClientConfigMockRecorder {
	return m.recorder
}

// SmtpHost mocks base method
func (m *MockEmailClientConfig) SmtpHost() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SmtpHost")
	ret0, _ := ret[0].(string)
	return ret0
}

// SmtpHost indicates an expected call of SmtpHost
func (mr *MockEmailClientConfigMockRecorder) SmtpHost() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SmtpHost", reflect.TypeOf((*MockEmailClientConfig)(nil).SmtpHost))
}

// SmtpPort mocks base method
func (m *MockEmailClientConfig) SmtpPort() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SmtpPort")
	ret0, _ := ret[0].(int)
	return ret0
}

// SmtpPort indicates an expected call of SmtpPort
func (mr *MockEmailClientConfigMockRecorder) SmtpPort() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SmtpPort", reflect.TypeOf((*MockEmailClientConfig)(nil).SmtpPort))
}

// Username mocks base method
func (m *MockEmailClientConfig) Username() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Username")
	ret0, _ := ret[0].(string)
	return ret0
}

// Username indicates an expected call of Username
func (mr *MockEmailClientConfigMockRecorder) Username() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Username", reflect.TypeOf((*MockEmailClientConfig)(nil).Username))
}

// Password mocks base method
func (m *MockEmailClientConfig) Password() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Password")
	ret0, _ := ret[0].(string)
	return ret0
}

// Password indicates an expected call of Password
func (mr *MockEmailClientConfigMockRecorder) Password() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Password", reflect.TypeOf((*MockEmailClientConfig)(nil).Password))
}

// InsecureSkipVerify mocks base method
func (m *MockEmailClientConfig) InsecureSkipVerify() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsecureSkipVerify")
	ret0, _ := ret[0].(bool)
	return ret0
}

// InsecureSkipVerify indicates an expected call of InsecureSkipVerify
func (mr *MockEmailClientConfigMockRecorder) InsecureSkipVerify() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsecureSkipVerify", reflect.TypeOf((*MockEmailClientConfig)(nil).InsecureSkipVerify))
}

// TempDir mocks base method
func (m *MockEmailClientConfig) TempDir() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TempDir")
	ret0, _ := ret[0].(string)
	return ret0
}

// TempDir indicates an expected call of TempDir
func (mr *MockEmailClientConfigMockRecorder) TempDir() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TempDir", reflect.TypeOf((*MockEmailClientConfig)(nil).TempDir))
}

// ValidGolaEmailDomain mocks base method
func (m *MockEmailClientConfig) ValidGolaEmailDomain() []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidGolaEmailDomain")
	ret0, _ := ret[0].([]string)
	return ret0
}

// ValidGolaEmailDomain indicates an expected call of ValidGolaEmailDomain
func (mr *MockEmailClientConfigMockRecorder) ValidGolaEmailDomain() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidGolaEmailDomain", reflect.TypeOf((*MockEmailClientConfig)(nil).ValidGolaEmailDomain))
}

// DefaultGolaEmailSender mocks base method
func (m *MockEmailClientConfig) DefaultGolaEmailSender() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DefaultGolaEmailSender")
	ret0, _ := ret[0].(string)
	return ret0
}

// DefaultGolaEmailSender indicates an expected call of DefaultGolaEmailSender
func (mr *MockEmailClientConfigMockRecorder) DefaultGolaEmailSender() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DefaultGolaEmailSender", reflect.TypeOf((*MockEmailClientConfig)(nil).DefaultGolaEmailSender))
}

// UnsupportedAttachmentExtensions mocks base method
func (m *MockEmailClientConfig) UnsupportedAttachmentExtensions() []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnsupportedAttachmentExtensions")
	ret0, _ := ret[0].([]string)
	return ret0
}

// UnsupportedAttachmentExtensions indicates an expected call of UnsupportedAttachmentExtensions
func (mr *MockEmailClientConfigMockRecorder) UnsupportedAttachmentExtensions() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnsupportedAttachmentExtensions", reflect.TypeOf((*MockEmailClientConfig)(nil).UnsupportedAttachmentExtensions))
}

// PermissibleTotalSizeOfAttachments mocks base method
func (m *MockEmailClientConfig) PermissibleTotalSizeOfAttachments() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PermissibleTotalSizeOfAttachments")
	ret0, _ := ret[0].(int)
	return ret0
}

// PermissibleTotalSizeOfAttachments indicates an expected call of PermissibleTotalSizeOfAttachments
func (mr *MockEmailClientConfigMockRecorder) PermissibleTotalSizeOfAttachments() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PermissibleTotalSizeOfAttachments", reflect.TypeOf((*MockEmailClientConfig)(nil).PermissibleTotalSizeOfAttachments))
}

// BaseTemplateFilePath mocks base method
func (m *MockEmailClientConfig) BaseTemplateFilePath() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BaseTemplateFilePath")
	ret0, _ := ret[0].(string)
	return ret0
}

// BaseTemplateFilePath indicates an expected call of BaseTemplateFilePath
func (mr *MockEmailClientConfigMockRecorder) BaseTemplateFilePath() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BaseTemplateFilePath", reflect.TypeOf((*MockEmailClientConfig)(nil).BaseTemplateFilePath))
}

// LogoUrls mocks base method
func (m *MockEmailClientConfig) LogoUrls() configuration.LogoUrls {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LogoUrls")
	ret0, _ := ret[0].(configuration.LogoUrls)
	return ret0
}

// LogoUrls indicates an expected call of LogoUrls
func (mr *MockEmailClientConfigMockRecorder) LogoUrls() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LogoUrls", reflect.TypeOf((*MockEmailClientConfig)(nil).LogoUrls))
}

// OtherUrls mocks base method
func (m *MockEmailClientConfig) OtherUrls() configuration.Urls {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OtherUrls")
	ret0, _ := ret[0].(configuration.Urls)
	return ret0
}

// OtherUrls indicates an expected call of OtherUrls
func (mr *MockEmailClientConfigMockRecorder) OtherUrls() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OtherUrls", reflect.TypeOf((*MockEmailClientConfig)(nil).OtherUrls))
}
