package tests

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetOutlookEventsNoGraphToken(t *testing.T) {
	w := sendRequest(router, "GET", "/activities/me/outlook", nil, &doNotDeleteUser.Id)
	assertResponse(t, w, http.StatusUnauthorized, nil)
	assert.NotNil(t, w.Body)
}

func TestGetOutlookEventsExpiredGraphToken(t *testing.T) {
	w := sendRequest(router, "GET", "/activities/me/outlook", nil, &doNotDeleteUser2.Id)
	assertResponse(t, w, http.StatusUnauthorized, nil)
	assert.NotNil(t, w.Body)
}
