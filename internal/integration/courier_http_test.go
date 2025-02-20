package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

type CreateCourierRequest struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	Phone  string `json:"phone"`
	Status string `json:"status"`
}

type CourierResponse struct {
	Id     int64  `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Phone  string `json:"phone"`
	Status string `json:"status"`
}

func (s *IntegrationTestSuite) TestCreateCourierHTTP() {
	tests := []struct {
		name           string
		request        *CreateCourierRequest
		setupFn        func() // для предварительной настройки
		expectedStatus int
		validateFn     func(*testing.T, *http.Response, error)
	}{
		{
			name: "Success/ValidData",
			request: &CreateCourierRequest{
				Name:   "Test Courier",
				Email:  "test@example.com",
				Phone:  "+1234567890",
				Status: "active",
			},
			expectedStatus: http.StatusOK,
			validateFn: func(t *testing.T, resp *http.Response, err error) {
				require.NoError(t, err)
				var courier CourierResponse
				err = json.NewDecoder(resp.Body).Decode(&courier)
				require.NoError(t, err)
				require.NotZero(t, courier.Id)
				require.Equal(t, "Test Courier", courier.Name)
			},
		},
		{
			name: "Error/InvalidJSON",
			request: &CreateCourierRequest{
				Email: "not_an_email",
			},
			expectedStatus: http.StatusBadRequest,
			validateFn: func(t *testing.T, resp *http.Response, err error) {
				require.NoError(t, err)
				var errResp map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&errResp)
				require.NoError(t, err)
				require.Contains(t, errResp, "error")
			},
		},
		{
			name: "Error/DuplicateEmail",
			request: &CreateCourierRequest{
				Name:   "Duplicate",
				Email:  "duplicate@example.com",
				Phone:  "+0987654321",
				Status: "active",
			},
			setupFn: func() {
				// Создаем курьера с тем же email
				s.createCourierViaGRPC(&CreateCourierRequest{
					Name:   "Original",
					Email:  "duplicate@example.com",
					Phone:  "+1111111111",
					Status: "active",
				})
			},
			expectedStatus: http.StatusConflict,
			validateFn: func(t *testing.T, resp *http.Response, err error) {
				require.NoError(t, err)
				var errResp map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&errResp)
				require.NoError(t, err)
				require.Contains(t, errResp["error"], "duplicate")
			},
		},
		{
			name:           "Error/EmptyRequest",
			request:        &CreateCourierRequest{},
			expectedStatus: http.StatusBadRequest,
			validateFn: func(t *testing.T, resp *http.Response, err error) {
				require.NoError(t, err)
				var errResp map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&errResp)
				require.NoError(t, err)
				require.Contains(t, errResp, "error")
			},
		},
		{
			name: "Success/LongValues",
			request: &CreateCourierRequest{
				Name:   s.generateString(255),
				Email:  fmt.Sprintf("%s@example.com", s.generateString(50)),
				Phone:  "+" + s.generateString(15),
				Status: "active",
			},
			expectedStatus: http.StatusOK,
			validateFn: func(t *testing.T, resp *http.Response, err error) {
				require.NoError(t, err)
				var courier CourierResponse
				err = json.NewDecoder(resp.Body).Decode(&courier)
				require.NoError(t, err)
				require.NotZero(t, courier.Id)
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			// Очищаем БД перед каждым тестом
			s.cleanupDB()

			// Выполняем setup если есть
			if tt.setupFn != nil {
				tt.setupFn()
			}

			// Готовим запрос
			body, err := json.Marshal(tt.request)
			s.Require().NoError(err)

			// Выполняем запрос
			resp, err := http.Post(
				fmt.Sprintf("http://localhost:%s/api/v1/couriers", s.httpPort),
				"application/json",
				bytes.NewBuffer(body),
			)

			// Проверяем статус
			if resp != nil {
				s.Require().Equal(tt.expectedStatus, resp.StatusCode)
				defer resp.Body.Close()
			}

			// Выполняем дополнительные проверки
			if tt.validateFn != nil {
				tt.validateFn(s.T(), resp, err)
			}
		})
	}
}

// Вспомогательные методы

func (s *IntegrationTestSuite) cleanupDB() {
	_, err := s.db.Exec("TRUNCATE couriers CASCADE")
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) generateString(length int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, length)
	for i := range b {
		b[i] = letters[s.r.Intn(len(letters))]
	}
	return string(b)
}

func (s *IntegrationTestSuite) createCourierViaGRPC(req *CreateCourierRequest) {
	// Используем уже существующий gRPC клиент для создания тестовых данных
	_, err := s.courier.Create(context.Background(), &proto.CreateCourierRequest{
		Name:   req.Name,
		Email:  req.Email,
		Phone:  req.Phone,
		Status: req.Status,
	})
	s.Require().NoError(err)
}
