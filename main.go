package main

//latest
import (
	"context"
	"fmt"
	"github.com/gobwas/httphead"
	"github.com/gobwas/ws"
	"github.com/portto/solana-go-sdk/client"
	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/program/token"
	"github.com/portto/solana-go-sdk/rpc"
	"github.com/portto/solana-go-sdk/types"
	"go-blockchain-expert/ethereum-test"
	"log"
	"net"
	"net/http"
	"runtime"
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
	test7()
	test5()
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

	fmt.Println(len(disc))
	fmt.Println(cap(disc))

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

var test5Ch = make(chan uint64, 1024)

func test5() {
	//defer ants.Release()
	//
	//var pool, _ = ants.NewPoolWithFunc(1000, func(i interface{}) {
	//	//fmt.Println(i)
	//	time.Sleep(time.Millisecond)
	//})
	//
	//defer pool.Release()
	//
	//go func() {
	//	var i = uint64(0)
	//	for {
	//		i++
	//		test5Ch <- i
	//		if 0 == i%2000 {
	//			fmt.Println(pool.Cap())
	//			fmt.Println(pool.Free())
	//			fmt.Println(pool.Running())
	//			fmt.Println(pool.Waiting())
	//			time.Sleep(time.Second * 5)
	//			fmt.Println("---------------------")
	//		}
	//	}
	//}()
	//
	//go func() {
	//	for {
	//		pool.Invoke(<-test5Ch)
	//	}
	//}()
	//
	////time.Sleep(time.Second * 3)
	////
	////ants.Submit(func() {
	////	fmt.Println("1")
	////})
	////ants.Submit(func() {
	////	fmt.Println("2")
	////})
	////ants.Submit(func() {
	////	fmt.Println("3")
	////})
	////ants.Submit(func() {
	////	fmt.Println("5")
	////})
	//
	//time.Sleep(time.Hour)
}

func test6() {
	ln, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		// handle error
	}

	// Prepare handshake header writer from http.Header mapping.
	header := ws.HandshakeHeaderHTTP(http.Header{
		"X-Go-Version": []string{runtime.Version()},
	})

	u := ws.Upgrader{
		OnHost: func(host []byte) error {
			if string(host) == "github.com" {
				return nil
			}
			return ws.RejectConnectionError(
				ws.RejectionStatus(403),
				ws.RejectionHeader(ws.HandshakeHeaderString(
					"X-Want-Host: github.com\r\n",
				)),
			)
		},
		OnHeader: func(key, value []byte) error {
			if string(key) != "Cookie" {
				return nil
			}
			ok := httphead.ScanCookie(value, func(key, value []byte) bool {
				// Check session here or do some other stuff with cookies.
				// Maybe copy some values for future use.
				return true
			})
			if ok {
				return nil
			}
			return ws.RejectConnectionError(
				ws.RejectionReason("bad cookie"),
				ws.RejectionStatus(400),
			)
		},
		OnBeforeUpgrade: func() (ws.HandshakeHeader, error) {
			return header, nil
		},
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		_, err = u.Upgrade(conn)
		if err != nil {
			log.Printf("upgrade error: %s", err)
		}
	}
}

//// FUarP2p5EnxD66vVDL4PWRoWMzA56ZVHG24hpEDFShEz
//var feePayer, _ = types.AccountFromBase58("4TMFNY9ntAn3CHzguSAvDNLPRoQTaK3sWbQQXdDXaE6KWRBLufGL6PJdsD2koiEe3gGmMdRK3aAw7sikGNksHJrN")
//
//// 9aE476sH92Vz7DMPyq5WLPkrKWivxeuTKEFKd2sZZcde
//var alice, _ = types.AccountFromBase58("4voSPg3tYuWbKzimpQK9EbXHmuyy5fUrtXvpLDMLkmY6TRncaTHAKGD8jUg3maB5Jbrd9CkQg4qjJMyN6sQvnEF2")
//
//var mintPubkey = common.PublicKeyFromString("F6tecPzBMF47yJ2EN6j2aGtE68yR5jehXcZYVZa6ZETo")
//
//var aliceTokenRandomTokenPubkey = common.PublicKeyFromString("HeCBh32JJ8DxcjTyc6q46tirHR8hd2xj3mGoAcQ7eduL")
//
//var aliceTokenATAPubkey = common.PublicKeyFromString("J1T6kAPowNFcxFh4pmwSqxQM9AitN7HwLyvxV2ZGfLf2")
//
//func test7() {
//	c := client.NewClient(rpc.DevnetRPCEndpoint)
//
//	res, err := c.GetLatestBlockhash(context.Background())
//	if err != nil {
//		log.Fatalf("get recent block hash error, err: %v\n", err)
//	}
//
//	tx, err := types.NewTransaction(types.NewTransactionParam{
//		Message: types.NewMessage(types.NewMessageParam{
//			FeePayer:        feePayer.PublicKey,
//			RecentBlockhash: res.Blockhash,
//			Instructions: []types.Instruction{
//				token.TransferChecked(token.TransferCheckedParam{
//					From:     aliceTokenRandomTokenPubkey,
//					To:       aliceTokenATAPubkey,
//					Mint:     mintPubkey,
//					Auth:     alice.PublicKey,
//					Signers:  []common.PublicKey{},
//					Amount:   1e8,
//					Decimals: 8,
//				}),
//			},
//		}),
//		Signers: []types.Account{feePayer, alice},
//	})
//	if err != nil {
//		log.Fatalf("failed to new tx, err: %v", err)
//	}
//
//	txhash, err := c.SendTransaction(context.Background(), tx)
//	if err != nil {
//		log.Fatalf("send raw tx error, err: %v\n", err)
//	}
//
//	log.Println("txhash:", txhash)
//}

// xD66vVDL4PWRoWMzA56ZVHG24hpEDFShEz
var feePayer, _ = types.AccountFromBase58("4ukSnzryvAVU98uDGQTXgTgfhxiN7uLZ5t21bBvvyRQAL8enozLmZYZHPBLRbdkj6xesnzkwYy2pWYDxxLPiTwfY")

// 9aE476sH92Vz7DMPyq5WLPkrKWivxeuTKEFKd2sZZcde
//var alice, _ = types.AccountFromBase58("4ukSnzryvAVU98uDGQTXgTgfhxiN7uLZ5t21bBvvyRQAL8enozLmZYZHPBLRbdkj6xesnzkwYy2pWYDxxLPiTwfY")

var mintPubkey = common.PublicKeyFromString("Es9vMFrzaCERmJfrF4H2FYD4KCoNkY11McCe8BenwNYB") //

var aliceTokenRandomTokenPubkey = common.PublicKeyFromString("7yLhAmGMMJvqDzyN13WjUVK1VRoZusqvA3UUnMvFWcek") //

//var aliceTokenATAPubkey = common.PublicKeyFromString("9nj6jK69KGU38q1G5wC822MpMMMryUyhAMyzTFLbtmV") //
var aliceTokenATAPubkey = common.PublicKeyFromString("7VhJLc53PH2T17sS46qHDaCYCR7ZcRSdPPpjNApjiJTh") //

func test7() {
	//fmt.Println("from", aliceTokenRandomTokenPubkey.ToBase58())
	//var from, _ = types.AccountFromBase58("4ukSnzryvAVU98uDGQTXgTgfhxiN7uLZ5t21bBvvyRQAL8enozLmZYZHPBLRbdkj6xesnzkwYy2pWYDxxLPiTwfY")
	//aliceTokenRandomTokenPubkey = from.PublicKey
	//fmt.Println("from", aliceTokenRandomTokenPubkey.ToBase58())
	//
	//fmt.Println("to", aliceTokenATAPubkey.ToBase58())
	//var to, _ = types.AccountFromBase58("5U9iwZfR9rreJSLXuHQvrray4MWuQTjHHrDMBVEcuxSgtdGjakd46Vw1nDhrXetJt3HM2waE4aMH4MhRare3zzcu")
	//aliceTokenATAPubkey = to.PublicKey
	//fmt.Println("to", aliceTokenATAPubkey.ToBase58())

	c := client.NewClient(rpc.MainnetRPCEndpoint)

	res, err := c.GetLatestBlockhash(context.Background())
	if err != nil {
		log.Fatalf("get recent block hash error, err: %v\n", err)
	}

	//tx, err := types.NewTransaction(types.NewTransactionParam{
	//	Message: types.NewMessage(types.NewMessageParam{
	//		FeePayer:        feePayer.PublicKey,
	//		RecentBlockhash: res.Blockhash,
	//		Instructions: []types.Instruction{
	//			token.TransferChecked(token.TransferCheckedParam{
	//				From:     aliceTokenRandomTokenPubkey,
	//				To:       aliceTokenATAPubkey,
	//				Mint:     mintPubkey,
	//				Auth:     feePayer.PublicKey,
	//				Signers:  []common.PublicKey{},
	//				Amount:   1,
	//				Decimals: 6,
	//			}),
	//		},
	//	}),
	//	Signers: []types.Account{feePayer},
	//})
	//if err != nil {
	//	log.Fatalf("failed to new tx, err: %v", err)
	//}
	//
	//txhash, err := c.SendTransaction(context.Background(), tx)
	tx, err := types.NewTransaction(types.NewTransactionParam{
		Message: types.NewMessage(types.NewMessageParam{
			FeePayer:        feePayer.PublicKey,
			RecentBlockhash: res.Blockhash,
			Instructions: []types.Instruction{
				token.Transfer(token.TransferParam{
					From:    aliceTokenRandomTokenPubkey,
					To:      aliceTokenATAPubkey,
					Auth:    feePayer.PublicKey,
					Signers: []common.PublicKey{},
					Amount:  222222,
					//From:     aliceTokenRandomTokenPubkey,
					//To:       aliceTokenATAPubkey,
					//Mint:     mintPubkey,
					//Auth:     feePayer.PublicKey,
					//Signers:  []common.PublicKey{},
					//Amount:   1,
					//Decimals: 6,
				}),
			},
		}),
		Signers: []types.Account{feePayer},
	})
	if err != nil {
		log.Fatalf("failed to new tx, err: %v", err)
	}

	txhash, err := c.SendTransaction(context.Background(), tx)
	if err != nil {
		log.Fatalf("send raw tx error, err: %v\n", err)
	}

	log.Println("txhash:", txhash)
}

func test8() {

}
