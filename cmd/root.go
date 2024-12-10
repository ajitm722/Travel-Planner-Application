package cmd

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"travel-planner/internal/flights"
	"travel-planner/internal/hotels"
	"travel-planner/internal/tracing"
	"travel-planner/internal/weather"

	"github.com/spf13/cobra"
)

var (
	city      string
	startDate string
	endDate   string
	minBudget int
	maxBudget int
)

var rootCmd = &cobra.Command{
	Use:   "travel-planner",
	Short: "Plan your travels with weather, flights, and hotel suggestions",
	Run: func(cmd *cobra.Command, args []string) {
		// Your current logic for CLI will remain unchanged.
	},
}

// HandlePlanRequest is the HTTP handler function that will be called to handle requests
func HandlePlanRequest(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	city = r.URL.Query().Get("city")
	startDate = r.URL.Query().Get("start-date")
	endDate = r.URL.Query().Get("end-date")
	minBudget, _ = strconv.Atoi(r.URL.Query().Get("min-budget"))
	maxBudget, _ = strconv.Atoi(r.URL.Query().Get("max-budget"))

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

	// Call weather API
	weatherInfo, err := weather.GetWeather(ctx, city, startDate, endDate)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching weather: %v", err), http.StatusInternalServerError)
		return
	}

	// Call flights API
	flightOptions, err := flights.GetFlights(ctx, city, startDate, endDate, minBudget, maxBudget)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching flights: %v", err), http.StatusInternalServerError)
		return
	}

	// Call hotels API
	hotelOptions, err := hotels.GetHotels(ctx, city, minBudget, maxBudget)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching hotels: %v", err), http.StatusInternalServerError)
		return
	}

	// Return the results in JSON format (or as plain text depending on preference)
	response := fmt.Sprintf(
		`{
			"weather": "%s",
			"flights": "%s",
			"hotels": "%s"
		}`, weatherInfo, flightOptions, hotelOptions)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(response))
}

// HTTP server setup
func init() {
	// This is the same as before for the CLI-based logic
	rootCmd.Flags().StringVarP(&city, "city", "c", "", "City to travel to")
	rootCmd.Flags().StringVarP(&startDate, "start-date", "s", "", "Start date of the travel (YYYY-MM-DD)")
	rootCmd.Flags().StringVarP(&endDate, "end-date", "e", "", "End date of the travel (YYYY-MM-DD)")
	rootCmd.Flags().IntVar(&minBudget, "min-budget", 0, "Minimum budget")
	rootCmd.Flags().IntVar(&maxBudget, "max-budget", 0, "Maximum budget")

	_ = rootCmd.MarkFlagRequired("city")
	_ = rootCmd.MarkFlagRequired("start-date")
	_ = rootCmd.MarkFlagRequired("end-date")
}
