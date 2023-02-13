package src

type CalcLine struct {
	input   chan *MatchResult
	outputs []func(data interface{})
}

func NewCalcLine(symbol string, lineScale LineScale, marketLine *MarketLine) *CalcLine {
	return &CalcLine{}
}

func (this *CalcLine) Start() error {
	return nil
}

func (this *CalcLine) Stop() error {
	return nil
}

func (this *CalcLine) Input(data interface{}) {
	this.input <- data.(*MatchResult)
}

func (this *CalcLine) Output(output func(data interface{})) {
	this.outputs = append(this.outputs, output)
}

func (this *CalcLine) next(data interface{}) {
	var output func(data interface{})
	for _, output = range this.outputs {
		output(data)
	}
}
