package main

import (
	"flag"

	"github.com/pouya-eghbali/alien-go/pkg/exec"
	"github.com/pouya-eghbali/alien-go/pkg/io"
)

var (
	tr     = exec.ExternalCommand("tr")
	lolcat = exec.ExternalCommand("lolcat")
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
