package application

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jmvdr-iscte/Orders/handler"
)

func loadRoutes() *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) { // http handler para o /
		w.WriteHeader(http.StatusOK)
	})

	router.Route("/orders", loadOrderRoutes) // faz set up a um sub router, o segundo paramtereo é a função na qual recebemos o sub router

	return router
}

func loadOrderRoutes(router chi.Router) { // como está a receber o sub router quaisquer routes que vamos fazer assign a esta função vão ter o perfixo /orders
	orderHandler := &handler.Order{} // fazemos desta variável um pointer ao guardarmos a referencia/memory address da instancia

	router.Post("/", orderHandler.Create)
	router.Get("/", orderHandler.List)
	router.Get("/{id}", orderHandler.GetById)
	router.Put("/{id}", orderHandler.UpdateById)
	router.Delete("/{id}", orderHandler.Delete)
}
