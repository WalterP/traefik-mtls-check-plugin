package traefik_mtls_check_plugin

import (
	"context"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
)

// mTlsCheckTypeName This is the name of the plugin. It is used internally by traefik to identify the plugin.
const (
	mTlsCheckTypeName = "mTlsCheck"
)

// mTlsCheckTypeBuilder This is the interface that is used by traefik to build the plugin.
type serviceBuilder interface {
	BuildHTTP(ctx context.Context, serviceName string) (http.Handler, error)
}

// Config This is the configuration struct. Here we are defining the configuration options that the user can set in the traefik.toml file.
type Config struct {
	ResponseCode int    `yaml:"responseCode" json:"responseCode"`
	CaCert       string `yaml:"caCert" json:"caCert"`
	CACertPath   string `yaml:"caCertPath" json:"caCertPath"`
	Message      string `yaml:"message" json:"message"`
}

// mTlsCheck This is the plugin struct. It will contain the next http.Handler in the chain, the name of the plugin, and the message to return.
type mTlsCheck struct {
	config *Config
	next   http.Handler
	name   string
}

// CreateConfig This function is used to create the default configuration for the plugin. It is called when the user does not.
// provide a configuration for the plugin in the traefik.toml file.
func CreateConfig() *Config {
	return &Config{
		Message:      "Not found",
		ResponseCode: 404,
		CaCert:       "",
		CACertPath:   "",
	}
}

// Init This function is used to initialize the plugin. It is called once when traefik starts.

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {

	logMessage := fmt.Sprintf("Creating tlsredirect plugin %q with config %+v\n", name, config)
	fmt.Printf(logMessage)

	// Return new plugin instance
	mtlsCheck := &mTlsCheck{
		config: config,
		next:   next,
		name:   name,
	}
	return mtlsCheck, nil
}

// validateCert validates a certificate against a given CA certificate.
// It returns true if the certificate is valid, false otherwise.
func validateCert(cert *x509.Certificate, caCert []byte) bool {
	// Create a new certificate pool and add the CA certificate to it.
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(caCert)

	// Verify the certificate using the certificate pool.
	_, err := cert.Verify(x509.VerifyOptions{
		Roots: certPool,
	})

	// If there is no error, the certificate is valid.
	if err == nil {
		fmt.Printf("Found valid certificate: %v\n", cert.Subject)
		return true
	}

	// If there is an error, print the error message.
	fmt.Sprintf("Error: %v\n", err)
	return false
}

// readCaCertFromFile reads a CA certificate from a file and returns it as a byte array.
func readCaCertFromFile(caCertPath string) []byte {
	caCert, err := os.ReadFile(caCertPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	return caCert
}

// ServeHTTP This is the main function of the plugin. It is called for each request that traefik receives.
func (p *mTlsCheck) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	valid := false

	if r.TLS != nil && len(r.TLS.PeerCertificates) > 0 {
		var caCert []byte
		if p.config.CACertPath != "" {
			caCert = readCaCertFromFile(p.config.CACertPath)
		} else {
			caCert = []byte(p.config.CaCert)
		}

		valid = validateCert(r.TLS.PeerCertificates[0], caCert)
	}
	if valid {
		p.next.ServeHTTP(w, r)
		return
	} else {
		w.WriteHeader(p.config.ResponseCode)
		message := fmt.Sprintf("Current message:= %v", p.config.Message)
		w.Write([]byte(message))
		return

	}
}
