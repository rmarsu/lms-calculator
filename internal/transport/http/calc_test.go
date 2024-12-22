package transport

import (
	"io"
	"lms-1/internal/service"
	"lms-1/pkg/calc"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCalculateHandler(t *testing.T) {
	tests := []struct {
		name         string
		requestBody  string
		expectedCode int
		expectedBody string
		method       string
	}{
		{
			name:         "Valid expression",
			requestBody:  `{"expression": "2 + 2"}`,
			expectedCode: http.StatusOK,
			expectedBody: `{"result":4}`,
			method:       http.MethodPost,
		},
		{
			name:         "Valid expression with minus",
			requestBody:  `{"expression": "2 - 42"}`,
			expectedCode: http.StatusOK,
			expectedBody: `{"result":-40}`,
			method:       http.MethodPost,
		},
		{
			name:         "Valid expression with minus 2",
			requestBody:  `{"expression": "2 - 42 - 55"}`,
			expectedCode: http.StatusOK,
			expectedBody: `{"result":-95}`,
			method:       http.MethodPost,
		},
		{
			name:         "Valid expression with minus 4",
			requestBody:  `{"expression": "2 -42 - 55"}`,
			expectedCode: http.StatusOK,
			expectedBody: `{"result":-95}`,
			method:       http.MethodPost,
		},
		{
			name:         "Valid expression with negative numbers",
			requestBody:  `{"expression": "-1 + 1"}`,
			expectedCode: http.StatusOK,
			expectedBody: `{"result":0}`,
			method:       http.MethodPost,
		},
		{
			name:         "Valid expression with negative numbers 2",
			requestBody:  `{"expression": "1 + (-1 + 2)"}`,
			expectedCode: http.StatusOK,
			expectedBody: `{"result":2}`,
			method:       http.MethodPost,
		},
		{
			name:         "Non valid expression with negative numbers 3",
			requestBody:  `{"expression": "1 + -1"}`,
			expectedCode: http.StatusUnprocessableEntity,
			expectedBody: `{"error":"некорректное выражение"}`,
			method:       http.MethodPost,
		},
		{
			name:         "Not an expression",
			requestBody:  `{"expression": "1 - "}`,
			expectedCode: http.StatusUnprocessableEntity,
			expectedBody: `{"error":"некорректное выражение"}`,
			method:       http.MethodPost,
		},
		{
			name:         "Valid expression with multiply",
			requestBody:  `{"expression": "2 * 12"}`,
			expectedCode: http.StatusOK,
			expectedBody: `{"result":24}`,
			method:       http.MethodPost,
		},
		{
			name:         "Valid expression with grouping",
			requestBody:  `{"expression": "(2 + 2) * 3"}`,
			expectedCode: http.StatusOK,
			expectedBody: `{"result":12}`,
			method:       http.MethodPost,
		},
		{
			name:         "Non valid group expression",
			requestBody:  `{"expression": "((((((((((1)))))))))"}`,
			expectedCode: http.StatusUnprocessableEntity,
			expectedBody: `{"error":"неверно расставленные скобки"}`,
			method:       http.MethodPost,
		},
		{
			name:         "Valid group expression",
			requestBody:  `{"expression": "((((((((((1))))))))))"}`,
			expectedCode: http.StatusOK,
			expectedBody: `{"result":1}`,
			method:       http.MethodPost,
		},
		{
			name:         "Valid expression",
			requestBody:  `{"expression": "1"}`,
			expectedCode: http.StatusOK,
			expectedBody: `{"result":1}`,
			method:       http.MethodPost,
		},
		{
			name:         "Empty expression",
			requestBody:  `{"expression": ""}`,
			expectedCode: http.StatusUnprocessableEntity,
			expectedBody: `{"error":"некорректное выражение"}`,
			method:       http.MethodPost,
		},
		{
			name:         "Invalid characters",
			requestBody:  `{"expression":"2 + a"}`,
			expectedCode: http.StatusUnprocessableEntity,
			expectedBody: `{"error":"некорректное выражение"}`,
			method:       http.MethodPost,
		},
		{
			name:         "Invalid characters 3",
			requestBody:  `{"meow":"2 + 2"}`,
			expectedCode: http.StatusUnprocessableEntity,
			expectedBody: `{"error":"некорректное выражение"}`,
			method:       http.MethodPost,
		},
		{
			name:         "Invalid characters2",
			requestBody:  `{"expression":"meow"}`,
			expectedCode: http.StatusUnprocessableEntity,
			expectedBody: `{"error":"некорректное выражение"}`,
			method:       http.MethodPost,
		},
		{
			name:         "Malformed JSON",
			requestBody:  `{"expression": 2 + 2`,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"error":"Ошибка при парсинге JSON"}`,
			method:       http.MethodPost,
		},
		{
			name:         "Valid expression",
			requestBody:  `{"expression": "2 + 2 - 22 + 432232 / 2 + (-1)"}`,
			expectedCode: http.StatusOK,
			expectedBody: `{"result":216097}`,
			method:       http.MethodPost,
		},
		{
			name:         "Division by zero",
			requestBody:  `{"expression": "2 / 0"}`,
			expectedCode: http.StatusUnprocessableEntity,
			expectedBody: `{"error":"деление на 0"}`,
			method:       http.MethodPost,
		},
		{
			name:         "Plus number",
			requestBody:  `{"expression": "+2 + 2"}`,
			expectedCode: http.StatusUnprocessableEntity,
			expectedBody: `{"error":"некорректное выражение"}`,
			method:       http.MethodPost,
		},
		{
			name:         "Brackets in brackets",
			requestBody:  `{"expression": "((2+3)*2)"}`,
			expectedCode: http.StatusOK,
			expectedBody: `{"result":10}`,
			method:       http.MethodPost,
		},
		{
			name:         "Multiplication by float",
			requestBody:  `{"expression": "2 * 1.5"}`,
			expectedCode: http.StatusOK,
			expectedBody: `{"result":3}`,
			method:       http.MethodPost,
		},
		{
			name:         "Division by float",
			requestBody:  `{"expression": "2 / 1.5"}`,
			expectedCode: http.StatusOK,
			expectedBody: `{"result":1.3333333333333333}`,
			method:       http.MethodPost,
		},
		{
			name:         "Extra JSON",
			requestBody:  `{"expression": "1+1", extra: "data"}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"error":"Ошибка при парсинге JSON"}`,
			method:       http.MethodPost,
		},
		{
			name:         "Non valid method",
			requestBody:  `{"expression": "2 / 2"}`,
			expectedCode: http.StatusMethodNotAllowed,
			expectedBody: `{"error":"Доступен только метод POST"}`,
			method:       http.MethodGet,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/api/v1/calculate", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			services := service.NewService(&service.Deps{
				Calc: calc.NewCalc(),
			})

			handler := NewHandler(services)

			handlerForTest := CheckPOSTMethod(handler.calculateHandler)
			handlerForTest(w, req)

			res := w.Result()
			body, _ := io.ReadAll(res.Body)

			if res.StatusCode != tt.expectedCode {
				t.Errorf("expected status %d, got %d", tt.expectedCode, res.StatusCode)
			}
			if strings.TrimSpace(string(body)) != tt.expectedBody {
				t.Errorf("expected body %s, got %s", tt.expectedBody, string(body))
			}
		})
	}
}
