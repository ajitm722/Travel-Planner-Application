package hotels

import (
	"context"
	"fmt"
	"log"
	"travel-planner/internal/tracing"

	"go.opentelemetry.io/otel/attribute"
)

func GetHotels(ctx context.Context, city string, minBudget, maxBudget int) (string, error) {
	_, span := tracing.StartSpan(ctx, "GetHotels")
	defer span.End()

	// Add attributes to the span
	span.SetAttributes(
		attribute.String("city", city),
		attribute.Int("minBudget", minBudget),
		attribute.Int("maxBudget", maxBudget),
	)

	// Simulate an error scenario
	if city == "" {
		err := fmt.Errorf("city cannot be empty")
		span.RecordError(err)
		return "", err
	}

	log.Printf("Fetching hotels in city: %s, budget: %d-%d", city, minBudget, maxBudget)

	// Simulate API Call
	result := fmt.Sprintf("Hotels in %s: $100-$200", city)
	return result, nil
}
