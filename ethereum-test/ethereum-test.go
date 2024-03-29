package ethereum_test

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"go-blockchain-expert/src"
	"math"
	"math/big"
	"sync"
	"time"
)

func Start() {
	{
		var token = common.HexToAddress("0x9364e119AD76e0346126aFcbDF5C9f0189500Cc5")
		fmt.Println(token)
		token = common.HexToAddress("0x9364e119AD76e0346126aFcbDF5C9f0189500Cc6")
		fmt.Println(token)

	}
	step4()
	step3()
	step2()

	var e *big.Int
	var v []byte
	var a string

	fmt.Println(e)
	fmt.Println(v)
	fmt.Println(a)

	e = new(big.Int).SetInt64(258)
	v = e.Bytes()
	a = fmt.Sprintf("%032b", 9)

	amount := big.NewFloat(0.001)

	var err error
	var privateKey *ecdsa.PrivateKey
	var transactOpts *bind.TransactOpts
	privateKey, err = crypto.HexToECDSA("4f87bc2f074d87fefd948cc2017ec4021035fa8c49d6786e1f2be6be6a991a4e")
	transactOpts, err = bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(128))

	//fmt.Println("main")
	//work.InitRedis(nil, "nil")
	rpcDial, err := rpc.Dial("http://52.77.179.179:8588")
	if err != nil {
		panic(err)
	}
	client := ethclient.NewClient(rpcDial)
	fmt.Println(client)

	token, err := src.NewIERC20(common.HexToAddress("0x9364e119AD76e0346126aFcbDF5C9f0189500Cc5"), client)
	if err != nil {
		panic(err)
	}
	//每个代币都会有相应的位数，例如eos是18位，那么我们转账的时候，需要在金额后面加18个0
	//decimal, err := token.Decimals(nil)
	decimal := 2
	if err != nil {
		panic(err)
	}
	//这是处理位数的代码段
	tenDecimal := big.NewFloat(math.Pow(10, float64(decimal)))
	convertAmount, _ := new(big.Float).Mul(tenDecimal, amount).Int(&big.Int{})
	//然后就可以转账到你需要接受的账户上了
	//toAddress 是接受eos的账户地址
	txs, err := token.Transfer(transactOpts, common.HexToAddress("0xbF395508fB2409dbD6c5C2f7cF8824AC79017939"), convertAmount)
	fmt.Println(txs.Hash().Hex())
	fmt.Println(txs)
	txs, err = token.Transfer(transactOpts, common.HexToAddress("0xbF395508fB2409dbD6c5C2f7cF8824AC79017939"), convertAmount)
	fmt.Println(txs.Hash().Hex())
	fmt.Println(txs)
}

func step2() {
	e := make(chan int, 5)
	e <- 1
	e <- 2
	e <- 3

	fmt.Println(<-e)
	fmt.Println(<-e)
	close(e)
	fmt.Println(<-e)
	fmt.Println(<-e)
}

const (
	mutexLocked  = 1 << iota
	mutexLocked2 = iota
	mutexLocked3 = 1 << 1
	mutexLocked4 = iota
)

func step3() {
	fmt.Println(mutexLocked)
	fmt.Println(mutexLocked2)
	fmt.Println(mutexLocked3)
	fmt.Println(mutexLocked4)

	xx := sync.Mutex{}
	xx.Lock()
}

func step4() {
	i := int64(0)
	//amount := 1000000000
	amount := int64(1000000000)
	xx := sync.Mutex{}

	ts := time.Now().UnixMilli()
	for i < amount {
		xx.Lock()
		xx.Unlock()
		i++
	}

	ts = time.Now().UnixMilli() - ts
	fmt.Println((amount * int64(1000)) / ts)
}
