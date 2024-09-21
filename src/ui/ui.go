package ui

import (
	"encoding/json"
	"l0/types"
	"net/http"

	"github.com/gin-gonic/gin"
)

func toJson(str string) (jsonMap map[string]any) {
    if err := json.Unmarshal([]byte(str), &jsonMap); err != nil{
        return map[string]any{"status": "fail"}
    }

    return jsonMap

}

func StartUI(ctx *types.LocalContext) {
    r := gin.New()

    r.Use(gin.Recovery())

    r.GET("/", func(c *gin.Context) {
        c.HTML(http.StatusOK, "form.html", nil)
    })

    r.POST("/submit", func(c *gin.Context) {

        input := c.PostForm("input_id")

        jsonMap := toJson(ctx.Cache.Get(input))

        c.JSON(http.StatusOK, jsonMap)
    })

    r.LoadHTMLGlob("ui/templates/*")
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
