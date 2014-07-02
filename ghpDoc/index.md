# xlNodeID_go

NodeID library for the xlattice_go project.

In Go, a NodeID is a slice of bytes used to uniquely identify an entity.  The 
length of the slice must be either 20 or 32 bytes.  The first is the length
of an SHA-1 hash and so can be used, for example, to store the SHA-1 hash of a
document.  In the XLattice project, 20 byte/160 bit values are always used 
to store such hashes.

SHA-1 is the Secure Hash Algorithm designed by the US National Security Agency.
It is a Federal Information Processing Standard published by the National 
Institute of Science and Technology, 
[FIPS180-4](http://csrc.nist.gov/publications/fips/fips180-4/fips-180-4.pdf)

In recent years many attempts have been made is crack SHA-1, that is to identify
two documents whose SHA-1 hash is the same.  At this time no such hash collision 
has been found.  Nevertheless SHA-1 has been declared obsolete (ie, 
cryptographically insecure) and so a competition was held to settle on a 
new hash algorithm, SHA-3.  The competition was won by 
[Keccak](http://noekeon.org/Keccak-implementation-3.2.pdf).  SHA-3 is
available in a number of varians.  We use Keccak-256, the 32-byte/256-bit
variant.  When XLattice documents refer to SHA-3, we mean Keccak-256.
