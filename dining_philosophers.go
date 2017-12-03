package main

import "fmt"
import "time"
import "sync"
import "math/rand"

func sleep() {
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
}

func philosopher(name string, l *sync.Mutex, r *sync.Mutex, q chan bool, log chan string) {
	log <- name + ":start"
	for i := 0; i < 100; i++ {
		sleep()
		l.Lock()
		log <- name + ":left"
		sleep()
		r.Lock()
		log <- name + ":right"
		sleep()
		log <- name + ":done"
		l.Unlock()
		r.Unlock()
	}
	log <- name + ":finish"
	q <- true
}

func philosopherL(name string, l *sync.Mutex, r *sync.Mutex, q chan bool, log chan string) {
	log <- name + ":start"
	for i := 0; i < 100; i++ {
		sleep()
		r.Lock()
		log <- name + ":right"
		sleep()
		l.Lock()
		log <- name + ":left"
		sleep()
		log <- name + ":done"
		l.Unlock()
		r.Unlock()
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

func main() {
	rand.Seed(time.Now().UnixNano())
	var f []*sync.Mutex
	for i := 0; i < 5; i++ {
		f = append(f, new(sync.Mutex))
	}
	q := make(chan bool)
	log := make(chan string)

	go monitor(log)

	go philosopher("A", f[0], f[1], q, log)
	go philosopher("B", f[1], f[2], q, log)
	go philosopher("C", f[2], f[3], q, log)
	go philosopher("D", f[3], f[4], q, log)
	go philosopher("E", f[4], f[0], q, log)

	for i := 0; i < 5; i++ {
		<-q
	}
}
