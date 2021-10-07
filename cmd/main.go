package main

import (
    "context"
    "crypto/tls"
    "github.com/allnash/moxie/models"
    guuid "github.com/google/uuid"
    "github.com/joho/godotenv"
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
    "golang.org/x/crypto/acme"
    "golang.org/x/crypto/acme/autocert"
    "gopkg.in/natefinch/lumberjack.v2"
    "net/http"
    "os"
    "os/signal"
    "time"
)

const AppEnvFilename = "/etc/moxie/app.env"

func load() error {
    return godotenv.Load(AppEnvFilename)
}

func main() {

    // Create Echo server
    e := echo.New()
    e.Logger.SetOutput(&lumberjack.Logger{
        Filename:   "/var/log/moxie/moxie.log",
        MaxSize:    100, // megabytes
        MaxBackups: 3,
        MaxAge:     28,   //days
        Compress:   true, // disabled by default
    })

    // Load ENV
    err := load()
    if err != nil {
        e.Logger.Error(err)
    }

    // Hosts
    hosts := map[string]*models.Host{}

    //-----
    // API
    //-----

    api := echo.New()
    api.Pre(middleware.HTTPSRedirect())
    api.Use(middleware.Logger())
    api.Use(middleware.Recover())
    api.Use(middleware.RequestIDWithConfig(middleware.RequestIDConfig{
        Generator: func() string {
            return customGenerator()
        },
    }))
    api.Use(middleware.GzipWithConfig(middleware.GzipConfig{
        Level: 5,
    }))
    api.Use(middleware.BodyLimit("10M"))
    // Add to Hosts
    hosts[os.Getenv("API_DOMAIN")] = &models.Host{Echo: api}

    api.GET("/", func(c echo.Context) error {
        return c.String(http.StatusOK, "API")
    })

    //------
    // Asset
    //------

    assets := echo.New()
    assets.Pre(middleware.HTTPSRedirect())
    assets.Use(middleware.Logger())
    assets.Use(middleware.Recover())
    assets.Use(middleware.GzipWithConfig(middleware.GzipConfig{
        Level: 5,
    }))
    assets.Use(expiresServerHeader)
    assets.Use(middleware.BodyLimit("10M"))
    assets.Use(middleware.StaticWithConfig(middleware.StaticConfig{
        Root:   os.Getenv("ASSET_DIRECTORY"),
        Browse: true,
    }))
    // Add to Hosts
    hosts[os.Getenv("ASSET_DOMAIN")] = &models.Host{Echo: assets}

    // Server
    e.Use(middleware.Recover())
    e.Use(middleware.Logger())
    e.Any("/*", func(c echo.Context) (err error) {
        req := c.Request()
        res := c.Response()
        host := hosts[req.Host]

        if host == nil {
            err = echo.ErrNotFound
        } else {
            host.Echo.ServeHTTP(res, req)
        }

        return
    })

    autoTLSManager := autocert.Manager{
        Prompt: autocert.AcceptTOS,
        // Cache certificates to avoid issues with rate limits (https://letsencrypt.org/docs/rate-limits)
        Cache:      autocert.DirCache("/var/www/.cache"),
        HostPolicy: autocert.HostWhitelist(os.Getenv("API_DOMAIN"), os.Getenv("ASSET_DOMAIN")),
    }

    // Start server
    go func() {
        s := http.Server{
            Addr:    ":443",
            Handler: e, // set Echo as handler
            TLSConfig: &tls.Config{
                //Certificates: nil, // <-- s.ListenAndServeTLS will populate this field
                GetCertificate: autoTLSManager.GetCertificate,
                NextProtos:     []string{acme.ALPNProto},
            },
            ReadTimeout: 30 * time.Second, // use custom timeouts
        }
        if err := s.ListenAndServeTLS("", ""); err != http.ErrServerClosed {
            e.Logger.Fatal(err)
            e.Logger.Fatal("shutting down the server")
        }
    }()

    // Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
    // Use a buffered channel to avoid missing signals as recommended for signal.Notify
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, os.Interrupt)
    <-quit
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    if err := e.Shutdown(ctx); err != nil {
        e.Logger.Fatal(err)
    }
}

func customGenerator() string {
    id := guuid.New()
    return id.String()
}

// ServerHeader middleware adds a `Server` header to the response.
func expiresServerHeader(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        c.Response().Header().Set("Cache-Control", "public, max-age=3600")
        return next(c)
    }
}
