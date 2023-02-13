package src

type PersistentTrade struct {
	input   chan *MarketTrade
	outputs []func(data interface{})
}

func NewPersistentTrade(symbol string) *PersistentTrade {
	return &PersistentTrade{}
}

func (this *PersistentTrade) Start() error {
	return nil
}

func (this *PersistentTrade) Stop() error {
	return nil
}

func (this *PersistentTrade) Input(data interface{}) {
	for 50 < len(this.input) {
		<-this.input
	}
	this.input <- data.(*MarketTrade)
}

func (this *PersistentTrade) Output(output func(data interface{})) {
}

func (this *PersistentTrade) next(data interface{}) {
	var output func(data interface{})
	for _, output = range this.outputs {
		output(data)
	}
}
