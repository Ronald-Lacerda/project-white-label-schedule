package shared

import (
	"math/rand"
	"sync"
	"time"

	"github.com/oklog/ulid/v2"
)

var (
	entropyMu sync.Mutex
	entropy   = rand.New(rand.NewSource(time.Now().UnixNano()))
)

// NewID gera um ULID único ordenável por tempo.
func NewID() string {
	entropyMu.Lock()
	id := ulid.MustNew(ulid.Timestamp(time.Now()), entropy)
	entropyMu.Unlock()
	return id.String()
}
