package campaign

import "time"

//untuk mewakili tabel yang ada dalam database
type Campaign struct {
	ID              int
	UserID          int
	Name            string
	SHotDescription string
	Perks           string
	BackerCount     int
	GoalAmount      int
	CurrentAmount   int
	Slug            string
	CreatedAt       time.Time
	UpdateAt        time.Time
	CampaignImages  []CampaignImage //relasi ke "campaignimages" pake penghubung payload (ntar ada di repository)
}

type CampaignImage struct {
	ID         int
	CampaignID int
	FileName   string
	IsPrimary  int
	CreatedAt  time.Time
	UpdateAt   time.Time
}
