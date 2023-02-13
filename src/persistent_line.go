package src

type PersistentLine struct {
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

}

func (this *PersistentLine) Output(output func(data interface{})) {

}
