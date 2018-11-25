// QUESTION SERVICE

package main

import (
	"github.com/valyala/fasthttp"
)

func main() {
	fasthttp.ListenAndServe(":8080", initRoutes().Handler)
}
