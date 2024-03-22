package main

import (
	"fmt"
	"slackwatch/backend/internal/app"
)

func main() {
	application, err := app.Initialize()
	if err != nil {
		fmt.Println("Error initializing application:", err)
		return
	}

	if err := application.Run(); err != nil {
		fmt.Println("Error running application:", err)
	}
}
