package gobr

import (
	"errors"
	"fmt"
	"sync"
)

type signalHub struct {
	rw        sync.RWMutex
	receivers map[string]chan struct{} //	requires the receiver to keep track of its key as the registering party.
}

func newSignalHub() *signalHub {
	h := signalHub{}
	h.receivers = make(map[string]chan struct{})
	return &h
}

func (h *signalHub) broadcast() { // therefore no blocking for any senders! 😁
	defer h.rw.Unlock()
	h.rw.Lock()
	for signature, channel := range h.receivers { // if map is empty then no signals get sent or lost.
		fmt.Println("broadcasting to", signature)
		channel <- struct{}{}
	}
}

func (h *signalHub) register(signature string) (chan struct{}, error) {
	defer h.rw.Unlock()
	h.rw.Lock()
	_, ok := h.receivers[signature]
	if ok {
		return nil, errors.New("signalHub.register() failed: receiver signature already exists")
	}

	h.receivers[signature] = make(chan struct{})
	return h.receivers[signature], nil
}

func (h *signalHub) deregister(signature string) {
	defer h.rw.Unlock()
	h.rw.Lock()
	delete(h.receivers, signature)
}
