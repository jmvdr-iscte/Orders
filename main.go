package main

import (
	"context"
	"fmt"

	"github.com/jmvdr-iscte/Orders/application"
)

func main() {
	app := application.New() // importa do meu package

	err := app.Start((context.TODO()))
	if err != nil {
		fmt.Println("failed to start app:", err)
	}
}
