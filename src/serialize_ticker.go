package src

import "time"

type SerializeTicker struct {
	marketTicker *MarketTicker
	lastTs       uint64
	interval     uint64
	input        chan *MarketTicker
	outputs      []func(data interface{})
}

func NewSerializeTicker(symbol string, marketLine *MarketLine) *SerializeTicker {
	return &SerializeTicker{}
}

func (this *SerializeTicker) Start() error {
	var tick = time.Tick(time.Duration(this.interval) * time.Millisecond)

	for {
		select {
		case this.marketTicker = <-this.input:
			this.process()
		case <-tick:
			this.process()
		}
	}
}

func (this *SerializeTicker) Stop() error {
	return nil
}

func (this *SerializeTicker) Input(data interface{}) {
	this.input <- data.(*MarketTicker)
}

func (this *SerializeTicker) Output(output func(data interface{})) {
	this.outputs = append(this.outputs, output)
}

func (this *SerializeTicker) process() {
	if nil == this.marketTicker || isExpire(this.marketTicker.Ts) {
		return
	}

	var ts = Now()
	if this.lastTs < ts {
		this.lastTs = ts + this.interval

		//TODO

		this.marketTicker = nil
	}
}
