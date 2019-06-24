package main

import (
	"bytes"
	"net/http"

	"github.com/BiteBit/protoc-gen-gin/generator"
	"github.com/BiteBit/protoc-gen-gin/tool"
	"github.com/golang/protobuf/proto"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"google.golang.org/genproto/googleapis/api/annotations"
)

type (
	// Generator godoc
	Generator struct {
		generator.Generator

		Param             map[string]string // Command-line parameters.
		PackageImportPath string            // Go import path of the package we're generating code for
		ImportPrefix      string            // String to prefix to imported package file names.
		ImportMap         map[string]string // Mapping from .proto file name to import path
		Pkg               map[string]string // The names under which we import support packages
	}
)

const (
	// GeneratorName godoc
	GeneratorName = "protoc-gen-box"
)

// New godoc
func New() *Generator {
	g := new(Generator)
	g.Buffer = new(bytes.Buffer)
	g.Request = new(plugin.CodeGeneratorRequest)
	g.Response = new(plugin.CodeGeneratorResponse)

	return g
}

// GenerateAllFiles godoc
func (g *Generator) GenerateAllFiles() {
	for _, file := range g.Request.GetProtoFile() {
		if !g.CanGenerate(g.FileName) {
			continue
		}

		g.Reset()
		g.SetUp(file)
		g.ImportDefault()
		g.GenCommentHead(GeneratorName)
		g.GenPackageName()
		g.GenImports()
		g.GenSvcInterface()
		g.GenSvcImplement()
		g.GenSvcRegister()

		g.Response.File = append(g.Response.File, &plugin.CodeGeneratorResponse_File{
			Name:    proto.String(tool.GeneratedFileName(g.FileName)),
			Content: proto.String(g.String()),
		})
	}
}

// ImportDefault 导入默认的包信息
func (g *Generator) ImportDefault() {
	g.AddImport("context")
	g.AddImport("github.com/gin-gonic/gin")
	g.AddImport("github.com/gin-gonic/gin/binding")
}

// GenSvcInterface 生成服务接口
func (g *Generator) GenSvcInterface() {
	for _, svc := range g.File.GetService() {
		g.F("type %s interface {", tool.ToServerInterfaceType(svc))

		for _, method := range svc.GetMethod() {
			inType := tool.ToType(g.FilePkg, method.GetInputType())
			outType := tool.ToType(g.FilePkg, method.GetOutputType())

			g.F("	%s(ctx context.Context, req *%s) (resp *%s, err error)", method.GetName(), inType, outType)
		}

		g.P("}")
	}
}

// GenSvcImplement 生成服务实现
func (g *Generator) GenSvcImplement() {
	g.P("")

	for _, svc := range g.File.GetService() {
		for _, method := range svc.GetMethod() {
			inType := tool.ToType(g.File.GetPackage(), method.GetInputType())

			g.F(`func %s%s(svc %s) func(ctx *gin.Context) {`, svc.GetName(), method.GetName(), tool.ToServerInterfaceType(svc))
			g.P("	return func(ctx *gin.Context) {")
			g.F("		req := &%s{}", inType)
			g.P(`		if err := ctx.ShouldBindWith(req, binding.Default(ctx.Request.Method, ctx.Request.Header.Get("Content-Type"))); err != nil {`)
			g.P("			ctx.JSON(400, err)")
			g.P("			ctx.Abort()")
			g.P("			return")
			g.P("		}")
			g.P("")
			g.F("		if resp, err := svc.%s(ctx, req); err != nil {", method.GetName())
			g.P("			ctx.JSON(500, err)")
			g.P("			ctx.Abort()")
			g.F("		} else {")
			g.P("			ctx.JSON(200, resp)")
			g.P("		}")
			g.P("	}")
			g.P("}")
			g.P("")
		}
	}
}

// GenSvcRegister 生成服务注册器
func (g *Generator) GenSvcRegister() {
	g.P("")
	for _, svc := range g.File.GetService() {
		name := tool.ToServerInterfaceType(svc)
		g.F(`func Register%s(engine *gin.Engine, server %s) {`, name, name)

		for _, method := range svc.GetMethod() {
			ext, _ := proto.GetExtension(method.GetOptions(), annotations.E_Http)
			rule := ext.(*annotations.HttpRule)

			var httpMethod string
			var pathPattern string
			switch pattern := rule.Pattern.(type) {
			case *annotations.HttpRule_Get:
				pathPattern = pattern.Get
				httpMethod = http.MethodGet
			case *annotations.HttpRule_Put:
				pathPattern = pattern.Put
				httpMethod = http.MethodPut
			case *annotations.HttpRule_Post:
				pathPattern = pattern.Post
				httpMethod = http.MethodPost
			case *annotations.HttpRule_Patch:
				pathPattern = pattern.Patch
				httpMethod = http.MethodPatch
			case *annotations.HttpRule_Delete:
				pathPattern = pattern.Delete
				httpMethod = http.MethodDelete
			default:
			}

			g.F("	engine.%s(\"%s\", %s%s(server))", httpMethod, tool.OpenAPI2Go(pathPattern), svc.GetName(), method.GetName())
		}

		g.P("}")
	}
}

// Reset godoc
func (g *Generator) Reset() {
	g.Generator.Reset()
	g.File = nil
	g.FileName = ""
	g.FilePkg = ""
}
