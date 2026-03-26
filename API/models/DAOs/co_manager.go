package DAOs

type CoManager struct {
	ProjectId int `json:"project_id" gorm:"type:bigint(20) unsigned;primaryKey;not null"`
	UserId    int `json:"user_id" gorm:"type:bigint(20) unsigned;primaryKey;not null"`

	Project Project `json:"project" gorm:"foreignKey:ProjectId;references:Id;constraint:OnDelete:CASCADE"`
	User    User    `json:"user" gorm:"foreignKey:UserId;references:Id;constraint:OnDelete:CASCADE"`
}
