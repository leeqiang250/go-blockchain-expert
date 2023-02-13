package src

type PersistentTicker struct {
}

func NewPersistentTicker(symbol string) *PersistentTicker {
	return &PersistentTicker{}
}

func (this *PersistentTicker) Start() error {
	return nil
}

func (this *PersistentTicker) Stop() error {
	return nil
}

func (this *PersistentTicker) Input(data interface{}) {

}

func (this *PersistentTicker) Output(output func(data interface{})) {

}
