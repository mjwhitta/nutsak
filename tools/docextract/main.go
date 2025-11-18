package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var (
	aliases   = regexp.MustCompile(`^Aliases:`)
	leadSlash = regexp.MustCompile(`^//\s*`)
	newType   = regexp.MustCompile(`^[A-Z-]+:`)
	pkg       = regexp.MustCompile(`package nutsak`)
	types     = regexp.MustCompile(`.*SEED TYPES.*`)
)

func main() {
	var e error
	var f *os.File
	var keep bool
	var line string
	var s *bufio.Scanner
	var sb strings.Builder

	if f, e = os.Open("nutsak.go"); e != nil {
		fmt.Println(e.Error())
		os.Exit(1)
	}
	defer func() {
		if e := f.Close(); e != nil {
			panic(e)
		}
	}()

	s = bufio.NewScanner(f)
	for s.Scan() {
		line = leadSlash.ReplaceAllString(s.Text(), "")

		switch {
		case pkg.MatchString(line):
			fmt.Print(strings.TrimSuffix(sb.String(), "\\n"))
			return
		case types.MatchString(line):
			keep = true
			continue
		}

		if keep {
			switch {
			case line == "":
				sb.WriteString("\\n")
			case aliases.MatchString(line), newType.MatchString(line):
				sb.WriteString(line)
			default:
				sb.WriteString(line + "\\n")
			}
		}
	}

	if s.Err() != nil {
		fmt.Println(s.Err().Error())
	}
}
