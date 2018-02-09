package nx_base

import (
	"encoding/gob"
	"time"
)

func init() {

	gob.Register(&NxAuction{})
	gob.Register([]NxAuction{})

	gob.Register(&NxChangename{})
	gob.Register([]NxChangename{})

	gob.Register(&NxChunks{})
	gob.Register([]NxChunks{})

	gob.Register(&NxConfig{})
	gob.Register([]NxConfig{})

	gob.Register(&NxDomains{})
	gob.Register([]NxDomains{})

	gob.Register(&NxDumpRank{})
	gob.Register([]NxDumpRank{})

	gob.Register(&NxGuilds{})
	gob.Register([]NxGuilds{})

	gob.Register(&NxLetter{})
	gob.Register([]NxLetter{})

	gob.Register(&NxPubdata{})
	gob.Register([]NxPubdata{})

	gob.Register(&NxResDumpinfo{})
	gob.Register([]NxResDumpinfo{})

	gob.Register(&NxRoleDumpinfo{})
	gob.Register([]NxRoleDumpinfo{})

	gob.Register(&NxRoles{})
	gob.Register([]NxRoles{})

	gob.Register(&NxRolesCrossing{})
	gob.Register([]NxRolesCrossing{})

	gob.Register(&NxSns{})
	gob.Register([]NxSns{})

	gob.Register(&NxSnsFeed{})
	gob.Register([]NxSnsFeed{})

	gob.Register(&NxSnsRelation{})
	gob.Register([]NxSnsRelation{})

	gob.Register(&NxSnslog{})
	gob.Register([]NxSnslog{})

	gob.Register(&NxTrade{})
	gob.Register([]NxTrade{})

	gob.Register(&Test{})
	gob.Register([]Test{})

}

type NxAuction struct {
	Name       string    `json:"name" xorm:"not null pk VARCHAR(128)"`
	Uid        string    `json:"uid" xorm:"not null VARCHAR(32)"`
	CreateTime time.Time `json:"create_time" xorm:"DATETIME"`
	Type       int       `json:"type" xorm:"INT(11)"`
	Record     string    `json:"record" xorm:"not null VARCHAR(128)"`
	Version    int       `json:"version" xorm:"INT(11)"`
	Amount     int       `json:"amount" xorm:"INT(11)"`
	Range      int       `json:"range" xorm:"INT(11)"`
	SaveTime   time.Time `json:"save_time" xorm:"DATETIME"`
	SaveData   []byte    `json:"save_data" xorm:"LONGBLOB"`
}

type NxChangename struct {
	Name     string    `json:"name" xorm:"not null pk VARCHAR(128)"`
	Uid      string    `json:"uid" xorm:"not null VARCHAR(32)"`
	NewName  string    `json:"new_name" xorm:"not null VARCHAR(128)"`
	SaveTime time.Time `json:"save_time" xorm:"DATETIME"`
}

type NxChunks struct {
	SerialNo   string    `json:"serial_no" xorm:"not null pk VARCHAR(64)"`
	Name       string    `json:"name" xorm:"index VARCHAR(128)"`
	ServerId   int       `json:"server_id" xorm:"index INT(11)"`
	SceneId    int       `json:"scene_id" xorm:"index INT(11)"`
	CreateTime time.Time `json:"create_time" xorm:"DATETIME"`
	Deleted    int       `json:"deleted" xorm:"INT(11)"`
	SaveTime   time.Time `json:"save_time" xorm:"DATETIME"`
	SaveData   []byte    `json:"save_data" xorm:"LONGBLOB"`
}

type NxConfig struct {
	Id         string `json:"id" xorm:"not null pk VARCHAR(128)"`
	Name       string `json:"name" xorm:"index VARCHAR(128)"`
	Type       int    `json:"type" xorm:"INT(11)"`
	Property   string `json:"property" xorm:"LONGTEXT"`
	Status     int    `json:"status" xorm:"INT(11)"`
	OpName     string `json:"op_name" xorm:"index VARCHAR(128)"`
	OpType     int    `json:"op_type" xorm:"INT(11)"`
	OpProperty string `json:"op_property" xorm:"LONGTEXT"`
	OpStatus   int    `json:"op_status" xorm:"INT(11)"`
	OpVersion  int    `json:"op_version" xorm:"INT(11)"`
	OpId       int    `json:"op_id" xorm:"INT(11)"`
}

type NxDomains struct {
	Name       string    `json:"name" xorm:"not null pk VARCHAR(128)"`
	Uid        string    `json:"uid" xorm:"not null unique VARCHAR(32)"`
	CreateTime time.Time `json:"create_time" xorm:"DATETIME"`
	Deleted    int       `json:"deleted" xorm:"INT(11)"`
	SaveTime   time.Time `json:"save_time" xorm:"DATETIME"`
	SaveData   []byte    `json:"save_data" xorm:"LONGBLOB"`
}

type NxDumpRank struct {
	RankName string    `json:"rank_name" xorm:"not null pk VARCHAR(32)"`
	Name     string    `json:"name" xorm:"not null default '' index VARCHAR(32)"`
	Uid      string    `json:"uid" xorm:"not null index VARCHAR(100)"`
	RoleJob  int       `json:"role_job" xorm:"INT(11)"`
	Nation   int       `json:"nation" xorm:"INT(11)"`
	Value    int       `json:"value" xorm:"INT(11)"`
	Rank     int       `json:"rank" xorm:"INT(11)"`
	Type     int       `json:"type" xorm:"INT(11)"`
	SaveTime time.Time `json:"save_time" xorm:"DATETIME"`
}

type NxGuilds struct {
	Name       string    `json:"name" xorm:"not null pk VARCHAR(128)"`
	Uid        string    `json:"uid" xorm:"not null unique VARCHAR(32)"`
	CreateTime time.Time `json:"create_time" xorm:"DATETIME"`
	Deleted    int       `json:"deleted" xorm:"INT(11)"`
	SaveTime   time.Time `json:"save_time" xorm:"DATETIME"`
	SaveData   []byte    `json:"save_data" xorm:"LONGBLOB"`
}

type NxLetter struct {
	SerialNo     string    `json:"serial_no" xorm:"not null pk VARCHAR(32)"`
	MsgSendTime  time.Time `json:"msg_send_time" xorm:"index DATETIME"`
	MsgTime      string    `json:"msg_time" xorm:"VARCHAR(64)"`
	MsgName      string    `json:"msg_name" xorm:"index VARCHAR(128)"`
	MsgUid       string    `json:"msg_uid" xorm:"index VARCHAR(32)"`
	MsgSource    string    `json:"msg_source" xorm:"index VARCHAR(128)"`
	MsgSourceUid string    `json:"msg_source_uid" xorm:"index VARCHAR(32)"`
	MsgType      int       `json:"msg_type" xorm:"index INT(11)"`
	MsgContent   string    `json:"msg_content" xorm:"VARCHAR(256)"`
	MsgAppendix  string    `json:"msg_appendix" xorm:"VARCHAR(4096)"`
}

type NxPubdata struct {
	Name       string    `json:"name" xorm:"not null pk VARCHAR(128)"`
	Uid        string    `json:"uid" xorm:"not null unique VARCHAR(32)"`
	CreateTime time.Time `json:"create_time" xorm:"DATETIME"`
	Deleted    int       `json:"deleted" xorm:"INT(11)"`
	SaveTime   time.Time `json:"save_time" xorm:"DATETIME"`
	SaveData   []byte    `json:"save_data" xorm:"LONGBLOB"`
}

type NxResDumpinfo struct {
	Id         string `json:"id" xorm:"not null VARCHAR(128)"`
	Name       string `json:"name" xorm:"VARCHAR(128)"`
	Type       int    `json:"type" xorm:"INT(11)"`
	Property   string `json:"property" xorm:"LONGTEXT"`
	Status     int    `json:"status" xorm:"INT(11)"`
	OpName     string `json:"op_name" xorm:"VARCHAR(128)"`
	OpType     int    `json:"op_type" xorm:"INT(11)"`
	OpProperty string `json:"op_property" xorm:"LONGTEXT"`
	OpStatus   int    `json:"op_status" xorm:"INT(11)"`
	OpVersion  int    `json:"op_version" xorm:"INT(11)"`
	QueryId    string `json:"query_id" xorm:"index VARCHAR(128)"`
	KeyId      int    `json:"key_id" xorm:"not null pk autoincr INT(11)"`
}

type NxRoleDumpinfo struct {
	RoleName            string    `json:"role_name" xorm:"not null default '' index VARCHAR(32)"`
	RoleUid             string    `json:"role_uid" xorm:"not null pk VARCHAR(100)"`
	RoleAccount         string    `json:"role_account" xorm:"index VARCHAR(100)"`
	Deleted             int       `json:"deleted" xorm:"INT(11)"`
	DeleteTime          time.Time `json:"delete_time" xorm:"DATETIME"`
	Sex                 int       `json:"sex" xorm:"INT(11)"`
	RoleJob             int       `json:"role_job" xorm:"INT(11)"`
	Nation              int       `json:"nation" xorm:"index INT(11)"`
	VipLevel            int       `json:"vip_level" xorm:"INT(11)"`
	RoleLevel           int       `json:"role_level" xorm:"INT(11)"`
	ExpToUpgrade        int64     `json:"exp_to_upgrade" xorm:"BIGINT(20)"`
	BattleAbility       int       `json:"battle_ability" xorm:"INT(11)"`
	MilitaryRank        int       `json:"military_rank" xorm:"INT(11)"`
	MilitaryExploitToal int       `json:"military_exploit_toal" xorm:"INT(11)"`
	MilitaryExploit     int       `json:"military_exploit" xorm:"INT(11)"`
	Exp                 int64     `json:"exp" xorm:"BIGINT(20)"`
	Str                 int       `json:"str" xorm:"INT(11)"`
	Ing                 int       `json:"ing" xorm:"INT(11)"`
	Dex                 int       `json:"dex" xorm:"INT(11)"`
	Sta                 int       `json:"sta" xorm:"INT(11)"`
	Hp                  int       `json:"hp" xorm:"INT(11)"`
	MaxHp               int       `json:"max_hp" xorm:"INT(11)"`
	MinAttack           int       `json:"min_attack" xorm:"INT(11)"`
	MaxAttack           int       `json:"max_attack" xorm:"INT(11)"`
	Defend              int       `json:"defend" xorm:"INT(11)"`
	Hit                 int       `json:"hit" xorm:"INT(11)"`
	Miss                int       `json:"miss" xorm:"INT(11)"`
	Crit                int       `json:"crit" xorm:"INT(11)"`
	Toug                int       `json:"toug" xorm:"INT(11)"`
	GuildName           string    `json:"guild_name" xorm:"not null default '' VARCHAR(32)"`
	GuildLevel          int       `json:"guild_level" xorm:"INT(11)"`
	CapitalRec          string    `json:"capital_rec" xorm:"LONGTEXT"`
	ToolBoxInfo         int       `json:"tool_box_info" xorm:"INT(4)"`
	ToolBox             string    `json:"tool_box" xorm:"LONGTEXT"`
	EquipStrengthenRec  string    `json:"equip_strengthen_rec" xorm:"LONGTEXT"`
	EquipBox            string    `json:"equip_box" xorm:"LONGTEXT"`
	SkillView           string    `json:"skill_view" xorm:"LONGTEXT"`
	PassiveskillRec     string    `json:"passiveSkill_rec" xorm:"LONGTEXT"`
	UpdateTime          time.Time `json:"update_time" xorm:"DATETIME"`
}

type NxRoles struct {
	Name        string    `json:"name" xorm:"not null pk VARCHAR(128)"`
	Uid         string    `json:"uid" xorm:"not null unique VARCHAR(32)"`
	Account     string    `json:"account" xorm:"not null index VARCHAR(128)"`
	CreateTime  time.Time `json:"create_time" xorm:"index DATETIME"`
	DeleteTime  time.Time `json:"delete_time" xorm:"DATETIME"`
	OpenTime    time.Time `json:"open_time" xorm:"DATETIME"`
	ServerId    int       `json:"server_id" xorm:"index INT(11)"`
	TotalSec    int       `json:"total_sec" xorm:"INT(11)"`
	Deleted     int       `json:"deleted" xorm:"INT(11)"`
	LetterMax   int       `json:"letter_max" xorm:"INT(11)"`
	SysFlags    int       `json:"sys_flags" xorm:"INT(11)"`
	Illegals    int       `json:"illegals" xorm:"INT(11)"`
	Online      int       `json:"online" xorm:"index INT(11)"`
	RoleIndex   int       `json:"role_index" xorm:"INT(11)"`
	SceneId     int       `json:"scene_id" xorm:"INT(11)"`
	Location    string    `json:"location" xorm:"VARCHAR(64)"`
	ResumeTime  time.Time `json:"resume_time" xorm:"index DATETIME"`
	ResumeInfo  string    `json:"resume_info" xorm:"VARCHAR(128)"`
	ChargeInfo  string    `json:"charge_info" xorm:"VARCHAR(128)"`
	SavePara1   string    `json:"save_para1" xorm:"VARCHAR(2048)"`
	StoreSerial int       `json:"store_serial" xorm:"INT(8)"`
	MoneyInfo   string    `json:"money_info" xorm:"VARCHAR(128)"`
	SaveTime    time.Time `json:"save_time" xorm:"index DATETIME"`
	SaveData    []byte    `json:"save_data" xorm:"LONGBLOB"`
	LazyData    []byte    `json:"lazy_data" xorm:"LONGBLOB"`
	SavePara2   string    `json:"save_para2" xorm:"VARCHAR(2048)"`
}

type NxRolesCrossing struct {
	Name     string    `json:"name" xorm:"not null pk VARCHAR(128)"`
	Uid      string    `json:"uid" xorm:"not null unique VARCHAR(32)"`
	Status   int       `json:"status" xorm:"not null INT(11)"`
	RoleName string    `json:"role_name" xorm:"not null unique(IX_index_name_server) VARCHAR(128)"`
	ServerId int       `json:"server_id" xorm:"not null unique(IX_index_name_server) INT(11)"`
	SceneId  int       `json:"scene_id" xorm:"INT(11)"`
	SaveTime time.Time `json:"save_time" xorm:"DATETIME"`
}

type NxSns struct {
	Uid         string    `json:"uid" xorm:"not null pk unique VARCHAR(32)"`
	CreateTime  time.Time `json:"create_time" xorm:"DATETIME"`
	Deleted     int       `json:"deleted" xorm:"INT(11)"`
	SaveTime    time.Time `json:"save_time" xorm:"DATETIME"`
	ActiveCount int       `json:"active_count" xorm:"INT(11)"`
	SaveData    []byte    `json:"save_data" xorm:"LONGBLOB"`
}

type NxSnsFeed struct {
	FeedId     string    `json:"feed_id" xorm:"not null pk VARCHAR(32)"`
	Uid        string    `json:"uid" xorm:"VARCHAR(32)"`
	TargetUid  string    `json:"target_uid" xorm:"VARCHAR(32)"`
	RangeMask  int       `json:"range_mask" xorm:"INT(11)"`
	Level      int       `json:"level" xorm:"index INT(11)"`
	Category   int       `json:"category" xorm:"INT(11)"`
	DescInfo   string    `json:"desc_info" xorm:"VARCHAR(1024)"`
	CommentAt  string    `json:"comment_at" xorm:"VARCHAR(32)"`
	GoodCount  int       `json:"good_count" xorm:"INT(11)"`
	BadCount   int       `json:"bad_count" xorm:"INT(11)"`
	CreateTime time.Time `json:"create_time" xorm:"index DATETIME"`
}

type NxSnsRelation struct {
	SerialNo   string    `json:"serial_no" xorm:"not null pk VARCHAR(32)"`
	Uid        string    `json:"uid" xorm:"index VARCHAR(32)"`
	TargetUid  string    `json:"target_uid" xorm:"index VARCHAR(32)"`
	RangeMask  int       `json:"range_mask" xorm:"INT(11)"`
	CreateTime time.Time `json:"create_time" xorm:"DATETIME"`
}

type NxSnslog struct {
	SerialNo   string    `json:"serial_no" xorm:"not null pk VARCHAR(32)"`
	LogTime    time.Time `json:"log_time" xorm:"DATETIME"`
	LogUid     string    `json:"log_uid" xorm:"VARCHAR(32)"`
	LogSource  int       `json:"log_source" xorm:"INT(11)"`
	LogType    int       `json:"log_type" xorm:"INT(11)"`
	LogContent string    `json:"log_content" xorm:"VARCHAR(256)"`
	LogComment string    `json:"log_comment" xorm:"VARCHAR(256)"`
}

type NxTrade struct {
	SerialNo   string    `json:"serial_no" xorm:"not null pk VARCHAR(32)"`
	SellerName string    `json:"seller_name" xorm:"VARCHAR(128)"`
	SellerUid  string    `json:"seller_uid" xorm:"VARCHAR(32)"`
	BuyerName  string    `json:"buyer_name" xorm:"VARCHAR(128)"`
	BuyerUid   string    `json:"buyer_uid" xorm:"VARCHAR(32)"`
	BuyerCount int       `json:"buyer_count" xorm:"INT(11)"`
	ItemType   int       `json:"item_type" xorm:"INT(11)"`
	ItemName   string    `json:"item_name" xorm:"VARCHAR(128)"`
	ItemInfo   string    `json:"item_info" xorm:"VARCHAR(2048)"`
	ItemCat1   string    `json:"item_cat1" xorm:"VARCHAR(128)"`
	ItemCat2   string    `json:"item_cat2" xorm:"VARCHAR(128)"`
	ItemCat3   string    `json:"item_cat3" xorm:"VARCHAR(128)"`
	Status     int       `json:"status" xorm:"INT(11)"`
	PriceMin   int       `json:"price_min" xorm:"INT(11)"`
	PriceMax   int       `json:"price_max" xorm:"INT(11)"`
	PriceDeal  int       `json:"price_deal" xorm:"INT(11)"`
	TimeBeg    time.Time `json:"time_beg" xorm:"DATETIME"`
	TimeEnd    time.Time `json:"time_end" xorm:"DATETIME"`
	TimeDeal   time.Time `json:"time_deal" xorm:"DATETIME"`
}

type Test struct {
	Aaa int     `json:"aaa" xorm:"not null pk unique INT(10)"`
	Bbb string  `json:"bbb" xorm:"VARCHAR(255)"`
	Ccc float32 `json:"ccc" xorm:"FLOAT"`
}
