package src

type PersistentLine struct {
	symbol  string
	input   chan *MarketLine
	outputs []func(data interface{})
}

func NewPersistentLine(symbol string, lineScale LineScale, marketLine *MarketLine) *PersistentLine {
	return &PersistentLine{}
}

func (this *PersistentLine) Start() error {
	var marketLine *MarketLine
	var cache = make(map[uint64]*MarketLine, 2)
	for {
		for 1 < len(this.input) {
			marketLine = <-this.input
			cache[marketLine.ID] = marketLine
		}

		marketLine = <-this.input
		cache[marketLine.ID] = marketLine

		for _, marketLine = range cache {
			SyncMarketLine(this.symbol, marketLine)
			cache[marketLine.ID] = nil
			delete(cache, marketLine.ID)
		}
	}
}

func (this *PersistentLine) Stop() error {
	return nil
}

func (this *PersistentLine) Input(data interface{}) {
	this.input <- data.(*MarketLine)
}

func (this *PersistentLine) Output(output func(data interface{})) {
}
