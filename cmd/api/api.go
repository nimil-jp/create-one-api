package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ken109/gin-jwt"
	"github.com/nimil-jp/gin-utils/http/middleware"
	"github.com/nimil-jp/gin-utils/http/router"

	"go-gin-ddd/config"
	"go-gin-ddd/driver/rdb"
	"go-gin-ddd/infrastructure/email"
	"go-gin-ddd/infrastructure/gcp"
	"go-gin-ddd/infrastructure/log"
	"go-gin-ddd/infrastructure/paypal"
	"go-gin-ddd/infrastructure/persistence"
	"go-gin-ddd/infrastructure/stripe"
	"go-gin-ddd/interface/handler"
	"go-gin-ddd/usecase"
)

func Run() {
	logger := log.Logger()

	err := jwt.SetUp(
		jwt.Option{
			Realm:            config.DefaultRealm,
			SigningAlgorithm: jwt.HS256,
			SecretKey:        []byte(config.Env.App.Secret),
		},
	)
	if err != nil {
		panic(err)
	}
	logger.Info("Succeeded in setting up JWT.")

	engine := gin.New()

	engine.GET("health", func(c *gin.Context) { c.Status(http.StatusOK) })

	engine.Use(middleware.Log(log.ZapLogger(), time.RFC3339, false))
	engine.Use(middleware.RecoveryWithLog(log.ZapLogger(), true))

	// cors
	engine.Use(middleware.Cors(nil))

	// dependencies injection
	// ----- infrastructure -----
	emailInfra := email.New()
	gcs := gcp.NewGcs()
	firebase := gcp.NewFirebase()
	paypalInfra := paypal.New()
	stripeInfra := stripe.New()

	// persistence
	userPersistence := persistence.NewUser()
	transactionPersistence := persistence.NewTransaction()
	articlePersistence := persistence.NewArticle()

	// ----- use case -----
	userUseCase := usecase.NewUser(userPersistence, transactionPersistence, articlePersistence, firebase, emailInfra, paypalInfra, stripeInfra)
	transactionUseCase := usecase.NewTransaction(transactionPersistence, userPersistence, stripeInfra)
	articleUseCase := usecase.NewArticle(articlePersistence)

	// ----- handler -----
	webhookHandler := handler.NewWebhook(userPersistence)
	signedURLHandler := handler.NewSignedURL(gcs)

	userHandler := handler.NewUser(userUseCase)
	transactionHandler := handler.NewTransaction(transactionUseCase)
	articleHandler := handler.NewArticle(articleUseCase)

	r := router.New(config.Env.App.Name, engine, rdb.Get)

	middleware.FirebaseSetup(firebase.AuthClient())

	r.Group("", []gin.HandlerFunc{middleware.FirebaseAuth(false)}, func(r *router.Router) {
		r.Group("webhook", nil, func(r *router.Router) {
			r.Post("paypal", webhookHandler.Paypal)
		})

		r.Group("users", nil, func(r *router.Router) {

			r.Get("", userHandler.Search)

			r.Group("by", nil, func(r *router.Router) {
				r.Get("id/:user_id", userHandler.GetBy("id"))
				r.Get("username/:username", userHandler.GetBy("username"))
			})

			r.Group(":user_id", nil, func(r *router.Router) {
				r.Get("following", userHandler.Following)
				r.Get("followers", userHandler.Followers)

				r.Get("supporting", userHandler.Supporting)
				r.Get("supporters", userHandler.Supporters)
				r.Get("supported-transactions", userHandler.SupportedTransactions)
				r.Get("articles", userHandler.Articles)
			})
		})

		r.Group("articles", nil, func(r *router.Router) {
			r.Get("", articleHandler.Search)
			r.Get(":id", articleHandler.GetByID)
		})
	})

	r.Group("", []gin.HandlerFunc{middleware.FirebaseAuth(true)}, func(r *router.Router) {
		r.Get("timeline", userHandler.Timeline)

		r.Group("signed-url", nil, func(r *router.Router) {
			r.Get("profile", signedURLHandler.Profile)
			r.Get("article", signedURLHandler.Article)
		})

		r.Group("users", nil, func(r *router.Router) {
			r.Post("", userHandler.Create)

			r.Get("me", userHandler.GetMe)

			r.Group(":user_id", nil, func(r *router.Router) {
				r.Put("", userHandler.Edit)
				r.Patch("", userHandler.Patch)
				r.Get("connect-paypal", userHandler.ConnectPaypal)
				r.Patch("connect-stripe", userHandler.ConnectStripe)

				r.Group("following", nil, func(r *router.Router) {
					r.Post(":target_user_id", userHandler.Follow(true))
					r.Delete(":target_user_id", userHandler.Follow(false))
				})

				r.Group("following", nil, func(r *router.Router) {
					r.Group("articles", nil, func(r *router.Router) {
						r.Get("", userHandler.FollowingArticles)
					})
				})

				r.Group("supporters", nil, func(r *router.Router) {
					r.Group("articles", nil, func(r *router.Router) {
						r.Get("", userHandler.SupportersArticles)
					})
				})

				r.Delete("", userHandler.Delete)
			})
		})

		r.Group("transactions", nil, func(r *router.Router) {
			r.Post("paypal-payment-intent", transactionHandler.CreateStripePaymentIntent)
			r.Post("", transactionHandler.Create)
		})

		r.Group("articles", nil, func(r *router.Router) {
			r.Post("", articleHandler.Create)
			r.Put(":id", articleHandler.Update)
			r.Delete(":id", articleHandler.Delete)
		})
	})

	logger.Info("Succeeded in setting up routes.")

	// serve
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.Env.Port),
		Handler: engine,
	}

	go func() {
		if err = srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	logger.Info("Succeeded in listen and serve.")

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("Server forced to shutdown: %+v", err)
	}

	logger.Info("Server exiting")
}
