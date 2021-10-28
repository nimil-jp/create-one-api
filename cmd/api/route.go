package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	jwt "github.com/ken109/gin-jwt"
	"github.com/pkg/errors"

	"go-gin-ddd/pkg/context"
	"go-gin-ddd/pkg/xerrors"
)

type router struct {
	g *gin.RouterGroup
}

func newRouter(engine *gin.Engine) *router {
	return &router{g: engine.Group("")}
}

func (r *router) group(relativePath string, handlers []gin.HandlerFunc, fn func(r *router)) {
	if handlers == nil {
		handlers = []gin.HandlerFunc{}
	}
	fn(&router{g: r.g.Group(relativePath, handlers...)})
}

type handlerFunc func(ctx context.Context, c *gin.Context) error

func (r *router) get(relativePath string, handlerFunc handlerFunc) {
	r.g.GET(relativePath, wrapperFunc(handlerFunc))
}

func (r *router) post(relativePath string, handlerFunc handlerFunc) {
	r.g.POST(relativePath, wrapperFunc(handlerFunc))
}

func (r *router) put(relativePath string, handlerFunc handlerFunc) {
	r.g.PUT(relativePath, wrapperFunc(handlerFunc))
}

func (r *router) patch(relativePath string, handlerFunc handlerFunc) {
	r.g.PATCH(relativePath, wrapperFunc(handlerFunc))
}

func (r *router) delete(relativePath string, handlerFunc handlerFunc) {
	r.g.DELETE(relativePath, wrapperFunc(handlerFunc))
}

func (r *router) options(relativePath string, handlerFunc handlerFunc) {
	r.g.OPTIONS(relativePath, wrapperFunc(handlerFunc))
}

func (r *router) head(relativePath string, handlerFunc handlerFunc) {
	r.g.HEAD(relativePath, wrapperFunc(handlerFunc))
}

func wrapperFunc(handlerFunc handlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx context.Context

		if userID, ok := jwt.GetClaim(c, "user_id"); ok {
			ctx = context.New(c.GetHeader("X-Request-Id"), uint(userID.(float64)))
		} else {
			ctx = context.New(c.GetHeader("X-Request-Id"), 0)
		}

		c.Writer.Header().Add("X-Request-Id", ctx.RequestID())

		err := handlerFunc(ctx, c)

		if err != nil {
			switch v := err.(type) {
			case *xerrors.Expected:
				if v.StatusOk() {
					return
				} else {
					c.JSON(v.StatusCode(), v.Message())
				}
			case *xerrors.Validation:
				c.JSON(http.StatusBadRequest, v)
			default:
				if gin.Mode() == gin.DebugMode {
					c.JSONP(http.StatusInternalServerError, map[string]string{"request_id": ctx.RequestID(), "error": v.Error()})
				} else {
					c.JSONP(http.StatusInternalServerError, map[string]string{"request_id": ctx.RequestID()})
				}
			}

			_ = c.Error(errors.Errorf("%+v", err))
		}
	}
}
