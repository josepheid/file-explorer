package api

import (
	"crypto/tls"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"

	"github.com/goteleport-interview/fs4/api/handlers"
	"github.com/goteleport-interview/fs4/api/internal/auth"
	"github.com/goteleport-interview/fs4/api/internal/middleware"
	"github.com/goteleport-interview/fs4/api/internal/sessions"
	"github.com/rs/cors"
)

// Server serves the directory browser API and webapp.
type Server struct {
	handler http.Handler
}

// NewServer creates a directory browser server.
// It serves webassets from the provided filesystem.
func NewServer(webassets fs.FS, rootPath string) (*Server, error) {
	mux := http.NewServeMux()
	s := &Server{handler: mux}

	auth := auth.New()
	session := sessions.New()

	// API routes
	mux.Handle("POST /api/v1/login", handlers.NewLoginHandler(auth, session))
	mux.Handle("POST /api/v1/logout", handlers.NewLogoutHandler(session))

	// Protected routes
	mux.Handle("GET /api/v1/browse", middleware.RequireAuth(session)(handlers.NewBrowseHandler(rootPath)))

	// web assets
	hfs := http.FS(webassets)
	files := http.FileServer(hfs)
	mux.Handle("/assets/", files)
	mux.Handle("/favicon.ico", files)

	// fall back to index.html for all unknown routes
	index, err := extractIndexHTML(hfs)
	if err != nil {
		return nil, err
	}
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write(index); err != nil {
			log.Println("failed to serve index.html", err)
		}
	}))

	s.handler = cors.Default().Handler(mux)

	return s, nil
}

func (s *Server) ListenAndServe(addr string) error {
	certFile := "./api/internal/certs/localhost.pem"
	keyFile := "./api/internal/certs/localhost-key.pem"

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return fmt.Errorf("failed to load TLS certificate: %w", err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12, // Minimum TLS 1.2
		MaxVersion:   tls.VersionTLS13, // Maximum TLS 1.3
		CipherSuites: []uint16{
			// TLS 1.3 cipher suites
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_CHACHA20_POLY1305_SHA256,

			// TLS 1.2 cipher suites (secure ones only)
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
		},
		PreferServerCipherSuites: true,
	}

	server := &http.Server{
		Addr:      addr,
		Handler:   s.handler,
		TLSConfig: tlsConfig,
	}

	return server.ListenAndServeTLS(certFile, keyFile)
}

func extractIndexHTML(fs http.FileSystem) ([]byte, error) {
	f, err := fs.Open("index.html")
	if err != nil {
		return nil, fmt.Errorf("could not open index.html: %w", err)
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("could not read index.html: %w", err)
	}

	return b, nil
}
