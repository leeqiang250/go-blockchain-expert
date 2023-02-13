package src

type SerializeLine struct {
	input   chan *MarketLine
	outputs []func(data interface{})
}

func NewSerializeLine(symbol string, marketLine *MarketLine) *SerializeLine {
	return &SerializeLine{}
}

func (this *SerializeLine) Start() error {
	for {
		var marketLine = <-this.input
		if !isExpire(marketLine.Ts) {

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
