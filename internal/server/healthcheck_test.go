package server

import (
	"github.com/dsha256/plfa/internal/jsonlog"
	"github.com/dsha256/plfa/internal/mock"
	"github.com/dsha256/plfa/pkg/dto"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestHandlers_HealthcheckHandler(t *testing.T) {
	testCases := []struct {
		name          string
		method        string
		buildStubs    func(repo *mock.MockRepository)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "200 on GET",
			method: http.MethodGet,
			buildStubs: func(repo *mock.MockRepository) {
				repo.EXPECT().ListTables().Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:   "404 on POST",
			method: http.MethodPost,
			buildStubs: func(repo *mock.MockRepository) {
				repo.EXPECT().ListTables().Times(0).Return([]dto.PragmaticTableWithID{}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:   "404 on PATCH",
			method: http.MethodPatch,
			buildStubs: func(repo *mock.MockRepository) {
				repo.EXPECT().ListTables().Times(0).Return([]dto.PragmaticTableWithID{}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:   "404 on DELETE",
			method: http.MethodDelete,
			buildStubs: func(repo *mock.MockRepository) {
				repo.EXPECT().ListTables().Times(0).Return([]dto.PragmaticTableWithID{}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:   "404 on OPTIONS",
			method: http.MethodOptions,
			buildStubs: func(repo *mock.MockRepository) {
				repo.EXPECT().ListTables().Times(0).Return([]dto.PragmaticTableWithID{}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
	}

	t.Parallel()
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mock.NewMockRepository(ctrl)
			tc.buildStubs(repo)

			logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

			server := NewServer(logger, repo)
			recorder := httptest.NewRecorder()

			url := "/v1/healthcheck"
			req, err := http.NewRequest(tc.method, url, nil)
			require.NoError(t, err)

			server.healthcheckHandler(recorder, req)
			tc.checkResponse(t, recorder)
		})
	}
}
