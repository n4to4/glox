package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: generate <outout directory>")
	}

	outputDir := os.Args[1]
	defineAst(outputDir, "Expr", []string{
		"Assign   : name Token, value Expr",
		"Binary   : left Expr, operator Token, right Expr",
		"Grouping : expression Expr",
		"Literal  : value interface{}",
		"Logical  : left Expr, operator Token, right Expr",
		"Unary    : operator Token, right Expr",
		"Variable : name Token",
	})

	defineAst(outputDir, "Stmt", []string{
		"Block      : statements []Stmt",
		"Expression : expression Expr",
		"If         : condition Expr, thenBranch *Stmt, elseBranch *Stmt",
		"Print      : expression Expr",
		"Var        : name Token, initializer *Expr",
		"While      : condition Expr, body Stmt",
	})
}

func defineAst(outputDir, baseName string, types []string) {
	path := fmt.Sprintf("%s/%s.go", outputDir, strings.ToLower(baseName))
	file, err := os.Create(path)
	if err != nil {
		log.Fatalf("cannot open file %q, %v\n", path, err)
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	generateAst(w, baseName, types)
	w.Flush()
}

func generateAst(w io.StringWriter, baseName string, types []string) {
	w.WriteString("package main\n")
	w.WriteString("\n")
	w.WriteString(fmt.Sprintf("type %s interface {\n", baseName))
	w.WriteString("\tTokenLiteral() string\n")
	w.WriteString(fmt.Sprintf("\t%sAcceptor\n", baseName))
	w.WriteString("}\n")

	w.WriteString("\n")
	generateVisitor(w, baseName, types)

	w.WriteString("\n")
	generateAcceptor(w, baseName)

	for _, t := range types {
		w.WriteString("\n")
		generateType(w, t)
		w.WriteString("\n")
		generateAccept(w, baseName, t)
	}
}

func generateVisitor(w io.StringWriter, baseName string, types []string) {
	w.WriteString(fmt.Sprintf("type %sVisitor interface {\n", baseName))
	for _, t := range types {
		splits := strings.Split(t, ":")
		typeName := strings.Trim(splits[0], " ")
		w.WriteString(fmt.Sprintf("\tVisit%s%s(%s %s) (interface{}, error)\n",
			typeName, baseName,
			strings.ToLower(baseName), typeName,
		))
	}
	w.WriteString("}\n")
}

func generateAcceptor(w io.StringWriter, baseName string) {
	w.WriteString(fmt.Sprintf("type %sAcceptor interface {\n", baseName))
	w.WriteString(fmt.Sprintf("\tAccept(v %sVisitor) (interface{}, error)\n", baseName))
	w.WriteString("}\n")
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

	w.WriteString(fmt.Sprintf(`func (x %s) TokenLiteral() string { return "" }`+"\n", typeName))
}

func generateAccept(w io.StringWriter, baseName, types string) {
	splits := strings.Split(types, ":")
	typeName := strings.Trim(splits[0], " ")

	w.WriteString(fmt.Sprintf("func (x %s) Accept(v %sVisitor) (interface{}, error) {\n", typeName, baseName))
	w.WriteString(fmt.Sprintf("\treturn v.Visit%s%s(x)\n", typeName, baseName))
	w.WriteString("}\n")
}
