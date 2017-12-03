package main

import "fmt"
import "time"
import "sync"
import "math/rand"

func sleep() {
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
}

func philosopher(name string, a *sync.Mutex, b *sync.Mutex, q chan bool, log chan string, req chan int, ack chan bool) {
	log <- name + ":start"
	for i := 0; i < 100; i++ {
		sleep()
		for {
			req <- -1
			if <- ack {
				break
			}
			sleep()
		}
		a.Lock()
		log <- name + ":left"
		sleep()
		b.Lock()
		req <- 1
		log <- name + ":right"
		sleep()
		log <- name + ":done"
		a.Unlock()
		b.Unlock()
	}
	log <- name + ":finish"
	q <- true
}

func philosopherL(name string, a *sync.Mutex, b *sync.Mutex, q chan bool, log chan string, req chan int, ack chan bool) {
	log <- name + ":start"
	for i := 0; i < 100; i++ {
		sleep()
		for {
			req <- -1
			if <- ack {
				break
			}
			sleep()
		}
		b.Lock()
		log <- name + ":right"
		sleep()
		a.Lock()
		log <- name + ":left"
		sleep()
		log <- name + ":done"
		a.Unlock()
		b.Unlock()
	}
	log <- name + ":finish"
	q <- true
}

func monitor(log chan string) {
	al := " "
	ar := " "
	bl := " "
	br := " "
	cl := " "
	cr := " "
	dl := " "
	dr := " "
	el := " "
	er := " "
	for {
		mes := <-log
		switch mes {
		case "A:left":
			al = "Y"
		case "A:right":
			ar = "Y"
		case "A:done":
			al = " "
			ar = " "
		case "B:left":
			bl = "Y"
		case "B:right":
			br = "Y"
		case "B:done":
			bl = " "
			br = " "
		case "C:left":
			cl = "Y"
		case "C:right":
			cr = "Y"
		case "C:done":
			cl = " "
			cr = " "
		case "D:left":
			dl = "Y"
		case "D:right":
			dr = "Y"
		case "D:done":
			dl = " "
			dr = " "
		case "E:left":
			el = "Y"
		case "E:right":
			er = "Y"
		case "E:done":
			el = " "
			er = " "
		}
		fmt.Printf("%sA%s  ", al, ar)
		fmt.Printf("%sB%s  ", bl, br)
		fmt.Printf("%sC%s  ", cl, cr)
		fmt.Printf("%sD%s  ", dl, dr)
		fmt.Printf("%sE%s  \n", el, er)
	}
}

func count(req chan int, ack chan bool){
	count := 5
	for {
		r := <- req
		if r == -1 {
			if count == 1 {
				ack <- false
				fmt.Println("Wait!!")
			} else {
				count--
				ack <- true
			}
		} else if r == 1 {
			count++
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	m := []*sync.Mutex{new(sync.Mutex), new(sync.Mutex), new(sync.Mutex), new(sync.Mutex), new(sync.Mutex)}
	q := make(chan bool)
	log := make(chan string)
	ack := make(chan bool)
	req := make(chan int)

	go monitor(log)
	go count(req, ack)

	go philosopher("A", m[0], m[1], q, log, req, ack)
	go philosopher("B", m[1], m[2], q, log, req, ack)
	go philosopher("C", m[2], m[3], q, log, req, ack)
	go philosopher("D", m[3], m[4], q, log, req, ack)
	go philosopher("E", m[4], m[0], q, log, req, ack)

	for i := 0; i < 5; i++ {
		<-q
	}
}
