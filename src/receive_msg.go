package src

type ReceiveMsg struct {
}

func NewReceiveMsg(symbol string, offset uint64) *ReceiveMsg {
	return &ReceiveMsg{}
}

func (this *ReceiveMsg) Start() error {
	return nil
}

func (this *ReceiveMsg) Stop() error {
	return nil
}

func (this *ReceiveMsg) Input(data interface{}) {

}

func (this *ReceiveMsg) Output(output func(data interface{})) {

}
