package main

import (
	"flag"

	cmdAgent "github.com/nicocha30/ligolo-ng/cmd/agent"
	cmdProxy "github.com/nicocha30/ligolo-ng/cmd/proxy"
	"github.com/nicocha30/ligolo-ng/types"
)

func main() {
	// Type of application to run
	var cmdType = flag.String("type", "proxy", "Type of the application to run")

	// Common flags
	var versionFlag = flag.Bool("version", false, "show the current version")
	var verbose = flag.Bool("v", false, "enable verbose mode")

	// Proxy command flags
	var listenInterface = flag.String("laddr", "0.0.0.0:11601", "listening address (prefix with https:// for websocket)")
	var enableAutocert = flag.Bool("autocert", false, "automatically request letsencrypt certificates, requires port 80 to be accessible")
	var enableSelfcert = flag.Bool("selfcert", false, "dynamically generate self-signed certificates")
	var certFile = flag.String("certfile", "certs/cert.pem", "TLS server certificate")
	var keyFile = flag.String("keyfile", "certs/key.pem", "TLS server key")
	var domainWhitelist = flag.String("allow-domains", "", "autocert authorised domains, if empty, allow all domains, multiple domains should be comma-separated.")
	var selfcertDomain = flag.String("selfcert-domain", "ligolo", "The selfcert TLS domain to use")
	var hideBanner = flag.Bool("nobanner", false, "don't show banner on startup")
	var configFile = flag.String("config", "", "the config file to use")
	var daemonMode = flag.Bool("daemon", false, "run as daemon mode (no CLI)")
	var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
	var memprofile = flag.String("memprofile", "", "write memory profile to `file`")

	// Agent command flags
	var ignoreCertificate = flag.Bool("ignore-cert", false, "ignore TLS certificate validation (dangerous), only for debug purposes")
	var acceptFingerprint = flag.String("accept-fingerprint", "", "accept certificates matching the following SHA256 fingerprint (hex format)")
	var retry = flag.Bool("retry", false, "auto-retry on error")
	var socksProxy = flag.String("proxy", "", "proxy URL address (http://admin:secret@127.0.0.1:8080)"+
		" or socks://admin:secret@127.0.0.1:8080")
	var serverAddr = flag.String("connect", "", "connect to proxy (domain:port)")
	var bindAddr = flag.String("bind", "", "bind to ip:port")
	var userAgent = flag.String("ua", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) "+
		"AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36", "HTTP User-Agent")

	flag.Parse() // Parse command line flags

	switch *cmdType {
	case "proxy", "server":
		flags := &types.ProxyFlags{
			Verbose:         verbose,
			ListenInterface: listenInterface,
			EnableAutocert:  enableAutocert,
			EnableSelfcert:  enableSelfcert,
			CertFile:        certFile,
			KeyFile:         keyFile,
			DomainWhitelist: domainWhitelist,
			SelfcertDomain:  selfcertDomain,
			VersionFlag:     versionFlag,
			HideBanner:      hideBanner,
			ConfigFile:      configFile,
			DaemonMode:      daemonMode,
			CPUProfile:      cpuprofile,
			MemProfile:      memprofile,
		}
		cmdProxy.Proxy(flags)
	case "agent", "client":
		flags := &types.AgentFlags{
			IgnoreCertificate: ignoreCertificate,
			AcceptFingerprint: acceptFingerprint,
			Verbose:           verbose,
			Retry:             retry,
			SocksProxy:        socksProxy,
			ServerAddr:        serverAddr,
			BindAddr:          bindAddr,
			UserAgent:         userAgent,
			VersionFlag:       versionFlag,
		}
		cmdAgent.Agent(flags)
	default:
		panic("Unknown command type: " + *cmdType)
	}
}
