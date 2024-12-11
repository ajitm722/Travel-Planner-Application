package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"travel-planner/internal/flights"
	"travel-planner/internal/hotels"
	"travel-planner/internal/tracing"
	"travel-planner/internal/weather"
)

// HandlePlanRequest is the HTTP handler function that will be called to handle requests
func HandlePlanRequest(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	city := r.URL.Query().Get("city")
	startDate := r.URL.Query().Get("start-date")
	endDate := r.URL.Query().Get("end-date")
	minBudget, _ := strconv.Atoi(r.URL.Query().Get("min-budget"))
	maxBudget, _ := strconv.Atoi(r.URL.Query().Get("max-budget"))

	// Validate required parameters
	if city == "" || startDate == "" || endDate == "" {
		http.Error(w, "Missing required parameters: city, start-date, end-date", http.StatusBadRequest)
		return
	}

	// Validate budget values
	if minBudget == 0 {
		minBudget = 100
	}
	if maxBudget == 0 {
		maxBudget = 1000
	}

	// Validate date formats
	if _, err := time.Parse("2006-01-02", startDate); err != nil {
		http.Error(w, "Invalid start-date format, expected YYYY-MM-DD", http.StatusBadRequest)
		return
	}
	if _, err := time.Parse("2006-01-02", endDate); err != nil {
		http.Error(w, "Invalid end-date format, expected YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	// Set up OpenTelemetry TracerProvider
	ctx := context.Background()
	tracerProvider, err := tracing.SetupTracerProvider()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to set up tracer provider: %v", err), http.StatusInternalServerError)
		return
	}
	defer tracerProvider.Shutdown(ctx)

	// Start root span
	ctx, span := tracing.StartSpan(ctx, "travel-planner")
	defer span.End()

	// Create channels to collect results
	weatherCh := make(chan string, 1)
	flightCh := make(chan string, 1)
	hotelCh := make(chan string, 1)
	errCh := make(chan error, 3)

	// Create a context with cancellation
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Fetch weather information asynchronously
	go func() {
		defer close(weatherCh)
		weatherInfo, err := weather.GetWeather(ctx, city, startDate, endDate)
		if err != nil {
			cancel() // Cancel other API calls
			errCh <- fmt.Errorf("error fetching weather: %v", err)
			return
		}
		weatherCh <- weatherInfo
	}()

	// Fetch flight options asynchronously
	go func() {
		defer close(flightCh)
		flightOptions, err := flights.GetFlights(ctx, city, startDate, endDate, minBudget, maxBudget)
		if err != nil {
			cancel() // Cancel other API calls
			errCh <- fmt.Errorf("error fetching flights: %v", err)
			return
		}
		flightCh <- flightOptions
	}()

	// Fetch hotel options asynchronously
	go func() {
		defer close(hotelCh)
		hotelOptions, err := hotels.GetHotels(ctx, city, minBudget, maxBudget)
		if err != nil {
			cancel() // Cancel other API calls
			errCh <- fmt.Errorf("error fetching hotels: %v", err)
			return
		}
		hotelCh <- hotelOptions
	}()

	// Collect results with error handling and timeout
	select {
	case err := <-errCh:
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	case <-time.After(10 * time.Second):
		http.Error(w, "Timeout while fetching data", http.StatusGatewayTimeout)
		return
	default:
		// Prepare the response
		response := map[string]interface{}{
			"weather": <-weatherCh,
			"flights": <-flightCh,
			"hotels":  <-hotelCh,
		}

		// Marshal the response into formatted JSON
		jsonResponse, err := json.MarshalIndent(response, "", "  ")
		if err != nil {
			http.Error(w, fmt.Sprintf("Error marshalling JSON response: %v", err), http.StatusInternalServerError)
			return
		}

		// Return the JSON response
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	}
}
