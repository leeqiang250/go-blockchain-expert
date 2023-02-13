package src

type CalcTicker struct {
}

func NewCalcTicker(symbol string, marketTicker *MarketTicker) *CalcTicker {
	return &CalcTicker{}
}

func (this *CalcTicker) Start() error {
	return nil
}

func (this *CalcTicker) Stop() error {
	return nil
}

func (this *CalcTicker) Input(data interface{}) {

}

func (this *CalcTicker) Output(output func(data interface{})) {

}
