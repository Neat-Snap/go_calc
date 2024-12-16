package application

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	calc "github.com/Neat-Snap/go_calc/pkg/calculation"
)

var calcFunc = calc.Calc

func mockCalc(expression string) (float64, error) {
	if expression == "2+2" || expression == "2%2B2" {
		return 4, nil
	}
	return 0, fmt.Errorf("invalid expression")
}

func TestExpressionHandler_POST(t *testing.T) {
	originalCalc := calcFunc
	calcFunc = mockCalc
	defer func() { calcFunc = originalCalc }()

	requestBody := Request{Expression: "2+2"}
	body, _ := json.Marshal(requestBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	ExpressionHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d; got %d", http.StatusOK, resp.StatusCode)
	}

	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response body: %s", err)
	}

	expected := "4.00"
	if response.Result != expected {
		t.Errorf("expected result %s; got %s", expected, response.Result)
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
	body, _ := io.ReadAll(w.Body)
	if string(body) != expected {
		t.Errorf("expected body %s; got %s", expected, string(body))
	}
}

func TestApplication_Run(t *testing.T) {
	app := New()
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	r, w, _ := os.Pipe()
	os.Stdin = r
	fmt.Fprintln(w, "exit")
	w.Close()

	err := app.Run()
	if err != nil {
		t.Errorf("Run() returned an error: %s", err)
	}
}

func TestApplication_StartServer(t *testing.T) {
	originalCalc := calcFunc
	calcFunc = mockCalc
	defer func() { calcFunc = originalCalc }()

	os.Setenv("PORT", "8080")
	app := New()
	go func() {
		app.StartServer()
	}()
	time.Sleep(1 * time.Second)

	client := &http.Client{}
	requestBody := Request{Expression: "2+2"}
	body, _ := json.Marshal(requestBody)

	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/api/v1/calculate", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("failed to create POST request: %s", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("failed to make POST request: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d; got %d", http.StatusOK, resp.StatusCode)
	}

	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response body: %s", err)
	}

	expected := "4.00"
	if response.Result != expected {
		t.Errorf("expected result %s; got %s", expected, response.Result)
	}
}

func equalJSON(a, b map[string]interface{}) bool {
	return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
}
