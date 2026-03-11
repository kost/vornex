package runner

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/gologger/levels"
	"github.com/kost/vornex/v2/pkg/port"
	"github.com/kost/vornex/v2/pkg/privileges"
	"github.com/kost/vornex/v2/pkg/protocol"
	"github.com/kost/vornex/v2/pkg/result"
	"github.com/kost/vornex/v2/pkg/scan"
	"github.com/stretchr/testify/require"
)

func TestTorProbeDiscovery(t *testing.T) {
	// 1. Check if Tor ControlPort is available, else skip
	conn, err := net.DialTimeout("tcp", "127.0.0.1:9051", 1*time.Second)
	if err != nil {
		t.Skip("Skipping Tor test because control port 127.0.0.1:9051 is not active: " + err.Error())
		return
	}
	conn.Close()

	privileges.IsPrivileged = true

	// 2. Configure Runner for Host Discovery Only
	options := Options{
		Host:              []string{"127.0.0.1", "2gzyxa5ihm7nsggfxnu52rck2vv4rvmdlkiu3zzui5du4xyclen53wid.onion", "invalid-fake-dead-onion-address.onion"},
		OnlyHostDiscovery: true,
		WithHostDiscovery: true,
		ScanType:          ConnectScan, // SynScan requires raw sockets
		ProbeTor:          "127.0.0.1:9051",
		TorPassword:       "", // Default local tests have no Tor password
		IPVersion:         []string{scan.IPv4},
	}
	options.configureHostDiscovery([]*port.Port{{Port: 80, Protocol: protocol.TCP}, {Port: 443, Protocol: protocol.TCP}, {Port: 8443, Protocol: protocol.TCP}})

	runner, err := NewRunner(&options)
	require.Nil(t, err, "could not create runner")
	defer runner.Close()

	// Capture outputs from runner
	discoveredHosts := make([]string, 0)
	runner.options.OnResult = func(hr *result.HostResult) {
		discoveredHosts = append(discoveredHosts, hr.Host)
	}

	gologger.DefaultLogger.SetMaxLevel(levels.LevelSilent) // Hide console output from interfering test runner

	// 3. Run Enumeration loop
	err = runner.RunEnumeration(context.Background())
	require.Nil(t, err, "failed to run enumeration")

	for discoveredIP := range runner.scanner.HostDiscoveryResults.GetIPs() {
		t.Logf("IP Detected in HostDiscoveryResults: %s", discoveredIP)
		discoveredHosts = append(discoveredHosts, discoveredIP)
	}

	// 4. Assertions
	require.Contains(t, discoveredHosts, "2gzyxa5ihm7nsggfxnu52rck2vv4rvmdlkiu3zzui5du4xyclen53wid.onion", "Expected Tor Project onion to be resolved alive")
	require.NotContains(t, discoveredHosts, "invalid-fake-dead-onion-address.onion", "Expected fake onion address to be dead")
}
