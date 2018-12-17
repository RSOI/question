// QUESTION SERVICE

package main

import (
	"fmt"
	"os"

	"github.com/RSOI/question/controller"
	"github.com/RSOI/question/database"
	"github.com/RSOI/question/utils"
	"github.com/valyala/fasthttp"
)

// PORT application port
const PORT = 8080

func main() {
	if len(os.Args) > 1 {
		utils.DEBUG = os.Args[1] == "debug"
	}
	utils.LOG("Launched in debug mode...")
	utils.LOG(fmt.Sprintf("Question service is starting on localhost: %d", PORT))

	controller.Init(database.Connect())
	fasthttp.ListenAndServe(fmt.Sprintf(":%d", PORT), initRoutes().Handler)
}
