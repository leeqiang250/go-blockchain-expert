package main

import (
	"fmt"
	"go-blockchain-expert/ethereum-test"
	"sync"
	"time"
)

func b() (i int) {
	defer func() {
		i++
		fmt.Println("defer2:", i)
	}()
	defer func() {
		i++
		fmt.Println("defer1:", i)
	}()
	return i //或者直接写成return
}

func c() *int {
	var i int
	defer func() {
		i++
		fmt.Println("defer2:", i)
	}()
	defer func() {
		i++
		fmt.Println("defer1:", i)
	}()
	return &i
}

func cc(data []int) {
	fmt.Printf("%p", &data)
	fmt.Println(len(data))

	for len(data) < 3 {
		data = append(data, len(data))
	}

	fmt.Printf("%p", &data)
	fmt.Println(len(data))
}

func main() {
	test4()

	ethereum_test.Start()

	data := make([]int, 0, 10)
	data = append(data, len(data))
	data = append(data, len(data))

	fmt.Printf("%p", &data)
	fmt.Println(len(data))

	cc(data)

	fmt.Printf("%p", &data)
	fmt.Println(len(data))

	fmt.Println("return:", *(c()))
	fmt.Println("return:", c())

	fmt.Println("return:", b())

	//ethash.HashimotoFull(nil, nil, 2)

	test3()

	go test1()

	time.Sleep(time.Hour)

	defer test0(7)
	defer test0(8)
	defer test0(9)
}

func test0(num int) {
	var dd = sync.NewCond(nil)
	dd.Signal()

	fmt.Println(num, time.Now())
}

func test1() {
	var mu sync.Mutex
	var cond = sync.NewCond(&mu)

	var i = 0
	for i < 100 {
		go func(i int) {
			cond.L.Lock()
			cond.Wait()
			fmt.Println(i)
			cond.L.Unlock()
		}(i)
		i++
	}

	time.Sleep(time.Second * 3)
	cond.Signal()

	time.Sleep(time.Second * 3)
	cond.Signal()

	time.Sleep(time.Second * 3)
	cond.Broadcast()

	time.Sleep(time.Second * 3)
	cond.Signal()

	time.Sleep(time.Second * 3)
	cond.Broadcast()

}

func test2() {
	var disc = make(chan uint64, 2)
	var closed = make(chan uint64, 2)

	go func() {
		var value, ok = <-disc
		fmt.Println(ok)
		value++
		closed <- value
	}()

	time.Sleep(time.Second * 3)

	go func() {
		select {
		case disc <- 3:
		case d, e := <-closed:
			{
				fmt.Println(d, e)
			}
		}
	}()

	time.Sleep(time.Hour)
}

func test3() {
	var disc = make(chan int64, 2)
	var closed = make(chan int64, 2)

	go func() {
		for {
			var random = time.Now().UnixMicro() % 3
			fmt.Println("----------------")
			fmt.Println(random)
			if 0 == random {
				disc <- time.Now().UnixMicro()
			} else if 1 == random {
				closed <- time.Now().UnixMicro()
			} else if 2 == random {

			}
			time.Sleep(time.Second * 2)
		}
	}()

	go func() {
		for {
			select {
			case result, ok := <-disc:
				{
					fmt.Println("disc", result, ok)
				}
			case result, ok := <-closed:
				{
					fmt.Println("closed", result, ok)
				}
			}
		}
	}()

	time.Sleep(time.Hour)
}

func test4() {
	var disc = make(chan int64, 20)
	var closed = make(chan int64, 20)

	disc <- 1
	disc <- 1
	disc <- 1

	closed <- 1
	closed <- 1
	closed <- 1

	fmt.Println("1")
	select {
	case result, ok := <-disc:
		{
			fmt.Println("disc", result, ok)
		}
	case result, ok := <-closed:
		{
			fmt.Println("closed", result, ok)
		}
	}
	fmt.Println("2")

}
