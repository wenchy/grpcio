// xtime represents server time with offset
package xtime

import (
	"time"
)

type offset struct {
	Seconds int64
	Nanos   int
}

var myoffset *offset

func init() {
	myoffset = new(offset)
}

// Now constructs a new time.Time from the current time and offset.
func Now() time.Time {
	return myoffset.asTime()
}

// Set the offset.
func Set(seconds int64, nanos int) {
	myoffset.Seconds = seconds
	myoffset.Nanos = nanos
}

// Reset the offset.
func Reset() {
	myoffset.Seconds = 0
	myoffset.Nanos = 0
}

// asTime converts x to a time.Time.
func (x *offset) asTime() time.Time {
	now := time.Now()
	return time.Unix(now.Unix()+x.Seconds, int64(now.Nanosecond()+x.Nanos)).UTC()
}
