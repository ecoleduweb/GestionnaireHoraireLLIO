package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCategories(t *testing.T) {
	w := sendRequest(router, "GET", "/categories", nil, doNotDeleteUser.Id)
	assertResponse(t, w, http.StatusOK, nil)
	assert.NotNil(t, w.Body)
}

func TestGetCategory(t *testing.T) {
	w := sendRequest(router, "GET", fmt.Sprintf("/category/%d", doNotDeleteCategory.Id), nil, doNotDeleteUser.Id)
	assertResponse(t, w, http.StatusOK, nil)
	assert.NotNil(t, w.Body)
}

func TestGetNotFoundCategorie(t *testing.T) {
	w := sendRequest(router, "GET", "/category/0", nil, doNotDeleteUser.Id)
	assertResponse(t, w, http.StatusNotFound, nil)
}
