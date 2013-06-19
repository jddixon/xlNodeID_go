package xlattice_go

import (
	"strings"
)

/**
 * An EndPoint is specified by a transport and an Address, including
 * the local part.  If the transport is TCP/IP, for example, the
 * Address includes the IP address and the port number.
 */

type EndPoint struct {
	transport *TransportI
	address   *AddressI
}

func NewEndPoint(t *TransportI, a *AddressI) *EndPoint {
	// XXX need some checks
	return &EndPoint{t, a}
}

func (e *EndPoint) getAddress() *AddressI {
	return e.address
}

func (e *EndPoint) getTransport() *TransportI {
	return e.transport
}

// func (e *EndPoint) clone() *EndPoint {
//     NOT IMPLEMENTED
// }

func (e *EndPoint) String() string {
	// e.transport is a pointer to something that satisfies
	//   the Transport interface and similarly for e.address
	var parts = []string{(*e.transport).String(), " ", (*e.address).String()}
	return strings.Join(parts, "")
}