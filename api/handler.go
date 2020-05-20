package api

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	rh "github.com/rinosukmandityo/maknews/repositories/helper"
	"github.com/rinosukmandityo/maknews/services/logic"
)

func RegisterHandler() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	newsRepo := rh.ChooseRepo()
	elasticRepo := rh.ElasticRepo()
	kafkaRepo := rh.KafkaConnection()
	redisRepo := rh.RedisRepo()

	newsSvc := logic.NewNewsService(newsRepo, redisRepo, elasticRepo, kafkaRepo)
	kafkaSvc := logic.NewKafkaService(kafkaRepo)

	go func() { // just assume that this is another service that register kafka topic
		kafkaSvc.ReadMessage(newsRepo, elasticRepo)
	}()

	registerNewsHandler(r, NewNewsHandler(newsSvc))

	return r
}

func registerNewsHandler(r *chi.Mux, handler NewsHandler) {
	// Subrouters:
	r.Route("/news", func(r chi.Router) {
		r.Post("/", handler.Post) // POST /news
		r.Get("/", handler.Get)   // GET /news?offset=0&limit=10
		// Subrouters:
		r.Route("/{id}", func(r chi.Router) {
			r.Use(handler.NewsCtx)
			r.Put("/", handler.Update)    // PUT /news/newsid01
			r.Delete("/", handler.Delete) // DELETE /news/newsid01
		})
	})
}
