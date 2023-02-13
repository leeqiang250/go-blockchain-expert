package src

type CalcTicker struct {
	marketTicker *MarketTicker
	input        chan *MatchResult
	outputs      []func(data interface{})
}

func NewCalcTicker(symbol string, marketTicker *MarketTicker) *CalcTicker {
	return &CalcTicker{}
}

func (this *CalcTicker) Start() error {

	this.next(this.marketTicker.DeepCopy())

	return nil
}

func (this *CalcTicker) Stop() error {
	return nil
}

func (this *CalcTicker) Input(data interface{}) {
	this.input <- data.(*MatchResult)
}

func (this *CalcTicker) Output(output func(data interface{})) {
	this.outputs = append(this.outputs, output)
}

func (this *CalcTicker) next(data interface{}) {
	var output func(data interface{})
	for _, output = range this.outputs {
		output(data)
	}
}
