package expr

import (
	"fmt"
	"github.com/seeadoog/expr/ast"
	"strconv"
)

type tokenV struct {
	tkn  string
	kind int
	num  float64
	x, y int
}

type tokenizer struct {
	next        func(c rune) error
	tokens      []tokenV
	tkn         []rune
	exp         []rune
	pos         int
	xy          int
	y           int
	currentStat int
}

func isVariableConstraint(s string) (any, bool) {
	switch s {
	case "true":
		return true, true
	case "false":
		return false, true
	case "nil":
		return nil, true
	}
	f, err := strconv.ParseFloat(s, 64)
	return f, err == nil

}

func (t *tokenizer) X() int {
	return t.pos - t.xy
}
func (t *tokenizer) Y() int {
	return t.y + 1
}

func parseTokenizer(exp string) ([]tokenV, error) {
	t := tokenizer{
		tokens: []tokenV{},
		exp:    []rune(exp),
	}
	t.next = t.statStart
	r := []rune(exp)
	for t.pos = 0; t.pos < len(r); t.pos++ {
		c := r[t.pos]
		err := t.next(c)
		switch c {
		case '\n':
			t.y++
			t.xy = t.pos
		case '\r':
		default:

		}

		if err != nil {
			return nil, fmt.Errorf("parse exp error as token error:%w '%v' at:%d:%d", err, exp, t.Y(), t.X())
		}
	}
	if len(t.tkn) > 0 {

		if t.currentStat == stateInString {
			return nil, fmt.Errorf("string is not closed: '%s'  at %d:%d", string(t.tkn), t.Y(), t.X())
		}
		t.tokens = append(t.tokens, tokenV{
			tkn: string(t.tkn),
			x:   t.X(),
			y:   t.Y(),
		})
	}
	return t.tokens, nil

}

func (t *tokenizer) appendToken(kind int, raw string) {
	if len(t.tkn) > 0 {

		seg := string(t.tkn)
		kd := t.getTknKind(seg)
		t.tokens = append(t.tokens, tokenV{
			tkn:  seg,
			kind: kd,
			x:    t.X() - 1,
			y:    t.Y(),
		})
		//t.tkn = t.tkn[:0]
		//t.tokens = append(t.tokens, tokenV{
		//	tkn:  string(t.tkn),
		//	kind: variables,
		//})
	}

	t.tokens = append(t.tokens, tokenV{
		tkn:  raw,
		kind: kind,
		x:    t.X(),
		y:    t.Y(),
	})
	t.tkn = t.tkn[:0]

}

func (t *tokenizer) getNext() (rune, bool) {
	t.pos++
	if t.pos >= len(t.exp) {
		return 0, false
	}
	return t.exp[t.pos], true
}

func (t *tokenizer) pre() (rune, bool) {
	if t.pos-1 < 0 {
		return 0, false
	}
	return t.exp[t.pos-1], true
}

func (t *tokenizer) getTknKind(seg string) int {
	kind := variables
	switch seg {
	case "or":
		kind = ast.ORR
	case "const":
		kind = ast.CONST
	case "in":
		kind = ast.IN
	case "as":
		kind = ast.AS
	}
	return kind
}

func (t *tokenizer) appendId() {
	if len(t.tkn) > 0 {
		seg := string(t.tkn)
		kind := t.getTknKind(seg)
		t.tokens = append(t.tokens, tokenV{
			tkn:  seg,
			kind: kind,
			x:    t.X(),
			y:    t.Y(),
		})
		t.tkn = t.tkn[:0]
	}
}

func (t *tokenizer) statStart(r rune) error {
	switch r {
	case '(', ')', ';', '{', '}', '[', ']', '%':
		t.appendToken(int(r), string(r))

	case '?':
		c, ok := t.getNext()
		if !ok {
			return fmt.Errorf("unexpected  eof after '?'")
		}
		if c == '?' {
			t.appendToken(ast.NONIL, "??")
			return nil
		}
		t.pos--
		t.appendToken(int(r), "?")
	case '#':
		t.next = func(c rune) error {
			switch c {
			case '\n', '\r':
				t.next = t.statStart
			}
			return nil
		}
	case ':':
		c, ok := t.getNext()
		if !ok {
			return fmt.Errorf("unexpected  eof after ':'")
		}
		if c == ':' {
			t.appendToken(ast.ACC, "::")
			return nil
		}
		t.pos--
		t.appendToken(int(r), ":")
	case '\'':
		t.next = t.statStringStart
	case '`':
		t.next = t.statStringStartWith('`')
	case '"':
		t.next = t.statStringStartWith('"')
	case ',':
		t.appendToken(',', ",")
	case ' ', '\t', '\n', '\r':
		t.appendId()

	case '+':
		c, ok := t.getNext()
		if !ok {
			return fmt.Errorf("unexpected  eof after '+'")
		}
		if c == '=' {
			t.appendToken(ast.ADDEQ, "+=")
			return nil
		}
		if c == '+' {
			t.appendToken(ast.ADDADD, "++")
			return nil
		}
		t.pos--
		t.appendToken(int(r), "+")
	case '*', '/', '^', '@':
		if r == '*' || r == '/' {
			c, ok := t.getNext()
			if !ok {
				return fmt.Errorf("unexpected eof after '%c'", r)
			}
			if c == '=' {
				if r == '*' {
					t.appendToken(ast.MULEQ, "*=")
				} else {
					t.appendToken(ast.DIVEQ, "/=")
				}
				return nil
			}
			t.pos--
		}
		t.appendToken(int(r), string(r))

	case '-':
		c, ok := t.getNext()
		if !ok {
			return fmt.Errorf("unexpected  eof after '-'")
		}
		if c == '>' {
			t.appendToken(ast.ACC, "->")
			return nil
		}
		if c == '-' {
			t.appendToken(ast.SUBSUB, "--")
			return nil
		}
		if c == '=' {
			t.appendToken(ast.SUBEQ, "-=")
			return nil
		}
		t.pos--
		t.appendToken(int(r), "-")
	case '!':
		c, ok := t.getNext()
		if !ok {
			return fmt.Errorf("unexpected  eof after '!'")
		}
		if c == '=' {

			c2, ok := t.getNext()
			if ok {
				if c2 == '=' {
					t.appendToken(ast.NOTEQT, "!==")
					return nil
				}
			}
			t.pos--
			t.appendToken(ast.NOTEQ, "!=")
			return nil
		}
		if c == '!' {
			t.appendToken(ast.NONIL, "!!")
			return nil
		}
		t.pos--
		t.appendToken(int(r), "!")
	case '=':
		//t.next = t.statParseEq
		c, ok := t.getNext()
		if !ok {
			return fmt.Errorf("unexpected  eof after '='")
		}
		if c == '=' {

			c2, ok := t.getNext()
			if ok {
				if c2 == '=' {
					t.appendToken(ast.EQT, "===")
					return nil
				}
			}
			t.pos--

			t.appendToken(ast.EQ, "==")
			return nil
		}
		if c == '>' {
			t.appendToken(ast.LAMB, "=>")
			return nil
		}
		t.pos--
		t.appendToken(int(r), "=")

	case '|':
		t.next = t.statParseOr
	case '&':
		t.next = t.statParseAND
	case '>':
		c, ok := t.getNext()
		if !ok {
			return fmt.Errorf("unexpected  eof after '>'")
		}
		if c == '=' {
			t.appendToken(ast.GTE, ">=")
			return nil
		}
		t.pos--
		t.appendToken(ast.GT, ">")
	case '<':
		c, ok := t.getNext()
		if !ok {
			return fmt.Errorf("unexpected  eof after '>'")
		}
		if c == '=' {
			t.appendToken(ast.LTE, "<=")
			return nil
		}
		t.pos--
		t.appendToken(ast.LT, "<")
	case '.':
		c, ok := t.getNext()
		if !ok {
			return fmt.Errorf("unexpected  eof after '.'")
		}
		if c == '.' {
			c, ok := t.getNext()
			if !ok {
				return fmt.Errorf("unexpected  eof after '..'")
			}
			if c == '.' {
				t.appendToken(ast.VARIADIC, "...")
				return nil
			}
			t.pos--
			return nil
		}
		t.pos--
		t.appendToken(ast.ACC, ".")
	default:
		t.tkn = append(t.tkn, r)
		if len(t.tkn) == 1 {
			if r >= '0' && r <= '9' {
				t.next = t.parseNumber
			}
		}

	}
	return nil
}

func pointNum(r []rune) int {
	s := 0
	for _, n := range r {
		if n == '.' {
			s++
		}
	}
	return s
}

func (t *tokenizer) parseNumber(c rune) error {
	if (c >= '0' && c <= '9') || c == '.' || c == 'x' {
		t.tkn = append(t.tkn, c)

		if pointNum(t.tkn) > 1 {
			return fmt.Errorf("parser invalid number: %s", string(t.tkn))
		}
		return nil
	}
	//for t.pos < len(t.exp) {
	//	c := t.exp[t.pos]
	//	t.pos++
	//	if (c >= '0' && c <= '9') || c == '.' || c == 'x' {
	//		t.tkn = append(t.tkn, c)
	//
	//		if pointNum(t.tkn) > 1 {
	//			return fmt.Errorf("parser invalid number: %s", string(t.tkn))
	//		}
	//		return nil
	//	} else {
	//		break
	//	}
	//}
	//
	//s := string(t.tkn)
	//t.tkn = t.tkn[:0]
	//var n float64
	//var err error
	//if pointNum(t.tkn) == 1 {
	//	n, err = strconv.ParseFloat(s, 64)
	//	if err != nil {
	//		return fmt.Errorf("parser invalid number: %s", s)
	//	}
	//
	//} else {
	//	var n1 int64
	//	n1, err = strconv.ParseInt(s, 0, 64)
	//	if err != nil {
	//		return fmt.Errorf("parser invalid number: %s", s)
	//	}
	//	n = float64(n1)
	//}

	//t.tokens = append(t.tokens, tokenV{
	//	tkn:  s,
	//	kind: number,
	//	num:  n,
	//})
	//if t.pos <= len(t.exp) {
	//	t.pos--
	//}
	t.pos--
	t.next = t.statStart
	return nil
}

func (t *tokenizer) statParseAND(r rune) error {
	if r != '&' {
		t.appendToken('&', "&")
		t.pos--
		t.next = t.statStart
		return nil
		//return errors.New("invalid token after & ")
	}
	t.appendToken(ast.AND, "&&")
	t.next = t.statStart
	return nil
}
func (t *tokenizer) statParseOr(r rune) error {
	if r != '|' {
		t.appendToken('|', "|")
		t.pos--
		t.next = t.statStart
		return nil
		//return errors.New("invalid token after | ")
	}
	t.appendToken(ast.OR, "||")
	t.next = t.statStart
	return nil
}

func (t *tokenizer) statStringStart(r rune) error {
	switch r {
	case '\'':
		t.tokens = append(t.tokens, tokenV{
			tkn:  string(t.tkn),
			kind: constant,
			x:    t.X(),
			y:    t.Y(),
		})
		t.tkn = t.tkn[:0]
		t.next = t.statStart
	case '\\':
		t.next = t.escapeNext(t.statStringStart)
	default:
		t.tkn = append(t.tkn, r)

	}
	return nil
}

const (
	statStart     = 0
	stateInString = 1
)

func (t *tokenizer) statStringStartWith(c rune) func(c rune) error {
	var fff func(rune) error
	fff = func(r rune) error {
		t.currentStat = stateInString
		switch r {
		case c:
			t.tokens = append(t.tokens, tokenV{
				tkn:  string(t.tkn),
				kind: constant,
			})
			t.tkn = t.tkn[:0]
			t.next = t.statStart
			t.currentStat = statStart
		case '\\':
			t.next = t.escapeNext(fff)
		default:
			t.tkn = append(t.tkn, r)
		}
		return nil
	}
	return fff
}

func (t *tokenizer) escapeNext(statFunc func(c rune) error) func(c rune) error {
	return func(c rune) error {
		switch c {
		case 'n':
			t.tkn = append(t.tkn, '\n')
		case 'r':
			t.tkn = append(t.tkn, '\r')
		case 'b':
			t.tkn = append(t.tkn, '\b')
		case 'v':
			t.tkn = append(t.tkn, '\v')
		case 'a':
			t.tkn = append(t.tkn, '\a')
		case 'f':
			t.tkn = append(t.tkn, '\f')
		case 't':
			t.tkn = append(t.tkn, '\t')
		default:
			t.tkn = append(t.tkn, c)
		}
		t.next = statFunc
		return nil
	}
}
