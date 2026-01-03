package proxy

import (
	"io"
	"math/rand/v2"
	"net"
	"time"

	"github.com/Jonathanthedeveloper/havoc.git/internal/state"
)

func Copy(state *state.HavocState, src net.Conn, dst net.Conn) error {

	// 32KB buffer
	buf := make([]byte, 32*1024)

	for {
		// Read
		read, readErr := src.Read(buf)

		if read > 0 {
			settings := state.GetChaos()

			// Drop
			if settings.DropRate > 0 && rand.Float64() < settings.DropRate {
				continue
			}

			// Latency and Jitter
			delay := settings.Latency

			if settings.Jitter > 0 {
				variance := time.Duration(rand.Float64() * float64(settings.Jitter))
				delay += variance
			}

			if delay > 0 {
				time.Sleep(delay)
			}

			// Write
			_, writeErr := dst.Write(buf[:read])

			if writeErr != nil {
				return writeErr
			}

		}

		if readErr != nil {

			if readErr == io.EOF {
				return nil
			}

			return readErr
		}
	}
}
