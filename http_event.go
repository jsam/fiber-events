package fiberEvents

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

type HTTPEvent struct {
	EventType           string         `json:"event_type"`
	RequestID           string         `json:"request_id"`
	HeaderXForwardedFor string         `json:"header_x_forwarded_for"`
	HeaderReferer       string         `json:"header_referer"`
	HeaderUserAgent     string         `json:"header_user_agent"`
	Path                string         `json:"path"`
	Protocol            string         `json:"protocol"`
	Status              int            `json:"status"`
	Method              string         `json:"method"`
	Route               string         `json:"route"`
	BytesReceived       int            `json:"bytes_received"`
	BytesResponded      int            `json:"bytes_responded"`
	StartTime           string         `json:"start_time"`
	EndTime             string         `json:"end_time"`
	Duration            int64          `json:"duration"`
	Logs                []*LogEvent    `json:"-"`
	Metrics             []*MetricEvent `json:"-"`
}

// Handler - logger will print JSON formatted logs onto STDOUT
func Handler(ctx *fiber.Ctx) error {
	start := time.Now()

	rid := string(ctx.Get("X-Request-ID"))
	if rid == "" {
		rid = utils.UUID()
	}

	event := HTTPEvent{
		EventType:           "http",
		RequestID:           rid,
		HeaderXForwardedFor: string(ctx.Get(fiber.HeaderXForwardedFor)),
		HeaderReferer:       string(ctx.Get(fiber.HeaderReferer)),
		HeaderUserAgent:     string(ctx.Get(fiber.HeaderUserAgent)),
		Path:                ctx.Path(),
		StartTime:           start.String(),
		Method:              string(ctx.Method()),
		Route:               string(ctx.Route().Path),
		BytesReceived:       ctx.Request().Header.ContentLength(),
		BytesResponded:      ctx.Response().Header.ContentLength(),
		Protocol:            string(ctx.Protocol()),
	}

	ctx.Locals("event", &event)
	ctx.Next()

	event.Status = ctx.Response().StatusCode()
	event.EndTime = time.Now().String()
	event.Duration = time.Since(start).Milliseconds()

	logStr, err := json.Marshal(event)
	log.Printf("%s", string(logStr))

	for _, logEvent := range event.Logs {
		logStr, err = json.Marshal(logEvent)
		log.Printf("%s", string(logStr))
	}

	for _, metricEvent := range event.Metrics {
		logStr, err = json.Marshal(metricEvent)
		log.Printf("%s", string(logStr))
	}

	return err
}

func New(config ...fiber.Config) fiber.Handler {
	log.SetFlags(0)
	return Handler
}
