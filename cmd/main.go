package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/litsynp/enums-in-go/internal/enum"
	"github.com/ztrue/tracerr"
)

func main() {
	router := http.NewServeMux()

	router.HandleFunc("GET /weekends/{weekend}", func(w http.ResponseWriter, r *http.Request) {
		weekendStr := r.PathValue("weekend")
		weekend, err := enum.DayOfWeekFromString(weekendStr)
		if err != nil {
			log.Default().Println(tracerr.Sprint(err))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			response := map[string]interface{}{
				"status":  "error",
				"message": fmt.Sprintf("Invalid weekend: %s", weekendStr),
			}
			json.NewEncoder(w).Encode(response)
			return
		}

		validate := validator.New()
		if err := validate.RegisterValidation("weekend", WeekendValidation); err != nil {
			log.Default().Println(tracerr.Sprint(err))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			response := map[string]interface{}{
				"status":  "error",
				"message": fmt.Sprintf("Invalid weekend: %s", weekendStr),
			}
			json.NewEncoder(w).Encode(response)
			return
		}

		req := struct {
			Day enum.DayOfWeek `validate:"weekend"`
		}{Day: weekend}

		err = validate.Struct(req)
		if err != nil {
			log.Default().Println(tracerr.Sprint(err))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			response := map[string]interface{}{
				"status":  "error",
				"message": fmt.Sprintf("%s is not a weekend", weekend),
			}
			json.NewEncoder(w).Encode(response)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]interface{}{
			"status":  "success",
			"weekend": req.Day.String(), // Assuming `String()` method returns the day as a string
		}
		json.NewEncoder(w).Encode(response)
	})

	router.HandleFunc("POST /weekends/validate", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Day enum.DayOfWeek `json:"day" validate:"required,weekend"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Default().Println(tracerr.Sprint(err))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			response := map[string]interface{}{
				"status":  "error",
				"message": fmt.Sprintf("Invalid request: %v", err),
			}
			json.NewEncoder(w).Encode(response)
			return
		}

		validate := validator.New()
		if err := validate.RegisterValidation("weekend", WeekendValidation); err != nil {
			log.Default().Println(tracerr.Sprint(err))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			response := map[string]interface{}{
				"status":  "error",
				"message": fmt.Sprintf("Invalid request: %v", err),
			}
			json.NewEncoder(w).Encode(response)
			return
		}

		err := validate.Struct(req)
		if err != nil {
			log.Default().Println(tracerr.Sprint(err))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			response := map[string]interface{}{
				"status":  "error",
				"message": fmt.Sprintf("%s is not a weekend", req.Day),
			}
			json.NewEncoder(w).Encode(response)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]interface{}{
			"status":  "success",
			"weekend": req.Day.String(), // Assuming `String()` method returns the day as a string
		}
		json.NewEncoder(w).Encode(response)
	})

	// Apply middleware
	mux := ChainMiddleware(router, LoggingMiddleware, RecoveryMiddleware)

	log.Println("Server started at :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

// Weekend validation function
func WeekendValidation(fl validator.FieldLevel) bool {
	dow, ok := fl.Field().Interface().(enum.DayOfWeek)
	if !ok {
		return false
	}
	return dow >= enum.DayOfWeekSaturday && dow <= enum.DayOfWeekSunday
}

// Middleware type
type Middleware func(http.Handler) http.Handler

// Custom ResponseWriter to capture status code
type logWriter struct {
	http.ResponseWriter
	code int
}

func (lw *logWriter) WriteHeader(code int) {
	lw.code = code
	lw.ResponseWriter.WriteHeader(code)
}

// Unwrap supports http.ResponseController.
func (lw *logWriter) Unwrap() http.ResponseWriter { return lw.ResponseWriter }

// Middleware chain function
func ChainMiddleware(h http.Handler, middlewares ...Middleware) http.Handler {
	for _, middleware := range middlewares {
		h = middleware(h)
	}
	return h
}

// Logging middleware
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lw := &logWriter{ResponseWriter: w}
		start := time.Now()
		next.ServeHTTP(lw, r)
		// Log request method, URL path, status code, and duration
		log.Printf("%s %s %d %s", r.Method, r.URL.Path, lw.code, time.Since(start))
	})
}

// Recovery middleware
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
