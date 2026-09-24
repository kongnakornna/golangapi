package http

import (
	"errors"
	"sync"

	"icmongolang/internal/middleware"

	"github.com/go-chi/chi/v5"
)

var once sync.Once

func MapMQTTRoutes(router chi.Router, h *MQTTHandler, mw *middleware.MiddlewareManager) error {
	var err error
	once.Do(func() { err = mapMQTTRoutesOnce(router, h, mw) })
	return err
}

func mapMQTTRoutesOnce(router chi.Router, h *MQTTHandler, mw *middleware.MiddlewareManager) error {
	if router == nil {
		return errors.New("mapMQTTRoutesOnce: router is nil")
	}
	if h == nil {
		return errors.New("mapMQTTRoutesOnce: handler is nil")
	}
	if mw == nil {
		return errors.New("mapMQTTRoutesOnce: middleware manager is nil")
	}

	rateLimit := mw.RateLimit()
	verifier := mw.Verifier(true)
	authenticator := mw.Authenticator()
	currentUser := mw.CurrentUser()
	activeUser := mw.ActiveUser()
	if rateLimit == nil || verifier == nil || authenticator == nil || currentUser == nil || activeUser == nil {
		return errors.New("mapMQTTRoutesOnce: one or more middleware functions returned nil")
	}

	router.Route("/mqtt", func(r chi.Router) {
		r.Use(rateLimit)
		r.Get("/subscriptions", h.Subscriptions)
		r.Get("/status", h.Status)
		r.Get("/gettopicdata", h.GetTopicData)
		r.Get("/devicecontrol", h.DeviceControl)

		r.Group(func(r chi.Router) {
			r.Use(verifier, authenticator, currentUser, activeUser)

			r.Post("/publish", h.Publish)
			r.Post("/subscribe", h.Subscribe)
			r.Post("/unsubscribe", h.Unsubscribe)
		})
	})
	return nil
}

/*
    curl -H "Authorization: Bearer <token>" "http://localhost:5000/api/mqtt/devicecontrol?topic=BAACTW02/CONTROL&message=ON"
	{
		"statuscode": 200,
		"code": 200,
		"topic_control": "BAACTW02/CONTROL",
		"topic_data": "BAACTW02/DATA",
		"message_sent": "ON",
		"payload": "25.5,60,1013",
		"data": ["25.5", "60", "1013"],
		"status": 1,
		"status_msg": "ON",
		"timestamp": "2026-06-08T18:30:00+07:00",
		"message": "Control sent to BAACTW02/CONTROL, response received",
		"message_th": "ส่งคำสั่งไปยัง BAACTW02/CONTROL และได้รับข้อมูลตอบกลับ",
		"from": "mqtt",
		"fetch_duration_ms": 245
	}
*/
