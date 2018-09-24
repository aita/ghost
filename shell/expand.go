package shell

import (
	"fmt"
	"strings"
	"unicode"
)

func expand(env *Environment, s string) (string, error) {
	needsExpand := strings.ContainsAny(s, `$"'\`)
	if !needsExpand {
		return s, nil
	}

	var err error
	switch s[0] {
	case '\'':
		s = s[1 : len(s)-1]

	case '"':
		s = s[1 : len(s)-1]
		s = expandEscape(s)
		s, err = expandDollar(env, s)
		if err != nil {
			return "", err
		}
	default:
		s = expandEscape(s)
		s, err = expandDollar(env, s)
		if err != nil {
			return "", err
		}
	}
	return s, nil
}

func expandDollar(env *Environment, src string) (string, error) {
	builder := strings.Builder{}
	for len(src) > 0 {
		index := strings.IndexRune(src, '$')
		if index < 0 {
			break
		}
		builder.WriteString(src[:index])
		src = src[index+1:]
		var name string
		if src[0] == '{' {
			last := strings.IndexRune(src, '}')
			if last < 0 {
				return "", fmt.Errorf("unbalanced {")
			}
			name = src[1:last]
			src = src[last+1:]
		} else {
			last := strings.IndexFunc(src, func(ch rune) bool {
				if unicode.IsSpace(ch) {
					return true
				}
				switch ch {
				case '$', '\'', '"':
					return true
				}
				return false
			})
			if last < 0 {
				last = len(src)
			}
			name = src[:last]
			src = src[last:]
		}
		val, _ := env.Get(name)
		builder.WriteString(val)
	}
	builder.WriteString(src)
	return builder.String(), nil
}

func expandEscape(src string) string {
	builder := strings.Builder{}
	for len(src) > 0 {
		index := strings.IndexRune(src, '\\')
		if index < 0 {
			break
		}
		builder.WriteString(src[:index])
		builder.WriteByte(src[index+1])
		src = src[index+2:]
	}
	builder.WriteString(src)
	return builder.String()
}
