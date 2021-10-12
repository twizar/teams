package ports_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twizar/common/pkg/dto"
	"github.com/twizar/teams/internal/application/service"
	"github.com/twizar/teams/internal/ports"
	"github.com/twizar/teams/test"
)

const (
	liverpoolID = "d6548941-53f1-4d27-ad3d-0286cf512af1"
	milanID     = "5d912b4e-4932-496d-b706-c22b58f76a21"
	sevillaID   = "b0f6d915-da69-4681-bd7e-d933dd599ab2"
	bayernID    = "418ca28d-af10-4fbb-8b10-6afd74a001b7"
)

func TestHTTPServer_APICalls(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		requestFactory func() *http.Request
		check          func(resultTeams []dto.Team, httpStatus int)
	}{
		{
			name: "get all teams",
			requestFactory: func() *http.Request {
				request, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/teams", http.NoBody)
				require.NoError(t, err)

				return request
			},
			check: func(resultTeams []dto.Team, httpStatus int) {
				assert.Equal(t, http.StatusOK, httpStatus)
				assert.Equal(t, 703, len(resultTeams))
			},
		},
		{
			name: "get teams by ID",
			requestFactory: func() *http.Request {
				request, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/teams", http.NoBody)
				require.NoError(t, err)

				q := request.URL.Query()
				ids := strings.Join([]string{liverpoolID, milanID, sevillaID, bayernID}, ",")
				q.Add("ids", ids)
				request.URL.RawQuery = q.Encode()

				return request
			},
			check: func(resultTeams []dto.Team, httpStatus int) {
				assert.Equal(t, http.StatusOK, httpStatus)
				assert.Equal(t, 4, len(resultTeams))
			},
		},
		{
			name: "search teams",
			requestFactory: func() *http.Request {
				request, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/teams/search", http.NoBody)
				require.NoError(t, err)
				q := request.URL.Query()
				q.Add("min-rating", "2")
				q.Add("order-by", "rating")
				q.Add("limit", "50")
				request.URL.RawQuery = q.Encode()

				return request
			},
			check: func(resultTeams []dto.Team, httpStatus int) {
				assert.Equal(t, http.StatusOK, httpStatus)
				assert.Equal(t, 50, len(resultTeams))
			},
		},
	}

	repo, clean := test.TeamRepoHelper(t)
	defer clean()

	teams := service.NewTeams(repo)
	server := ports.NewHTTPServer(teams)
	router := ports.ConfigureRouter(server)

	for _, testCase := range tests {
		testCase := testCase

		writer := httptest.NewRecorder()
		router.ServeHTTP(writer, testCase.requestFactory())
		body, err := io.ReadAll(writer.Body)
		require.NoError(t, err)

		var result []dto.Team
		err = json.Unmarshal(body, &result)
		require.NoError(t, err)

		testCase.check(result, writer.Code)
	}
}
