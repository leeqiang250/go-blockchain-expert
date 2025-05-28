package src

import (
	"context"
	"encoding/json"
	"github.com/panjf2000/ants/v2"
	"strconv"
	"sync"
)

type CalcTickers struct {
	wg           sync.WaitGroup
	pool         *ants.PoolWithFunc
	marketTicker *MarketTicker
	input        chan *MatchResult
	outputs      []func(data interface{})
}

func NewCalcTickers(symbol string, marketTicker *MarketTicker) *CalcTickers {
	var calcTickers = CalcTickers{}

	calcTickers.pool, _ = ants.NewPoolWithFunc(10, func(i interface{}) {
		//calcTickers.myFunc(i)
		calcTickers.wg.Done()
	})

	return &calcTickers
}

func (this *CalcTickers) Start() error {
	//var tick = time.Tick(time.Duration(3) * time.Millisecond)

	return nil
}

func (this *CalcTickers) Stop() error {
	return nil
}

func (this *CalcTickers) Input(data interface{}) {
	this.input <- data.(*MatchResult)
}

func (this *CalcTickers) Output(output func(data interface{})) {
	this.outputs = append(this.outputs, output)
}

func (this *CalcTickers) process() {
	var err error
	var seq string
	var symbol string
	var redisTicker string
	var marketTicker *MarketTicker
	var redisTickers map[string]string
	var redisTickerSeqs map[string]string
	redisTickers, err = client.HGetAll(context.Background(), MARKET_TICKER).Result()
	if nil != err {
		return
	}

	redisTickerSeqs, err = client.HGetAll(context.Background(), MARKET_TICKER_PUSH_SEQ).Result()
	if nil != err {
		return
	}

	var marketTickers = make(map[string]*MarketTicker, len(redisTickers))

	for symbol, redisTicker = range redisTickers {
		this.wg.Add(1)
		_ = this.pool.Invoke(redisTicker)
	}

	this.wg.Wait()

	var marketTickerSeqs = make(map[string]uint64, len(redisTickerSeqs))
	for symbol, seq = range redisTickerSeqs {
		marketTickerSeqs[symbol], _ = strconv.ParseUint(seq, 10, 64)
	}

	for symbol, marketTicker = range marketTickers {
		if marketTickerSeqs[symbol] < marketTicker.Seq {
			marketTickerSeqs[symbol] = marketTicker.Seq
		} else {
			marketTickers[symbol] = nil
			delete(marketTickers, symbol)
		}
	}

	if 0 < len(marketTickers) {
		this.next(marketTickers)
		client.HMSet(context.Background(), MARKET_TICKER_PUSH_SEQ, marketTickerSeqs)
	}

}

func (this *CalcTickers) next(data interface{}) {
	var output func(data interface{})
	for _, output = range this.outputs {
		output(data)
	}
}

func (this *CalcTickers) myFunc(redisTicker string) {
	var marketTicker = MarketTicker{}
	var err = json.Unmarshal([]byte(redisTicker), &marketTicker)
	if nil != err {
		return
	}

	this.wg.Done()
}
