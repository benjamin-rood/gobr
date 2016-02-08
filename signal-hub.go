package gobr

import (
	"errors"
	"fmt"
	"sync"
)

type SignalHub struct {
	rw        sync.RWMutex
	receivers map[string]chan struct{} //	requires the receiver to keep track of its key ("signature") as the registering party.
}

func NewSignalHub() *SignalHub {
	h := SignalHub{}
	h.receivers = make(map[string]chan struct{})
	return &h
}

func (h *SignalHub) Broadcast() { // therefore no blocking for any senders! 😁
	defer h.rw.Unlock()
	h.rw.Lock()
	for signature, channel := range h.receivers { // if map is empty then no signals get sent or lost.
		fmt.Println("broadcasting to", signature)
		channel <- struct{}{}
	}
}

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

func (h *SignalHub) Deregister(signature string) {
	defer h.rw.Unlock()
	h.rw.Lock()
	delete(h.receivers, signature)
}
