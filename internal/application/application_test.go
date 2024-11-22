package application

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	calc "github.com/Neat-Snap/go_calc/pkg/calculation"
)

var calcFunc = calc.Calc

func mockCalc(expression string) (float64, error) {
	if expression == "2+2" || expression == "2%2B2" {
		return 4, nil
	}
	return 0, fmt.Errorf("invalid expression")
}

func TestExpressionHandler_GET(t *testing.T) {
	originalCalc := calcFunc
	calcFunc = mockCalc
	defer func() { calcFunc = originalCalc }()
	req := httptest.NewRequest(http.MethodGet, "/expression?expression=2%2B2", nil)
	w := httptest.NewRecorder()

	ExpressionHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d; got %d", http.StatusOK, resp.StatusCode)
	}

	expected := "4.00"
	body := w.Body.String()
	if body != expected {
		t.Errorf("expected body %s; got %s", expected, body)
	}
}

func TestExpressionHandler_POST(t *testing.T) {
	originalCalc := calcFunc
	calcFunc = mockCalc
	defer func() { calcFunc = originalCalc }()
	requestBody := Response{Expression: "2+2"}
	body, _ := json.Marshal(requestBody)

	req := httptest.NewRequest(http.MethodPost, "/expression", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	ExpressionHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d; got %d", http.StatusOK, resp.StatusCode)
	}

	expected := "4.00"
	responseBody := w.Body.String()
	if responseBody != expected {
		t.Errorf("expected body %s; got %s", expected, responseBody)
	}
}

func TestLoggingMiddleware(t *testing.T) {
	handler := LoggingMiddleWare(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d; got %d", http.StatusOK, resp.StatusCode)
	}

	expected := "test response"
	body := w.Body.String()
	if body != expected {
		t.Errorf("expected body %s; got %s", expected, body)
	}
}

func TestApplication_Run(t *testing.T) {
	app := New()
	err := app.Run()
	if err != nil {
		t.Errorf("Run() returned an error: %s", err)
	}
}

func TestApplication_StartServer(t *testing.T) {
	originalCalc := calcFunc
	calcFunc = mockCalc
	defer func() { calcFunc = originalCalc }()
	app := New()
	go func() {
		app.StartServer()
	}()
	defer func() {
	}()
	resp, err := http.Get("http://localhost:8080/expression?expression=2%2B2")
	if err != nil {
		t.Fatalf("failed to make GET request: %s", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d; got %d", http.StatusOK, resp.StatusCode)
	}
	expected := "4.00"
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %s", err)
	}
	if buf.String() != expected {
		t.Errorf("expected response body %s; got %s", expected, buf.String())
	}
}
