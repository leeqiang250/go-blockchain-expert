package src

type PullMatchResult struct {
}

func NewPullMatchResult(symbol string, offset uint64) *PullMatchResult {
	return &PullMatchResult{}
}

func (this *PullMatchResult) Start() error {
	return nil
}

func (this *PullMatchResult) Stop() error {
	return nil
}

func (this *PullMatchResult) Input(data interface{}) {

}

func (this *PullMatchResult) Output(output func(data interface{})) {

}
