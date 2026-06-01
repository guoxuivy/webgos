package models

type Department struct {
	BaseFields
	Name     string       `gorm:"size:50;unique;not null" json:"name"`
	ParentID int          `gorm:"column:parent_id;default:0" json:"parent_id"`
	LeaderID int          `gorm:"column:leader_id;default:0" json:"leader_id"`
	Remark   string       `gorm:"size:200" json:"remark"`
	Status   int          `gorm:"type:tinyint;default:1" json:"status"`
	Sort     int          `gorm:"column:sort;default:0" json:"sort"`
	Leader   *User        `gorm:"foreignKey:LeaderID" json:"leader"`
	Children []Department `gorm:"-" json:"children"`
}

func (Department) TableName() string {
	return "departments"
}
