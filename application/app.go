package application

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

type App struct { // guarda as dependencias da aplicação
	router http.Handler  // para ficar desassociado do chi
	rdb    *redis.Client // guardar o client de redis
}

func New() *App {
	app := &App{ // cria uma instancia da app e vaz assign à variavel
		rdb: redis.NewClient(&redis.Options{}), // gere o client internament
	}

	app.loadRoutes()
	return app
}

func (a *App) Start(ctx context.Context) error { // criar um receiver, é como se fosse o owner do metodo, recebe um pointer para a instancia da app
	server := &http.Server{
		Addr:    ":3000",
		Handler: a.router,
	}
	err := a.rdb.Ping(ctx).Err() //:= realiza assign e inicialização
	if err != nil {
		return fmt.Errorf("failed to connect to redis: %w", err)
	}

	defer func() { // criamos uma função anonima para poder apanhar os erros, não funciona chamar no fim pois a função principal devolve um erro
		if err := a.rdb.Close(); err != nil {
			fmt.Println("failed to close redis", err)
		}
	}()

	fmt.Println("Starting server")

	ch := make(chan error, 1) // criar um channel para a goroutine voltar ao main thread, cada channel tem o seu assigned type que será enviado pelas threafds. O segundo parametro simboliza o tamanho do buffer a enviar. Num cannal bufferd owritter vai estar blockeado até que o sinal seja lido, num unbuffered channel pode continuar a escrever até chegar ao tamanho do buffer.

	go func() { // GO ROUTINE :Inicializa outra thread e garante que não bloqueia a nossa main thread
		err = server.ListenAndServe()
		if err != nil {
			ch <- fmt.Errorf("failed to start server: %w", err)
		}
		close(ch)
	}()

	select { // é um switch para channels
	case err := <-ch: // vai captar todos os valores deste canal para esta variável é o receiver. O receiver vai parar a execução do código até o channel estar fechado. pode-se verificar também por outro parametro o $open se o channel está aberto ou fechado
		return err
	case <-ctx.Done(): // devolve um channel se o channel foi cancelado
		timeout, cancel := context.WithTimeout(context.Background(), time.Second*10) // 10 segundos para todos os processos terminarem
		defer cancel()

		return server.Shutdown(timeout)
	}
}
