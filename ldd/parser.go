package ldd

import(
	"text/scanner"
	"fmt"
)

/*
linux-vdso.so.1 (0x00007ffdbd15e000)
libreadline.so.6 => /usr/lib/libreadline.so.6 (0x00007fc5d3591000)
libncursesw.so.6 => /usr/lib/libncursesw.so.6 (0x00007fc5d3324000)
libdl.so.2 => /usr/lib/libdl.so.2 (0x00007fc5d3120000)
libc.so.6 => /usr/lib/libc.so.6 (0x00007fc5d2d7f000)
/lib64/ld-linux-x86-64.so.2 (0x00007fc5d37db000)
*/

type Parser struct {
	scanner scanner.Scanner

	items   chan item
}

type itemType int

const(
	itemError itemType = iota
	itemFilename
)

type item struct {
	typ      itemType
	val      string
}

type stateFn func(*Parser) stateFn

var lastSoNameParsed string

func parseHexAddress( p * Parser ) stateFn {
	p.scanner.Scan()


	if p.scanner.TokenText()[0:2] != "0x" {
		return p.errorf( "Token \"%s\" is not a hex address, should start with \"0x\", not \"%s\"", p.scanner.TokenText(), p.scanner.TokenText()[0:2] )
	}

	for _,c := range p.scanner.TokenText()[2:] {
		if ! ( c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z' || c >= '0' && c <= '9' ) {
			return p.errorf( "Token \"%s\" is not a hex address, char '%c' is not a hex character", p.scanner.TokenText(), c )
		}
	}

	p.scanner.Scan()
	if p.scanner.TokenText() != ")" {
		return p.errorf( "Expected character \")\", got \"%s\"", p.scanner.TokenText() )
	}

	p.emit( lastSoNameParsed )

	return parseSoName
}

func parseArrow( p * Parser ) stateFn {
	p.scanner.Scan()

	if p.scanner.TokenText() != ">" {
		return p.errorf( "Expected '>', but got \"%s\"", p.scanner.TokenText() )
	}

	return parseSoName
}

func parseSoName( p * Parser ) stateFn {
	var tok rune
	var text string = ""
	var nextState stateFn

	for tok != scanner.EOF {
		tok = p.scanner.Scan()

		switch p.scanner.TokenText() {
			case "(":
				nextState = parseHexAddress
			case "=":
				nextState = parseArrow
			default:
				text += p.scanner.TokenText()
		}

		if nextState != nil {
			lastSoNameParsed = text

			return nextState
		}
	}

	if tok == scanner.EOF {
		return nil
	}

	return p.errorf( "Parser state should not end up here" )
}

func NewParser( scanner scanner.Scanner ) (* Parser, chan item) {
	p := & Parser {
		scanner: scanner,
		items: make(chan item),
	}

	go p.run()
	return p, p.items
}

func (p *Parser) errorf( format string, args ... interface{} ) stateFn {
	p.items <- item{
		itemError,
		fmt.Sprintf(format, args...),
	}

	return nil
}

func (p *Parser) emit( dependency string ) {
	p.items <- item{
		itemFilename,
		dependency,
	}
}

func (p *Parser) run() {

	for state := parseSoName; state != nil; {
		state = state( p )
	}

	close( p.items )
}

func (p *Parser) Next() string {
	var tok rune
	for tok != scanner.EOF {
		tok = p.scanner.Scan()
		fmt.Println("At position", p.scanner.Pos(), ":", p.scanner.TokenText())
		fmt.Println( tok )
	}

	return ""
}
