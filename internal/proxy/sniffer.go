package proxy

import (
	"bytes"
	"net"
	"strings"
	"sync"
)

// Sniffer wraps a net.Conn and peeks at the initial bytes to detect HTTP metadata.
type Sniffer struct {
	net.Conn
	buffer    bytes.Buffer
	mu        sync.Mutex
	captured  bool
	Method    string
	Path      string
	Status    string
	isRequest bool
}

// NewSniffer creates a new Sniffer. isRequest should be true for client->server connections.
func NewSniffer(conn net.Conn, isRequest bool) *Sniffer {
	return &Sniffer{
		Conn:      conn,
		isRequest: isRequest,
	}
}

func (s *Sniffer) Read(p []byte) (n int, err error) {
	n, err = s.Conn.Read(p)
	if n > 0 {
		s.mu.Lock()
		defer s.mu.Unlock()

		if !s.captured {
			// Append read bytes to our internal buffer
			s.buffer.Write(p[:n])

			// Try to parse what we have
			data := s.buffer.String()

			// We only need the first line typically
			if newlineIdx := strings.Index(data, "\n"); newlineIdx != -1 {
				firstLine := data[:newlineIdx]
				parts := strings.Fields(firstLine)

				if len(parts) >= 2 {
					if s.isRequest {
						// Request: GET /path HTTP/1.1
						// We want method and path
						s.Method = parts[0]
						if len(parts) > 1 {
							s.Path = parts[1]
						}
					} else {
						// Response: HTTP/1.1 200 OK
						// We want the status code
						s.Status = parts[1]
					}
					s.captured = true
					// Free the buffer memory as we don't need it anymore
					s.buffer.Reset()
				}
			} else if s.buffer.Len() > 4096 {
				// If we haven't found a newline in 4KB, give up to save memory
				s.captured = true
				s.buffer.Reset()
			}
		}
	}
	return n, err
}
