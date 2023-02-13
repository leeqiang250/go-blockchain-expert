package src

type CalcTrade struct {
	input   chan *MatchResult
	outputs []func(data interface{})
}

func NewCalcTrade(symbol string, marketTrade *MarketTrade) *CalcTrade {
	return &CalcTrade{}
}

func (this *CalcTrade) Start() error {
	return nil
}

func (this *CalcTrade) Stop() error {
	return nil
}

func (this *CalcTrade) Input(data interface{}) {
	this.input <- data.(*MatchResult)
}

func (this *CalcTrade) Output(output func(data interface{})) {
	this.outputs = append(this.outputs, output)
}

func (this *CalcTrade) next(data interface{}) {
	var output func(data interface{})
	for _, output = range this.outputs {
		output(data)
	}
}
