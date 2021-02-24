package mg

import (
	"github.com/globalsign/mgo/bson"
	"time"
)

//mongo官方库暂不稳定，使用mgo

type creditOP struct {
	Count  int32     `bson:"count"`
	Total  int64     `bson:"total"`
	LastAt time.Time `bson:"last_at"`
}

type UserStats struct {
	ID         uint64    `bson:"_id"`
	RegisterAt time.Time `bson:"register_at"`
	Channel    string    `bson:"chn"`
	IP         string    `bson:"ip"`
	AID        string    `bson:"aid"`
	TotalGain  int64     `bson:"total_gain"`
	Recharge   creditOP  `bson:"recharge"`
	Withdraw   creditOP  `bson:"withdraw"`
}

type DailyStats struct {
	ID        uint64   `bson:"_id"`
	Channel   string   `bson:"chn"`
	IP        string   `bson:"ip"`
	AID       string   `bson:"aid"`
	TotalGain int64    `bson:"total_gain"`
	Recharge  creditOP `bson:"recharge"`
	Withdraw  creditOP `bson:"withdraw"`
}

type IMTemplate struct {
	ID         bson.ObjectId `bson:"_id"`
	Type       []string      `bson:"type"`
	Template   string        `bson:"template"`
	MinAmount  int64         `bson:"min_amount"`
	DisablePKG []string      `bson:"disable_pkg"`
	DisableCHN []string      `bson:"disable_chn"`
}
