package tests

import (
	"llio-api/models/enums"
	"net/http"
	"testing"
)

func TestExportExcelAsAdmin(t *testing.T) {
	w := sendRequest(router, "GET", "/report/excel", nil, enums.Administrator)

	assertResponse(t, w, http.StatusOK, nil)

	if w.Body == nil || w.Body.Len() == 0 {
		t.Fatal("Excel file is empty")
	}

	// Vérifie le Content-Type
	if w.Header().Get("Content-Type") != "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" {
		t.Fatalf("Expected Content-Type Excel, got %s", w.Header().Get("Content-Type"))
	}
}

func TestExportExcelAsNoAdministratorForbidden(t *testing.T) {
	w := sendRequest(router, "GET", "/report/excel", nil, enums.Employee)

	assertResponse(t, w, http.StatusForbidden, nil)
}

func TestExportExcelUnauthenticated(t *testing.T) {
	w := sendRequest(router, "GET", "/report/excel", nil, 999) 

	assertResponse(t, w, http.StatusUnauthorized, nil)
}