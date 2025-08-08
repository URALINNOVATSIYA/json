package json

import (
	"fmt"
	"strconv"
	"unicode"
	"unicode/utf16"
	"unicode/utf8"
)

type Json struct {
	json             []byte
	cursor           int
	jsonLen          int
	buffer           []byte
	listBuilder      func() any
	listItemPusher   func(any, int, any) any
	objectBuilder    func() any
	objectItemPusher func(any, string, any) any
}

func Bytes(json []byte) *Json {
	return &Json{
		json:    json,
		cursor:  0,
		jsonLen: len(json),
		buffer:  make([]byte, 0, 128),
	}
}

func (j *Json) WithListBuilder(listBuilder func() any, listItemPusher func(any, int, any) any) *Json {
	j.listBuilder = listBuilder
	j.listItemPusher = listItemPusher
	return j
}

func (j *Json) WithObjectBuilder(objectBuilder func() any, objectItemPusher func(any, string, any) any) *Json {
	j.objectBuilder = objectBuilder
	j.objectItemPusher = objectItemPusher
	return j
}

func (j *Json) Decode() (v any, err error) {
	v, err = j.parseValue()
	if err != nil {
		return nil, err
	}
	if j.cursor >= j.jsonLen {
		return v, nil
	}
	return nil, syntaxError(j.cursor)
}

func (j *Json) parseValue() (any, error) {
	var v any
	var err error

	j.skipWhitespaces()
	for j.cursor < j.jsonLen {
		switch j.json[j.cursor] {
		case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			v, err = j.parseNumber()
		case '"':
			v, err = j.parseString()
		case 't':
			v, err = j.parseTrue()
		case 'f':
			v, err = j.parseFalse()
		case 'n':
			v, err = j.parseNull()
		case '[':
			v, err = j.parseList()
		case '{':
			v, err = j.parseObject()
		default:
			return nil, syntaxError(j.cursor)
		}

		if err != nil {
			return nil, err
		}

		j.skipWhitespaces()
		return v, nil
	}
	return nil, syntaxError(j.cursor)
}

func (j *Json) skipWhitespaces() {
	for j.cursor < j.jsonLen {
		switch j.json[j.cursor] {
		case ' ', '\t', '\r', '\n':
			j.cursor++
			continue
		default:
			return
		}
	}
}

func (j *Json) parseList() (any, error) {
	j.cursor++
	var v any
	var err error
	var key int
	l := make([]any, 0)
	var cl any

	if j.listBuilder != nil {
		cl = j.listBuilder()
	}

	j.skipWhitespaces()
	expectValue := false
	for j.cursor < j.jsonLen {
		switch j.json[j.cursor] {
		case ',':
			if key == 0 || expectValue {
				return nil, syntaxError(j.cursor)
			}
			expectValue = true
			j.cursor++
			continue
		case ']':
			if expectValue {
				return nil, syntaxError(j.cursor)
			}
			j.cursor++
			if j.listBuilder == nil {
				return l, nil
			}
			return cl, nil
		default:
			v, err = j.parseValue()
			if err != nil {
				return nil, err
			}

			if j.listBuilder != nil {
				cl = j.listItemPusher(cl, key, v)
			} else {
				l = append(l, v)
			}

			expectValue = false
			key++
			continue
		}
	}

	return nil, syntaxError(j.cursor)
}

func (j *Json) parseObject() (any, error) {
	var k string
	var v any
	var err error
	var counter int
	expectItem := false
	m := j.makeObject()

	j.cursor++
	j.skipWhitespaces()
	for j.cursor < j.jsonLen {
		switch j.json[j.cursor] {
		case '"':
			k, err = j.parseString()
			if err != nil {
				return nil, err
			}

			j.skipWhitespaces()
			if j.cursor >= j.jsonLen || j.json[j.cursor] != ':' {
				return nil, syntaxError(j.cursor)
			}

			j.cursor++
			v, err = j.parseValue()
			if err != nil {
				return nil, err
			}

			if j.objectBuilder != nil {
				m = j.objectItemPusher(m, k, v)
			} else {
				m.(map[string]any)[k] = v
			}

			counter++
			expectItem = false
			continue
		case ',':
			if counter == 0 || expectItem {
				return nil, syntaxError(j.cursor)
			}
			expectItem = true
			j.cursor++
			j.skipWhitespaces()
			continue
		case '}':
			if expectItem {
				return nil, syntaxError(j.cursor)
			}
			j.cursor++
			return m, nil
		default:
			return nil, syntaxError(j.cursor)
		}
	}
	return nil, syntaxError(j.cursor)
}

func (j *Json) makeObject() any {
	if j.objectBuilder == nil {
		return map[string]any{}
	}

	return j.objectBuilder()
}

func (j *Json) parseNumber() (any, error) {
	if j.json[j.cursor] == '-' && (j.json[j.cursor+1] < 48 || j.json[j.cursor+1] > 57) {
		return nil, syntaxError(j.cursor)
	}

	asFloat := false
	i := j.cursor + 1
	for ; i < j.jsonLen; i++ {
		switch j.json[i] {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			continue
		case '.':
			if asFloat {
				return nil, syntaxError(i)
			}
			asFloat = true
			continue
		}
		break
	}

	if j.jsonLen >= i+2 && (j.json[i] == 'e' || j.json[i] == 'E') &&
		(j.json[i+1] == '+' || j.json[i+1] == '-') &&
		j.json[i+2] > 47 && j.json[i+2] < 58 {
		for i = i + 2; i < j.jsonLen; i++ {
			if j.json[i] > 47 && j.json[i] < 58 {
				continue
			}
			break
		}
	}

	var n any
	n, _ = strconv.ParseFloat(string(j.json[j.cursor:i]), 64)
	j.cursor = i
	return n, nil
}

func (j *Json) parseString() (string, error) {
	defer func() {
		j.buffer = j.buffer[:0]
	}()

	j.cursor++
	for ; j.cursor < j.jsonLen; j.cursor++ {
		switch j.json[j.cursor] {
		case '"':
			j.cursor++
			return string(j.buffer), nil
		case '\\':
			j.cursor++
			switch j.json[j.cursor] {
			case '"', '\\', '/', '\'':
				j.buffer = append(j.buffer, j.json[j.cursor])
			case 'r':
				j.buffer = append(j.buffer, '\r')
			case 'n':
				j.buffer = append(j.buffer, '\n')
			case 't':
				j.buffer = append(j.buffer, '\t')
			case 'b':
				j.buffer = append(j.buffer, '\b')
			case 'f':
				j.buffer = append(j.buffer, '\f')
			case 'u':
				var r rune
				var err error
				r, err = j.parseUnicode()
				if err != nil {
					return "", err
				}
				j.buffer = utf8.AppendRune(j.buffer, r)
			}
		default:
			j.buffer = append(j.buffer, j.json[j.cursor])
		}
	}

	return "", syntaxError(j.cursor)
}

func (j *Json) parseUnicode() (rune, error) {
	j.cursor++

	var r, r2 rune
	var err error
	r, err = j.hexToRune()
	if err != nil {
		return 0, err
	}

	if utf16.IsSurrogate(r) {
		if j.json[j.cursor+1] != '\\' || j.json[j.cursor+2] != 'u' {
			return unicode.ReplacementChar, nil
		}
		j.cursor += 3
		r2, err = j.hexToRune()
		if err != nil {
			return 0, err
		}

		r = utf16.DecodeRune(r, r2)
	}

	return r, nil
}

func (j *Json) hexToRune() (rune, error) {
	if j.jsonLen < j.cursor+4 {
		return 0, syntaxError(j.cursor)
	}

	var r rune
	for _, c := range j.json[j.cursor : j.cursor+4] {
		switch {
		case '0' <= c && c <= '9':
			c = c - '0'
		case 'a' <= c && c <= 'f':
			c = c - 'a' + 10
		case 'A' <= c && c <= 'F':
			c = c - 'A' + 10
		default:
			return 0, syntaxError(j.cursor)
		}
		r = r*16 + rune(c)
	}

	j.cursor += 3
	return r, nil
}

func (j *Json) parseTrue() (any, error) {
	if j.jsonLen < j.cursor+4 {
		return nil, syntaxError(j.cursor)
	}
	if j.json[j.cursor+1] == 'r' && j.json[j.cursor+2] == 'u' && j.json[j.cursor+3] == 'e' {
		j.cursor += 4
		return true, nil
	}
	return nil, syntaxError(j.cursor)
}

func (j *Json) parseFalse() (any, error) {
	if j.jsonLen < j.cursor+5 {
		return nil, syntaxError(j.cursor)
	}
	if j.json[j.cursor+1] == 'a' && j.json[j.cursor+2] == 'l' && j.json[j.cursor+3] == 's' && j.json[j.cursor+4] == 'e' {
		j.cursor += 5
		return false, nil
	}
	return nil, syntaxError(j.cursor)
}

func (j *Json) parseNull() (any, error) {
	if j.jsonLen < j.cursor+4 {
		return nil, syntaxError(j.cursor)
	}
	if j.json[j.cursor+1] == 'u' && j.json[j.cursor+2] == 'l' && j.json[j.cursor+3] == 'l' {
		j.cursor += 4
		return nil, nil
	}
	return nil, syntaxError(j.cursor)
}

func syntaxError(cursor int) error {
	return fmt.Errorf("syntax error at position %d", cursor)
}
