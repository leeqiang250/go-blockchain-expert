package src

type SerializeTrade struct {
	input   chan *MarketTrade
	outputs []func(data interface{})
}

func NewSerializeTrade(symbol string) *SerializeTrade {
	return &SerializeTrade{}
}

func (this *SerializeTrade) Start() error {
	for {
		var marketTrade = <-this.input
		if !isExpire(marketTrade.Ts) {

		}
	}
}

func (this *SerializeTrade) Stop() error {
	return nil
}

func (this *SerializeTrade) Input(data interface{}) {
	this.input <- data.(*MarketTrade)
}

func (this *SerializeTrade) Output(output func(data interface{})) {
	this.outputs = append(this.outputs, output)
}

func (this *SerializeTrade) next(data interface{}) {
	var output func(data interface{})
	for _, output = range this.outputs {
		output(data)
	}
}
