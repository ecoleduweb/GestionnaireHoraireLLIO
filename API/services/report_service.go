package services

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"llio-api/repositories"
	"strconv"
	"time"
)

func GenerateExportCSV(from string, to string) (*bytes.Buffer, error) {
    activities, err := repositories.GetAllForExport(from, to)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch activities: %w", err)
    }

    buf := &bytes.Buffer{}
    writer := csv.NewWriter(buf)

    writer.Write([]string{
        "activity_id", "activity_name", "activity_description", "start_time", "end_time",
        "category_id", "category_name",
        "project_id", "project_name",
        "user_id", "user_fisrt_name", "user_last_name",
    })

    for _, a := range activities {
        writer.Write([]string{
            strconv.Itoa(a.Id), a.Name, a.Description,
            a.StartDate.Format(time.RFC3339), a.EndDate.Format(time.RFC3339),
            strconv.Itoa(a.CategoryId), a.Category.Name,
            strconv.Itoa(a.ProjectId), a.Project.Name,
            strconv.Itoa(a.UserId), a.User.FirstName, a.User.LastName,
        })
    }

    writer.Flush()
    if err := writer.Error(); err != nil {
        return nil, fmt.Errorf("csv write failed: %w", err)
    }

    return buf, nil
}
