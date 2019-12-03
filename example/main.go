package main

import (
	"fmt"
	"heartbeatCheck"
	"time"
)

func main() {
	var testcount int
	var timeout bool
	checkFn := func() bool {
		fmt.Printf("checkFn: testcount=%d\n", testcount)
		if testcount > 3 {
			return false
		}
		return true
	}
	timeoutCallback := func() interface{} {
		fmt.Printf("timeout testconunt:%d\n", testcount)
		timeout = true
		return testcount
	}

	r := heartbeatCheck.New(checkFn, timeoutCallback, time.Second*3, time.Second, 3)
	go r.Run()
	for {
		time.Sleep(time.Second)
		testcount++ //make checkFn return false
		if timeout {
			fmt.Println(r.Show())
			r.WaitForClose()
			fmt.Println("r.Run() over")
			break
		}
	}

	//testing reuse regularcheck
	testcount = 0
	timeout = false
	ok := r.RunWithGoroutine()
	if !ok {
		panic(" reuse regularcheck fail")
	}
	fmt.Println("reuse ok")
	for {
		time.Sleep(time.Second)
		testcount++ //make checkFn return false
		if timeout {
			fmt.Println(r.Show())
			r.WaitForClose()
			fmt.Println("2222: r.Run() over")
			break
		}
	}

	//test sync
	testcount = 0
	timeout = false
	ok = r.RunWithGoroutine()
	if !ok {
		panic("reuse regularcheck fail")
	}
	fmt.Println("reuse ok")

	ok = r.RunWithGoroutine()
	if ok {
		panic(" multi Run regularcheck fail")
	}
	fmt.Println("ok: one regularcheck Running is not reentrant")
	for {
		time.Sleep(time.Second)
		testcount++ //make checkFn return false
		if timeout {
			fmt.Println(r.Show())
			r.WaitForClose()
			fmt.Println("333: r.Run() over")
			break
		}
	}
}
