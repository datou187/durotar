package maestro

//BaseFields 基本字段，所有track必须包含
type BaseFields struct {
	EventID    string //事件ID, 比如bet_order
	UserID     uint64 //用户id
	IsRobot    bool
	SubType    string //游戏代号,如zjh
	ThirdType  string //子玩法或房间类型(房间号)
	Limit      string //门槛, 0.01/1
	SeqNO      string //下注旗(局)号或唯一标志性某一局游戏id
	OutTransID string //订单号(optional)
}
