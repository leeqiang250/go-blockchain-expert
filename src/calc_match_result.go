package src

import (
	"encoding/json"
	"fmt"
)

type MatchResult struct {
	input   chan []byte
	outputs []func(data interface{})
}

func NewCalcMatchResult(symbol string) *MatchResult {
	return &MatchResult{}
}

func (this *MatchResult) Start() error {
	for {
		var data = <-this.input
		var dd = &struct {
		}{}
		var err = json.Unmarshal(data, dd)
		fmt.Println(err)
		this.next(dd)
	}
	return nil
}

func (this *MatchResult) Stop() error {
	return nil
}

func (this *MatchResult) Input(data interface{}) {
	this.input <- data.([]byte)
}

func (this *MatchResult) Output(output func(data interface{})) {
	this.outputs = append(this.outputs, output)
}

func (this *MatchResult) next(data interface{}) {
	var output func(data interface{})
	for _, output = range this.outputs {
		output(data)
	}
}
