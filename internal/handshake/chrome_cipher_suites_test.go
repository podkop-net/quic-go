package handshake

import (
	"slices"
	"testing"

	utls "github.com/metacubex/utls"
)

// Chrome's order of the suites follows the AES hardware of the CPU, ChaCha20 first without it; the
// spec of the QUIC ClientHello takes the order of the CPU it runs on.
func TestChromeCipherSuites(t *testing.T) {
	withAES := []uint16{utls.TLS_AES_128_GCM_SHA256, utls.TLS_AES_256_GCM_SHA384, utls.TLS_CHACHA20_POLY1305_SHA256}
	withoutAES := []uint16{utls.TLS_CHACHA20_POLY1305_SHA256, utls.TLS_AES_128_GCM_SHA256, utls.TLS_AES_256_GCM_SHA384}
	if got := chromeCipherSuites(true); !slices.Equal(got, withAES) {
		t.Fatalf("with AES hardware: %x, want %x", got, withAES)
	}
	if got := chromeCipherSuites(false); !slices.Equal(got, withoutAES) {
		t.Fatalf("without AES hardware: %x, want %x", got, withoutAES)
	}
	spec := chromeQUICClientHelloSpec(nil).CipherSuites
	if !slices.Equal(spec, chromeCipherSuites(hasAESGCMHardwareSupport)) {
		t.Fatalf("the spec offers %x on a CPU whose AES-GCM hardware is %v", spec, hasAESGCMHardwareSupport)
	}
	t.Logf("this CPU: AES-GCM hardware %v, suites %x", hasAESGCMHardwareSupport, spec)
}
