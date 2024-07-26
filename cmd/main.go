package main

import (
	"context"
	"errors"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	"github.com/tusmasoma/go-tech-dojo/infra/mysql"
	"github.com/tusmasoma/go-tech-dojo/infra/redis"
	"github.com/tusmasoma/go-tech-dojo/interfaces/handler"
	"github.com/tusmasoma/go-tech-dojo/interfaces/middleware"
	"github.com/tusmasoma/go-tech-dojo/pkg/log"
	"github.com/tusmasoma/go-tech-dojo/usecase"

	_ "github.com/go-sql-driver/mysql"
)

const (
	ReadTimeout                   = 5 * time.Second
	WriteTimeout                  = 10 * time.Second
	IdleTimeout                   = 15 * time.Second
	GracefulShutdownTimeout       = 5 * time.Second
	PreflightCacheDurationSeconds = 300
)

func main() {
	var addr string
	// .envファイルから環境変数を読み込む
	if err := godotenv.Load(); err != nil {
		log.Info("No .env file found", log.Ferror(err))
	}
	flag.StringVar(&addr, "addr", ":8083", "tcp host:port to connect")
	flag.Parse()

	Serve(addr)
}

func Serve(addr string) {
	mainCtx, cancelMain := context.WithCancel(context.Background())
	defer cancelMain()

	db, err := mysql.NewMySQLDB(mainCtx)
	if err != nil {
		log.Error("Failed to connect to DB", log.Ferror(err))
		return
	}
	defer db.Close()

	client := redis.NewRedisClient(mainCtx)

	transactionRepo := mysql.NewTransactionRepository(db)
	userRepo := mysql.NewUserRepository(db)
	userCollectionRepo := mysql.NewUserCollectionRepository(db)
	collectionRepo := mysql.NewCollectionRepository(db)
	collectionCacheRepo := redis.NewCollectionRepository(client)
	userUseCase := usecase.NewUserUseCase(userRepo, transactionRepo, userCollectionRepo, collectionRepo, collectionCacheRepo)
	userHandler := handler.NewUserHandler(userUseCase)
	authMiddleware := middleware.NewAuthMiddleware()

	/* ===== URLマッピングを行う ===== */
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Origin"},
		ExposedHeaders:   []string{"Link", "Authorization"},
		AllowCredentials: false,
		MaxAge:           PreflightCacheDurationSeconds,
	}))

	r.Route("/api", func(r chi.Router) {
		r.Route("/user", func(r chi.Router) {
			r.Post("/create", userHandler.CreateUser)
			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.Authenticate)
				r.Get("/get", userHandler.GetUser)
				r.Put("/update", userHandler.UpdateUser)
			})
		})
		r.Route("/collection", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				r.Use(authMiddleware.Authenticate)
				r.Get("/list", userHandler.ListUserCollections)
			})
		})
	})

	/* ===== サーバの設定 ===== */
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  ReadTimeout,
		WriteTimeout: WriteTimeout,
		IdleTimeout:  IdleTimeout,
	}
	/* ===== サーバの起動 ===== */
	log.Info("Server running...")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt, os.Kill)
	defer stop()

	go func() {
		if err = srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("Server failed", log.Ferror(err))
		}
	}()

	<-ctx.Done()
	log.Info("Server stopping...")

	tctx, cancel := context.WithTimeout(context.Background(), GracefulShutdownTimeout)
	defer cancel()

	if err = srv.Shutdown(tctx); err != nil {
		log.Error("Failed to shutdown http server", log.Ferror(err))
	}
	log.Info("Server exited")
}
