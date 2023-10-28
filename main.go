package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/jmvdr-iscte/Orders/application"
)

func main() {
	app := application.New(application.LoadConfig()) // importa do meu package

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt) // só devemos usar o Contect.Background no intuito de gerar outros contexts. pois usar em conjunto com outros pode gerar problemas de concurrencia
	defer cancel()                                                          // simboliza que a função cancel apenas deve ser chamada no fim da função em que está
	err := app.Start(ctx)
	if err != nil {
		fmt.Println("failed to start app:", err)
	}
}
