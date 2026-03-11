package runner

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

// isOnionAlive connects to the Tor ControlPort, authenticates, and checks if
// the specified hidden service descriptor (.onion) is alive on the Tor network.
func isOnionAlive(onionAddress string, controlPort string, authCookie string) (bool, error) {
	conn, err := net.DialTimeout("tcp", controlPort, 5*time.Second)
	if err != nil {
		return false, fmt.Errorf("failed to connect to tor control port %s: %w", controlPort, err)
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)

	// Authenticate
	authCmd := fmt.Sprintf("AUTHENTICATE \"%s\"\r\n", authCookie)
	_, err = fmt.Fprint(conn, authCmd)
	if err != nil {
		return false, fmt.Errorf("failed to send authenticate command: %w", err)
	}

	authReply, err := reader.ReadString('\n')
	if err != nil {
		return false, fmt.Errorf("failed to read authenticate reply: %w", err)
	}
	if !strings.HasPrefix(authReply, "250") {
		return false, fmt.Errorf("authentication failed, expected 250 but got: %s", authReply)
	}

	// Subscribe to HS_DESC events to get the fetch results
	_, err = fmt.Fprintf(conn, "SETEVENTS HS_DESC\r\n")
	if err != nil {
		return false, fmt.Errorf("failed to send SETEVENTS command: %w", err)
	}

	// Read the 250 OK for SETEVENTS
	eventsReply, err := reader.ReadString('\n')
	if err != nil || !strings.HasPrefix(eventsReply, "250") {
		return false, fmt.Errorf("failed to set events: %v, reply: %s", err, eventsReply)
	}

	// Prepare HSFETCH payload, we must strip .onion suffix according to spec
	onionAddress = strings.TrimSuffix(strings.TrimSpace(onionAddress), ".onion")

	// The syntax is HSFETCH <address> where address is the v3 descriptor ID (without .onion)
	_, err = fmt.Fprintf(conn, "HSFETCH %s\r\n", onionAddress)
	if err != nil {
		return false, fmt.Errorf("failed to send HSFETCH command: %w", err)
	}

	// Read the 250 OK for HSFETCH (it means the fetch started, not that it succeeded)
	fetchStartReply, err := reader.ReadString('\n')
	if err != nil || !strings.HasPrefix(fetchStartReply, "250") {
		return false, fmt.Errorf("HSFETCH failed to start: %v, reply: %s", err, fetchStartReply)
	}

	// Read asynchronously for 650 HSFETCH reply or 5xx error
	// Set a deadline since HSFETCH waits for directory authorities to respond
	// This can take a while on the Tor network
	conn.SetReadDeadline(time.Now().Add(35 * time.Second))

	for {
		reply, err := reader.ReadString('\n')
		if err != nil {
			return false, fmt.Errorf("failed during HS_DESC read: %w", err)
		}

		reply = strings.TrimSpace(reply)

		if strings.HasPrefix(reply, "650 HS_DESC") {
			// Look for our specific onion address in the event
			if strings.Contains(reply, onionAddress) {
				if strings.Contains(reply, "RECEIVED") {
					return true, nil
				} else if strings.Contains(reply, "FAILED") {
					return false, nil
				}
			}
		}
	}
}
