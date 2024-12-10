package cmd

import (
	"context"
	"fmt"
	"os"

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
		// Set up OpenTelemetry TracerProvider
		ctx := context.Background()
		tracerProvider, err := tracing.SetupTracerProvider()
		if err != nil {
			fmt.Println("Failed to set up tracer provider:", err)
			os.Exit(1)
		}
		defer func() {
			_ = tracerProvider.Shutdown(ctx)
		}()

		// Start root span
		ctx, span := tracing.StartSpan(ctx, "travel-planner")
		defer span.End()

		// Call weather API
		fmt.Println("Fetching weather information...")
		weatherInfo, err := weather.GetWeather(ctx, city, startDate, endDate)
		if err != nil {
			fmt.Println("Error fetching weather:", err)
			return
		}
		fmt.Printf("Weather Information: %s\n", weatherInfo)

		// Call flights API
		fmt.Println("Fetching flight options...")
		flightOptions, err := flights.GetFlights(ctx, city, startDate, endDate, minBudget, maxBudget)
		if err != nil {
			fmt.Println("Error fetching flights:", err)
			return
		}
		fmt.Printf("Flight Options: %s\n", flightOptions)

		// Call hotels API
		fmt.Println("Fetching hotel options...")
		hotelOptions, err := hotels.GetHotels(ctx, city, minBudget, maxBudget)
		if err != nil {
			fmt.Println("Error fetching hotels:", err)
			return
		}
		fmt.Printf("Hotel Options: %s\n", hotelOptions)
	},
}

func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		return err
	}
	return nil
}

func init() {
	rootCmd.Flags().StringVarP(&city, "city", "c", "", "City to travel to")
	rootCmd.Flags().StringVarP(&startDate, "start-date", "s", "", "Start date of the travel (YYYY-MM-DD)")
	rootCmd.Flags().StringVarP(&endDate, "end-date", "e", "", "End date of the travel (YYYY-MM-DD)")
	rootCmd.Flags().IntVar(&minBudget, "min-budget", 0, "Minimum budget")
	rootCmd.Flags().IntVar(&maxBudget, "max-budget", 0, "Maximum budget")

	_ = rootCmd.MarkFlagRequired("city")
	_ = rootCmd.MarkFlagRequired("start-date")
	_ = rootCmd.MarkFlagRequired("end-date")
}
