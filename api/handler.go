package api

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/rinosukmandityo/maknews/logic"
	rh "github.com/rinosukmandityo/maknews/repositories/helper"
)

func RegisterHandler() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	newsSvc := logic.NewNewsService(rh.ChooseRepo())
	elasticSvc := logic.NewElasticService(rh.ElasticRepo())
	kafkaSvc := logic.NewKafkaService()
	redisSvc := logic.NewRedisService(rh.RedisRepo())

	go func() {
		kafkaSvc.ReadMessage(newsSvc, elasticSvc)
	}()

	registerNewsHandler(r, NewNewsHandler(newsSvc, elasticSvc, kafkaSvc, redisSvc))

	return r
}

func registerNewsHandler(r *chi.Mux, handler NewsHandler) {
	r.Get("/news", handler.Get)
	r.Post("/news", handler.Post)
	r.Post("/update", handler.Update)
	r.Post("/delete", handler.Delete)
}
