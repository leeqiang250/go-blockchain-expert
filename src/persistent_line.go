package src

type PersistentLine struct {
	input   chan *MarketLine
	outputs []func(data interface{})
}

func NewPersistentLine(symbol string, lineScale LineScale, marketLine *MarketLine) *PersistentLine {
	return &PersistentLine{}
}

func (this *PersistentLine) Start() error {
	return nil
}

func (this *PersistentLine) Stop() error {
	return nil
}

func (this *PersistentLine) Input(data interface{}) {
	this.input <- data.(*MarketLine)
}

func (this *PersistentLine) Output(output func(data interface{})) {
}

func (this *PersistentLine) next(data interface{}) {
	var output func(data interface{})
	for _, output = range this.outputs {
		output(data)
	}
}
