package src

type PersistentTrade struct {
}

func NewPersistentTrade(symbol string) *PersistentTrade {
	return &PersistentTrade{}
}

func (this *PersistentTrade) Start() error {
	return nil
}

func (this *PersistentTrade) Stop() error {
	return nil
}

func (this *PersistentTrade) Input(data interface{}) {

}

func (this *PersistentTrade) Output(output func(data interface{})) {

}
