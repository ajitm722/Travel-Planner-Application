package hotels

import (
	"context"
	"fmt"
	"time"
	"travel-planner/internal/tracing"

	"go.opentelemetry.io/otel/attribute"
	"golang.org/x/exp/rand"
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

	// Simulate API Call with a random sleep between 1-5 seconds
	sleepDuration := time.Duration(rand.Intn(5)+1) * time.Second // Random duration between 1 and 5 seconds
	time.Sleep(sleepDuration)

	// Simulate API Call result
	result := fmt.Sprintf("Hotels in %s: $100-$200", city)

	return result, nil
}
