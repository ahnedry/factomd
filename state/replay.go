// Copyright 2015 FactomProject Authors. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package state

import (
	"fmt"
	"sync"
	"time"

	"github.com/FactomProject/factomd/common/interfaces"
	"github.com/FactomProject/factomd/common/primitives"
)

const numBuckets = 27

var _ = time.Now()
var _ = fmt.Print

type Replay struct {
	mutex    sync.Mutex
	buckets  [numBuckets]map[[32]byte]int
	basetime int // hours since 1970
	center   int // Hour of the current time.
}

// Remember that Unix time is in seconds since 1970.  This code
// wants to be handed time in seconds.
func hours(unix int64) int {
	return int(unix / 60 / 60)
}

// Returns false if the hash is too old, or is already a
// member of the set.  Timestamp is in seconds.
func (r *Replay) Valid(mask int, hash [32]byte, timestamp interfaces.ITimestamp, systemtime interfaces.ITimestamp) (index int, valid bool) {
	timeSeconds := timestamp.GetTimeSeconds()
	systemTimeSeconds := systemtime.GetTimeSeconds()
	// Check the timestamp to see if within 12 hours of the system time.  That not valid, we are
	// just done without any added concerns.
	if timeSeconds-systemTimeSeconds > 60*60*12 || systemTimeSeconds-timeSeconds > 60*60*12 {
		return -1, false
	}

	now := hours(systemTimeSeconds)

	r.mutex.Lock()
	defer r.mutex.Unlock()

	// We don't let the system clock go backwards.  likely an attack if it does.
	if now < r.center {
		now = r.center
	}

	if r.center == 0 {
		r.center = now
		r.basetime = now - (numBuckets / 2)
	}
	for r.center < now {
		copy(r.buckets[:], r.buckets[1:])
		r.buckets[numBuckets-1] = nil
		r.center++
		r.basetime++
	}

	t := hours(timeSeconds)
	index = t - r.basetime
	if index < 0 || index >= numBuckets {
		fmt.Println("dddd Timestamp false on time:", index)
		return 0, false
	}

	if r.buckets[index] == nil {
		r.buckets[index] = make(map[[32]byte]int)
	} else {
		v, _ := r.buckets[index][hash]
		if v&mask > 0 {
			return index, false
		}
	}
	return index, true
}

// Checks if the timestamp is valid.  If the timestamp is too old or
// too far into the future, then we don't consider it valid.  Or if we
// have seen this hash before, then it is not valid.  To that end,
// this code remembers hashes tested in the past, and rejects the
// second submission of the same hash.
func (r *Replay) IsTSValid(mask int, hash interfaces.IHash, timestamp interfaces.ITimestamp) bool {
	return r.IsTSValid_(mask, hash.Fixed(), timestamp, primitives.NewTimestampNow())
}

// To make the function testable, the logic accepts the current time
// as a parameter.  This way, the test code can manipulate the clock
// at will.
func (r *Replay) IsTSValid_(mask int, hash [32]byte, timestamp interfaces.ITimestamp, now interfaces.ITimestamp) bool {

	if index, ok := r.Valid(mask, hash, timestamp, now); ok {
		r.mutex.Lock()
		defer r.mutex.Unlock()
		// Mark this hash as seen
		r.buckets[index][hash] = r.buckets[index][hash] | mask
		return true
	}

	return false
}
