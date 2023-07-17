package main

import (
	"fmt"
	"os"

	"github.com/imkcat/catchat/internal/modules/app"
)

func main() {
	appInstance, err := app.NewApp()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := appInstance.CliApp.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}
