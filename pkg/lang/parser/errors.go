package parser

import (
	"fmt"
	"strings"

	"github.com/pouya-eghbali/posh/pkg/lang/parser/types"
)

func PrintError(code string, failedAt *types.Node) {
	pos := (*failedAt).GetPos()
	image := (*failedAt).GetImage()
	lines := strings.Split(code, "\n")
	line := lines[pos.Line]
	fmt.Println(line)
	fmt.Print(strings.Repeat(" ", pos.Column-len(image)))
	fmt.Println("^")
	fmt.Print(strings.Repeat(" ", pos.Column-len(image)))
	fmt.Printf("Error at %d:%d\n", pos.Line, pos.Column)
}
