package configuration

type ConfigData struct {
	Environment        string                       `json:"environment"`
	Email              Email                        `json:"email"`
	OpenTracingUrl     string                       `json:"open_tracing_url"`
	TracingServiceName string                       `json:"tracing_service_name" binding:"required"`
	TracingOCAgentHost string                       `json:"tracing_oc_agent_host" binding:"required"`
	LogLevel           string                       `json:"log_level" binding:"required"`
}

type Email struct {
	SmtpHost                         string   `json:"smtp_host"`
	SmtpPort                         int      `json:"smtp_port"`
	Username                         string   `json:"username"`
	InsecureSkipVerify               bool     `json:"insecure_skip_verify"`
	ValidMensuvadiEmailDomains       []string `json:"valid_mensuvadi_email_domains"`
	DefaultMensuvadiEmailSender      string   `json:"default_mensuvadi_email_sender"`
	UnsupportedAttachmentExtensions  []string `json:"unsupported_attachment_extensions"`
	PermissibleAttachmentSizeInBytes int      `json:"permissible_attachment_size_in_bytes"`
	BaseTemplateFilePath             string   `json:"base_template_file_path"`
	LogoUrls                         LogoUrls `json:"logo_urls"`
}

type LogoUrls struct {
	Mensuvadi       string `json:"mensuvadi"`
	Facebook        string `json:"facebook"`
	Instagram       string `json:"instagram"`
	Twitter         string `json:"twitter"`
	LinkedIn        string `json:"linkedin"`
	DownloadIOS     string `json:"download_ios"`
	DownloadAndroid string `json:"download_android"`
	HelpCenter      string `json:"help_center_url"`
}
