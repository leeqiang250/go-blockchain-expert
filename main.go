package main

//latest
import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gobwas/httphead"
	"github.com/gobwas/ws"
	"github.com/portto/solana-go-sdk/client"
	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/program/token"
	"github.com/portto/solana-go-sdk/rpc"
	"github.com/portto/solana-go-sdk/types"
	"github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"
	"go-blockchain-expert/ethereum-test"
	"go-blockchain-expert/src"
	"log"
	"math"
	"net"
	"net/http"
	"runtime"
	"strconv"
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

var Client *redis.Client

func main() {

	Client = redis.NewClient(&redis.Options{
		Addr:     "localhost:32768",
		Username: "default",
		Password: "redispw", // no password set
		DB:       0,         // use default DB
	})

	go src.NewTaskRedis(Client, "if (redis.call('get', KEYS[1]) == ARGV[1]) then redis.call('expire', KEYS[1], ARGV[2]) return 1 else if (1 == redis.call('setnx', KEYS[1], ARGV[1])) then redis.call('expire', KEYS[1], ARGV[2]) return 1 else return 0 end return 0 end", time.Second, "keyUniqueName", "nodeUniqueName", func() {
		fmt.Println("nodeUniqueName")
	}).Run()

	src.NewTaskRedis(Client, "if (redis.call('get', KEYS[1]) == ARGV[1]) then redis.call('expire', KEYS[1], ARGV[2]) return 1 else if (1 == redis.call('setnx', KEYS[1], ARGV[1])) then redis.call('expire', KEYS[1], ARGV[2]) return 1 else return 0 end return 0 end", time.Second, "keyUniqueName", "nodeUniqueName2", func() {
		fmt.Println("nodeUniqueName2")
	}).Run()

	src.XXX()

	ethereum_test.Start()
	//src.Test()

	//fmt.Println(strconv.FormatUint(math.MaxUint64, 10))
	//var dd = decimal.NewFromFloat(math.MaxFloat64).String()
	{
		var data = decimal.NewFromFloat(33.3333333333333333333)
		fmt.Println(data)
		fmt.Println(data.String())
	}

	test13()
	test12()
	test11()
	test10()
	test8()
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
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:32768",
		Username: "default",
		Password: "redispw", // no password set
		DB:       0,         // use default DB
	})

	//test9(rdb)

	fmt.Println(rdb.Ping(context.Background()).Result())

	//fmt.Println(rdb.Set(context.Background(), "a", "a", time.Hour))
	//fmt.Println(rdb.Get(context.Background(), "a"))
	//
	//rdb.FlushAll(context.Background())

	//fmt.Println(rdb.ScriptKill(context.Background()))
	//fmt.Println(rdb.ScriptFlush(context.Background()))

	//fmt.Println(rdb.ScriptLoad(context.Background(), "if redis.call('get', KEYS[1]) == ARGV[1] then return redis.call('del', KEYS[1]) else return 0 end"))
	//
	//fmt.Println(rdb.ScriptExists(context.Background(), "e9f69f2beb755be68b5e456ee2ce9aadfbc4ebf4"))
	//fmt.Println(rdb.ScriptExists(context.Background(), "e9f69f2beb755be68b5e456ee2ce9aadfbc4ebf5"))
	//
	//fmt.Println(rdb.EvalSha(context.Background(), "e9f69f2beb755be68b5e456ee2ce9aadfbc4ebf4", []string{"a-k"}, "a-v"))
	//
	//fmt.Println(rdb.Eval(context.Background(), "if redis.call('get', KEYS[1]) == ARGV[1] then return redis.call('del', KEYS[1]) else return 0 end", []string{"a-k"}, "a-v"))
	//
	//fmt.Println(rdb)

	var key = "xxx"
	var symbol = "btcusdt"
	var value100 = "100"
	var value150 = "150"
	var value200 = "200"

	{
		var script = "local seq = redis.call('hget', KEYS[2], ARGV[1]) if (seq) then if (tonumber(ARGV[2]) > tonumber(seq)) then redis.call('lpush', KEYS[1], ARGV[3]) redis.call('hset', KEYS[2], ARGV[1], ARGV[2]) if (5 < redis.call('llen', KEYS[1])) then redis.call('rpop', KEYS[1]) end return 3 else return 2 end  else redis.call('lpush', KEYS[1], ARGV[3]) redis.call('hset', KEYS[2], ARGV[1], ARGV[2]) return 1 end"
		var stringCmd = rdb.ScriptLoad(context.Background(), script)
		fmt.Println(stringCmd)
		fmt.Println(rdb.ScriptExists(context.Background(), stringCmd.Val()))

		var i = uint64(0)
		for i < 1000 {
			i++

			var trade = src.MarketTrade{
				Seq:    i,
				Ts:     uint64(time.Now().UnixMilli()),
				Side:   src.BUY,
				Price:  decimal.NewFromInt(int64(i)),
				Vol:    decimal.NewFromInt(int64(i)),
				Amount: decimal.NewFromInt(int64(i)),
				Offset: i,
			}

			var json, _ = json.Marshal(trade)
			fmt.Println(string(json))
			fmt.Println(rdb.EvalSha(context.Background(), stringCmd.Val(), []string{fmt.Sprintf(src.MARKET_TRADE, symbol), fmt.Sprintf(src.MARKET_TRADE_SEQ, symbol)}, symbol, trade.Seq, json))
			fmt.Println(rdb.LLen(context.Background(), fmt.Sprintf(src.MARKET_TRADE, symbol)))
			fmt.Println(rdb.LRange(context.Background(), fmt.Sprintf(src.MARKET_TRADE, symbol), 0, 1000).Val())
		}
	}

	rdb.HSet(context.Background(), key, symbol, value100)

	var script0 = "return redis.call('hget', '" + key + "', '" + symbol + "')"
	var stringCmd0 = rdb.ScriptLoad(context.Background(), script0)
	fmt.Println(stringCmd0)
	fmt.Println(rdb.ScriptExists(context.Background(), stringCmd0.Val()))
	fmt.Println(rdb.EvalSha(context.Background(), stringCmd0.Val(), nil))

	var script1 = "return redis.call('hget', KEYS[1], KEYS[2])"
	var stringCmd1 = rdb.ScriptLoad(context.Background(), script1)
	fmt.Println(stringCmd1)
	fmt.Println(rdb.ScriptExists(context.Background(), stringCmd1.Val()))
	fmt.Println(rdb.EvalSha(context.Background(), stringCmd1.Val(), []string{key, symbol}))
	rdb.HSet(context.Background(), key, symbol, value200)
	fmt.Println(rdb.EvalSha(context.Background(), stringCmd1.Val(), []string{key, symbol}))

	rdb.HSet(context.Background(), key, symbol, value100)
	fmt.Println(rdb.HGet(context.Background(), key, symbol))

	var script2 = "if(tonumber(KEYS[3])>tonumber(redis.call('hget', KEYS[1], KEYS[2]))) then return 1 else return 2 end"
	//script2 = "if(KEYS[3] > redis.call('hget', KEYS[1], KEYS[2])) then return 1 else return 2 end"
	var stringCmd2 = rdb.ScriptLoad(context.Background(), script2)
	fmt.Println(stringCmd2)
	fmt.Println(rdb.ScriptExists(context.Background(), stringCmd2.Val()))

	rdb.HSet(context.Background(), key, symbol, value100)
	fmt.Println(rdb.EvalSha(context.Background(), stringCmd2.Val(), []string{key, symbol, value150}))
	rdb.HSet(context.Background(), key, symbol, value200)
	fmt.Println(rdb.EvalSha(context.Background(), stringCmd2.Val(), []string{key, symbol, value150}))

	fmt.Println(rdb.EvalSha(context.Background(), stringCmd2.Val(), []string{key, symbol, "9"}), "9")
	fmt.Println(rdb.EvalSha(context.Background(), stringCmd2.Val(), []string{key, symbol, "10"}), "10")
	fmt.Println(rdb.EvalSha(context.Background(), stringCmd2.Val(), []string{key, symbol, "20"}), "20")
	fmt.Println(rdb.EvalSha(context.Background(), stringCmd2.Val(), []string{key, symbol, "100"}), "100")
	fmt.Println(rdb.EvalSha(context.Background(), stringCmd2.Val(), []string{key, symbol, "199"}), "199")
	fmt.Println(rdb.EvalSha(context.Background(), stringCmd2.Val(), []string{key, symbol, "200"}), "200")
	fmt.Println(rdb.EvalSha(context.Background(), stringCmd2.Val(), []string{key, symbol, "201"}), "201")
	fmt.Println(rdb.EvalSha(context.Background(), stringCmd2.Val(), []string{key, symbol, "300"}), "300")
	fmt.Println(rdb.EvalSha(context.Background(), stringCmd2.Val(), []string{key, symbol, "1300"}), "1300")

	var cmd1 = rdb.EvalSha(context.Background(), stringCmd2.Val(), []string{key, symbol + "x", "9"})
	fmt.Println(cmd1)
	fmt.Println(cmd1.Val())
	fmt.Println(cmd1.Result())
	fmt.Println(cmd1.Bool())
	fmt.Println(rdb.EvalSha(context.Background(), stringCmd2.Val(), []string{key, symbol + "x", "1300"}), "1300")

	var script3 = "return tonumber(KEYS[1])"
	var stringCmd3 = rdb.ScriptLoad(context.Background(), script3)
	fmt.Println(stringCmd3)
	fmt.Println(rdb.ScriptExists(context.Background(), stringCmd3.Val()))

	var mat = strconv.FormatUint(math.MaxUint64, 10)
	fmt.Println(mat)
	mat = strconv.FormatUint(math.MaxUint64/uint64(2), 10)
	fmt.Println(mat)
	fmt.Println(rdb.EvalSha(context.Background(), stringCmd3.Val(), []string{"20"}))
	fmt.Println(rdb.EvalSha(context.Background(), stringCmd3.Val(), []string{mat}))
	//18446744073709551615
	//9223372036854775807
	var total = math.MaxUint64 / uint64(2)
	fmt.Println(total)

	total = total / uint64(100)
	fmt.Println(total) //每年

	total = total / uint64(365)
	fmt.Println(total) //每天

	total = total / uint64(24)
	fmt.Println(total) //每小时

	total = total / uint64(60)
	fmt.Println(total) //每分钟

	total = total / uint64(60)
	fmt.Println(total) //每秒

}

func test9(client *redis.Client, seq string) {
	client.FlushAll(context.Background())

	var key_seq = "market::overview::seq"
	var key = "market::overview"
	var symbol = "btcusdt"
	//var value100 = "100"
	//var value150 = "150"
	//var value200 = "200"

	var script = "local seq = redis.call('hget', KEYS[1], KEYS[3]) if (seq) then return 1 else redis.call('hset', KEYS[1], KEYS[3], ARGV[1]) return 0 end"
	//script = "local seq = redis.call('hget', KEYS[1], KEYS[2]) return seq"
	var stringCmd = client.ScriptLoad(context.Background(), script)
	var result, err = stringCmd.Result()
	if nil != err {
		fmt.Println(err)
		return
	}
	fmt.Println(result)

	fmt.Println(client.HDel(context.Background(), key, symbol))
	fmt.Println(client.EvalSha(context.Background(), result, []string{key_seq, key, symbol}, seq))
	fmt.Println(client.HGet(context.Background(), key, symbol))
	fmt.Println(client.HSet(context.Background(), key, symbol, "31232"))
	fmt.Println(client.EvalSha(context.Background(), result, []string{key_seq, key, symbol}, seq))
	fmt.Println(client.HGet(context.Background(), key, symbol))

}

func test10() {
	var ch = make(chan uint64, 2000)

	go func() {
		var i = uint64(0)
		for {
			i++
			ch <- i
		}
	}()

	//交易释放
	for {
		for 100 < len(ch) {
			<-ch
		}
		//<-ch
		fmt.Println(len(ch), <-ch)
		time.Sleep(time.Second)
	}
}

func test11() {
	var ch = make(chan uint64, 2000)

	go func() {
		var i = uint64(0)
		for {
			i++
			ch <- i
		}
	}()

	for {
		for 1 < len(ch) {
			<-ch
		}
		//<-ch
		fmt.Println(len(ch), <-ch)
		time.Sleep(time.Second)
	}
}

func test12() {
	var ch = make(chan int64, 2000)

	go func() {
		for {
			ch <- time.Now().UnixMilli()
		}
	}()

	var ts = int64(0)
	var pending = int64(0)
	var ticker = time.Tick(time.Millisecond * 100)

	for {
		for 1 < len(ch) {
			<-ch
		}

		select {
		case data := <-ch:
			if ts+100 < time.Now().UnixMilli() {
				fmt.Println("xxxxxx", len(ch), time.Now().UnixMilli()-data)
				time.Sleep(time.Second)
				ts = time.Now().UnixMilli()
			} else {
				pending = data
			}
		case <-ticker:
			if 0 < pending {
				fmt.Println("ticker", len(ch), time.Now().UnixMilli()-pending)
				time.Sleep(time.Second)
				pending = 0
			}
		}
	}
}

func test13() {
	var ch = make(chan *src.MarketLine, 2000)

	go func() {
		var seq = uint64(0)
		for {
			seq++

			ch <- &src.MarketLine{
				Seq:     seq,
				ID:      uint64(time.Now().UnixMilli() / 1000),
				Open:    decimal.Zero,
				Close:   decimal.Zero,
				High:    decimal.Zero,
				Low:     decimal.Zero,
				Vol:     decimal.Zero,
				Amount:  decimal.Zero,
				Count:   0,
				IsClose: false,
				Offset:  0,
			}

			var dd, _ = json.Marshal(&src.MarketLine{
				Seq:     seq,
				ID:      uint64(time.Now().UnixMilli() / 1000),
				Open:    decimal.Zero,
				Close:   decimal.Zero,
				High:    decimal.NewFromFloat(3333333.3333333333),
				Low:     decimal.Zero,
				Vol:     decimal.Zero,
				Amount:  decimal.Zero,
				Count:   0,
				IsClose: false,
				Offset:  0,
			})
			fmt.Println(string(dd))
		}
	}()

	var pending = make(map[uint64]*src.MarketLine, 10)
	var ticker = time.Tick(time.Millisecond * 100)

	for {
		select {
		case data := <-ch:
			pending[data.ID] = data
		case <-ticker:
			for _, data := range pending {
				pending[data.ID] = nil
				delete(pending, data.ID)
				fmt.Println(data)
			}
		}
	}
}

//evalsha c76866510ad8135b2e990a1f9a4d5cd8203f5566 3 xxx btcusdt 9: 1 9
//evalsha c76866510ad8135b2e990a1f9a4d5cd8203f5566 3 xxx btcusdt 10: 2 10
//evalsha c76866510ad8135b2e990a1f9a4d5cd8203f5566 3 xxx btcusdt 20: 2 20
//evalsha c76866510ad8135b2e990a1f9a4d5cd8203f5566 3 xxx btcusdt 100: 2 100
//evalsha c76866510ad8135b2e990a1f9a4d5cd8203f5566 3 xxx btcusdt 199: 2 199
//evalsha c76866510ad8135b2e990a1f9a4d5cd8203f5566 3 xxx btcusdt 200: 2 200
//evalsha c76866510ad8135b2e990a1f9a4d5cd8203f5566 3 xxx btcusdt 201: 1 201
//evalsha c76866510ad8135b2e990a1f9a4d5cd8203f5566 3 xxx btcusdt 300: 1 300
//evalsha c76866510ad8135b2e990a1f9a4d5cd8203f5566 3 xxx btcusdt 1300: 2 1300

//evalsha 596c46727b8c2265cb7aee8d05acc7439013909b 3 xxx btcusdt 9: 2 9
//evalsha 596c46727b8c2265cb7aee8d05acc7439013909b 3 xxx btcusdt 10: 2 10
//evalsha 596c46727b8c2265cb7aee8d05acc7439013909b 3 xxx btcusdt 20: 2 20
//evalsha 596c46727b8c2265cb7aee8d05acc7439013909b 3 xxx btcusdt 100: 2 100
//evalsha 596c46727b8c2265cb7aee8d05acc7439013909b 3 xxx btcusdt 199: 2 199
//evalsha 596c46727b8c2265cb7aee8d05acc7439013909b 3 xxx btcusdt 200: 2 200
//evalsha 596c46727b8c2265cb7aee8d05acc7439013909b 3 xxx btcusdt 201: 1 201
//evalsha 596c46727b8c2265cb7aee8d05acc7439013909b 3 xxx btcusdt 300: 1 300
//evalsha 596c46727b8c2265cb7aee8d05acc7439013909b 3 xxx btcusdt 1300: 1 1300
