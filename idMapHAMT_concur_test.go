package nodeID

// xlNodeID_go/idMap_concur_test.go

import (
	"bytes"
	"fmt"
	. "gopkg.in/check.v1"
	"time"
)

var _ = fmt.Print

// Gets run as the Kth goroutine out of M.  Each goroutine is responsible
// for inserting N/M values into the map.  When the values have been
// inserted, signals 'done'.  N must not be prime relative to M; K must
// be in the range [0,M).
//
func (s *XLSuite) doHamtOneWay(c *C, keys [][]byte, N, M, K int,
	m *IDMapHAMT, done chan bool) {

	quota := N / M
	start := K * quota
	end := (K + 1) * quota

	for i := start; i < end; i++ {
		_ = m.Insert(keys[i], &keys[i])
	}
	done <- true

}

// Given N values, create M goroutines, each inserting N/M values
func (s *XLSuite) doHamtMWayTest(c *C, keys [][]byte, N, M int) {
	m, err := NewIDMapHAMT(5, 18) // XXX WAS 5,5
	c.Assert(err, IsNil)

	chans := make([]chan bool, M)
	for k := 0; k < M; k++ {
		chans[k] = make(chan bool)
	}
	for k := 0; k < M; k++ {
		go s.doHamtOneWay(c, keys, N, M, k, m, chans[k])
	}
	time.Sleep(time.Millisecond)

	for k := 0; k < M; k++ {
		<-chans[k]
	}
	// verify that the keys are present in the map
	for i := 0; i < N; i++ {
		value, err := m.Find(keys[i])
		c.Assert(err, IsNil)

		val := value.(*[]byte) // XXX sometimes fails, gets nil
		c.Assert(bytes.Equal(*val, keys[i]), Equals, true)

	}
	// BEGIN STATS CHECKS -------------------------------------------
	// entryCount, mapCount, deepest := m.Size()
	entryCount, _, _ := m.Size()
	c.Assert(entryCount, Equals, uint(N))

	// DEBUG XXX These comments apply to idMap test
	// Quite inefficient: for 256K entries, 61832 maps; for 1M, 97486 maps.
	// Max depth seen: 5.
	//
	//fmt.Printf("entryCount %7d, mapCount %5d, deepest %2d\n",
	//	entryCount, mapCount, deepest)
	// END
}

// XXX in effect, comment this test out
func (s *XLSuite) TestConcurrentHamtInsertions(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_CONCURRENT_HAMT_INSERTIONS")
	}

	// build an array of N random-ish K-byte keys
	K := 32
	N := 1024 * 1024
	t0 := time.Now()
	keys := makeSomeKeys(N, K)
	t1 := time.Now()
	deltaT := t1.Sub(t0)
	fmt.Printf("setup time for %d  %d-byte keys: %v\n", N, K, deltaT)

	t0 = time.Now()
	s.doHamtMWayTest(c, keys, N, 4)  // do 4-fold test
	s.doHamtMWayTest(c, keys, N, 8)  // do 8-fold test
	s.doHamtMWayTest(c, keys, N, 16) // do 16-fold test
	s.doHamtMWayTest(c, keys, N, 32) // do 32-fold test
	s.doHamtMWayTest(c, keys, N, 64) // do 64-fold test
	t1 = time.Now()
	deltaT = t1.Sub(t0)
	fmt.Printf("run time for %d concurrent insertion tests: %v\n", 5, deltaT)
	// typically about 6.7 sec
}
