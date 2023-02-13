package src

type SerializeLine struct {
	cache   map[uint64]*MarketLine
	input   chan *MarketLine
	outputs []func(data interface{})
}

func NewSerializeLine(symbol string, marketLine *MarketLine) *SerializeLine {
	return &SerializeLine{}
}

func (this *SerializeLine) Start() error {
	var marketLine *MarketLine
	for {
		for 1 < len(this.input) {
			marketLine = <-this.input
			if !isExpire(marketLine.Ts) {
				this.cache[marketLine.ID] = marketLine
			}
		}

		marketLine = <-this.input
		if !isExpire(marketLine.Ts) {
			this.cache[marketLine.ID] = marketLine
		}

		for _, marketLine = range this.cache {
			this.cache[marketLine.ID] = nil
			delete(this.cache, marketLine.ID)
		}
	}
}

func (this *SerializeLine) Stop() error {
	return nil
}

func (this *SerializeLine) Input(data interface{}) {
	this.input <- data.(*MarketLine)
}

func (this *SerializeLine) Output(output func(data interface{})) {
	this.outputs = append(this.outputs, output)
}

func (this *SerializeLine) next(data interface{}) {
	var output func(data interface{})
	for _, output = range this.outputs {
		output(data)
	}
}
