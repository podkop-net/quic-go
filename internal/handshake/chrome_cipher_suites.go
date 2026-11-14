package handshake

import (
	"runtime"

	utls "github.com/metacubex/utls"
	"golang.org/x/sys/cpu"
)

// hasAESGCMHardwareSupport tells a CPU on which AES-GCM is fast, as crypto/tls tells it (Go 1.26,
// cipher_suites.go, kept in sync with crypto/internal/fips140/aes/gcm.supportsAESGCM): the AEAD of
// a QUIC connection is crypto/tls's, so its speed is what the order of the suites must follow.
var hasAESGCMHardwareSupport = cpu.X86.HasAES && cpu.X86.HasPCLMULQDQ && cpu.X86.HasSSE41 && cpu.X86.HasSSSE3 ||
	cpu.ARM64.HasAES && cpu.ARM64.HasPMULL ||
	cpu.S390X.HasAES && cpu.S390X.HasAESCTR && cpu.S390X.HasGHASH ||
	runtime.GOARCH == "ppc64" || runtime.GOARCH == "ppc64le"

// chromeCipherSuites are the TLS 1.3 suites Chrome offers, in its order for a CPU with or without
// AES hardware: BoringSSL puts ChaCha20-Poly1305 first when the CPU has none
// (ssl_write_client_cipher_list, EVP_has_aes_hardware), as Chrome does on a phone or a router
// without the AES instructions. A server takes the first suite of the client it supports, so the
// order decides the cipher: AES-GCM in software costs such a CPU about twice what ChaCha20 does.
func chromeCipherSuites(aesHardware bool) []uint16 {
	if aesHardware {
		return []uint16{utls.TLS_AES_128_GCM_SHA256, utls.TLS_AES_256_GCM_SHA384, utls.TLS_CHACHA20_POLY1305_SHA256}
	}
	return []uint16{utls.TLS_CHACHA20_POLY1305_SHA256, utls.TLS_AES_128_GCM_SHA256, utls.TLS_AES_256_GCM_SHA384}
}
