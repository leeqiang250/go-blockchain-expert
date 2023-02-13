package src

type PersistentTrade struct {
	symbol  string
	input   chan *MarketTrade
	outputs []func(data interface{})
}

func NewPersistentTrade(symbol string) *PersistentTrade {
	return &PersistentTrade{}
}

func (this *PersistentTrade) Start() error {
	for {
		for 50 < len(this.input) {
			<-this.input
		}
		for {
			SyncMarketTrade(this.symbol, <-this.input)
		}
	}
}

func (this *PersistentTrade) Stop() error {
	return nil
}

func (this *PersistentTrade) Input(data interface{}) {
	this.input <- data.(*MarketTrade)
}

func (this *PersistentTrade) Output(output func(data interface{})) {
}
