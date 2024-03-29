package generator

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
)

type (
	// Generator 基础的生成器
	Generator struct {
		*bytes.Buffer
		Request  *plugin.CodeGeneratorRequest  // The input.
		Response *plugin.CodeGeneratorResponse // The output.

		Imports  []string                        // 程序导入
		File     *descriptor.FileDescriptorProto // 正在处理的协议文件描述信息
		FileName string                          // 正在处理的文件名
		FilePkg  string                          // 正在处理的包名
	}
)

// SetUp 初始化
func (g *Generator) SetUp(file *descriptor.FileDescriptorProto) {
	g.File = file
	g.FileName = file.GetName()
	g.FilePkg = file.GetPackage()
}

// GenCommentHead 生成注释头
func (g *Generator) GenCommentHead(generator string) {
	g.F("// Code generated by %s. DO NOT EDIT.", generator)
	g.F("// source: %s", g.FileName)
	g.P("")
}

// GenPackageName 生成包名
func (g *Generator) GenPackageName() {
	g.F("package %s", g.FilePkg)
	g.P("")
}

// GenImports 生成导入信息
func (g *Generator) GenImports() {
	if len(g.Imports) == 0 {
		return
	}

	g.P("import (")
	for _, ipt := range g.Imports {
		g.F("	\"%s\"", ipt)
	}
	g.P(")")
	g.P("")
}

// AddImport 增加导入包信息
func (g *Generator) AddImport(pkg ...string) {
	g.Imports = append(g.Imports, pkg...)
}

// CanGenerate 是否需要生成文件
func (g *Generator) CanGenerate(name string) bool {
	ok := false

	for _, f := range g.Request.GetFileToGenerate() {
		if f == name {
			ok = true
			break
		}
	}

	return ok
}

// Reset godoc
func (g *Generator) Reset() {
	g.Buffer.Reset()
}

// P godoc
func (g *Generator) P(strs ...string) {
	for _, str := range strs {
		g.F(str)
	}
}

// F godoc
func (g *Generator) F(format string, args ...interface{}) {
	content := fmt.Sprintf(format, args...)

	g.WriteString(content)
	g.WriteString("\n")
}

// Error reports a problem, including an error, and exits the program.
func (g *Generator) Error(err error, msgs ...string) {
	s := strings.Join(msgs, " ") + ":" + err.Error()
	log.Print("protoc error:", s)
	os.Exit(1)
}

// Fail reports a problem and exits the program.
func (g *Generator) Fail(msgs ...string) {
	s := strings.Join(msgs, " ")
	log.Print("protoc error:", s)
	os.Exit(1)
}
