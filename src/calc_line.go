package src

type CalcLine struct {
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

}

func (this *CalcLine) Output(output func(data interface{})) {
}
