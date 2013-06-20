package xlattice_go

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"errors"
	"hash"
)

/**
 * A Node is uniquely identified by a NodeID and can satisfy an
 * identity test constructed using its public key.  That is, it
 * can prove that it holds the private key materials corresponding
 * to the public key.
 *
 * @author Jim Dixon
 */
type Node struct {
	nodeID      *NodeID         // public
	key         *rsa.PrivateKey // private
	pubkey      *rsa.PublicKey  // public
	endPoints   []*EndPoint
	overlays    []*OverlayI
	peers       []*Peer
	connections []*ConnectionI
}

func NewNewNode(id *NodeID) (*Node, error) {
	// XXX create default 2K bit RSA key
	return NewNode(id, nil, nil, nil, nil)
}

func NewNode(id *NodeID, key *rsa.PrivateKey, e *[]*EndPoint, p *[]*Peer,
	c *[]*ConnectionI) (*Node, error) {

	if id == nil {
		err := errors.New("IllegalArgument: nil NodeID")
		return nil, err
	}
	nodeID := (*id).Clone() // we use this copy, not the nodeID passed

	///////////////////////////////
	// XXX STUB: switch on key type
	// * extract the public key
	// * build the DigSigner
	///////////////////////////////
	if key == nil {
		k, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			return nil, err
		}
		key = k
	}

	var endPoints []*EndPoint // an empty slice
	var overlays []*OverlayI
	if e != nil {
		count := len(*e)
		for i := 0; i < count; i++ {
			endPoints = append(endPoints, (*e)[i])
			// XXX get the overlay from the endPoint
			// overlays = append(overlays, (*o)[i])
		}
	}
	var peers []*Peer // an empty slice
	if p != nil {
		count := len(*p)
		for i := 0; i < count; i++ {
			peers = append(peers, (*p)[i])
		}
	}
	var cnxs []*ConnectionI // another empty slice
	if c != nil {
		count := len(*c)
		for i := 0; i < count; i++ {
			cnxs = append(cnxs, (*c)[i])
		}
	}
	q := new(Node)
	(*q).nodeID = nodeID // the clone
	(*q).key = key
	(*q).pubkey = &(*key).PublicKey
	(*q).endPoints = endPoints
	(*q).overlays = overlays
	(*q).peers = peers
	(*q).connections = cnxs
	return q, nil
}

// IDENTITY /////////////////////////////////////////////////////////
func (n *Node) GetNodeID() *NodeID {
	return n.nodeID
}
func (n *Node) GetPublicKey() *rsa.PublicKey {
	// XXX This should be a copy or a serialization
	return n.pubkey
}

// Returns an instance of a DigSigner which can be run in a separate
// goroutine.  This allows the Node to calculate more than one
// digital signature at the same time.
//
// XXX would prefer that *DigSigner be returned
func (n *Node) getSigner() *signer {
	return newSigner(n.key)
}

// OVERLAYS /////////////////////////////////////////////////////////
func (n *Node) addOverlay(o *OverlayI) error {
	if o == nil {
		return errors.New("IllegalArgument: nil Overlay")
	}
	n.overlays = append(n.overlays, o)
	return nil
}

/**
 * @return a count of the number of overlays the peer can be
 *         accessed through
 */
func (n *Node) SizeOverlays() int {
	return len(n.overlays)
}

/** @return how to access the peer (transport, protocol, address) */
func (n *Node) GetOverlay(x int) *OverlayI {
	return n.overlays[x]
}

// PEERS ////////////////////////////////////////////////////////////
func (n *Node) addPeer(o *Peer) error {
	if o == nil {
		return errors.New("IllegalArgument: nil Peer")
	}
	n.peers = append(n.peers, o)
	return nil
}

/**
 * @return a count,  the number of peers
 */
func (n *Node) SizePeers() int {
	return len(n.peers)
}
func (n *Node) GetPeer(x int) *Peer {
	return n.peers[x]
}

// CONNECTORS ///////////////////////////////////////////////////////
func (n *Node) addConnectionI(c *ConnectionI) error {
	if c == nil {
		return errors.New("IllegalArgument: nil ConnectionI")
	}
	n.connections = append(n.connections, c)
	return nil
}

/** @return a count of known Connections for this Peer */
func (n *Node) SizeConnections() int {
	return len(n.connections)
}

/**
 * Return a ConnectionI, an Address-Protocol pair identifying
 * an Acceptor for the Peer.  Connections are arranged in order
 * of preference, with the zero-th ConnectionI being the most
 * preferred.
 *
 * XXX Could as easily return an EndPoint.
 *
 * @return the Nth Connection
 */
func (n *Node) GetConnection(x int) *ConnectionI {
	return n.connections[x]
}

// EQUAL ////////////////////////////////////////////////////////////
func (n *Node) Equal(any interface{}) bool {
	if any == n {
		return true
	}
	if any == nil {
		return false
	}
	switch v := any.(type) {
	case *Node:
		_ = v
	default:
		return false
	}
	other := any.(*Node) // type assertion
	// THINK ABOUT publicKey.equals(any.publicKey)
	if n.nodeID == other.nodeID {
		return true
	}
	if n.nodeID.Length() != other.nodeID.Length() {
		return false
	}
	myVal := n.nodeID.Value()
	otherVal := other.nodeID.Value()
	for i := 0; i < n.nodeID.Length(); i++ {
		if myVal[i] != otherVal[i] {
			return false
		}
	}
	return false
}
func (n *Node) String() string {
	return "NOT IMPLEMENTED"
}

// DIG SIGNER ///////////////////////////////////////////////////////

type signer struct {
	key    *rsa.PrivateKey
	digest hash.Hash
}

func newSigner(key *rsa.PrivateKey) *signer {
	// XXX some validation, please
	h := sha1.New()
	ds := signer{key: key, digest: h}
	return &ds
}
func (s *signer) Algorithm() string {
	return "SHA1+RSA" // XXX NOT THE PROPER NAME
}
func (s *signer) Length() int {
	return 42 // XXX NOT THE PROPER VALUE
}
func (s *signer) Update(chunk []byte) {
	s.digest.Write(chunk)
}

// XXX NOTE CHANGE IN INTERFACE
func (s *signer) Sign() ([]byte, error) {
	h := s.digest.Sum(nil)
	sig, err := rsa.SignPKCS1v15(rand.Reader, s.key, crypto.SHA1, h)
	return sig, err
}

func (s *signer) String() string {
	return "NOT IMPLEMENTED"
}
