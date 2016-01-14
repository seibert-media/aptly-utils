package architecture

import "bytes"

type Architecture string

const (
	ALL     = Architecture("all")
	I386    = Architecture("i386")
	AMD64   = Architecture("amd64")
	DEFAULT = AMD64
)

func Join(architectures []Architecture, sep string) string {
	b := bytes.NewBufferString("")
	first := true
	for _, a := range architectures {
		if first {
			first = false
		} else {
			b.WriteString(sep)
		}
		b.WriteString(string(a))
	}
	return string(b.Bytes())
}

func Parse(architectures []string) []Architecture {
	var result []Architecture
	for _, name := range architectures {
		result = append(result, Architecture(name))
	}
	return result
}
