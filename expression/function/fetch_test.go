package function

import (
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/niming-dev/ddd-demo/go-common/expression"
)

func Test_FetchGet(t *testing.T) {
	res, err := fetch{}.Call(nil, []*expression.Data{
		expression.NewString("GET"),                        // method
		expression.NewString("http://localhost:9876/test"), // url
		expression.NewString(""),                           // encoding: json or form, default is form
		expression.NewStruct(map[string]*expression.Data{
			"a": expression.NewInt(10),
			"b": expression.NewBool(true),
			"c": expression.NewString("cstring"),
		}),
	})
	if nil != err {
		t.Fatal(err)
	}

	t.Log(res.String())
}

func Test_FetchPostJson(t *testing.T) {
	res, err := fetch{}.Call(nil, []*expression.Data{
		expression.NewString("POST"),                       // method
		expression.NewString("http://localhost:9876/test"), // url
		expression.NewString("json"),                       // encoding: json or form, default is form
		expression.NewStruct(map[string]*expression.Data{
			"a": expression.NewInt(10),
			"b": expression.NewBool(true),
			"c": expression.NewString("cstring"),
		}),
	})
	if nil != err {
		t.Fatal(err)
	}

	t.Log(res.String())
}

func Test_FetchPostForm(t *testing.T) {
	res, err := fetch{}.Call(nil, []*expression.Data{
		expression.NewString("POST"),                       // method
		expression.NewString("http://localhost:9876/test"), // url
		expression.NewString("form"),                       // encoding: json or form, default is form
		expression.NewStruct(map[string]*expression.Data{
			"a": expression.NewInt(10),
			"b": expression.NewBool(true),
			"c": expression.NewString("cstring"),
		}),
	})
	if nil != err {
		t.Fatal(err)
	}

	t.Log(res.String())
}

func init() {
	r := gin.Default()

	r.GET("/test", func(c *gin.Context) {

		pa := c.Query("a")
		pb := c.Query("b")
		pc := c.Query("c")

		c.JSON(http.StatusOK, map[string]string{"a": pa, "b": pb, "c": pc})
	})

	r.POST("/test", func(c *gin.Context) {

		reqBytes, err := io.ReadAll(c.Request.Body)
		if nil != err {
			c.String(http.StatusBadRequest, fmt.Sprintf("%v", err))
		}

		c.Writer.Write(reqBytes)
		c.Writer.WriteHeader(http.StatusOK)
	})

	go r.Run(":9876")
}
