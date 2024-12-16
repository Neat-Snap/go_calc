package application

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	calc "github.com/Neat-Snap/go_calc/pkg/calculation"
)

type Config struct {
	Port string
}

func NewConfig() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "1081"
		log.Printf("Port not specified, defaulting to %s", port)
	}
	return &Config{
		Port: port,
	}
}

type Application struct {
	config *Config
}

func New() *Application {
	config := NewConfig()
	return &Application{config: config}
}

func (app *Application) Run() error {
	log.SetFlags(log.Ldate | log.Ltime)
	log.Println("Starting the app without a server")
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("Error occured while reading input string: %s\n", err)
	}

	if text == "exit" {
		log.Println("Application was successfully closed")
		return nil
	}
	result, err := calc.Calc(text)

	if err != nil {
		log.Printf("Error occured while calculating the result: %s\n", err)
	} else {
		log.Printf("%s = %f", text, result)
	}
	return nil
}

type Request struct {
	Expression string `json:"expression"`
}

type Response struct {
	Result string `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

func ExpressionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "Only POST method is allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid request format"}`, http.StatusUnprocessableEntity)
		log.Printf("Error parsing request body: %s\n", err)
		return
	}

	result, err := calc.Calc(req.Expression)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		resp := Response{Error: "Expression is not valid"}
		json.NewEncoder(w).Encode(resp)
		log.Printf("Invalid expression: %s. Error: %s\n", req.Expression, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	resp := Response{Result: fmt.Sprintf("%.2f", result)}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		log.Printf("Error encoding response: %s\n", err)
	}
}

type CustomResponseWriter struct {
	http.ResponseWriter
	Body       *bytes.Buffer
	StatusCode int
}

func (crw *CustomResponseWriter) Write(b []byte) (int, error) {
	crw.Body.Write(b)
	return crw.ResponseWriter.Write(b)
}

func (crw *CustomResponseWriter) WriteHeader(statusCode int) {
	crw.StatusCode = statusCode
	crw.ResponseWriter.WriteHeader(statusCode)
}

func LoggingMiddleWare(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		crw := &CustomResponseWriter{
			ResponseWriter: w,
			Body:           bytes.NewBuffer(nil),
			StatusCode:     http.StatusOK,
		}

		handlerFunc(crw, r)
		log.Printf("Response Status: %d, Response Body: %s\n", crw.StatusCode, crw.Body.String())
	}
}

func (app *Application) StartServer() {
	mux := http.NewServeMux()
	handler := LoggingMiddleWare(ExpressionHandler)
	mux.Handle("/api/v1/calculate", handler)

	port := app.config.Port
	log.Printf("Starting server on port %s", port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", app.config.Port), mux)
	if err != nil {
		log.Printf("Error occurred while running the server: %s", err)
	}
}
