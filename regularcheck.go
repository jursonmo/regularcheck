package heartbeatCheck

import (
	"fmt"
	"sync"
	"time"
)

type RegularCheck struct {
	sync.Mutex
	checkFn                  func() bool
	timeoutCallback          func() interface{}
	checkIntv, failCheckIntv time.Duration
	maxFail, fail            int
	running, closing         bool
	closeWait                chan struct{}
	hasDryCloseChan          bool
	t                        *time.Timer
}

func New(checkFn func() bool, timeoutCallback func() interface{}, checkIntv, failCheckIntv time.Duration, maxFail int) *RegularCheck {
	return &RegularCheck{checkFn: checkFn,
		timeoutCallback: timeoutCallback,
		checkIntv:       checkIntv,
		failCheckIntv:   failCheckIntv,
		maxFail:         maxFail,
		closeWait:       make(chan struct{}, 1),
	}
}

func (r *RegularCheck) RunWithGoroutine() bool {
	synch := make(chan bool, 1)
	f := func() {
		r.run(synch)
	}
	go f()
	return <-synch
}

func notify(ch chan bool, b bool) {
	if ch == nil || cap(ch) == 0 {
		return
	}
	select {
	case ch <- b:
	default:
	}
}

func (r *RegularCheck) Run() bool {
	return r.run(nil)
}

func (r *RegularCheck) run(res chan bool) bool {
	r.Lock()
	if r.running {
		r.Unlock()
		notify(res, false)
		return false
	}
	r.running = true //avoid other goroutine r.Run()

	//make sure dry closeWait chan when r is reused.
	if !r.hasDryCloseChan {
		select {
		case <-r.closeWait:
		default:
		}
	}
	r.hasDryCloseChan = false //reset
	r.Unlock()

	notify(res, true)
	r.t = time.NewTimer(r.checkIntv)
	for {
		<-r.t.C
		if r.closing {
			break
		}
		if r.checkFn() {
			if !r.rstTimerToNormalCheck() {
				break
			}
			r.fail = 0 //reset fail count
			continue
		}

		//fail
		r.fail++
		if r.fail > r.maxFail {
			r.timeoutCallback()
			break
		}
		//check fail, so need to reset to failCheckIntv
		if !r.rstTimerToFailCheck() {
			break
		}
	}

	//reset source
	r.Lock()
	defer r.Unlock()
	// if !r.t.Stop() {
	// 	fmt.Printf("r.t.Stop() fail\n")
	// 	//<-r.t.C
	// 	select {
	// 	case <-r.t.C: //does't mean sendTime have been done, so here maybe can not dry the channel
	// 	default:
	// 	}
	// 	fmt.Printf("r.t.Stop() over\n")
	// }
	r.t.Stop() //we will create a new timer when we r.Run() again, so here we don't care whether Stop() return true or false,

	r.running = false
	r.closing = false
	r.fail = 0

	//notify CloseWithWait
	select {
	case r.closeWait <- struct{}{}:
	default:
	}

	fmt.Printf("check over, %s\n", r.Show())
	return true
}

func (r *RegularCheck) rstTimerToNormalCheck() bool {
	r.Lock()
	defer r.Unlock()
	if r.closing {
		return false
	}
	r.t.Reset(r.checkIntv) //reset to normal check interval
	return true
}

func (r *RegularCheck) rstTimerToFailCheck() bool {
	r.Lock()
	defer r.Unlock()
	if r.closing {
		return false
	}
	r.t.Reset(r.failCheckIntv) //reset to failCheckIntv
	return true
}

func (r *RegularCheck) Close() {
	r.Lock()
	defer r.Unlock()
	if !r.running {
		return
	}
	r.closing = true
	r.t.Reset(time.Duration(100))
}

func (r *RegularCheck) CloseWithWait() {
	r.Close()
	r.WaitForClose()
}

//不要WaitForClose 和 reuse 在不同的goroutie 里操作，这会有竞争问题， 如果reuse 时先清空closeWait，
func (r *RegularCheck) WaitForClose() {
	if r.hasDryCloseChan {
		return
	}
	select {
	case <-r.closeWait:
		r.hasDryCloseChan = true
	}
}

func (r *RegularCheck) Show() string {
	return fmt.Sprintf("maxFail:%d, fail:%d, running:%v, closing:%v, len(closewait):%d, hasDryCloseChan:%v",
		r.maxFail, r.fail, r.running, r.closing, len(r.closeWait), r.hasDryCloseChan)
}
