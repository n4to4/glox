package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {}

func defineAst(outputDir, baseName string, types []string) {
	path := fmt.Sprintf("%s/%s.go", outputDir, baseName)
	file, err := os.Create(path)
	if err != nil {
		log.Fatalf("cannot open file %q, %v\n", path, err)
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	generateAst(w, baseName, types)
}

func generateAst(w io.StringWriter, baseName string, types []string) {
	w.WriteString("package lox\n")
	w.WriteString("\n")
	w.WriteString(fmt.Sprintf("type %s interface {\n", baseName))
	w.WriteString("\tTokenLiteral() string\n")
	w.WriteString("}\n")

	for _, t := range types {
		w.WriteString("\n")
		generateType(w, t)
	}
}

func generateType(w io.StringWriter, types string) {
	splits := strings.Split(types, ":")
	typeName := strings.Trim(splits[0], " ")
	allFields := strings.Trim(splits[1], " ")
	fieldNames := strings.Split(allFields, ", ")

	w.WriteString(fmt.Sprintf("type %s struct {\n", typeName))
	for _, field := range fieldNames {
		w.WriteString(fmt.Sprintf("\t%s\n", field))
	}
	w.WriteString("}\n")
	w.WriteString("\n")

	w.WriteString(fmt.Sprintf(`func (x *%s) TokenLiteral() string { return "" }`+"\n", typeName))
}
