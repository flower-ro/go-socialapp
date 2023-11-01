package idgenerate

import (
	"fmt"
	"math/rand"
	"sync/atomic"
	"time"
)

type IdGenerator struct {
	id int64
}

func (i *IdGenerator) Reset() {
	atomic.StoreInt64(&i.id, 0)
}

func (i *IdGenerator) GenID() int64 {
	return atomic.AddInt64(&i.id, 1)
}

func (i *IdGenerator) GetCurrentID() int64 {
	return i.id
}

func CreateGuid() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x%x", uint32(time.Now().Unix()), b[4:8])
}
