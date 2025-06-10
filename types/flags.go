package types

// AgentFlags holds the configuration for the agent application.
type AgentFlags struct {
	IgnoreCertificate *bool   // Ignore TLS certificate validation
	AcceptFingerprint *string // Accept certificates matching the given SHA256 fingerprint (hex format)
	Verbose           *bool   // Enable verbose mode
	Retry             *bool   // Auto-retry on error
	SocksProxy        *string // Proxy URL address (http://admin:secret@
	ServerAddr        *string // Connect to proxy (domain:port)
	BindAddr          *string // Bind to ip:port
	UserAgent         *string // HTTP User-Agent
	VersionFlag       *bool   // Show the current version
}

// ProxyFlags holds the configuration for the proxy application.
type ProxyFlags struct {
	Verbose         *bool   // Enable verbose mode
	ListenInterface *string // Listening address
	EnableAutocert  *bool   // Enable automatic Letsencrypt certificates
	EnableSelfcert  *bool   // Enable self-signed certificates
	CertFile        *string // TLS certificate file path
	KeyFile         *string // TLS key file path
	DomainWhitelist *string // Allowed domains for autocert
	SelfcertDomain  *string // Domain for self-signed cert
	VersionFlag     *bool   // Show version
	HideBanner      *bool   // Hide startup banner
	ConfigFile      *string // Config file path
	DaemonMode      *bool   // Run in daemon mode
	CPUProfile      *string // CPU profile output file
	MemProfile      *string // Memory profile output file
}
