package ui

import (
	"encoding/json"
	local_context "l0/types/local_context"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
)

func toJson(str string) (jsonMap map[string]any) {
	if err := json.Unmarshal([]byte(str), &jsonMap); err != nil {
		return map[string]any{"status": "fail"}
	}

	return jsonMap

}

func StartUI(ctx *local_context.LocalContext) {
	r := gin.New()

	r.Use(gin.Recovery())

	r.GET("/", func(c *gin.Context) {
		var ls []string
		for _, it := range ctx.Cache.GetAll() {
			ls = append(ls, it.Id)
		}
        sort.Strings(ls)

		c.HTML(http.StatusOK, "form.html", gin.H{
			"items": ls,
		})
	})

	r.POST("/submit", func(c *gin.Context) {

		input := c.PostForm("input_id")

		jsonStr, err := ctx.Cache.Get(input)

		if err != nil {
			c.String(http.StatusOK, "Order with id '%s' does not exist", input)
			return
		}

		c.Data(http.StatusOK, "application/json", []byte(*jsonStr))
	})

	r.LoadHTMLGlob("ui/templates/*")
	r.Run()
}
