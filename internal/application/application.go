package application

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
		port = "8080"
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

type Response struct {
	Expression string `json:"expression"`
}

func ExpressionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		expression := r.URL.Query().Get("expression")

		result, err := calc.Calc(expression)
		if err != nil {
			http.Error(w, "Error occurred while calculating the expression", http.StatusBadRequest)
			log.Printf("Error occurred while calculating the expression: %s. Error: %s\n", expression, err)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte(fmt.Sprintf("%.2f", result)))
		if err != nil {
			http.Error(w, "Error occurred while sending the response", http.StatusInternalServerError)
			log.Printf("Error occurred while sending the response: %s\n", err)
			return
		}
	} else if r.Method == http.MethodPost {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error occurred while reading request body", http.StatusInternalServerError)
			log.Printf("Error occurred while reading request body: %s\n", err)
			return
		}
		defer r.Body.Close()

		var response Response
		if err := json.Unmarshal(body, &response); err != nil {
			http.Error(w, "Error occurred while parsing JSON element", http.StatusBadRequest)
			log.Printf("Error occurred while parsing JSON element: %s\n", err)
			return
		}

		result, err := calc.Calc(response.Expression)
		if err != nil {
			http.Error(w, "Error occurred while calculating the expression", http.StatusBadRequest)
			log.Printf("Error occurred while calculating the expression: %s\n", err)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte(fmt.Sprintf("%.2f", result)))
		if err != nil {
			http.Error(w, "Error occurred while sending the response", http.StatusInternalServerError)
			log.Printf("Error occurred while sending the response: %s\n", err)
			return
		}
	} else {
		http.Error(w, "Only POST and GET methods are allowed", http.StatusMethodNotAllowed)
		return
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
	mux.Handle("/expression", handler)

	port := app.config.Port
	log.Printf("Starting server on port %s", port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", app.config.Port), mux)
	if err != nil {
		log.Printf("Error occurred while running the server: %s", err)
	}
}
