package regularcheck

import (
	"testing"
	"time"
)

func TestRegularcheck(t *testing.T) {
	var testcount int
	//var timeout bool
	checkFn := func() bool {
		if testcount > 5 {
			return false
		}
		return true
	}
	timeoutCallback := func() interface{} {
		t.Logf("over testconunt:%d\n", testcount)
		//timeout = true
		return testcount
	}
	//it will timeout after 3+1*3 = 6 seconds
	timeoutSecond := time.Duration(3 + 1*3)
	r := New(checkFn, timeoutCallback, time.Second*3, time.Second, 3)
	go r.Run()
	start := time.Now()
	for {
		time.Sleep(time.Second)
		testcount = 6 //make checkFn return false

		r.WaitForClose()
		t.Logf("r.Run() over, elapsed second:%d \n", time.Since(start)/time.Second)
		if time.Since(start)/time.Second != timeoutSecond {
			t.Fatalf("fail:elapsed second:%d != %d\n", time.Since(start)/time.Second, timeoutSecond)
		}
		break
	}

}

func TestRegularcheckReuse(t *testing.T) {
	var testcount int
	//var timeout bool
	checkFn := func() bool {
		if testcount > 5 {
			return false
		}
		return true
	}
	timeoutCallback := func() interface{} {
		t.Logf("over testconunt:%d\n", testcount)
		//timeout = true
		return testcount
	}
	//it will timeout after 3+1*3 = 6 seconds
	timeoutSecond := time.Duration(3 + 1*3)
	r := New(checkFn, timeoutCallback, time.Second*3, time.Second, 3)
	go r.Run()
	start := time.Now()
	for {
		time.Sleep(time.Second)
		testcount = 6 //make checkFn return false

		r.WaitForClose()
		t.Logf("r.Run() over, elapsed second:%d \n", time.Since(start)/time.Second)
		if time.Since(start)/time.Second != timeoutSecond {
			t.Fatalf("fail:elapsed second:%d != %d\n", time.Since(start)/time.Second, timeoutSecond)
		}
		break
	}
	t.Log("first time check over ,now going to checking reuse case")
	//testing reuse
	testcount = 0
	ok := r.RunWithGoroutine()
	if !ok {
		t.Fatal(" reuse regularcheck fail")
	}
	t.Log("reuse running is ok")
	start = time.Now()
	for {
		time.Sleep(time.Second)
		testcount = 6 //make checkFn return false

		r.WaitForClose()
		t.Logf("r.Run() over, elapsed second:%d \n", time.Since(start)/time.Second)
		if time.Since(start)/time.Second != timeoutSecond {
			t.Fatalf("fail:elapsed second:%d != %d\n", time.Since(start)/time.Second, timeoutSecond)
		}
		break
	}
	t.Log("success, reuse check is ok")
}
