package nodeID

import (
	"fmt"
	xr "github.com/jddixon/rnglib_go"
	xu "github.com/jddixon/xlUtil_go"
	. "gopkg.in/check.v1"
)

func (s *XLSuite) TestBadNodeIDs(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_BAD_NODE_IDS")
	}
	c.Assert(false, Equals, IsValidID(nil))
	candidate := make([]byte, xu.SHA1_BIN_LEN-1)
	c.Assert(false, Equals, IsValidID(candidate))
	candidate = make([]byte, xu.SHA1_BIN_LEN)
	c.Assert(true, Equals, IsValidID(candidate))
	candidate = make([]byte, xu.SHA1_BIN_LEN+1)
	c.Assert(false, Equals, IsValidID(candidate))
	candidate = make([]byte, xu.SHA2_BIN_LEN)
	c.Assert(true, Equals, IsValidID(candidate))
}
func (s *XLSuite) TestThisAndThat(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_THIS_AND_THAT")
	}
	var err error
	rng := xr.MakeSimpleRNG()
	v1 := make([]byte, xu.SHA1_BIN_LEN)
	rng.NextBytes(v1)
	v2 := make([]byte, xu.SHA1_BIN_LEN)
	rng.NextBytes(v2)
	id1, err := NewNodeID(v1)
	c.Assert(err, Equals, nil)
	id2, err := NewNodeID(v2)
	c.Assert(err, Equals, nil)
	c.Assert(id1, Not(Equals), id2)

	v1a := id1.Value()
	v2a := id2.Value()

	// XXX test framework cannot handle these comparisons
	// c.Assert(v1, Not(DeepEquals), v1a)		// 'Deep' is for desperation
	// c.Assert(v2, Not(Equals), v2a)

	// XXX not sure that test framework results are meaningful
	c.Assert(&v1, Not(Equals), &v1a)
	c.Assert(&v2, Not(Equals), &v2a)

	c.Assert(xu.SHA1_BIN_LEN, Equals, len(v1a))
	c.Assert(xu.SHA1_BIN_LEN, Equals, len(v2a))
	for i := 0; i < xu.SHA1_BIN_LEN; i++ {
		c.Assert(v1[i], Equals, v1a[i])
		c.Assert(v2[i], Equals, v2a[i])
	}
	c.Assert(false, Equals, id1.Equal(nil))
	c.Assert(true, Equals, id1.Equal(id1))
	c.Assert(false, Equals, id1.Equal(id2)) // FOO
}

func (s *XLSuite) TestComparator(c *C) {
	var err error
	if VERBOSITY > 0 {
		fmt.Println("TEST_COMPARATOR")
	}
	rng := xr.MakeSimpleRNG()
	v1 := make([]byte, xu.SHA1_BIN_LEN)
	rng.NextBytes(v1)
	v3 := make([]byte, xu.SHA2_BIN_LEN)
	rng.NextBytes(v3)
	id1, err := NewNodeID(v1) // SHA1
	c.Assert(err, Equals, nil)
	id3, err := NewNodeID(v3) // SHA3
	c.Assert(err, Equals, nil)

	// make clones which will sort before and after v1
	v1a := make([]byte, xu.SHA1_BIN_LEN) // sorts AFTER v1
	for i := 0; i < xu.SHA1_BIN_LEN; i++ {
		if v1[i] == 255 {
			v1a[i] = 255
		} else {
			v1a[i] = v1[i] + 1
			break
		}
	}
	v1b := make([]byte, xu.SHA1_BIN_LEN) // sorts BEFORE v1
	for i := 0; i < xu.SHA1_BIN_LEN; i++ {
		if v1[i] == 0 {
			v1b[i] = 0
		} else {
			v1b[i] = v1[i] - 1
			break
		}
	}
	id1a, err := NewNodeID(v1a)
	c.Assert(err, Equals, nil)
	id1b, err := NewNodeID(v1b)
	c.Assert(err, Equals, nil)

	result, err := id1.Compare(id1) // self
	c.Assert(0, Equals, result)
	c.Assert(err, IsNil)

	klone, err := id1.Clone() // identical copy
	c.Assert(err, IsNil)
	result, err = id1.Compare(klone)
	c.Assert(0, Equals, result)
	c.Assert(err, IsNil)

	result, err = id1.Compare(nil) // nil comparand
	c.Assert(err, Not(IsNil))      // NOT

	result, err = id1.Compare(id3)
	c.Assert(err, Not(IsNil)) // different lengths	// NOT

	result, err = id1.Compare(id1a)
	c.Assert(-1, Equals, result) // id1 < id1a
	c.Assert(err, IsNil)

	result, err = id1.Compare(id1b) // id1 > id1b
	c.Assert(1, Equals, result)
	c.Assert(err, IsNil)

	result, err = id1a.Compare(id1b) // id1a > id1b
	c.Assert(1, Equals, result)
	c.Assert(err, IsNil)
}

func (s *XLSuite) makeANodeID(c *C, rng *xr.PRNG) (id *NodeID) {
	var length int
	if rng.NextBoolean() {
		length = xu.SHA1_BIN_LEN
	} else {
		length = xu.SHA2_BIN_LEN
	}
	data := make([]byte, length)
	rng.NextBytes(data)
	id, err := New(data)
	c.Assert(err, IsNil)
	return id
}
func (s *XLSuite) TestSameNodeID(c *C) {
	rng := xr.MakeSimpleRNG()
	id1 := s.makeANodeID(c, rng)
	c.Assert(SameNodeID(id1, id1), Equals, true)

	id2 := s.makeANodeID(c, rng)
	c.Assert(SameNodeID(id1, id2), Equals, false)
}
