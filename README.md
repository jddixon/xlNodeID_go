xlNodeID_go
===========

The NodeID library for the xlattice_go project.

## NodeID

A nodeID is a slice of bytes used to uniquely identify an entity.  The 
length of the slice must be either 20 or 32 bytes.  

The first is the length of an SHA1 hash and so can be used, for example, 
to store the SHA1 hash of a document.  

**SHA1**, also known as SHA-1, is the Secure Hash Algorithm 
designed by the US National Security Agency.
It is a Federal Information Processing Standard published by the National 
Institute of Science and Technology, 
[FIPS180-4](http://csrc.nist.gov/publications/fips/fips180-4/fips-180-4.pdf)

In recent years many attempts have been made is crack SHA1, that is to 
identify two documents whose SHA1 hash is the same.  At this time no such
hash collision has been found.  Nevertheless SHA1 has been declared obsolete 
(ie, cryptographically insecure) and so there has been a general shift towards
more secure variants of the hash, one of which is SHA256, which produces a
256-bit/32-byte hash.  For simplicity XLattice documents generally refer
to this as **SHA2**.

Over the last several yerasa competition was held to settle on a 
new and still more secure hash algorithm, SHA3.  The competition was won by 
[Keccak](http://noekeon.org/Keccak-implementation-3.2.pdf).  

SHA3 is available in a number of variants.  We use Keccak-256, the 
32-byte/256-bit variant.  When XLattice documents refer to **SHA3**, we mean 
Keccak-256.  

At this time SHA3 is not quite stabilized and so XLattice prefers SHA256.
In other words, 32-byte XLattice nodeIDs are conventionally
SHA2 hashes.

## IDMaps

Because nodeIDs are used to uniquely identify objects within the XLattice
universe, it is convenient to have utilities which can store and rapidly
retrieve nodeIDs and pointers to the associated objects.

An idMap is a trie structure which stores keys and pointers to their associated
values.  A key is a byte slice (and assumed to be a nodeID).  A value can
be anything at all.

Each idMap has a maximum depth.  This is the length of the longest key can
can be stored in the idMap.  By default the depth is 32, the length of an
SHA2 or SHA3 hash.

When a new key is inserted into an idMap, the software searches down the 
existing tree.  At each level there is a 256-cell Map.  If at level N the
cell indexed by byte N is free, the new key is entered into the Map at 
that level.  If the cell is occupied by a key which differs in the next
byte, a new Map is created, with the new key and the existing key in their
respecitive cells.  If the keys have equal next bytes, insertion recurses
down the tree, until a difference is found or the maximum depth is reached.

Despite the use of coarse locks, idMap
Insert and Find operations compare favorbly with that found using the 
[HAMT](http://en.wikipedia.org/wiki/Hash_array_mapped_trie) algorighm as
implemented in [hamt_go](http://jddixon.github.io/hamt_go).  

## On-line Documentation

More information on the **xlNodeID_go** project can be found [here](https://jddixon.github.io/xlNodeID_go)

## On-line Documentation
More information on the **xlNodeID_go** project can be found 
[here](https://jddixon.github.io/xlNodeID_go)
