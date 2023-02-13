package src

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"
	"gocloud.dev/pubsub"
	"log"
	"math"
	"time"
)

//CREATE TABLE `ex_snapshot_ticker` (
//`id` bigint unsigned NOT NULL COMMENT 'id',
//`symbol` varchar(100) NOT NULL COMMENT '交易对',
//`idx` bigint unsigned NOT NULL COMMENT '所属时间刻度的起始时间戳',
//`open` varchar(64) NOT NULL COMMENT '开盘价',
//`close` varchar(64) NOT NULL COMMENT '收盘价',
//`high` varchar(64) NOT NULL COMMENT '最高价',
//`low` varchar(64) NOT NULL COMMENT '最低价',
//`vol` varchar(64) NOT NULL COMMENT '成交额',
//`amount` varchar(64) NOT NULL COMMENT '成交量',
//`count` int unsigned NOT NULL COMMENT '交易数量',
//`price_change` varchar(100) NOT NULL COMMENT '涨跌额',
//`price_change_percent` varchar(100) NOT NULL COMMENT '涨跌幅(22.22%:0.2222，-22.22%:0.2222)',
//`is_open` tinyint unsigned NOT NULL COMMENT '是否开盘',
//`ts` bigint unsigned NOT NULL COMMENT 'Taker的时间戳',
//`offset` bigint unsigned NOT NULL COMMENT 'Offset',
//`seq` bigint unsigned NOT NULL COMMENT 'Seq',
//`updated_at` timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '修改时间',
//`created_at` timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
//PRIMARY KEY (`id`),
//UNIQUE KEY `unique_index_symbol` (`symbol`) USING BTREE
//) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Ticker快照';
//
//CREATE TABLE `ex_snapshot_trade` (
//`id` bigint unsigned NOT NULL COMMENT 'id',
//`symbol` varchar(100) NOT NULL COMMENT '交易对',
//`ts` bigint unsigned NOT NULL COMMENT 'Taker的时间戳',
//`offset` bigint unsigned NOT NULL COMMENT 'Offset',
//`seq` bigint unsigned NOT NULL COMMENT 'Seq',
//`updated_at` timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '修改时间',
//`created_at` timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
//PRIMARY KEY (`id`),
//UNIQUE KEY `unique_index_symbol` (`symbol`) USING BTREE
//) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='Trade快照';
//
//CREATE TABLE `ex_line_btcusdt` (
//`id` bigint unsigned NOT NULL COMMENT 'id',
//`idx` bigint unsigned NOT NULL COMMENT '所属时间刻度的起始时间戳',
//`open` varchar(64) NOT NULL COMMENT '开盘价',
//`close` varchar(64) NOT NULL COMMENT '收盘价',
//`high` varchar(64) NOT NULL COMMENT '最高价',
//`low` varchar(64) NOT NULL COMMENT '最低价',
//`vol` varchar(64) NOT NULL COMMENT '成交额',
//`amount` varchar(64) NOT NULL COMMENT '成交量',
//`count` int unsigned NOT NULL COMMENT '交易数量',
//`is_close` tinyint unsigned NOT NULL COMMENT '是否关闭',
//`offset` bigint unsigned NOT NULL COMMENT 'Offset',
//`seq` bigint unsigned NOT NULL COMMENT 'Seq',
//`scale_type` enum('1min','5min','15min','30min','1hour','4hour','12hour','1day','1week','1month') NOT NULL COMMENT '时间刻度',
//`updated_at` timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '修改时间',
//`created_at` timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
//PRIMARY KEY (`id`),
//UNIQUE KEY `unique_index_scale_type_idx` (`scale_type`,`idx`) USING BTREE
//) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='非1分钟K线';
//
//CREATE TABLE `ex_line_1min_btcusdt` (
//`id` bigint unsigned NOT NULL COMMENT 'id',
//`idx` bigint unsigned NOT NULL COMMENT '所属时间刻度的起始时间戳',
//`open` varchar(64) NOT NULL COMMENT '开盘价',
//`close` varchar(64) NOT NULL COMMENT '收盘价',
//`high` varchar(64) NOT NULL COMMENT '最高价',
//`low` varchar(64) NOT NULL COMMENT '最低价',
//`vol` varchar(64) NOT NULL COMMENT '成交额',
//`amount` varchar(64) NOT NULL COMMENT '成交量',
//`count` int unsigned NOT NULL COMMENT '交易数量',
//`is_close` tinyint unsigned NOT NULL COMMENT '是否关闭',
//`offset` bigint unsigned NOT NULL COMMENT 'Offset',
//`seq` bigint unsigned NOT NULL COMMENT 'Seq',
//`updated_at` timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '修改时间',
//`created_at` timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
//PRIMARY KEY (`id`),
//UNIQUE KEY `unique_index_idx` (`idx`) USING BTREE
//) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='1分钟K线';

type LineScale string

const (
	MIN_1   LineScale = "1min"
	MIN_5   LineScale = "5min"
	MIN_15  LineScale = "15min"
	MIN_30  LineScale = "30min"
	HOUR_1  LineScale = "1hour"
	HOUR_4  LineScale = "4hour"
	HOUR12  LineScale = "12hour"
	DAY_1   LineScale = "1day"
	WEEK_1  LineScale = "1week"
	MONTH_1 LineScale = "1month"
)

type Event string

const (
	SUB Event = "sub" //触发事件类型-推送
	REQ Event = "req" //触发事件类型-请求
)

type Side string

const (
	BUY  Side = "buy"  //成交方向-买
	SELL Side = "sell" //成交方向-卖
)

type MarketTicker struct {
	ID                 uint64          `json:"i"`   //所属时间刻度的起始时间戳，UTC 0
	Open               decimal.Decimal `json:"o"`   //开盘价
	Close              decimal.Decimal `json:"c"`   //收盘价
	High               decimal.Decimal `json:"h"`   //最高价
	Low                decimal.Decimal `json:"l"`   //最低价
	Vol                decimal.Decimal `json:"v"`   //24H成交额
	Amount             decimal.Decimal `json:"a"`   //24H成交量
	Count              uint32          `json:"co"`  //交易数量
	PriceChange        decimal.Decimal `json:"pc"`  //涨跌额
	PriceChangePercent decimal.Decimal `json:"pcp"` //涨跌幅(22.22%:0.2222, -22.22%:0.2222)
	IsOpen             bool            `json:"is"`  //是否开盘
	Ts                 uint64          `json:"t"`   //Taker的时间戳，UTC 0
	Offset             uint64          `json:"of"`  //Offset
	Seq                uint64          `json:"s"`   //数据版本-Taker的SeqId
}

type MarketLine struct {
	ID      uint64          `json:"i"`  //所属时间刻度的起始时间戳，UTC 0
	Open    decimal.Decimal `json:"o"`  //开盘价
	Close   decimal.Decimal `json:"c"`  //收盘价
	High    decimal.Decimal `json:"h"`  //最高价
	Low     decimal.Decimal `json:"l"`  //最低价
	Vol     decimal.Decimal `json:"v"`  //成交额
	Amount  decimal.Decimal `json:"a"`  //成交量
	Count   uint32          `json:"co"` //交易数量
	IsClose bool            `json:"is"` //是否关闭
	Ts      uint64          `json:"t"`  //Taker的时间戳，UTC 0
	Offset  uint64          `json:"of"` //Offset
	Seq     uint64          `json:"s"`  //数据版本-Taker的SeqId
}

type MarketTrade struct {
	Side   Side            `json:"si"` //Taker的交易方向
	Price  decimal.Decimal `json:"p"`  //Maker的成交价格
	Vol    decimal.Decimal `json:"v"`  //Maker的成交额
	Amount decimal.Decimal `json:"a"`  //Maker的成交量
	Ts     uint64          `json:"t"`  //Taker的时间戳，UTC 0
	Offset uint64          `json:"of"` //Offset
	Seq    uint64          `json:"s"`  //数据版本-Taker的SeqId
}

var client *redis.Client

func isExpire(ts uint64) bool {
	return false
}

func (this *MarketTicker) DeepCopy() *MarketTicker {
	return nil
}

func (this *MarketTrade) DeepCopy() *MarketTrade {
	return nil
}

func (this *MarketLine) DeepCopy() *MarketLine {
	return &MarketLine{
		ID:      this.ID,
		Open:    this.Open,
		Close:   this.Close,
		High:    this.High,
		Low:     this.Low,
		Vol:     this.Vol,
		Amount:  this.Amount,
		Count:   this.Count,
		IsClose: this.IsClose,
		Ts:      this.Ts,
		Offset:  this.Offset,
		Seq:     this.Seq,
	}
}
func Test() {
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:32768",
		Username: "default",
		Password: "redispw", // no password set
		DB:       0,         // use default DB
	})
	var i = uint64(0)
	for i < 100 {
		i++

		//SyncTrade("btcusdt", &MarketTrade{
		//	Side:   BUY,
		//	Price:  decimal.NewFromInt(int64(i)),
		//	Vol:    decimal.NewFromInt(int64(i)),
		//	Amount: decimal.NewFromInt(int64(i)),
		//	Ts:     uint64(time.Now().UnixMilli()) / 1000,
		//	Offset: i,
		//	Seq:    i,
		//})
		//
		//SyncTrade("ethusdt", &MarketTrade{
		//	Side:   BUY,
		//	Price:  decimal.NewFromInt(int64(i)),
		//	Vol:    decimal.NewFromInt(int64(i)),
		//	Amount: decimal.NewFromInt(int64(i)),
		//	Ts:     uint64(time.Now().UnixMilli()) / 1000,
		//	Offset: i,
		//	Seq:    i,
		//})

		//SyncTicker("btcusdt", &MarketTicker{
		//	ID:                 i,
		//	Open:               decimal.NewFromInt(int64(i)),
		//	Close:              decimal.NewFromInt(int64(i)),
		//	High:               decimal.NewFromInt(int64(i)),
		//	Low:                decimal.NewFromInt(int64(i)),
		//	Vol:                decimal.NewFromInt(int64(i)),
		//	Amount:             decimal.NewFromInt(int64(i)),
		//	Count:              uint(i),
		//	PriceChange:        decimal.NewFromInt(int64(i)),
		//	PriceChangePercent: decimal.NewFromInt(int64(i)),
		//	Ts:                 uint64(time.Now().UnixMilli()) / 1000,
		//	Offset:             i,
		//	Seq:                i,
		//})
		//
		//SyncTicker("ethusdt", &MarketTicker{
		//	ID:                 i,
		//	Open:               decimal.NewFromInt(int64(i)),
		//	Close:              decimal.NewFromInt(int64(i)),
		//	High:               decimal.NewFromInt(int64(i)),
		//	Low:                decimal.NewFromInt(int64(i)),
		//	Vol:                decimal.NewFromInt(int64(i)),
		//	Amount:             decimal.NewFromInt(int64(i)),
		//	Count:              uint(i),
		//	PriceChange:        decimal.NewFromInt(int64(i)),
		//	PriceChangePercent: decimal.NewFromInt(int64(i)),
		//	Ts:                 uint64(time.Now().UnixMilli()) / 1000,
		//	Offset:             i,
		//	Seq:                i,
		//})

		fmt.Println(math.MaxUint16)
		fmt.Println(math.MaxUint32)

		SyncMarketLine("btcusdt", &MarketLine{
			ID:      (uint64(time.Now().UnixMilli()) / 6000) * 60,
			Open:    decimal.NewFromInt(int64(i)),
			Close:   decimal.NewFromInt(int64(i)),
			High:    decimal.NewFromInt(int64(i)),
			Low:     decimal.NewFromInt(int64(i)),
			Vol:     decimal.NewFromInt(int64(i)),
			Amount:  decimal.NewFromInt(int64(i)),
			Count:   uint32(i),
			IsClose: false,
			Ts:      uint64(time.Now().UnixMilli()),
			Offset:  i,
			Seq:     i,
		})

		SyncMarketLine("ethusdt", &MarketLine{
			ID:      (uint64(time.Now().UnixMilli()) / 6000) * 60,
			Open:    decimal.NewFromInt(int64(i)),
			Close:   decimal.NewFromInt(int64(i)),
			High:    decimal.NewFromInt(int64(i)),
			Low:     decimal.NewFromInt(int64(i)),
			Vol:     decimal.NewFromInt(int64(i)),
			Amount:  decimal.NewFromInt(int64(i)),
			Count:   uint32(i),
			IsClose: false,
			Ts:      uint64(time.Now().UnixMilli()),
			Offset:  i,
			Seq:     i,
		})
	}
}

//TODO 行情 UI 截图

const (
	MARKET_SYMBOL_SEQ string = "%s_seq"
)

const (
	MARKET_TICKER_SCRIPT   string = "local seq = redis.call('hget', KEYS[1], ARGV[1]) if (seq) then if (tonumber(ARGV[2]) > tonumber(seq)) then redis.call('hset', KEYS[1], ARGV[1], ARGV[2]) redis.call('hset', KEYS[1], ARGV[3], ARGV[4]) return 3 else return 2 end else redis.call('hset', KEYS[1], ARGV[1], ARGV[2]) redis.call('hset', KEYS[1], ARGV[3], ARGV[4]) return 1 end"
	MARKET_TICKER          string = "market::ticker"
	MARKET_TICKER_PUSH_SEQ string = "market::ticker::push::seq"
)

func SyncMarketTicker(symbol string, marketTicker *MarketTicker) {
	//TODO 只需要加载一次
	var sha1 = client.ScriptLoad(context.Background(), MARKET_TICKER_SCRIPT).Val()

	var data, err = json.Marshal(marketTicker)
	if nil != err {
		return
	}

	var cmd = client.EvalSha(
		context.Background(),
		sha1, []string{MARKET_TICKER},
		fmt.Sprintf(MARKET_SYMBOL_SEQ, symbol),
		marketTicker.Seq,
		symbol,
		data,
	)

	fmt.Println(cmd)
}

//local seq = redis.call('hget', KEYS[1], ARGV[1])
//if (seq)
//then
//if (tonumber(ARGV[2]) > tonumber(seq))
//then
//redis.call('hset', KEYS[1], ARGV[1], ARGV[2])
//redis.call('hset', KEYS[1], ARGV[3], ARGV[4])
//return 3
//else
//return 2
//end
//else
//redis.call('hset', KEYS[1], ARGV[1], ARGV[2])
//redis.call('hset', KEYS[1], ARGV[3], ARGV[4])
//return 1
//end
//
//操作返回值
//1:初始化数据
//2:低版本过期数据
//3:高版本有效数据

const (
	MARKET_LINE_SCRIPT string = "local seq = redis.call('hget', KEYS[1], ARGV[1]) if (seq) then if (tonumber(ARGV[2]) > tonumber(seq)) then redis.call('hset', KEYS[1], ARGV[1], ARGV[2]) redis.call('hset', KEYS[1], ARGV[3], ARGV[4]) return 3 else return 2 end else redis.call('hset', KEYS[1], ARGV[1], ARGV[2]) redis.call('hset', KEYS[1], ARGV[3], ARGV[4]) return 1 end"
	MARKET_LINE        string = "market::line::%s"
)

func SyncMarketLine(symbol string, marketLine *MarketLine) {
	//TODO 只需要加载一次
	var sha1 = client.ScriptLoad(context.Background(), MARKET_LINE_SCRIPT).Val()

	var data, err = json.Marshal(marketLine)
	if nil != err {
		return
	}

	client.EvalSha(
		context.Background(),
		sha1,
		[]string{fmt.Sprintf(MARKET_LINE, symbol)},
		fmt.Sprintf(MARKET_SYMBOL_SEQ, symbol),
		marketLine.Seq,
		marketLine.ID,
		data,
	)
}

//local seq = redis.call('hget', KEYS[1], ARGV[1])
//if (seq)
//then
//if (tonumber(ARGV[2]) > tonumber(seq))
//then
//redis.call('hset', KEYS[1], ARGV[1], ARGV[2])
//redis.call('hset', KEYS[1], ARGV[3], ARGV[4])
//return 3
//else
//return 2
//end
//else
//redis.call('hset', KEYS[1], ARGV[1], ARGV[2])
//redis.call('hset', KEYS[1], ARGV[3], ARGV[4])
//return 1
//end
//
//操作返回值
//1:初始化数据
//2:低版本过期数据
//3:高版本有效数据
//
//
//
//
//
//
//
//
//
//
//
//
//

//ticker 3秒一次
//盘口 1 秒一次

func xx() {
	// PRAGMA: This example is used on gocloud.dev; PRAGMA comments adjust how it is shown and can be ignored.
	// PRAGMA: On gocloud.dev, add a blank import: _ "gocloud.dev/pubsub/kafkapubsub"
	// PRAGMA: On gocloud.dev, hide lines until the next blank line.
	ctx := context.Background()

	// pubsub.OpenSubscription creates a *pubsub.Subscription from a URL.
	// The host + path are used as the consumer group name.
	// The "topic" query parameter sets one or more topics to subscribe to.
	// The set of brokers must be in an environment variable KAFKA_BROKERS.

	subscription, err := pubsub.OpenSubscription(ctx,
		"kafka://my-group?topic=my-topic")
	if err != nil {
		log.Fatal(err)
	}
	defer subscription.Shutdown(ctx)

	var msg, _ = subscription.Receive(ctx)
	fmt.Println(msg)
	msg.Ack()
}

func xxx() {
	// PRAGMA: This example is used on gocloud.dev; PRAGMA comments adjust how it is shown and can be ignored.
	// PRAGMA: On gocloud.dev, add a blank import: _ "gocloud.dev/pubsub/kafkapubsub"
	// PRAGMA: On gocloud.dev, hide lines until the next blank line.
	ctx := context.Background()

	// pubsub.OpenTopic creates a *pubsub.Topic from a URL.
	// The host + path are the topic name to send to.
	// The set of brokers must be in an environment variable KAFKA_BROKERS.
	topic, err := pubsub.OpenTopic(ctx, "kafka://my-topic")
	if err != nil {
		log.Fatal(err)
	}
	defer topic.Shutdown(ctx)
	topic.Send(ctx, &pubsub.Message{
		LoggableID: "",
		Body:       nil,
		Metadata:   nil,
		BeforeSend: nil,
		AfterSend:  nil,
	})
}

const (
	MARKET_TRADE_SCRIPT        = "local seq = redis.call('hget', KEYS[2], ARGV[1]) if (seq) then if (tonumber(ARGV[2]) > tonumber(seq)) then redis.call('lpush', KEYS[1], ARGV[3]) redis.call('hset', KEYS[2], ARGV[1], ARGV[2]) if (100 < redis.call('llen', KEYS[1])) then redis.call('rpop', KEYS[1]) end return 3 else return 2 end else redis.call('lpush', KEYS[1], ARGV[3]) redis.call('hset', KEYS[2], ARGV[1], ARGV[2]) if (100 < redis.call('llen', KEYS[1])) then redis.call('rpop', KEYS[1]) end return 1 end"
	MARKET_TRADE        string = "market::trade::%s"
	MARKET_TRADE_SEQ    string = "market::trade::seq"
)

func SyncMarketTrade(symbol string, marketTrade *MarketTrade) {
	//TODO 只需要加载一次
	var sha1 = client.ScriptLoad(context.Background(), MARKET_TRADE_SCRIPT).Val()

	var data, err = json.Marshal(marketTrade)
	if nil != err {
		return
	}

	var cmd = client.EvalSha(
		context.Background(),
		sha1,
		[]string{
			fmt.Sprintf(MARKET_TRADE, symbol),
			MARKET_TRADE_SEQ,
		},
		symbol,
		marketTrade.Seq,
		data,
	)

	fmt.Println(cmd)
}

//local seq = redis.call('hget', KEYS[2], ARGV[1])
//if (seq)
//then
//if (tonumber(ARGV[2]) > tonumber(seq))
//then
//redis.call('lpush', KEYS[1], ARGV[3])
//redis.call('hset', KEYS[2], ARGV[1], ARGV[2])
//if (100 < redis.call('llen', KEYS[1]))
//then
//redis.call('rpop', KEYS[1])
//end
//return 3
//else
//return 2
//end
//else
//redis.call('lpush', KEYS[1], ARGV[3])
//redis.call('hset', KEYS[2], ARGV[1], ARGV[2])
//if (100 < redis.call('llen', KEYS[1]))
//then
//redis.call('rpop', KEYS[1])
//end
//return 1
//end
//操作返回值
//1:初始化数据
//2:低版本过期数据
//3:高版本有效数据

func XXX() {
	var data = &MarketLine{
		ID:      0,
		Open:    decimal.Zero,
		Close:   decimal.Zero,
		High:    decimal.Zero,
		Low:     decimal.Zero,
		Vol:     decimal.Zero,
		Amount:  decimal.Zero,
		Count:   0,
		IsClose: false,
		Ts:      0,
		Offset:  0,
		Seq:     0,
	}
	var ch = make(chan *MarketLine, 100)

	go func() {
		for {
			data.Open = data.Open.Add(decimal.NewFromInt(1))
			//time.Sleep(time.Second)
			data.Close = data.Close.Add(decimal.NewFromInt(1))

			ch <- data.DeepCopy()
		}
	}()

	for {
		var aa, _ = json.Marshal(<-ch)
		var bb = &MarketLine{}
		json.Unmarshal(aa, bb)

		if bb.Open.Equal(bb.Close) {
			fmt.Println(string(aa))
		} else {
			fmt.Println(string(aa))
		}
	}

}

func Now() uint64 {
	return 2
}
