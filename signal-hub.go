package gobr

import (
	"errors"
	"fmt"
	"sync"
)

// SignalHub is for simple non-blocking to signal receiver channels individually and broadcasting to all receivers
type SignalHub struct {
	rw        sync.RWMutex
	receivers map[string]chan struct{} //	requires the receiver to keep track of its key ("signature") as the registering party.
}

// NewSignalHub returns a pointer to a new instance, with the sync.RWMutex set to the default (Unlocked).
func NewSignalHub() *SignalHub {
	h := SignalHub{}
	h.receivers = make(map[string]chan struct{})
	return &h
}

// Broadcast safely sends a signal to all receivers. By using invoking the method *as a goroutine* instead of sending channel directly, we avoid the chance of the sender blocking while waiting for the receiver to retreive the message sent on the channel.
// This also means that we can Broadcast "as-if" there are processes waiting to receive the synchronisation signal, even though in some cases there might not be.
func (h *SignalHub) Broadcast() {
	defer h.rw.Unlock()
	h.rw.Lock()
	for signature, channel := range h.receivers { // if map is empty then no signals get sent or lost.
		fmt.Println("broadcasting to", signature)
		channel <- struct{}{}
	}
}

// Signal 'safely' sends a signal message (struct{}{}) without blocking the sender *if invoked as a goroutine*.
// Otherwise you need to check for the error if you want to confirm receipt.
func (h *SignalHub) Signal(signature string) error {
	defer h.rw.Unlock()
	h.rw.Lock()
	_, ok := h.receivers[signature]
	if !ok {
		return errors.New("SignalHub.Signal() failed: receiver with key (signature string passed to method) does not exist")
	}
	h.receivers[signature] <- struct{}{}
	return nil
}

// Register safely adds an entry to the receivers map, using the passed string signature as the map key.
func (h *SignalHub) Register(signature string) (chan struct{}, error) {
	defer h.rw.Unlock()
	h.rw.Lock()
	_, ok := h.receivers[signature]
	if ok {
		return nil, errors.New("SignalHub.register() failed: receiver signature already exists")
	}

	h.receivers[signature] = make(chan struct{})
	return h.receivers[signature], nil
}

// Deregister safely removes an entry from the receivers map with key "signature".
func (h *SignalHub) Deregister(signature string) {
	defer h.rw.Unlock()
	h.rw.Lock()
	delete(h.receivers, signature)
}
