package captain

import (
	"net"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDatagram(t *testing.T) {
	t.Run("invalid address", func(t *testing.T) {
		_, err := NewDatagram("localhost", nil)
		require.Error(t, err)
	})
	t.Run("domain host name to long", func(t *testing.T) {
		_, err := NewDatagram("localhostlocalhostlocalhostlocalhostlocalhostlocalhostlocalhostlocalhostlocalhostlocalhost"+
			"localhostlocalhostlocalhostlocalhostlocalhostlocalhostlocalhostlocalhostlocalhostlocalhost"+
			"localhostlocalhostlocalhostlocalhostlocalhostlocalhostlocalhostlocalhostlocalhostlocalhost:8080", nil)
		require.Error(t, err)
	})

	t.Run("domain", func(t *testing.T) {
		datagram, err := NewDatagram("localhost:8080", []byte{1, 2, 3})
		require.NoError(t, err)
		require.Equal(t, Datagram{
			0, AddrSpec{
				FQDN:     "localhost",
				Port:     8080,
				AddrType: ATYPDomain,
			},
			[]byte{1, 2, 3},
		}, datagram)
		require.Equal(t, []byte{0, ATYPDomain, 9, 'l', 'o', 'c', 'a', 'l', 'h', 'o', 's', 't', 0x1f, 0x90}, datagram.Header())
		require.Equal(t, []byte{0, ATYPDomain, 9, 'l', 'o', 'c', 'a', 'l', 'h', 'o', 's', 't', 0x1f, 0x90, 1, 2, 3}, datagram.Bytes())

	})
	t.Run("ipv4", func(t *testing.T) {
		datagram, err := NewDatagram("127.0.0.1:8080", []byte{1, 2, 3})
		require.NoError(t, err)
		require.Equal(t, Datagram{
			0, AddrSpec{
				IP:       net.IPv4(127, 0, 0, 1),
				Port:     8080,
				AddrType: ATYPIPv4,
			},
			[]byte{1, 2, 3},
		}, datagram)
		require.Equal(t, []byte{0, ATYPIPv4, 127, 0, 0, 1, 0x1f, 0x90}, datagram.Header())
		require.Equal(t, []byte{0, ATYPIPv4, 127, 0, 0, 1, 0x1f, 0x90, 1, 2, 3}, datagram.Bytes())
	})
	t.Run("ipv6", func(t *testing.T) {
		datagram, err := NewDatagram("[::1]:8080", []byte{1, 2, 3})
		require.NoError(t, err)
		require.Equal(t, Datagram{
			0, AddrSpec{
				IP:       net.IPv6loopback,
				Port:     8080,
				AddrType: ATYPIPv6,
			},
			[]byte{1, 2, 3},
		}, datagram)
		require.Equal(t, []byte{0, ATYPIPv6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0x1f, 0x90}, datagram.Header())
		require.Equal(t, []byte{0, ATYPIPv6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0x1f, 0x90, 1, 2, 3}, datagram.Bytes())

	})
}

func TestParseDatagram(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		args    args
		wantDa  Datagram
		wantErr bool
	}{
		{
			"IPv4",
			args{[]byte{0, ATYPIPv4, 127, 0, 0, 1, 0x1f, 0x90, 1, 2, 3}},
			Datagram{
				0, AddrSpec{
					IP:       net.IPv4(127, 0, 0, 1),
					Port:     8080,
					AddrType: ATYPIPv4,
				},
				[]byte{1, 2, 3},
			},
			false,
		},
		{
			"IPv6",
			args{[]byte{0, ATYPIPv6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0x1f, 0x90, 1, 2, 3}},
			Datagram{
				0, AddrSpec{
					IP:       net.IPv6loopback,
					Port:     8080,
					AddrType: ATYPIPv6,
				},
				[]byte{1, 2, 3},
			},
			false,
		},
		{
			"FQDN",
			args{[]byte{0, ATYPDomain, 9, 'l', 'o', 'c', 'a', 'l', 'h', 'o', 's', 't', 0x1f, 0x90, 1, 2, 3}},
			Datagram{
				0, AddrSpec{
					FQDN:     "localhost",
					Port:     8080,
					AddrType: ATYPDomain,
				},
				[]byte{1, 2, 3},
			},
			false,
		},
		{
			"invalid address type",
			args{[]byte{0, 0x02, 127, 0, 0, 1, 0x1f, 0x90}},
			Datagram{},
			true,
		},
		{
			"less min length",
			args{[]byte{0, ATYPIPv4, 127, 0, 0, 1, 0x1f}},
			Datagram{},
			true,
		},
		{
			"less domain length",
			args{[]byte{0, ATYPDomain, 10, 127, 0, 0, 1, 0x1f, 0x09}},
			Datagram{},
			true,
		},
		{
			"less ipv6 length",
			args{[]byte{0, ATYPIPv6, 127, 0, 0, 1, 0x1f, 0x09}},
			Datagram{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDa, err := ParseDatagram(tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDatagram() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && !reflect.DeepEqual(gotDa, tt.wantDa) {
				t.Errorf("ParseDatagram() gotDa = %v, want %v", gotDa, tt.wantDa)
			}
		})
	}
}
