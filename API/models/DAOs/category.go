package DAOs

import (
	"time"
)

type Category struct {
	Id          int        `json:"id" gorm:"primaryKey;autoIncrement;not null"`
	Name        string     `json:"name" gorm:"type:varchar(50);not null; uniqueIndex:idx_name_project"`
	Description string     `json:"description" gorm:"type:text;not null"`
	Billable    bool       `json:"billable" gorm:"type:boolean;not null;default:false"`
	Activities  []Activity `json:"activities" gorm:"foreignKey:CategoryId;references:Id"`
	CreatedAt   time.Time  `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `json:"updatedAt" gorm:"autoUpdateTime"`

	// Clés étrangères
	ProjectId int `json:"projectId" gorm:"uniqueIndex:idx_name_project"`

	// Relations
	Project Project `json:"project" gorm:"foreignKey:ProjectId;references:Id"`
}
