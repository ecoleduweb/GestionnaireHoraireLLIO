package DTOs

import "time"

type DetailedActivityDTO struct {
	Id          int       `json:"id"`
	Name        string    `json:"name" binding:"required,max=50"`
	Description string    `json:"description" binding:"required"`
	StartDate   time.Time `json:"startDate" binding:"required"`
	EndDate     time.Time `json:"endDate" binding:"required"`
	UserId      int       `json:"userId" `
	ProjectId   int       `json:"projectId" `
	ProjectName string    `json:"projectName" `
	CategoryId  int       `json:"categoryId" `
}
