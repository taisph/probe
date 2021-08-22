package addrprobe

import (
	"bytes"
)

type NetworkAddress struct {
	Network string
	Address string
}

func FromString(address string) NetworkAddress {
	b := []byte(address)
	i := bytes.IndexByte(b, ':')

	if bytes.Equal(b[0:i], []byte("unix")) {
		return NetworkAddress{Network: "unix", Address: string(b[i+1:])}
	}

	l := bytes.LastIndexByte(b, ':')
	if i == l {
		return NetworkAddress{Network: "tcp", Address: address}
	}

	return NetworkAddress{Network: string(b[0:i]), Address: string(b[i+1:])}
}
