package parser

import (
	"path"
	"strings"

	"github.com/pouya-eghbali/posh/pkg/lang/parser/rules"
	"github.com/pouya-eghbali/posh/pkg/lang/parser/types"
	"github.com/pouya-eghbali/posh/pkg/lang/parser/utils"
)

func CompileMainFile(inputPath string, outputName string, astOutput bool) error {
	temp, err := utils.CreateTempDir()
	if err != nil {
		return err
	}

	baseDir := path.Dir(inputPath)
	filePath := strings.TrimPrefix(inputPath, baseDir+"/")
	posh := types.NewPoshFile(filePath, baseDir, temp, "main", map[string]types.CompiledFile{})
	err = utils.CompilePoshFile(posh, rules.MatchPosh)

	if err != nil {
		return err
	}

	return utils.CompileTempDir(temp, outputName)
}
