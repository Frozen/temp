package main

import (
	"fmt"
	"github.com/pkg/errors"
)

type TokenType int

func (a TokenType) String() string {
	switch a {
	case TokenTypeIdent:
		return "TokenTypeIdent"
	case TokenTypeSpace:
		return "TokenTypeSpace"
	case TokenTypeInt:
		return "TokenTypeInt"
	case TokenTypeComparsion:
		return "TokenTypeComparsion"
	case TokenTypeDoubleQuotedString:
		return "TokenTypeDoubleQuotedString"
	case TokenTypeSingleQuotedString:
		return "TokenTypeSingleQuotedString"
	case TokenTypeMath:
		return "TokenTypeMath"
	case TokenTypeFloat:
		return "TokenTypeFloat"
	default:
		//fmt.Println("unknown type ", int(a))
		//panic("unknown type")
		return "unknown type"
	}
}

const (
	TokenTypeIdent TokenType = iota + 1
	TokenTypeSpace
	TokenTypeDoubleQuotedString
	TokenTypeSingleQuotedString
	TokenTypeMath
	TokenTypeComparsion
	TokenTypeInt
	TokenTypeFloat
)

type ParseError interface {
}

func NewParseError(line, pos int, mess string) {

}

type Token struct {
	Type  TokenType
	Line  int
	Pos   int
	Value []byte
}

func tokenaze(bts []byte) ([]Token, error) {

	var tokens []Token

	//i := 0
	iter := NewIterator(bts)

	for {

		next := iter.NextNonShift()
		if next.eof {
			return tokens, nil
		}

		fmt.Println("next", next.pos, string(next.b))

		switch {
		case isAlpha(next.b):
			tokens = append(tokens, readIdent(iter))
		case isSpace(next.b):
			token := readSpace(iter)
			fmt.Println(token)
			tokens = append(tokens, token)
		case isDoubleQuoter(next.b):
			token := readDoubleQuoterString(iter)
			//fmt.Println("readDoubleQuoterString:", token)
			tokens = append(tokens, token)
		case isSingleQuoter(next.b):
			tokens = append(tokens, readSingleQuoterString(iter))
		case isMathChar(next.b):
			token, err := readMath(iter)
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, token)

		case isNum(next.b):
			tokens = append(tokens, readNum(iter))

		default:
			//panic("nn")
			fmt.Println("unknown char", string(next.b))
			return tokens, nil
		}
	}
}

func readIdent(iter *Iterator) Token {

	b := iter.Next()
	start := b.pos
	for {

		b := iter.Next()
		if b.eof || !or(b.b, isAlphaNum, char('_')) {
			return Token{
				Type:  TokenTypeIdent,
				Value: iter.Bytes()[start : b.pos+1],
			}
		}
	}
}

func readSpace(iter *Iterator) Token {
	start := iter.Next()

	for {
		b := iter.NextNonShift()
		if b.eof || !isSpace(b.b) {
			return Token{
				Type:  TokenTypeSpace,
				Value: iter.Bytes()[start.pos:b.pos],
			}
		}
		iter.Next()
	}
}

func readDoubleQuoterString(iter *Iterator) Token {
	start := iter.Next()

	for {

		b := iter.Next()
		if b.eof || isDoubleQuoter(b.b) {
			return Token{
				Type:  TokenTypeDoubleQuotedString,
				Value: iter.Bytes()[start.pos : b.pos+1],
			}
		}
	}
}

func readSingleQuoterString(iter *Iterator) Token {
	start := iter.Next()

	for {
		b := iter.Next()
		if b.eof || isSingleQuoter(b.b) {
			return Token{
				Type:  TokenTypeSingleQuotedString,
				Value: iter.Bytes()[start.pos : b.pos+1],
			}
		}
	}
}

func readMath(iter *Iterator) (Token, error) {
	start := iter.Next()

	//fmt.Println("start", start)

	next := iter.NextNonShift()
	if !isMathChar(next.b) {
		return Token{
			Type:  TokenTypeMath,
			Value: iter.Bytes()[start.pos : start.pos+1],
		}, nil
	}

	next = iter.Next()

	switch s := string([]byte{start.b, next.b}); s {
	case "==", ">=", "<=": // not equal??
		return Token{
			Type:  TokenTypeComparsion,
			Value: iter.Bytes()[start.pos : next.pos+1],
		}, nil
	default:
		return Token{}, errors.Errorf("invalid character %s", s)
	}
}

func readNum(iter *Iterator) Token {
	start := iter.Next()

	for {

		//fmt.Println("aa")

		next := iter.NextNonShift()
		if next.eof || !isNum(next.b) {
			return Token{
				Type:  TokenTypeInt,
				Value: iter.Bytes()[start.pos : next.pos+1],
			}
		}
		iter.Next()
	}
}

func isAlpha(b byte) bool {
	return b >= 'a' && b <= 'z' || b >= 'A' && b <= 'Z'
}

func isNum(b byte) bool {
	return b >= '0' && b <= '9'
}

func isAlphaNum(b byte) bool {
	return isAlpha(b) || isNum(b)
}

func char(is byte) func(b byte) bool {
	return func(b byte) bool {
		return b == is
	}
}

func or(b byte, predicats ...func(b byte) bool) bool {
	for _, f := range predicats {
		if f(b) {
			return true
		}
	}
	return false
}

func isSpace(b byte) bool {
	switch b {
	case ' ', '\t', '\n':
		return true
	default:
		return false
	}
}

func isDoubleQuoter(b byte) bool {
	return b == '"'
}

func isSingleQuoter(b byte) bool {
	return b == '\''
}

func isMathChar(b byte) bool {
	switch b {
	case '=', '>', '<', '+', '-':
		return true
	default:
		return false
	}
}

type Iterator struct {
	bts       []byte
	index     int
	lastIndex int
}

type Next struct {
	b   byte
	pos int
	eof bool
}

func NewIterator(bts []byte) *Iterator {
	return &Iterator{
		bts:       bts,
		index:     0,
		lastIndex: len(bts) - 1,
	}
}

func (a *Iterator) Next() Next {
	a.index += 1
	if a.index <= a.lastIndex {
		return Next{b: a.bts[a.index-1], pos: a.index - 1}
	}
	return Next{
		b:   0,
		pos: a.index - 1,
		eof: true,
	}
}

func (a *Iterator) Bytes() []byte {
	return a.bts
}

func (a *Iterator) NextNonShift() Next {
	if a.index <= a.lastIndex {
		return Next{
			b:   a.bts[a.index],
			pos: a.index,
		}
	}
	return Next{
		b:   0,
		pos: a.index,
		eof: true,
	}
}

var s = `let x = 6`

//var s = ` "4444"`

func main() {

	//s := &S{}

	//s.ReadFrom("")

	rs, _ := tokenaze([]byte(s))
	fmt.Println("tokens: ", len(rs))
	for _, row := range rs {
		fmt.Printf("%s %s\n", row.Type, string(row.Value))
	}

}
