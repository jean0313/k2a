package gen

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"

	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

var (
	Separator = fmt.Sprintf("%v", os.PathSeparator)
)

type GenContex struct {
	ChannelName         string
	MessageName         string
	OperationName       string
	SchemaName          string
	MessageSchema       string
	MessageSchemaFormat string

	GCtx *GlobalContext
}

type GlobalContext struct {
	Title           string
	Version         string
	AppDescription  string
	ServersUrl      string
	L3Domain        string
	ApplicationCode string
	Module          string

	Group               string
	Artifact            string
	Description         string
	PackageName         string
	ReleaseVersion      string
	CustomCommitMessage string

	AsyncAPIFile string
	DestDir      string
}

func NewGlobalContext() *GlobalContext {
	return &GlobalContext{
		DestDir: "output",

		Version:     "1.0.0",
		Group:       "com.example",
		Artifact:    "sample-demo",
		PackageName: "demo",

		ReleaseVersion: "1.0.0-RC-1",
	}
}

type Project struct {
	dirs  map[string]string
	ctxes []GenContex
	gCtx  *GlobalContext
}

func NewProject(ctxex []GenContex, gCtx *GlobalContext) *Project {
	return &Project{
		dirs:  make(map[string]string),
		ctxes: ctxex,
		gCtx:  gCtx,
	}
}

func (p *Project) init() {
	pkgPath := filepath.Join(strings.Split(p.gCtx.Group, ".")...)
	root := filepath.Join(p.gCtx.DestDir, p.gCtx.Artifact)

	// project root
	p.dirs["root"] = root

	// test base dir
	p.dirs["test"] = filepath.Join(root, "src", "test")
	p.dirs["java"] = filepath.Join(root, "src", "main", "java")

	// resources base dir
	p.dirs["resources"] = filepath.Join(root, "src", "main", "resources")
	p.dirs["api"] = filepath.Join(root, "src", "main", "resources", "api")
	p.dirs["avro"] = filepath.Join(root, "src", "main", "resources", "avro")

	// java package base dir
	base := filepath.Join(root, "src", "main", "java", pkgPath)
	p.dirs["base"] = base

	p.dirs["processing"] = filepath.Join(base, "consumer", "processing")
	p.dirs["config"] = filepath.Join(base, "config")
	p.dirs["dao"] = filepath.Join(base, "dao")
	p.dirs["mapper"] = filepath.Join(base, "data", "mapper")
	p.dirs["domain"] = filepath.Join(base, "domain")
	p.dirs["integration"] = filepath.Join(base, "integration")
	p.dirs["model"] = filepath.Join(base, "model")
	p.dirs["processing"] = filepath.Join(base, "processing")
	p.dirs["response"] = filepath.Join(base, "response")
	p.dirs["service"] = filepath.Join(base, "service")
	p.dirs["validator"] = filepath.Join(base, "validator")

	for _, path := range p.dirs {
		os.MkdirAll(path, os.ModePerm)
	}

	p.createIgnoreFile()
}

func (p *Project) createFile(ctx *GenContex, tmplFile string, dirName string, nameFn func(ctx *GenContex) string) error {
	path := p.dirs[dirName]
	javaPath := filepath.Join(path, nameFn(ctx))
	return renderTemplate(tmplFile, javaPath, ctx)
}

func (p *Project) createPom() error {
	return p.renderGlobalTemplateFile("pom.tmpl", "root", "pom.xml")
}

func (p *Project) createIgnoreFile() {
	createFile(filepath.Join(p.dirs["root"], ".gitignore"),
		`.idea/
target/
.nvm/
*.swp

*.log

*.class

*.jar
*.war
*.nar
*.ear
*.zip
*.tar.gz
*.rar`)
}

func (p *Project) createApplication() error {
	return p.renderGlobalTemplateFile("Application.tmpl", "base", "Application.java")
}

func (p *Project) createCommonFiles() error {
	if err := p.renderGlobalTemplateFile("CommonConsumerValidator.tmpl", "validator", "CommonConsumerValidator.java"); err != nil {
		return err
	}

	if err := p.renderGlobalTemplateFile("CommonProducerValidator.tmpl", "validator", "CommonProducerValidator.java"); err != nil {
		return err
	}

	return nil
}

func (p *Project) renderGlobalTemplateFile(tmplFile string, dirName string, targetFile string) error {
	path := p.dirs[dirName]
	targetPath := filepath.Join(path, targetFile)
	return renderTemplateFn(tmplFile, targetPath, p.gCtx)
}

func (p *Project) createFiles() {
	var err error
	for i := 0; i < len(p.ctxes); i++ {
		v := p.ctxes[i]
		zap.L().Info("generate file for context", zap.Any("contex", v))

		if v.OperationName == "publish" {
			if err = p.createFile(&v, "ProducerConfig.tmpl", "config", func(ctx *GenContex) string {
				return fmt.Sprintf("%vProducerConfig.java", capitalize(ctx.ChannelName))
			}); err != nil {
				panic(err)
			}

			if err = p.createFile(&v, "AbstractProducerResponse.tmpl", "response", func(ctx *GenContex) string {
				return fmt.Sprintf("Abstract%vProducerResponse.java", capitalize(ctx.ChannelName))
			}); err != nil {
				panic(err)
			}

			if err = p.createFile(&v, "PublisherService.tmpl", "service", func(ctx *GenContex) string {
				return fmt.Sprintf("%vPublishService.java", capitalize(ctx.ChannelName))
			}); err != nil {
				panic(err)
			}

			if err = p.createFile(&v, "PublisherServiceImpl.tmpl", "service", func(ctx *GenContex) string {
				return fmt.Sprintf("%vPublishServiceImpl.java", capitalize(ctx.ChannelName))
			}); err != nil {
				panic(err)
			}

			if err = p.createFile(&v, "ProducerValidator.tmpl", "validator", func(ctx *GenContex) string {
				return fmt.Sprintf("%vProducerValidator.java", capitalize(ctx.ChannelName))
			}); err != nil {
				panic(err)
			}
		}

		if v.OperationName == "subscribe" {
			if err = p.createFile(&v, "ConsumerConfig.tmpl", "config", func(ctx *GenContex) string {
				return fmt.Sprintf("%vConsumerConfig.java", capitalize(ctx.ChannelName))
			}); err != nil {
				panic(err)
			}

			if err = p.createFile(&v, "ConsumerConfig.tmpl", "config", func(ctx *GenContex) string {
				return fmt.Sprintf("%vConsumerProcessor.java", capitalize(ctx.ChannelName))
			}); err != nil {
				panic(err)
			}

			if err = p.createFile(&v, "ConsumerProcessor.tmpl", "processing", func(ctx *GenContex) string {
				return fmt.Sprintf("%vConsumerProcessor.java", capitalize(ctx.ChannelName))
			}); err != nil {
				panic(err)
			}

			if err = p.createFile(&v, "ConsumerService.tmpl", "service", func(ctx *GenContex) string {
				return fmt.Sprintf("%vConsumerService.java", capitalize(ctx.ChannelName))
			}); err != nil {
				panic(err)
			}

			if err = p.createFile(&v, "ConsumerValidator.tmpl", "validator", func(ctx *GenContex) string {
				return fmt.Sprintf("%vConsumerValidator.java", capitalize(ctx.ChannelName))
			}); err != nil {
				panic(err)
			}
		}
	}
}

func (p *Project) generateSchemas(ctx GenContex) {
	// TODO avro schema

	// TODO json schema

}

func generateProjectStructure(ctxes []GenContex, gCtx *GlobalContext) {
	project := NewProject(ctxes, gCtx)
	project.init()
	project.createPom()
	project.createApplication()
	project.createFiles()
	project.createCommonFiles()

	for _, v := range ctxes {
		project.generateSchemas(v)
	}
}

func renderTemplate(tmplFile string, targetFile string, ctx *GenContex) error {
	f, err := os.Create(targetFile)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	tmpl := getTemplate(tmplFile)
	err = tmpl.Execute(f, ctx)
	return err
}

func createFile(filePath string, content string) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	_, err = f.WriteString(content)
	if err != nil {
		return err
	}
	return nil
}

func renderTemplateFn(tmplFile string, targetFile string, ctx any) error {
	f, err := os.Create(targetFile)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	tmpl := getTemplate(tmplFile)
	err = tmpl.Execute(f, ctx)
	return err
}

func getTemplate(tmplFile string) *template.Template {
	path := filepath.Join("templates", tmplFile)
	tmpl, err := template.New(tmplFile).Funcs(template.FuncMap{
		"capitalize":   capitalize,
		"uncapitalize": uncapitalize,
	}).ParseFiles(path)

	if err != nil {
		panic(err)
	}
	return tmpl
}

func capitalize(str string) string {
	runes := []rune(str)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

func uncapitalize(str string) string {
	runes := []rune(str)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

func parseContexts(apiFile string) []GenContex {
	// version 2.0 - 2.6, start from channel
	// channel -> sub/pub -> message -> payload

	bs, err := os.ReadFile(apiFile)
	if err != nil {
		panic(err)
	}

	var v interface{}
	err = yaml.Unmarshal(bs, &v)
	if err != nil {
		panic(err)
	}
	v = convertMapI2MapS(v)

	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	api := AsyncAPIModel{}
	err = json.Unmarshal(data, &api)
	if err != nil {
		panic(err)
	}

	msgNameToSchemaNameMap := make(map[string]string)
	if messages := api.Components.Messages; messages != nil {
		for k, v := range messages {
			_, schemaName := extractMessageNameAndSchemaName(&v, msgNameToSchemaNameMap)
			msgNameToSchemaNameMap[k] = schemaName
		}
	}

	var ctxs []GenContex
	for k, v := range api.Channels {
		ctx := GenContex{
			ChannelName: k,
		}

		if v.Publish != nil {
			ctx.OperationName = "publish"
			ctx.MessageName, ctx.SchemaName = extractMessageNameAndSchemaName(v.Publish.Message, msgNameToSchemaNameMap)
		}

		if v.Subscribe != nil {
			ctx.OperationName = "subscribe"
			ctx.MessageName, ctx.SchemaName = extractMessageNameAndSchemaName(v.Subscribe.Message, msgNameToSchemaNameMap)
		}

		ctxs = append(ctxs, ctx)
	}

	for _, v := range ctxs {
		zap.L().Debug("context", zap.Any("ctx", v))
	}

	for k, v := range msgNameToSchemaNameMap {
		zap.L().Debug("message to schema mapping", zap.String(k, v))
	}

	return ctxs
}

func extractMessageNameAndSchemaName(message *Message, mapping map[string]string) (string, string) {
	if message.Reference != nil {
		ref := message.Reference.Ref
		if strings.HasPrefix(ref, "#/") {
			msgName := strings.Split(ref, "/")[3]
			if schemaName, ok := mapping[msgName]; ok {
				return msgName, schemaName
			}
			return msgName, ""
		}
	} else if message.MessageEntity != nil {
		return message.MessageEntity.Name, extractSchemaName(message.MessageEntity)
	}
	return "", ""
}

// only search for schema name in message entity
func extractSchemaName(messageEntity *MessageEntity) string {
	if messageEntity != nil && messageEntity.Payload != nil {
		payload := messageEntity.Payload
		if iref, ok := payload["$ref"]; ok {
			ref := fmt.Sprintf("%v", iref)
			if strings.HasPrefix(ref, "#/") {
				return strings.Split(ref, "/")[3]
			}
		} else {
			if name, ok := payload["name"]; ok {
				return fmt.Sprintf("%v", name)
			} else {
				return messageEntity.Name
			}
		}
	}
	return ""
}

func debug(api AsyncAPIModel) {
	zap.L().Debug("----------- Channels ------------")
	for k, v := range api.Channels {
		zap.L().Debug("channel", zap.Any(k, v))
		if v.Publish != nil {
			zap.L().Debug("		channel message ref", zap.Any("ref", v.Publish.Message.Reference))
			zap.L().Debug("		channel message payload", zap.Any("payload", v.Publish.Message.MessageEntity))
		}

		if v.Subscribe != nil {
			zap.L().Debug("		channel message ref", zap.Any("ref", v.Subscribe.Message.Reference))
			zap.L().Debug("		channel message payload", zap.Any("payload", v.Subscribe.Message.MessageEntity))
		}
	}

	zap.L().Debug("----------- Components.Messages ------------")
	for k, v := range api.Components.Messages {
		zap.L().Debug("component message", zap.Any(k, v))
		if v.MessageEntity != nil {
			zap.L().Debug("		compnent message payload", zap.Any("payload", v.MessageEntity.Payload))
		}
	}

	zap.L().Debug("----------- Components.Schemas ------------")
	for k, v := range api.Components.Schemas {
		zap.L().Debug("schema", zap.Any(k, v))
	}
}

func convertMapI2MapS(v interface{}) interface{} {
	switch x := v.(type) {
	case map[interface{}]interface{}:
		m := map[string]interface{}{}

		for k, v2 := range x {
			switch k2 := k.(type) {
			case string: // Fast check if it's already a string
				m[k2] = convertMapI2MapS(v2)
			default:
				m[fmt.Sprint(k)] = convertMapI2MapS(v2)
			}
		}

		v = m

	case []interface{}:
		for i, v2 := range x {
			x[i] = convertMapI2MapS(v2)
		}

	case map[string]interface{}:
		for k, v2 := range x {
			x[k] = convertMapI2MapS(v2)
		}
	}

	return v
}

func Gen(gCtx *GlobalContext) {
	if gCtx.AsyncAPIFile == "" {
		zap.L().Warn("async api spec file should not be empty!!")
		return
	}
	ctxes := parseContexts(gCtx.AsyncAPIFile)
	for i := 0; i < len(ctxes); i++ {
		ctxes[i].GCtx = gCtx
	}

	generateProjectStructure(ctxes, gCtx)
}
