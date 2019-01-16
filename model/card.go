package model

import "time"

const (
	ChinaUnicom = iota
	ChinaMobile
	ChinaTelecom
)

const (
	AuthUnauthed = iota
	AuthAudited
	AuthUnpass
	AuthPass
)

const (
	CardReadied = iota
	CardActivated
	CardDeactivated
	CardRetired
)

type Card struct {
	Id                 int
	Network            int
	Operator           int
	Source             int
	Msisdn             string `gorm:"unique_index:uiq_msisdn"`
	Iccid              string `gorm:"unique_index:uiq_iccid"`
	Imsi               string
	Imei               string
	CheckSum           string `mapstructure:"check_sum"`
	Used               float64
	Total              float64
	TotalUsed          float64
	MonthUsed          float64
	LastUsed           float32
	LastMonthUsed      float64
	VoiceUsed          int
	VoiceTotal         int
	VoiceTotalUsed     int
	VoiceMonthUsed     int
	VoiceLastUsed      int
	VoiceLastMonthUsed int
	Status             int
	Name               string
	IdentNum           string
	Mobile             string
	AuthStatus         int
	IsAuthed           bool
	EntAuth            bool
	RenewalAmount      float64
	RenewalCount       int
	Rate               float64
	ActivatedAt        *time.Time
	ExpiredAt          *time.Time
	LastActivatedAt    *time.Time
	UsedAt             *time.Time
	AuthorizedAt       *time.Time
	CardPlans          []CardPlan
	Orders             []Order
	Partners           []Partner `gorm:"many2many:partner_cards"`
	Users              []User    `gorm:"many2many:user_cards"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

func (c Card) TableName() string {
	return "cards"
}

func (c *Card) FullIccid() string {
	if c.Operator == ChinaUnicom && len(c.Iccid) == 19 {
		iccids := []byte(c.Iccid)
		var num0, num1, num2, num3 int
		for i := 0; i < 10; i++ {
			num0 += 2 * int(iccids[i*2]-'0')
			if int(iccids[i*2]-'0')-4 > 0 {
				num1 += 9
			}
		}
		for i := 1; i < 10; i++ {
			num2 += int(iccids[i*2-1] - '0')
		}
		num3 = 10 - (num0-num1+num2)%10
		if num3 == 10 {
			num3 = 0
		}
		iccids = append(iccids, byte(num3+'0'))
		return string(iccids)
	}
	return c.Iccid
}

func (c *Card) IsUnlimited() bool {
	for _, plan := range c.CardPlans {
		if plan.ExpiredAt == nil || time.Now().Before(*plan.ExpiredAt) {
			if plan.Plan.IsUnlimited {
				return true
			}
		}
	}
	return false
}

type Auth struct {
	Id         int
	CardId     int   `gorm:"index:idx_card_id"`
	Card       *Card `json:",omitempty"`
	UserId     int   `gorm:"index:idx_user_id"`
	UserIp     string
	Name       string
	IdentNum   string
	Mobile     string
	Cover      string
	Back       string
	People     string
	Status     int `gorm:"type:tinyint"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	AuthTracks []AuthTrack
}

func (a Auth) TableName() string {
	return "auths"
}

type AuthTrack struct {
	Id        int
	AuthId    int `gorm:"index:idx_auth_id"`
	OldStatus int
	Status    int
	Remark    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (at AuthTrack) TableName() string {
	return "auth_tracks"
}
