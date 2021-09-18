package doc

import (
	"embed"
	"html/template"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/wenchy/grpcio/internal/atom"
)

//go:embed template static
var f embed.FS

func init() {
	// data, _ := f.ReadFile("static/swagger/client_test_service.swagger.json")
	// fmt.Println(string(data))
}

type controller struct {
}

// InitRouter init this controller's router.
func InitRouter() {
	html := template.Must(template.ParseFS(f, "template/*"))
	atom.GinEngine.SetHTMLTemplate(html)

	c := &controller{}
	atom.GinEngine.GET("/", c.checkHealth)
	atom.GinEngine.GET("/docs/test", c.getTest)
	atom.GinEngine.GET("/docs", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusMovedPermanently, "/docs/index")
	})
	atom.GinEngine.GET("/docs/index", c.getIndex)
	atom.GinEngine.GET("/docs/index/:service", c.getService)
	atom.GinEngine.GET("/docs/files/static/swagger/:name", c.getSwagger)
}

// Trainer type will be used in the program
type Trainer struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	City string `json:"city"`
}

func (c *controller) getTest(ctx *gin.Context) {
	atom.Log.Info("request of getTest")
	var trainer Trainer
	trainer.Name = "webot"
	trainer.Age = 1
	trainer.City = "Shenzhen"
	ctx.JSON(http.StatusOK, trainer)
}

type Doc struct {
	Name string
	Path string
}

func (c *controller) checkHealth(ctx *gin.Context) {
	ctx.String(http.StatusOK, "OK")
}

func (c *controller) getIndex(ctx *gin.Context) {
	docs := make([]Doc, 0)
	entries, err := f.ReadDir("static/swagger")
	if err != nil {
		atom.Log.Errorf("read dir failed: %v", err)
		ctx.String(http.StatusInternalServerError, "read dir failed: %v", err)
		return
	}
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), "json") {
			service := strings.Split(entry.Name(), ".")[0]
			doc := Doc{Name: service, Path: "./index/" + service}
			docs = append(docs, doc)
		}
	}
	atom.Log.Debugf("docs: %+v", docs)
	ctx.HTML(http.StatusOK, "index.tpl.html", gin.H{
		"Docs": docs,
	})
}

func (c *controller) getService(ctx *gin.Context) {
	atom.Log.Info("request of getService")
	service := ctx.Param("service")
	ctx.HTML(http.StatusOK, "redoc.tpl.html", gin.H{
		"Name": service,
		"Path": "../files/static/swagger/" + service + ".swagger.json",
	})
}

func (c *controller) getSwagger(ctx *gin.Context) {
	atom.Log.Info("request of getSwagger")
	name := ctx.Param("name")
	data, err := f.ReadFile("static/swagger/" + name)
	if err != nil {
		atom.Log.Errorf("read file failed: %v", err)
		ctx.String(http.StatusInternalServerError, "read file failed: %v", err)
		return
	}
	ctx.Data(http.StatusOK, "application/json; charset=utf-8", data)
}
