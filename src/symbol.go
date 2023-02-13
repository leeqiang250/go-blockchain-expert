package src

import "math"

func InitSymbol(symbol string, marketTrade *MarketTrade, marketTicker *MarketTicker, marketLines map[LineScale]*MarketLine) {
	var offset = uint64(math.MaxUint16)

	if nil != marketTrade {
		offset = Min(offset, marketTrade.Offset)
	}
	if nil != marketTicker {
		offset = Min(offset, marketTrade.Offset)
	}
	for _, marketLine := range marketLines {
		offset = Min(offset, marketLine.Offset)
	}
	if offset == uint64(math.MaxUint16) {
		offset = 0
	}

	var persistentTrade = NewPersistentTrade(symbol)
	go persistentTrade.Start()

	var calcTrade = NewCalcTrade(symbol, marketTrade.DeepCopy())
	calcTrade.Output(persistentTrade.Input)
	go calcTrade.Start()

	var persistentTicker = NewPersistentTicker(symbol)
	go persistentTicker.Start()

	var calcTicker = NewCalcTicker(symbol, marketTicker.DeepCopy())
	calcTicker.Output(persistentTicker.Input)
	go calcTicker.Start()

	var match = NewSerializeMatchResult(symbol)
	match.Output(calcTrade.Input)
	match.Output(calcTicker.Input)
	for lineScale, marketLine := range marketLines {
		var persistentLine = NewPersistentLine(symbol, lineScale, marketLine.DeepCopy())
		go persistentLine.Start()

		var calcLine = NewCalcLine(symbol, lineScale, marketLine.DeepCopy())
		calcLine.Output(persistentLine.Input)
		go calcLine.Start()

		match.Output(calcLine.Input)
	}
	go match.Start()

	var receiveMsg = NewReceiveMsg(symbol, offset)
	receiveMsg.Output(match.Input)
	go receiveMsg.Start()
}

func InitTickers() {

}
