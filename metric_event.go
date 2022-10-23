package fiberEvents

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

type MetricEvent struct {
	EventType string                 `json:"event_type"`
	Name      string                 `json:"name"`
	Timestamp string                 `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

func AddMetric(ctx *fiber.Ctx, name string, data map[string]interface{}) {
	event := ctx.Locals("event").(*HTTPEvent)

	metric := &MetricEvent{
		EventType: "metric",
		Name:      name,
		Timestamp: time.Now().String(),
		Data:      data,
	}
	event.Metrics = append(event.Metrics, metric)
}
