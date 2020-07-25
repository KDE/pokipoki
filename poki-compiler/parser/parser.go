package parser

import (
	"io/ioutil"
	"log"
	"path"
	"regexp"
	"strings"
	"text/scanner"
)

type customScanner struct {
	scanner.Scanner
}

func (s *customScanner) ScanNoEOF() rune {
	ret := s.Scanner.Scan()
	if ret == scanner.EOF {
		log.Fatalf("%s: Unexpected EOF", s.Position)
	}
	return ret
}

func (s *customScanner) ScanIdent() string {
	s.ScanNoEOF()
	if !isIdent(s.TokenText()) {
		log.Fatalf("%s: '%s' is not a valid identifier", s.Position, s.TokenText())
	}
	return s.TokenText()
}

func (s *customScanner) ScanName() string {
	s.ScanNoEOF()
	if !isName(s.TokenText()) {
		log.Fatalf("%s: '%s' is not a valid name", s.Position, s.TokenText())
	}
	return s.TokenText()
}

func (s *customScanner) ScanExpecting(str string) string {
	s.ScanNoEOF()
	if s.TokenText() != str {
		log.Fatalf("%s: Was expecting '%s', got '%s'", s.Position, str, s.TokenText())
	}
	return s.TokenText()
}

func (s *customScanner) ScanToEOL() (ret []string) {
	for s.Peek() != '\n' {
		s.ScanNoEOF()
		ret = append(ret, s.TokenText())
	}
	return
}

var isIdent = regexp.MustCompile(`^[a-z][a-zA-Z]+$`).MatchString
var isName = regexp.MustCompile(`^[A-Z][a-zA-Z]+$`).MatchString

// ParseFile returns a parsed PokiPokiDocument from a file
func ParseFile(file string) (PokiPokiDocument, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return PokiPokiDocument{}, nil
	}

	doku := PokiPokiDocument{
		Objects: map[string]PokiPokiObject{},
	}

	var s customScanner
	s.Init(strings.NewReader(string(data)))
	s.Filename = path.Base(file)

	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		if s.TokenText() != "object" {
			log.Fatalf("%s: not an object definition", s.Position)
		}

		obj := PokiPokiObject{}
		obj.Name = s.ScanName()

		s.ScanExpecting("{")

		for {
			next := s.Scan()
			if next == '}' {
				break
			}

			prop := PokiPokiProperty{}
			if isName(s.TokenText()) {
				obj.Children = append(obj.Children, s.TokenText())
				continue
			} else if !isIdent(s.TokenText()) {
				log.Fatalf("%s: '%s' is not a valid identifier", s.Position, s.TokenText())
			}

			prop.Name = s.TokenText()
			prop.Type = s.ScanToEOL()

			obj.Properties = append(obj.Properties, prop)
		}

		doku.Objects[obj.Name] = obj
	}

	return doku, nil
}
