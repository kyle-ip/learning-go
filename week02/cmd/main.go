package cmd

import "week02/internal/router"

func main() {
	r := router.New()
	r.Run(":8080")
}
