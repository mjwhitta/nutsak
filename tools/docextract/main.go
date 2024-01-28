package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
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
	var out string
	var s *bufio.Scanner

	if f, e = os.Open("nutsak.go"); e != nil {
		fmt.Println(e.Error())
		os.Exit(1)
	}
	defer f.Close()

	s = bufio.NewScanner(f)
	for s.Scan() {
		line = leadSlash.ReplaceAllString(s.Text(), "")

		switch {
		case pkg.MatchString(line):
			fmt.Print(out[:len(out)-2])
			return
		case types.MatchString(line):
			keep = true
			continue
		}

		if keep {
			switch {
			case line == "":
				out += "\\n"
			case aliases.MatchString(line):
				out += line
			case newType.MatchString(line):
				out += line
			default:
				out += line + "\\n"
			}
		}
	}

	if s.Err() != nil {
		fmt.Println(s.Err().Error())
		os.Exit(2)
	}
}
