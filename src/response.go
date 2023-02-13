package src

type Response struct {
	Code    uint8       `json:"c"`  //响应码，
	Ts      uint64      `json:"t"`  //响应时间：消息生成时间
	Event   Event       `json:"e"`  //触发事件类型
	Channel string      `json:"ch"` //频道名称
	Data    interface{} `json:"d"`  //数据域
}
