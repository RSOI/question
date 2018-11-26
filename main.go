// QUESTION SERVICE

package main

import (
	"fmt"

	"github.com/RSOI/question/database"
	"github.com/valyala/fasthttp"
)

func main() {
	const PORT = 8080
	database.Connect()
	fmt.Println("Question service is starting on localhost:", PORT)
	fasthttp.ListenAndServe(fmt.Sprintf(":%d", PORT), initRoutes().Handler)
}
