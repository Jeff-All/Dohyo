// You can edit this code!
// Click here and start typing.
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Jeff-All/Dohyo/app"
)

func main() {
	fmt.Println("starting Dohyo")
	log.SetOutput(os.Stdout)
	app.ExecApp()
}
