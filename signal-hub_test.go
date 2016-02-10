package gobr

// These tests really depend on panicking or not... so I guess I need to grab
// any panic and use it as the error for t.Errorf... but I can't be bothered
// yet!
// Currently it meets the requirements of `go test`, because if there is
// a panic, then test will FAIL and refer to the function where the panic
// occured as the test which failed.

import (
	"math/rand"
	"testing"
	"time"
)

func TestSignalHub(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	quit := make(chan struct{})
	sh := NewSignalHub()

	for i := 'A'; i < 'G'; i++ { // i.e. 'A', 'B", ··· , 'F'
		go func(signature string, hub *SignalHub) {
			defer hub.Deregister(signature)
			ch, _ := hub.Register(signature)
			for {
				select {
				case <-quit:
					return
				case <-ch:
					// fmt.Println(signature + " received signal")
				}
			}
		}(string(i), sh)
	}

	time.Sleep(25 * time.Millisecond)
	// fmt.Println("starting..")
	sh.Broadcast(false)
	// keeps main from exiting until this function literal returns
	func(hub *SignalHub) {
		rt := randomTicking(5, 100)
		signature := string(0x41 + rand.Intn(6)) // randomly signal a receiver with string key directly.
		timer := time.NewTimer(1 * time.Second)
		for {
			select {
			case <-rt:
				// fmt.Println("sending to", signature)
				hub.Signal(signature, false)
				signature = string(0x41 + rand.Intn(6)) // randomly pick rew receiver
			case <-timer.C:
				// fmt.Println("quitting...")
				close(quit)
				time.Sleep(25 * time.Millisecond)
				return
			}
		}
	}(sh)
}

func randomTicking(min, max int) chan struct{} {
	rt := make(chan struct{})
	go func(rt chan struct{}) {
		d := time.Millisecond * time.Duration(min+rand.Intn(max-min))
		// fmt.Println("rt burst in", d)
		for {
			select {
			case <-time.After(d):
				rt <- struct{}{}
				d = time.Millisecond * time.Duration(min+rand.Intn(max-min))
				// fmt.Println("rt burst in", d)
			}
		}
	}(rt)
	return rt
}

func TestWaitForSignalOnce(t *testing.T) {
	quit := make(chan struct{})
	sh := NewSignalHub()
	go func() {
		WaitForSignalOnce("TEST_WAIT_ONCE", sh)
		close(quit)
	}()
	time.Sleep(5 * time.Millisecond)
	sh.Signal("TEST_WAIT_ONCE", false)
	<-quit // blocks until quit closed
}

func TestWaitForBroadcastOnce(t *testing.T) {
	quit := make(chan struct{})
	sh := NewSignalHub()
	go func() {
		WaitForSignalOnce("TEST_WAIT_ONCE", sh)
		close(quit)
	}()
	time.Sleep(5 * time.Millisecond)
	sh.Broadcast(false)
	<-quit
}
