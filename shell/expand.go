package shell

import (
	"strings"
)

func expand(env *Environment, s string) string {
	needsExpand := strings.ContainsAny(s, `$"'\`)
	if !needsExpand {
		return s
	}

	switch s[0] {
	case '\'':
		s = s[1 : len(s)-1]

	case '"':
		s = s[1 : len(s)-1]
		s = expandEscape(s)
		s = expandDollar(env, s)

	default:
		s = expandEscape(s)
		s = expandDollar(env, s)
	}
	return s
}

func expandDollar(env *Environment, src string) string {
	builder := strings.Builder{}
	for len(src) > 0 {
		index := strings.IndexRune(src, '$')
		if index < 0 {
			break
		}
		builder.WriteString(src[:index])
		src = src[index+1:]

		last := strings.IndexAny(src, `$'"`+" \t\n\v")
		if last < 0 {
			last = len(src)
		}
		val, _ := env.Get(src[:last])
		builder.WriteString(val)
		src = src[last:]
	}
	builder.WriteString(src)
	return builder.String()
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
