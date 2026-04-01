package DAOs

import "llio-api/models/enums"

type User struct {
	Id           int            `json:"id" gorm:"primaryKey;autoIncrement;not null"`
	FirstName    string         `json:"first_name" gorm:"type:varchar(50);not null"`
	LastName     string         `json:"last_name" gorm:"type:varchar(50);not null"`
	Email        string         `json:"email" gorm:"type:varchar(255);not null;unique;index"`
	Role         enums.UserRole `json:"role" gorm:"type:enum(2,1,0);not null;default:0"`
	RefreshToken *string        `json:"-" gorm:"type:text"`
	Activities   []Activity     `json:"activities" gorm:"foreignKey:UserId;references:Id"`
	Projects     []Project      `json:"projects" gorm:"foreignKey:ManagerId;references:Id"`
}
