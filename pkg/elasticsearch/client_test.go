package elasticsearch

import (
	"crypto/tls"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCustomCA(t *testing.T) {
	customCA := `-----BEGIN CERTIFICATE-----
MIIC2jCCAcKgAwIBAgIBATANBgkqhkiG9w0BAQsFADAeMRwwGgYDVQQDExNsb2dn
aW5nLXNpZ25lci10ZXN0MB4XDTE5MDUyMzE0NDkxM1oXDTI0MDUyMTE0NDkxNFow
HjEcMBoGA1UEAxMTbG9nZ2luZy1zaWduZXItdGVzdDCCASIwDQYJKoZIhvcNAQEB
BQADggEPADCCAQoCggEBAKjkJ6SM7sGcclLsDO34xr6ybhkw3xf2deCVFj/Gr2qW
TbjpYAULYNPRaaorkdFm7VNvLtpxlQtILJKx40Mi0aBfGw0imPLG03pFf0wSxa0F
RWZP4NBZ9vtt0+xir77d7BhI4oTHQYzY9JxqtvaOFCFlKpF96vO+DXap4s94/tB2
mY19a1SKePTCd514dAhi0LhJ4p9zk5t8AKH77kAVmZByf0L2sVYTSXPalbSbPlEm
mvTaOV8vkll+Iri3wOose87CR0squv5DX6so3Hznso6MPpq245I9yucyjRRTAquA
3g7z2N7pagF30U5FuEgTOVx8lyvJ/UFwqmPi//KJ6qMCAwEAAaMjMCEwDgYDVR0P
AQH/BAQDAgKkMA8GA1UdEwEB/wQFMAMBAf8wDQYJKoZIhvcNAQELBQADggEBAGEB
LMprANS8xWLLI7C0UzAJsXIhC8/12aT5vOw3LU/4+Pn3Zoy+0Avm+DlM1YsLZiig
Hv8rBHh7TuDhhPeeK/GWHAoM0I8kQVgClOREga0/8etfpMcq/e6AHVDxeXLE+L7o
wSIwNDVOwetlYVHFka0PguEM76Ms70Gia7RQU5GtUsIeB9zrwG7dkz4wWlS2eOQ3
MIL3TlNQxJGDYc34ourZ/QfRXnnzkw7tKH1Cp1SkmEKkfSpqXf9cfq4x2eJHfBw7
7cD3KHT422uXoSE8Tff/B76PNzpzlhE3UNOcLOdBSA1Ur71x8ryi0waxesQVonlr
18FP8YQpqqQy9QdtPno=
-----END CERTIFICATE-----`

	tlsConfig := &tls.Config{}

	err := addRootCA([]byte(customCA), tlsConfig)
	assert.NoError(t, err)
	assert.NotNil(t, tlsConfig.RootCAs)
}
