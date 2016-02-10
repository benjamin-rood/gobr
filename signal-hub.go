package gobr

import (
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/benjamin-rood/gobr"
)

// SignalHub is for simple non-blocking to signal receiver channels individually and broadcasting to all receivers
type SignalHub struct {
	rw        sync.RWMutex
	receivers map[string]chan struct{} //	requires the receiver to keep track of its key ("signature") as the registering party.
}

// NewSignalHub returns a pointer to a new instance, with the sync.RWMutex set to the default (Unlocked).
func NewSignalHub() *SignalHub {
	sh := SignalHub{}
	sh.receivers = make(map[string]chan struct{})
	return &sh
}

// Broadcast safely sends a signal to all receivers. By using invoking the method *as a goroutine* instead of sending channel directly, we avoid the chance of the sender blocking while waiting for the receiver to retreive the message sent on the channel.
// This also means that we can Broadcast "as-if" there are processes waiting to receive the synchronisation signal, even though in some cases there might not be.
func (sh *SignalHub) Broadcast() {
	defer sh.rw.Unlock()
	sh.rw.Lock()
	for signature, channel := range sh.receivers { // if map is empty then no signals get sent or lost.
		fmt.Println("broadcasting to", signature)
		channel <- struct{}{}
	}
}

// Signal 'safely' sends a signal message (struct{}{}) without blocking the sender *if invoked as a goroutine*.
// Otherwise you need to check for the error if you want to confirm receipt.
func (sh *SignalHub) Signal(signature string) error {
	defer sh.rw.Unlock()
	sh.rw.Lock()
	_, ok := sh.receivers[signature]
	if !ok {
		return errors.New("SignalHub.Signal() failed: receiver with key (signature string passed to method) does not exist")
	}
	sh.receivers[signature] <- struct{}{}
	return nil
}

// Register safely adds an entry to the receivers map, using the passed string signature as the map key.
func (sh *SignalHub) Register(signature string) (chan struct{}, bool) {
	_ = "breakpoint" // godebug
	sh.rw.Lock()
	defer sh.rw.Unlock()
	_, clash := sh.receivers[signature]
	if clash {
		log.Println("SignalHub.register() failed: receiver signature already exists")
		return nil, clash
	}

	sh.receivers[signature] = make(chan struct{})
	return sh.receivers[signature], false
}

// Deregister safely removes an entry from the receivers map with key "signature".
func (sh *SignalHub) Deregister(signature string) {
	defer sh.rw.Unlock()
	sh.rw.Lock()
	delete(sh.receivers, signature) // It's safe to do this even if the key is already absent from the map.
}

// WaitForSignalOnce - single-use operation in one tidy function:
// registers, blocks until signalled received once once, then deregisters on return.
func WaitForSignalOnce(s string, ts *gobr.SignalHub) {
	defer ts.Deregister(s)
	signal, ok := ts.Register(s)
	if !ok {
		return
	}
	<-signal // block while waiting for signal to be received.
}
