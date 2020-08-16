package cs

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thinkgos/go-socks5"
)

var key = `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAwZquZqQbc6TaZyaa0UV5XRqDe7FY6BNhk7FxFMvwPyQ0jSj9
T3dfmBNkLEbdKwOEk3frMG5o0zl5ZbXj+B+24KQ5v0PBVjLHyJpzd8bpkq3W/eAh
WIKihY7Xsxr2sES7j9WTt+6KIXbMEx2IIKDaONcVCXH51hIhp1qqZwBtVIykdUn3
LwzDibGqp4RKABDy9CxY3x8alPPYbT0aBf4f60U7YPlI1/k7QPkXg+DLlog+utsn
eovCe33VFT5IOszKVxPUFGqzxqbvNMgUFc5eron2SCHKUjaryve0jdUY1jniNupb
B4902aw+hRGero6FsfZkBsNiV2SUgeG/+5oR1QIDAQABAoIBAQCx2ZLUn3TIa2xm
zcPy8stmh/C5NFXj+8nrj1m+LQpqNqw/8KOi2JpsbYPcWMzbssObZNIdD5AkWev3
T3w4d4ncG4Eg/vEgak21Lo1cPtJa+G9DkR2Q3ZDG+E2WLvLnQny6yQyGLw+dZjBa
bwqaTqmpBYxBvP4xdT6NKnDXZkEJJQBG8mO5bRM6oZZpp9LidtodlU4daoxIzvbf
lEPUZkuKOsLkeOiM2icXuU9SSZEExOs/ig5tgLEHdHmVhnpvAQr75ukO/ImZOyw1
Ne7AbC6XkiRJpoh2Oe63o04fORBI/O5HeJNvJXPtuxoar7WIVbZMqqhkurjqtz5l
cjKp/zsBAoGBAOvytWUCApxSRoGifKjeCpjpMAfbXpFDIF7iR7zUNx3Zy7vfTTm4
FzSbaT87YZpJ8GqssYimmQKRAI33fQUM0bDxiKSZkakSerWELTgAThr8BN1e6hfc
ONhVkDKAVlBYc3ksXN1FrmfKSi/YnBAwEtWKKeYNN75svKwN1RsS4TsRAoGBANIO
vSSJqphCKio/XFBqj2Ozu5UFe/MVVC6XT38SvoVbdAiRJqeIgoErn7N+qo0RpjNj
TaMDk6R6564/0sgdR8iZxQ/9Cy5ujWQF8jedk4XLc6WXi9BXmoHlAfmLyhD9wujc
ZUUefQsBZ+i4J1CmVovu/DbhZYzue3EzkP1NnEKFAoGASI8ZDXjyyJPcrt0DLQMr
ix6a8K+bg1x7RfKcUQuJ75octyfSnd9o83qfgRyHxWTblFKLPhTNlSZ2XzIutjDd
A2cjuEqpqq7OIagGJ+SgIFhEPreDkdbdfFnDwGQLJyYsTKVB4aIeIjjpW5FnXOsL
v7N/cwm5jMvvsZGHaY4CyaECgYAlUwMew+txJIiTezCvBVA3Og+Buji9B7QulypD
/ROnZImooAoLSMFPrG2zGjW53UH37ZQ0/AS2/DPAjYypjDJeHZyba64Z8QDknf3d
Df3Rj0YcTWJFgdtta0C/k6wy+rQwZkEEWBeF5hkNi/NIbFYChVOBeOlvckyy36PK
roiudQKBgQDma8xa1OhcbhXQGL+UVY30BKihabjAN2OAN4Ukx+9kKgzoGQPPSTFa
in10BwKpf9b95yqqViF6VKb+NSOBe2Kdyxx5PWnGKkGNSdGoan+urh7m4NJSbkAi
rFVx8YeFEzQM36IsGYUwKWVoB9EhN5ig+q0Ac4MndnhNDs1ktq8hrg==
-----END RSA PRIVATE KEY-----`

var crt = `-----BEGIN CERTIFICATE-----
MIIDTzCCAjegAwIBAgIBATANBgkqhkiG9w0BAQsFADBJMQswCQYDVQQGEwJaTTES
MBAGA1UEChMJdjJ1YjFqLmNnMRIwEAYDVQQLEwl2MnViMWouY2cxEjAQBgNVBAMT
CXYydWIxai5jZzAeFw0yMDAzMzExMzQ1MDdaFw0zMDAzMjkxNDQ1MDdaMEkxCzAJ
BgNVBAYTAlpNMRIwEAYDVQQKEwl2MnViMWouY2cxEjAQBgNVBAsTCXYydWIxai5j
ZzESMBAGA1UEAxMJdjJ1YjFqLmNnMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIB
CgKCAQEAwZquZqQbc6TaZyaa0UV5XRqDe7FY6BNhk7FxFMvwPyQ0jSj9T3dfmBNk
LEbdKwOEk3frMG5o0zl5ZbXj+B+24KQ5v0PBVjLHyJpzd8bpkq3W/eAhWIKihY7X
sxr2sES7j9WTt+6KIXbMEx2IIKDaONcVCXH51hIhp1qqZwBtVIykdUn3LwzDibGq
p4RKABDy9CxY3x8alPPYbT0aBf4f60U7YPlI1/k7QPkXg+DLlog+utsneovCe33V
FT5IOszKVxPUFGqzxqbvNMgUFc5eron2SCHKUjaryve0jdUY1jniNupbB4902aw+
hRGero6FsfZkBsNiV2SUgeG/+5oR1QIDAQABo0IwQDAOBgNVHQ8BAf8EBAMCAqQw
HQYDVR0lBBYwFAYIKwYBBQUHAwEGCCsGAQUFBwMCMA8GA1UdEwEB/wQFMAMBAf8w
DQYJKoZIhvcNAQELBQADggEBAHl5zBhdfN2oxUsxjlmfaOLfRIDa6wEAyeWqasr0
BW1ZP+ehhpvQMxG9xXjTlbBHnj34W7fTkzvrj9xc4mU61tMugifbIWnzXIPWrVeu
xTQivw6iVmYckUBhoI6WiHuYv+NOi2h72uWLmv0JDfG1NFddFBccOIzQd4iTO+zi
ufrg3IrbJx+7EnA87vXGdZVItgz92HoQF3HPfeXzzSFMjNmxEJKNP1IU7VmlPSUv
0F9sF0wukMiOGUQ0tXeYv3ArHqEfwtF5H9OH5RCuspFFMx6qPyAc1Ccjs73GLJ8I
TL44tBTU3E0Bl+fyBSRkAXbVVTcYsxTeHsSuYm3pARTpKsw=
-----END CERTIFICATE-----`

func TestTcpTls_Forward_Direct(t *testing.T) {
	for _, single := range []bool{true, false} {
		func() {
			// server
			srv := &TCPTlsServer{
				Addr:   "127.0.0.1:0",
				CaCert: nil,
				Cert:   []byte(crt),
				Key:    []byte(key),
				Single: single,
				Status: make(chan error, 1),
				Handler: HandlerFunc(func(inconn net.Conn) {
					buf := make([]byte, 20)
					n, err := inconn.Read(buf)
					if !assert.NoError(t, err) {
						return
					}
					assert.Equal(t, "ping", string(buf[:n]))
					_, err = inconn.Write([]byte("pong"))
					if !assert.NoError(t, err) {
						return
					}
				}),
			}
			go func() { _ = srv.ListenAndServe() }()
			require.NoError(t, <-srv.Status)
			defer srv.Close()

			// client
			d := &TCPTlsDialer{
				CaCert: []byte(crt),
				Cert:   []byte(crt),
				Key:    []byte(key),
				Single: single,
			}
			if !single {
				d.CaCert = nil
			}

			cli, err := d.DialTimeout(srv.LocalAddr(), 5*time.Second)
			require.NoError(t, err)
			defer cli.Close()

			_, err = cli.Write([]byte("ping"))
			require.NoError(t, err)
			b := make([]byte, 20)
			n, err := cli.Read(b)
			require.NoError(t, err)
			require.Equal(t, "pong", string(b[:n]))
		}()
	}
}

func TestJumper_socks5_tls(t *testing.T) {
	for _, single := range []bool{true, false} {
		func() {
			// server
			srv := &TCPTlsServer{
				Addr:   "127.0.0.1:0",
				CaCert: nil,
				Cert:   []byte(crt),
				Key:    []byte(key),
				Single: single,
				Status: make(chan error, 1),
				Handler: HandlerFunc(func(inconn net.Conn) {
					buf := make([]byte, 20)
					n, err := inconn.Read(buf)
					if !assert.NoError(t, err) {
						return
					}
					assert.Equal(t, "ping", string(buf[:n]))
					_, err = inconn.Write([]byte("pong"))
					if !assert.NoError(t, err) {
						return
					}
				}),
			}
			go func() { _ = srv.ListenAndServe() }()
			require.NoError(t, <-srv.Status)
			defer srv.Close()

			// start socks5 proxy server
			cator := &socks5.UserPassAuthenticator{Credentials: socks5.StaticCredentials{"user": "password"}}
			proxySrv := socks5.NewServer(
				socks5.WithAuthMethods(
					[]socks5.Authenticator{
						new(socks5.NoAuthAuthenticator),
						cator,
					}),
			)
			proxyLn, err := net.Listen("tcp", "127.0.0.1:0")
			require.NoError(t, err)
			defer proxyLn.Close() // nolint: errcheck

			go func() {
				proxySrv.Serve(proxyLn) // nolint: errcheck
			}()

			time.Sleep(time.Millisecond * 100)
			// t.Logf("proxy server address: %v", proxyLn.Addr().String())

			// start jumper to socks5
			proxyURL := "socks5://" + "user:password@" + proxyLn.Addr().String()
			pURL, err := ParseProxyURL(proxyURL)
			require.NoError(t, err)
			// t.Logf("socks5 proxy url: %v", proxyURL)

			jumptcp := &TCPTlsDialer{
				CaCert:  []byte(crt),
				Cert:    []byte(crt),
				Key:     []byte(key),
				Single:  single,
				Forward: Socks5{pURL},
			}
			conn, err := jumptcp.DialTimeout(srv.LocalAddr(), time.Second)
			require.NoError(t, err)
			defer conn.Close() // nolint: errcheck
			_, err = conn.Write([]byte("ping"))
			require.NoError(t, err)
			b := make([]byte, 20)
			n, err := conn.Read(b)
			require.NoError(t, err)
			require.Equal(t, "pong", string(b[:n]))
		}()
	}
}
