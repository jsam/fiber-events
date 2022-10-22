package gofiberJsonLog

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Record struct {
	IP                  string `json:"ip"`
	HeaderXForwardedFor string `json:"header_x_forwarded_for"`
	Port                string `json:"port"`
	Path                string `json:"path"`
	URL                 string `json:"url"`
	Agent               string `json:"agent"`
	HeaderReferer       string `json:"header_referer"`
	Protocol            string `json:"protocol"`
	Status              int    `json:"status"`
	Method              string `json:"method"`
	Route               string `json:"route"`
	BytesReceived       int    `json:"bytes_received"`
	BytesResponded      int    `json:"bytes_responded"`
	StartTime           string `json:"start_time"`
	EndTime             string `json:"end_time"`
	Duration            int64  `json:"duration"`
}

// Log - logger will print JSON formatted logs onto STDOUT
func Log(ctx *fiber.Ctx) {
	start := time.Now()
	logger := Record{
		IP:                  ctx.IP(),
		HeaderXForwardedFor: string(ctx.Get(fiber.HeaderXForwardedFor)),
		Port:                ctx.Port(),
		Path:                ctx.Path(),
		URL:                 ctx.OriginalURL(),
		StartTime:           start.String(),
		Method:              string(ctx.Method()),
		Route:               string(ctx.Route().Path),
		BytesReceived:       ctx.Request().Header.ContentLength(),
		BytesResponded:      ctx.Response().Header.ContentLength(),
		Agent:               string(ctx.Get(fiber.HeaderUserAgent)),
		HeaderReferer:       string(ctx.Get(fiber.HeaderReferer)),
		Protocol:            string(ctx.Protocol()),
	}

	ctx.Next()

	logger.Status = ctx.Response().StatusCode()
	logger.EndTime = time.Now().String()
	logger.Duration = time.Since(start).Milliseconds()

	logStr, _ := json.Marshal(logger)
	log.Printf("%s", string(logStr))
}
