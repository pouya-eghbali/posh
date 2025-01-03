package lexer

import (
	"fmt"
	"regexp"

	"github.com/pouya-eghbali/posh/pkg/lang/parser/types"
)

type Pattern struct {
	Name string
	Re   *regexp.Regexp
}

// Name enum:

var patterns = []Pattern{
	{Name: "WHITESPACE", Re: regexp.MustCompile(`^(\s+)`)},
	{Name: "COMMENT", Re: regexp.MustCompile(`^(#.*\n)`)},
	{Name: "KEYWORD", Re: regexp.MustCompile(`^(fn|if|else|elif|and|or|not|return|true|false|import|from|as|for|in|break|continue)\s`)},
	{Name: "IDENTIFIER", Re: regexp.MustCompile(`^([a-zA-Z_][a-zA-Z0-9-_]*)`)},
	{Name: "PUNCTUATOR", Re: regexp.MustCompile(`^([{}()[\]<>,.;+-/*%=|!])`)},
	{Name: "INTEGER", Re: regexp.MustCompile(`^([0-9]+)`)},
	{Name: "FLOAT", Re: regexp.MustCompile(`^([0-9]+\.[0-9]+)`)},
	{Name: "STRING", Re: regexp.MustCompile(`^"(\\.|[^"]*)"`)},
}

func Lex(code string) ([]types.Node, error) {
	offset := 0
	line := 0
	column := 0

	tokens := []types.Node{}

	for {
		if offset >= len(code) {
			break
		}

		for _, pattern := range patterns {
			if match := pattern.Re.FindString(code[offset:]); match != "" {
				if pattern.Name != "WHITESPACE" && pattern.Name != "COMMENT" {

					image := match
					if pattern.Name == "KEYWORD" {
						image = match[:len(match)-1]
					}

					tokens = append(tokens, &types.TokenNode{
						BaseNode: types.BaseNode{Type: pattern.Name},
						Pos:      &types.Pos{Line: line, Column: column},
						Image:    image,
					})
				}

				if pattern.Name == "WHITESPACE" {
					for _, c := range match {
						if c == '\n' {
							line++
							column = 0
						} else {
							column++
						}
					}
				} else {
					for _, c := range match {
						if c == '\n' {
							line++
							column = 0
						} else {
							column++
						}
					}
				}

				offset += len(match)
				break
			} else if pattern.Name == "STRING" {
				return nil, fmt.Errorf("invalid string at %d:%d", line, column)
			}
		}
	}

	return tokens, nil
}
