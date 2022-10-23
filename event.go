package fiberEvents

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

type Event struct {
	*sync.Mutex `json:"-"`

	RequestID           string `json:"request_id"`
	HeaderXForwardedFor string `json:"header_x_forwarded_for"`
	HeaderReferer       string `json:"header_referer"`
	HeaderUserAgent     string `json:"header_user_agent"`
	Path                string `json:"path"`
	Protocol            string `json:"protocol"`
	Status              int    `json:"status"`
	Method              string `json:"method"`
	Route               string `json:"route"`
	BytesReceived       int    `json:"bytes_received"`
	BytesResponded      int    `json:"bytes_responded"`
	StartTime           string `json:"start_time"`
	EndTime             string `json:"end_time"`
	Duration            int64  `json:"duration"`

	Tracelog []*Log                   `json:"tracelog"`
	Events   []map[string]interface{} `json:"events"`
}

func (e *Event) AddEvent(name string, value interface{}) {
	e.Events = append(e.Events, map[string]interface{}{
		"name":  name,
		"value": value,
	})
}

// Handler - logger will print JSON formatted logs onto STDOUT
func Handler(ctx *fiber.Ctx) error {
	start := time.Now()

	rid := string(ctx.Get("X-Request-ID"))
	if rid == "" {
		rid = utils.UUID()
	}

	event := Event{
		Mutex:               &sync.Mutex{},
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

	return err
}

func New(config ...fiber.Config) fiber.Handler {
	log.SetFlags(0)
	return Handler
}
