package services

import (
	"fmt"
	"io"
	"llio-api/repositories"

	"github.com/xuri/excelize/v2"
)

func GenerateExcel(w io.Writer) error {
    activities, err := repositories.GetAllForExport()
    if err != nil {
        return fmt.Errorf("failed to fetch activities: %w", err)
    }

	f := excelize.NewFile()
	sheet := "Sheet1"

	headers := []string{
		"Prénom", "Nom",
		"Projet", "Unique ID",
		"Activité", "Catégorie",
		"Date de début", "Date de fin","Temps passé (h)",
	}

	for i, h := range headers {
		col := string(rune('A' + i))
		f.SetCellValue(sheet, col+"1", h)
	}

	row := 2

	for _, a := range activities {
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), a.User.FirstName)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), a.User.LastName)

		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), a.Project.Name)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), a.Project.UniqueId)

		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), a.Name)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), a.Category.Name)

		f.SetCellValue(sheet, fmt.Sprintf("G%d", row), a.StartDate.Format("2006-01-02 15:04"))
		f.SetCellValue(sheet, fmt.Sprintf("H%d", row), a.EndDate.Format("2006-01-02 15:04"))
        f.SetCellValue(sheet, fmt.Sprintf("I%d", row), a.TimeSpent)

		row++
	}

	return f.Write(w)
}
