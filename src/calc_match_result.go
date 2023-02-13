package src

import (
	"encoding/json"
)

type CalcMatchResult struct {
	input   chan []byte
	outputs []func(data interface{})
}

func NewCalcMatchResult(symbol string) *CalcMatchResult {
	return &CalcMatchResult{}
}

func (this *CalcMatchResult) Start() error {
	for {
		var matchResult = &MatchResult{}
		var err = json.Unmarshal(<-this.input, matchResult)
		if nil != err {
			this.next(matchResult)
		}
	}
}

func (this *CalcMatchResult) Stop() error {
	return nil
}

func (this *CalcMatchResult) Input(data interface{}) {
	this.input <- data.([]byte)
}

func (this *CalcMatchResult) Output(output func(data interface{})) {
	this.outputs = append(this.outputs, output)
}

func (this *CalcMatchResult) next(data interface{}) {
	var output func(data interface{})
	for _, output = range this.outputs {
		output(data)
	}
}
