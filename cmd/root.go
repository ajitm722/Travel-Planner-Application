package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

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

	if city == "" || startDate == "" || endDate == "" {
		http.Error(w, "Missing required parameters: city, start-date, end-date", http.StatusBadRequest)
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

	// Fetch weather information asynchronously
	go func() {
		weatherInfo, err := weather.GetWeather(ctx, city, startDate, endDate)
		if err != nil {
			errCh <- fmt.Errorf("error fetching weather: %v", err)
			return
		}
		weatherCh <- weatherInfo
	}()

	// Fetch flight options asynchronously
	go func() {
		flightOptions, err := flights.GetFlights(ctx, city, startDate, endDate, minBudget, maxBudget)
		if err != nil {
			errCh <- fmt.Errorf("error fetching flights: %v", err)
			return
		}
		flightCh <- flightOptions
	}()

	// Fetch hotel options asynchronously
	go func() {
		hotelOptions, err := hotels.GetHotels(ctx, city, minBudget, maxBudget)
		if err != nil {
			errCh <- fmt.Errorf("error fetching hotels: %v", err)
			return
		}
		hotelCh <- hotelOptions
	}()

	// Check for errors from any of the API calls
	select {
	case err := <-errCh:
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	default:
		// No errors, prepare the response
		response := map[string]interface{}{
			"weather": <-weatherCh,
			"flights": <-flightCh,
			"hotels":  <-hotelCh,
		}

		// Marshal the response into formatted JSON
		jsonResponse, err := json.MarshalIndent(response, "", "  ")
		if err != nil {
			// Handle error if JSON marshalling fails
			http.Error(w, fmt.Sprintf("Error marshalling JSON response: %v", err), http.StatusInternalServerError)
			return
		}

		// Return the JSON response
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	}
}
