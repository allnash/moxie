package main

import (
    "context"
    "fmt"
    "github.com/allnash/moxie/config"
    "github.com/allnash/moxie/models"
    "github.com/ilyakaznacheev/cleanenv"
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
    "gopkg.in/natefinch/lumberjack.v2"
    "net/http"
    "net/url"
    "os"
    "os/signal"
    "time"
)

const AppYamlFilename = "app.yaml"
const SSL_PORT = "8443"


func load() config.Config {
	var cfg config.Config
	// read configuration from the file and environment variables
	if err := cleanenv.ReadConfig(AppYamlFilename, &cfg); err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	return cfg
}

var hosts = map[string]*models.Host{}

func main() {

	// Load ENV
	cfg := load()

	// Hosts
	for _, service := range cfg.Services {
        // Service Target
        tenant := echo.New()
        var targets []*middleware.ProxyTarget
        tenant.Logger.SetOutput(&lumberjack.Logger{
            Filename:   cfg.Logfile,
            MaxSize:    100, // megabytes
            MaxBackups: 3,
            MaxAge:     28,   //days
            Compress:   true, // disabled by default
        })
        // Service Config
		if service.Type == "proxy" {
            // Web endpoint
			urlS, err := url.Parse(service.EgressUrl)
			if err != nil {
				tenant.Logger.Fatal(err)
			}
			targets = append(targets, &middleware.ProxyTarget{
				URL: urlS,
			})
			tenant.Use(middleware.Proxy(middleware.NewRoundRobinBalancer(targets)))
			tenant.GET("/*", func(c echo.Context) error {
				return c.String(http.StatusOK, "Tenant:"+c.Request().Host)
			})
			hosts[service.IngressUrl] = &models.Host{Echo: tenant}
		} else if service.Type == "static" {
            // Static endpoint
            tenant.Use(middleware.GzipWithConfig(middleware.GzipConfig{
                Level: 5,
            }))
            tenant.Use(expiresServerHeader)
            tenant.Use(middleware.BodyLimit("10M"))
            tenant.Use(middleware.StaticWithConfig(middleware.StaticConfig{
                Root:   service.EgressUrl,
                Browse: true,
                HTML5:  true,
            }))
            // Add to Hosts
            hosts[service.IngressUrl] = &models.Host{Echo: tenant}
        }
	}

	//---------
	// ROOT
	//---------
	server := echo.New()
	server.Use(middleware.Logger())
	server.Use(middleware.Recover())
	server.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "{\"success\":\"ok\"}")
	})
	hosts[cfg.StatusPage] = &models.Host{Echo: server}



	// Server
	e := echo.New()
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
	// 4 Terabyte limit
	e.Use(middleware.BodyLimit("4T"))

	// Start server with Graceful Shutdown
	go func() {
		if err := e.StartTLS(":"+SSL_PORT, "server.crt", "server.key"); err != nil && err != http.ErrServerClosed {
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

// ServerHeader middleware adds a `Server` header to the response.
func expiresServerHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("Cache-Control", "public, max-age=3600")
		return next(c)
	}
}
