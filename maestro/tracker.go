//Package maestro 数据平台接入sdk
package maestro

import (
	"encoding/json"
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	tracker *zap.Logger
)

const (
	//PATH 默认track路径
	PATH = "/var/log/maestro/track.json"
)

func fillBaseFields(base *BaseFields) []zap.Field {
	var fields []zap.Field
	fields = append(fields, zap.String("_id", fmt.Sprintf("%v", uuid.NewV1())))
	fields = append(fields, zap.String("_event_id", base.EventID))
	fields = append(fields, zap.Uint64("_user_id", base.UserID))
	if base.IsRobot {
		fields = append(fields, zap.Int32("user_type", 1))
	} else {
		fields = append(fields, zap.Int32("user_type", 0))
	}
	fields = append(fields, zap.String("sub_type", base.SubType))
	fields = append(fields, zap.String("third_type", base.ThirdType))
	fields = append(fields, zap.String("limit", base.Limit))
	fields = append(fields, zap.String("seq_no", base.SeqNO))
	if base.OutTransID != "" {
		fields = append(fields, zap.String("out_trans_id", base.OutTransID))
	}
	return fields
}

// Track 自定义事件
func Track(base *BaseFields, fields ...zap.Field) {
	for _, f := range fillBaseFields(base) {
		fields = append(fields, f)
	}
	tracker.Info("", fields...)
}

// TrackOrder 订单结算事件, 所有游戏都需要, 该事件在每一局游戏只能发送一次
func TrackOrder(base *BaseFields, betAmount float64, awardAmount float64, taxAmount float64) {
	var fields []zap.Field
	base.EventID = "bet_order"
	fields = fillBaseFields(base)
	fields = append(fields, zap.Float64("bet_amount", betAmount))
	fields = append(fields, zap.Float64("award_amount", awardAmount))
	fields = append(fields, zap.Float64("tax_amount", taxAmount))
	tracker.Info("", fields...)
}

// TrackTurnBet 轮次对战场下注行为事件
func TrackTurnBet(base *BaseFields, turn int32, coin float64, coinNum int32, amount float64) {
	var fields []zap.Field
	base.EventID = "turn_bet"
	fields = fillBaseFields(base)
	fields = append(fields, zap.Int32("turn", turn))
	fields = append(fields, zap.Float64("coin", coin))
	fields = append(fields, zap.Int32("coin_num", coinNum))
	fields = append(fields, zap.Float64("amount", amount))
	tracker.Info("", fields...)
}

// TrackHundredsBet 百人场下注行为事件
func TrackHundredsBet(base *BaseFields, region int32, coin float64, coinNum int32, amount float64) {
	var fields []zap.Field
	base.EventID = "hundreds_bet"
	fields = fillBaseFields(base)
	fields = append(fields, zap.Int32("region", region))
	fields = append(fields, zap.Float64("coin", coin))
	fields = append(fields, zap.Int32("coin_num", coinNum))
	fields = append(fields, zap.Float64("amount", amount))
	tracker.Info("", fields...)
}

// TrackCattleBet 血拼牛牛下注行为事件
func TrackCattleBet(base *BaseFields, multiple float64, baseAmount float64, bankerMultiple float64) {
	var fields []zap.Field
	base.EventID = "cattle_bet"
	fields = fillBaseFields(base)
	fields = append(fields, zap.Float64("multiple", multiple))
	fields = append(fields, zap.Float64("base_bet_amount", baseAmount))
	fields = append(fields, zap.Float64("banker_multiple", bankerMultiple))
	tracker.Info("", fields...)
}

// TrackFishHit 捕鱼命中行为事件
func TrackFishHit(base *BaseFields, fishMap map[int32]int32) {
	var fields []zap.Field
	base.EventID = "fishing_hit"
	fields = fillBaseFields(base)
	s, err := json.Marshal(fishMap)
	if err != nil {
		return
	}
	fields = append(fields, zap.String("fish_types", string(s)))
	tracker.Info("", fields...)
}

// TrackLordBet 斗地主倍数行为事件
func TrackLordBet(base *BaseFields, multiple float64, baseAmount float64, bankerMultiple float64, bombMultiple float64) {
	var fields []zap.Field
	base.EventID = "landlord_bet"
	fields = fillBaseFields(base)
	fields = append(fields, zap.Float64("multiple", multiple))
	fields = append(fields, zap.Float64("baseAmount", baseAmount))
	fields = append(fields, zap.Float64("bankerMultiple", bankerMultiple))
	fields = append(fields, zap.Float64("bombMultiple", bombMultiple))
	tracker.Info("", fields...)
}

// TrackTax 税收事件
func TrackTax(userID uint64, gameID string, projectID string, tax float64) {
	var fields []zap.Field
	fields = append(fields, zap.String("_id", fmt.Sprintf("%v", uuid.NewV1())))
	fields = append(fields, zap.String("_event_id", "tax"))
	fields = append(fields, zap.Uint64("_user_id", userID))
	fields = append(fields, zap.String("game_id", gameID))
	fields = append(fields, zap.String("project_id", projectID))
	fields = append(fields, zap.Float64("tax", tax))
	tracker.Info("", fields...)
}

func maestroTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendInt64(t.Unix())
}

// InitTracker 初始化
func InitTracker(path string) {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "_event_time"
	encoderCfg.EncodeTime = maestroTimeEncoder
	if path == "" {
		path = PATH
	}
	rawJSON := []byte(fmt.Sprintf(`{
		  "level": "info",
		  "encoding": "json",
		  "outputPaths": ["%s"],
		  "encoderConfig": {
			"levelEncoder": "uppercase"
		  },
		  "disableCaller": true
		}`, path))
	var cfg zap.Config
	var err error
	if err = json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}
	cfg.EncoderConfig = encoderCfg
	tracker, _ = cfg.Build()
}

//CloseTracker 关闭服务器之前调用，同步缓冲区
func CloseTracker() {
	_ = tracker.Sync()
}
