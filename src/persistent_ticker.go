package src

type PersistentTicker struct {
	symbol  string
	input   chan *MarketTicker
	outputs []func(data interface{})
}

func NewPersistentTicker(symbol string) *PersistentTicker {
	return &PersistentTicker{}
}

func (this *PersistentTicker) Start() error {
	for {
		var marketTicker *MarketTicker
		for 0 < len(this.input) {
			marketTicker = <-this.input
		}
		
		SyncTicker(this.symbol, marketTicker)
	}
}

func (this *PersistentTicker) Stop() error {
	return nil
}

func (this *PersistentTicker) Input(data interface{}) {
	this.input <- data.(*MarketTicker)
}

func (this *PersistentTicker) Output(output func(data interface{})) {
}

func (this *PersistentTicker) next(data interface{}) {
	var output func(data interface{})
	for _, output = range this.outputs {
		output(data)
	}
}
