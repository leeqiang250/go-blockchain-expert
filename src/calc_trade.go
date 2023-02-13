package src

type CalcTrade struct {
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

}

func (this *CalcTrade) Output(output func(data interface{})) {

}
