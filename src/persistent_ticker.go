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
		for 1 < len(this.input) {
			<-this.input
		}
		SyncMarketTicker(this.symbol, <-this.input)
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
