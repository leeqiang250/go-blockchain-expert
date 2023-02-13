package src

import (
	"encoding/json"
)

type SerializeMatchResult struct {
	input   chan []byte
	outputs []func(data interface{})
}

func NewSerializeMatchResult(symbol string) *SerializeMatchResult {
	return &SerializeMatchResult{}
}

func (this *SerializeMatchResult) Start() error {
	for {
		var matchResult = &MatchResult{}
		var err = json.Unmarshal(<-this.input, matchResult)
		if nil != err {
			this.next(matchResult)
		}
	}
}

func (this *SerializeMatchResult) Stop() error {
	return nil
}

func (this *SerializeMatchResult) Input(data interface{}) {
	this.input <- data.([]byte)
}

func (this *SerializeMatchResult) Output(output func(data interface{})) {
	this.outputs = append(this.outputs, output)
}

func (this *SerializeMatchResult) next(data interface{}) {
	var output func(data interface{})
	for _, output = range this.outputs {
		output(data)
	}
}
