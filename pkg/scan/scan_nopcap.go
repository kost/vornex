//go:build nopcap

package scan

import (
	"errors"
	"time"

	"github.com/kost/vornex/v2/pkg/port"
	"github.com/kost/vornex/v2/pkg/protocol"
)

func init() {
}

func buildListenHandler() (*ListenHandler, error) {
	return nil, errors.New("pcap is disabled")
}

func ICMPWriteWorker() {
}

func EthernetWriteWorker() {
}

func TransportWriteWorker() {
}

func SendAsyncPkg(listenHandler *ListenHandler, ip string, p *port.Port, pkgFlag PkgFlag) {
}

func sendAsyncTCP4(listenHandler *ListenHandler, ip string, p *port.Port, pkgFlag PkgFlag) {
}

func sendAsyncUDP4(listenHandler *ListenHandler, ip string, p *port.Port, pkgFlag PkgFlag) {
}

func sendAsyncTCP6(listenHandler *ListenHandler, ip string, p *port.Port, pkgFlag PkgFlag) {
}

func sendAsyncUDP6(listenHandler *ListenHandler, ip string, p *port.Port, pkgFlag PkgFlag) {
}

func (l *ListenHandler) ICMPReadWorker4() {
}

func (l *ListenHandler) ICMPReadWorker6() {
}

func (l *ListenHandler) TcpReadWorker4() {
}

func (l *ListenHandler) TcpReadWorker6() {
}

func (l *ListenHandler) UdpReadWorker4() {
}

func (l *ListenHandler) UdpReadWorker6() {
}

func SetupHandlerUnix(interfaceName, bpfFilter string, protocols ...protocol.Protocol) error {
	return errors.New("pcap is disabled")
}

func TransportReadWorker() {
}

func CleanupHandlersUnix() {
}

func SetupHandlers() error {
	return errors.New("pcap is disabled")
}

func SetupHandler(interfaceName string) error {
	return errors.New("pcap is disabled")
}

func ACKPort(listenHandler *ListenHandler, dstIP string, port int, timeout time.Duration) (bool, error) {
	return false, errors.New("pcap is disabled")
}
