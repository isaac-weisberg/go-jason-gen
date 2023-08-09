package main

import (
	"fmt"
	"strings"
)

type customBuilder struct {
	builder strings.Builder
}

func (customBuilder *customBuilder) WriteString(values ...string) {
	for _, val := range values {
		customBuilder.builder.WriteString(val)
	}
}

func (customBuilder *customBuilder) WriteLine(values ...string) {
	for _, val := range values {
		customBuilder.builder.WriteString(val)
	}
	customBuilder.builder.WriteByte('\n')
}

func (customBuilder *customBuilder) WriteLineFI(indentation int, format string, args ...any) {
	for i := 0; i < indentation; i++ {
		customBuilder.builder.WriteRune('\t')
	}
	customBuilder.builder.WriteString(fmt.Sprintf(format, args...))
	customBuilder.builder.WriteRune('\n')
}

func (customBuilder *customBuilder) WriteLineIndent(indentation int, values ...string) {
	for i := 0; i < indentation; i++ {
		customBuilder.builder.WriteRune('\t')
	}
	for _, val := range values {
		customBuilder.builder.WriteString(val)
	}
	customBuilder.builder.WriteByte('\n')
}

func (customBuilder *customBuilder) String() string {
	return customBuilder.builder.String()
}
