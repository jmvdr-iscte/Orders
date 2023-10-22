package application

import (
	"context"
	"fmt"
	"net/http"
)

type App struct { // guarda as dependencias da aplicação
	router http.Handler // para ficar desassociado do chi
}

func New() *App {
	app := &App{ // cria uma instancia da app e vaz assign à variavel
		router: loadRoutes(),
	}
	return app
}

func (a *App) Start(ctx context.Context) error { // criar um receiver, é como se fosse o owner do metodo, recebe um pointer para a instancia da app
	server := &http.Server{
		Addr:    ":3000",
		Handler: a.router,
	}
	err := server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
	return nil
}
