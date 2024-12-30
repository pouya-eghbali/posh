package main

import (
	"flag"

	"github.com/pouya-eghbali/alien-go/pkg/exec"
	"github.com/pouya-eghbali/alien-go/pkg/io"
)

var (
	lolcat = exec.ExternalCommand("lolcat")
	tr     = exec.ExternalCommand("tr")
	echo   = exec.ExternalCommand("echo")
)

func main() {
	var name string
	flag.StringVar(&name, "name", "", "")
	flag.Parse()
	message := io.Format("Hello, %s!", name)
	result := lolcat(tr(echo(&exec.RunContext{}, message), "[:lower:]", "[:upper:]"), "-f").Wait().ToString()
	io.Line(result)
}
