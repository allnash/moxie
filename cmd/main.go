package main

import (
	"context"
	"fmt"
	"github.com/allnash/moxie/config"
	"github.com/allnash/moxie/ipfilter"
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

const AppYamlFilename = "/etc/moxie/app.yaml"

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
	server.Use(middleware.Recover())
	server.GET("/status", func(c echo.Context) error {
		return c.String(http.StatusOK, "{\"success\":\"ok\"}")
	})
	hosts[cfg.StatusHost+":"+cfg.ProxyListenPort] = &models.Host{Echo: server}

	// Server
	e := echo.New()
	e.Use(middleware.Logger())
	e.Logger.SetOutput(&lumberjack.Logger{
		Filename:   cfg.Logfile,
		MaxSize:    100, // megabytes
		MaxBackups: 3,
		MaxAge:     28,   //days
		Compress:   true, // disabled by default
	})
	e.Use(ipfilter.MiddlewareWithConfig(ipfilter.Config{
		Skipper: middleware.DefaultSkipper,
		BlackList: []string{
			"104.244.100.0/24",
			"104.244.101.0/24",
			"104.244.102.0/24",
			"104.244.103.0/24",
			"130.254.100.0/24",
			"130.254.101.0/24",
			"130.254.102.0/24",
			"130.254.103.0/24",
			"130.254.104.0/24",
			"130.254.105.0/24",
			"130.254.106.0/24",
			"130.254.107.0/24",
			"130.254.108.0/24",
			"130.254.109.0/24",
			"130.254.110.0/24",
			"130.254.111.0/24",
			"130.254.112.0/24",
			"130.254.113.0/24",
			"130.254.114.0/24",
			"130.254.115.0/24",
			"130.254.116.0/24",
			"130.254.117.0/24",
			"130.254.118.0/24",
			"130.254.119.0/24",
			"130.254.120.0/24",
			"130.254.121.0/24",
			"130.254.122.0/24",
			"130.254.123.0/24",
			"130.254.124.0/24",
			"130.254.125.0/24",
			"130.254.126.0/24",
			"130.254.127.0/24",
			"130.254.96.0/24",
			"130.254.97.0/24",
			"130.254.98.0/24",
			"130.254.99.0/24",
			"130.44.200.0/24",
			"130.44.201.0/24",
			"130.44.202.0/24",
			"130.44.203.0/24",
			"147.53.113.0/24",
			"147.53.114.0/24",
			"147.53.115.0/24",
			"147.53.116.0/24",
			"147.53.118.0/24",
			"147.53.119.0/24",
			"147.53.120.0/24",
			"147.53.121.0/24",
			"147.53.122.0/24",
			"147.53.123.0/24",
			"147.53.124.0/24",
			"147.53.125.0/24",
			"147.53.126.0/24",
			"147.53.127.0/24",
			"154.201.32.0/24",
			"154.201.33.0/24",
			"154.201.34.0/24",
			"154.201.35.0/24",
			"154.201.36.0/24",
			"154.201.37.0/24",
			"154.201.38.0/24",
			"154.201.39.0/24",
			"154.201.40.0/24",
			"154.201.41.0/24",
			"154.201.42.0/24",
			"154.201.43.0/24",
			"154.201.44.0/24",
			"154.201.45.0/24",
			"154.201.46.0/24",
			"154.201.47.0/24",
			"154.201.56.0/24",
			"154.201.57.0/24",
			"154.201.58.0/24",
			"154.201.59.0/24",
			"154.201.60.0/24",
			"154.201.61.0/24",
			"154.201.62.0/24",
			"154.201.63.0/24",
			"154.202.100.0/24",
			"154.202.101.0/24",
			"154.202.102.0/24",
			"154.202.103.0/24",
			"154.202.104.0/24",
			"154.202.105.0/24",
			"154.202.106.0/24",
			"154.202.107.0/24",
			"154.202.108.0/24",
			"154.202.109.0/24",
			"154.202.110.0/24",
			"154.202.111.0/24",
			"154.202.112.0/24",
			"154.202.113.0/24",
			"154.202.114.0/24",
			"154.202.115.0/24",
			"154.202.116.0/24",
			"154.202.117.0/24",
			"154.202.118.0/24",
			"154.202.119.0/24",
			"154.202.120.0/24",
			"154.202.121.0/24",
			"154.202.122.0/24",
			"154.202.123.0/24",
			"154.202.124.0/24",
			"154.202.125.0/24",
			"154.202.126.0/24",
			"154.202.127.0/24",
			"154.202.96.0/24",
			"154.202.97.0/24",
			"154.202.98.0/24",
			"154.202.99.0/24",
			"154.83.10.0/24",
			"154.83.11.0/24",
			"154.83.36.0/24",
			"154.83.37.0/24",
			"154.83.38.0/24",
			"154.83.39.0/24",
			"154.83.40.0/24",
			"154.83.41.0/24",
			"154.83.42.0/24",
			"154.83.43.0/24",
			"154.83.44.0/24",
			"154.83.45.0/24",
			"154.83.46.0/24",
			"154.83.47.0/24",
			"154.83.8.0/24",
			"154.83.9.0/24",
			"154.84.132.0/24",
			"154.84.133.0/24",
			"154.84.134.0/24",
			"154.84.135.0/24",
			"154.84.139.0/24",
			"154.84.140.0/24",
			"154.84.142.0/24",
			"154.84.143.0/24",
			"156.225.10.0/24",
			"156.225.11.0/24",
			"156.225.12.0/24",
			"156.225.13.0/24",
			"156.225.14.0/24",
			"156.225.15.0/24",
			"156.225.8.0/24",
			"156.225.9.0/24",
			"156.227.10.0/24",
			"156.227.13.0/24",
			"156.227.14.0/24",
			"156.227.15.0/24",
			"156.227.9.0/24",
			"156.239.48.0/24",
			"156.239.49.0/24",
			"156.239.50.0/24",
			"156.239.51.0/24",
			"156.239.52.0/24",
			"156.239.53.0/24",
			"156.239.54.0/24",
			"156.239.55.0/24",
			"156.239.63.0/24",
			"158.62.208.0/24",
			"158.62.209.0/24",
			"158.62.210.0/24",
			"158.62.211.0/24",
			"158.62.216.0/24",
			"158.62.217.0/24",
			"158.62.218.0/24",
			"158.62.219.0/24",
			"158.62.220.0/24",
			"158.62.221.0/24",
			"158.62.222.0/24",
			"158.62.223.0/24",
			"185.100.215.0/24",
			"185.139.27.0/24",
			"185.93.32.0/24",
			"192.171.88.0/24",
			"194.50.243.0/24",
			"199.101.136.0/24",
			"199.101.137.0/24",
			"199.101.142.0/24",
			"199.101.143.0/24",
			"207.254.88.0/24",
			"207.254.89.0/24",
			"207.254.90.0/24",
			"207.254.91.0/24",
			"207.254.92.0/24",
			"207.254.93.0/24",
			"207.254.94.0/24",
			"207.254.95.0/24",
			"208.52.181.0/24",
			"208.52.183.0/24",
			"208.89.240.0/24",
			"208.89.241.0/24",
			"208.89.242.0/24",
			"208.89.243.0/24",
			"45.128.78.0/24",
			"45.129.126.0/24",
			"45.141.178.0/24",
			"45.141.179.0/24",
			"45.199.129.0/24",
			"45.199.130.0/24",
			"45.199.134.0/24",
			"45.199.136.0/24",
			"45.199.137.0/24",
			"45.199.138.0/24",
			"45.92.246.0/24",
			"50.114.110.0/24",
			"50.114.111.0/24",
			"5.180.152.0/24",
			"67.216.236.0/24",
			"67.216.237.0/24",
			"69.58.64.0/20",
			"69.58.64.0/24",
			"69.58.65.0/24",
			"69.58.66.0/24",
			"69.58.67.0/24",
			"69.58.68.0/24",
			"69.58.69.0/24",
			"69.58.70.0/24",
			"69.58.71.0/24",
			"69.58.72.0/24",
			"69.58.73.0/24",
			"69.58.74.0/24",
			"69.58.75.0/24",
			"69.58.76.0/24",
			"69.58.77.0/24",
			"69.58.78.0/24",
			"69.58.79.0/24",
			"69.58.88.0/24",
			"69.58.89.0/24",
			"69.58.90.0/24",
			"69.58.91.0/24",
			"72.14.92.0/24",
			"72.14.93.0/24",
			"72.14.94.0/24",
			"72.14.95.0/24",
		},
		BlockByDefault: false,
	}))
	e.Any("/*", func(c echo.Context) (err error) {
		req := c.Request()
		res := c.Response()
		host := hosts[req.Host]
		if host == nil {
			err = echo.ErrNotFound
			e.Logger.Info("Resource Not found - " + req.Host)
		} else {
			host.Echo.ServeHTTP(res, req)
		}

		return
	})
	// 4 Terabyte limit
	e.Use(middleware.BodyLimit("4T"))

	// Start server with Graceful Shutdown WITH CERT
	//go func() {
	//	if err := e.StartTLS(":"+cfg.SSLPort,
	//		"/etc/moxie/ssl/server.crt",
	//		"/etc/moxie/ssl/server.key"); err != nil && err != http.ErrServerClosed {
	//		e.Logger.Fatal("shutting down the server")
	//	}
	//}()

	// Start server with Graceful Shutdown WITHOUT CERT
	go func() {
		if err := e.Start(":9000"); err != nil && err != http.ErrServerClosed {
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
