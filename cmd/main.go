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
			"104.164.113.0/24",
			"104.164.163.0/24",
			"104.164.183.0/24",
			"104.164.35.0/24",
			"104.165.123.0/24",
			"104.165.127.0/24",
			"104.165.169.0/24",
			"104.165.232.0/24",
			"104.165.86.0/24",
			"104.165.92.0/24",
			"104.234.138.0/24",
			"104.234.143.0/24",
			"104.234.168.0/24",
			"104.252.131.0/24",
			"104.252.143.0/24",
			"104.252.19.0/24",
			"104.252.28.0/24",
			"104.252.30.0/24",
			"104.253.196.0/24",
			"107.186.7.0/24",
			"107.186.76.0/24",
			"136.0.36.0/24",
			"138.229.104.0/22",
			"138.229.108.0/24",
			"138.229.109.0/24",
			"138.229.110.0/24",
			"138.229.111.0/24",
			"138.229.96.0/21",
			"139.180.0.0/21",
			"139.180.224.0/24",
			"139.180.225.0/24",
			"139.180.226.0/24",
			"139.180.227.0/24",
			"139.180.228.0/22",
			"139.60.101.0/24",
			"141.11.10.0/24",
			"141.11.2.0/24",
			"141.11.23.0/24",
			"141.11.24.0/24",
			"141.11.25.0/24",
			"141.11.252.0/24",
			"141.11.253.0/24",
			"141.11.29.0/24",
			"141.11.3.0/24",
			"141.11.46.0/24",
			"141.11.47.0/24",
			"141.11.66.0/24",
			"141.11.67.0/24",
			"141.11.76.0/24",
			"141.11.77.0/24",
			"141.11.78.0/24",
			"141.11.79.0/24",
			"141.164.80.0/24",
			"141.164.81.0/24",
			"141.164.82.0/24",
			"141.164.83.0/24",
			"141.164.84.0/24",
			"141.164.85.0/24",
			"141.164.86.0/24",
			"141.164.87.0/24",
			"141.164.88.0/24",
			"141.164.89.0/24",
			"141.164.90.0/24",
			"141.164.91.0/24",
			"141.164.92.0/24",
			"141.164.93.0/24",
			"141.164.94.0/24",
			"141.164.95.0/24",
			"142.111.142.0/24",
			"142.111.152.0/24",
			"142.147.104.0/24",
			"142.147.105.0/24",
			"142.147.106.0/24",
			"142.147.107.0/24",
			"142.147.108.0/24",
			"142.147.109.0/24",
			"142.147.110.0/24",
			"142.147.111.0/24",
			"142.252.145.0/24",
			"142.252.215.0/24",
			"142.252.37.0/24",
			"147.185.107.0/24",
			"147.185.162.0/24",
			"147.53.112.0/24",
			"147.53.117.0/24",
			"147.92.52.0/22",
			"148.59.146.0/24",
			"148.59.184.0/24",
			"148.59.185.0/24",
			"149.20.240.0/24",
			"149.20.241.0/24",
			"149.20.242.0/24",
			"149.20.243.0/24",
			"149.20.244.0/24",
			"149.20.245.0/24",
			"149.20.246.0/24",
			"149.20.247.0/24",
			"152.44.100.0/22",
			"152.44.104.0/24",
			"152.44.105.0/24",
			"152.44.106.0/24",
			"152.44.107.0/24",
			"152.44.108.0/24",
			"152.44.109.0/24",
			"152.44.110.0/24",
			"152.44.111.0/24",
			"152.44.96.0/24",
			"152.44.97.0/24",
			"152.44.98.0/23",
			"154.16.122.0/24",
			"154.83.32.0/24",
			"154.83.33.0/24",
			"154.83.34.0/24",
			"154.83.35.0/24",
			"154.84.128.0/24",
			"154.84.129.0/24",
			"154.84.130.0/24",
			"154.84.131.0/24",
			"154.84.136.0/24",
			"154.84.137.0/24",
			"154.84.138.0/24",
			"154.84.141.0/24",
			"156.227.11.0/24",
			"156.227.12.0/24",
			"156.227.8.0/24",
			"156.239.32.0/24",
			"156.239.33.0/24",
			"156.239.34.0/24",
			"156.239.35.0/24",
			"156.239.36.0/24",
			"156.239.37.0/24",
			"156.239.38.0/24",
			"156.239.39.0/24",
			"156.239.40.0/24",
			"156.239.41.0/24",
			"156.239.42.0/24",
			"156.239.43.0/24",
			"156.239.44.0/24",
			"156.239.45.0/24",
			"156.239.46.0/24",
			"156.239.47.0/24",
			"156.239.56.0/24",
			"156.239.57.0/24",
			"156.239.58.0/24",
			"156.239.59.0/24",
			"156.239.60.0/24",
			"156.239.61.0/24",
			"156.239.62.0/24",
			"156.248.100.0/24",
			"156.248.101.0/24",
			"156.248.102.0/24",
			"156.248.103.0/24",
			"156.248.64.0/24",
			"156.248.65.0/24",
			"156.248.66.0/24",
			"156.248.67.0/24",
			"156.248.68.0/24",
			"156.248.69.0/24",
			"156.248.70.0/24",
			"156.248.71.0/24",
			"156.248.96.0/24",
			"156.248.97.0/24",
			"156.248.98.0/24",
			"156.248.99.0/24",
			"156.252.23.0/24",
			"158.62.213.0/24",
			"162.223.122.0/24",
			"162.244.144.0/21",
			"163.5.129.0/24",
			"165.140.11.0/24",
			"166.0.217.0/24",
			"166.0.219.0/24",
			"166.0.223.0/24",
			"166.88.129.0/24",
			"166.88.213.0/24",
			"166.88.220.0/24",
			"166.88.244.0/24",
			"166.88.58.0/24",
			"166.88.68.0/24",
			"167.160.64.0/21",
			"167.160.72.0/21",
			"167.160.72.0/24",
			"167.160.73.0/24",
			"167.160.74.0/24",
			"167.160.75.0/24",
			"167.160.76.0/24",
			"167.160.77.0/24",
			"167.160.78.0/24",
			"167.160.79.0/24",
			"168.245.143.0/24",
			"168.245.206.0/24",
			"168.91.10.0/24",
			"168.91.11.0/24",
			"168.91.12.0/24",
			"168.91.13.0/24",
			"168.91.14.0/24",
			"168.91.15.0/24",
			"168.91.32.0/24",
			"168.91.33.0/24",
			"168.91.34.0/24",
			"168.91.35.0/24",
			"168.91.36.0/24",
			"168.91.37.0/24",
			"168.91.38.0/24",
			"168.91.39.0/24",
			"168.91.40.0/23",
			"168.91.42.0/24",
			"168.91.43.0/24",
			"168.91.44.0/24",
			"168.91.45.0/24",
			"168.91.46.0/24",
			"168.91.47.0/24",
			"168.91.8.0/24",
			"168.91.9.0/24",
			"170.199.224.0/24",
			"170.199.225.0/24",
			"170.199.226.0/24",
			"170.199.227.0/24",
			"170.199.228.0/24",
			"170.199.229.0/24",
			"170.199.230.0/24",
			"170.199.231.0/24",
			"172.121.241.0/24",
			"172.121.255.0/24",
			"172.121.99.0/24",
			"172.252.10.0/24",
			"172.252.125.0/24",
			"172.252.133.0/24",
			"172.252.232.0/24",
			"172.252.233.0/24",
			"172.252.24.0/24",
			"172.252.56.0/24",
			"172.252.58.0/24",
			"172.81.112.0/22",
			"172.96.80.0/24",
			"172.96.81.0/24",
			"172.96.82.0/24",
			"172.96.83.0/24",
			"172.96.84.0/24",
			"172.96.85.0/24",
			"172.96.86.0/24",
			"172.96.87.0/24",
			"172.96.88.0/24",
			"172.96.89.0/24",
			"172.96.90.0/24",
			"172.96.91.0/24",
			"172.96.92.0/24",
			"172.96.93.0/24",
			"172.96.94.0/24",
			"172.96.95.0/24",
			"173.214.192.0/22",
			"173.245.93.0/24",
			"179.61.225.0/24",
			"181.214.229.0/24",
			"185.129.108.0/24",
			"185.129.109.0/24",
			"185.182.65.0/24",
			"185.35.78.0/24",
			"185.77.249.0/24",
			"185.92.46.0/24",
			"188.214.232.0/24",
			"188.214.233.0/24",
			"191.101.100.0/24",
			"191.101.102.0/24",
			"192.171.80.0/24",
			"192.171.81.0/24",
			"192.171.82.0/24",
			"192.171.83.0/24",
			"192.171.84.0/24",
			"192.171.85.0/24",
			"192.171.86.0/24",
			"192.171.87.0/24",
			"192.171.89.0/24",
			"192.171.90.0/24",
			"192.171.91.0/24",
			"192.171.92.0/24",
			"192.171.93.0/24",
			"192.171.94.0/24",
			"192.171.95.0/24",
			"192.177.109.0/24",
			"192.177.128.0/24",
			"192.177.129.0/24",
			"192.177.130.0/24",
			"192.177.131.0/24",
			"192.177.132.0/24",
			"192.177.133.0/24",
			"192.177.134.0/24",
			"192.177.135.0/24",
			"192.177.136.0/24",
			"192.177.137.0/24",
			"192.177.138.0/24",
			"192.177.139.0/24",
			"192.177.140.0/24",
			"192.177.141.0/24",
			"192.177.142.0/24",
			"192.177.143.0/24",
			"192.177.144.0/24",
			"192.177.145.0/24",
			"192.177.146.0/24",
			"192.177.147.0/24",
			"192.177.148.0/24",
			"192.177.149.0/24",
			"192.177.150.0/24",
			"192.177.151.0/24",
			"192.177.152.0/24",
			"192.177.153.0/24",
			"192.177.154.0/24",
			"192.177.155.0/24",
			"192.177.156.0/24",
			"192.177.157.0/24",
			"192.177.158.0/24",
			"192.177.159.0/24",
			"192.177.161.0/24",
			"192.177.164.0/24",
			"192.177.167.0/24",
			"192.177.172.0/24",
			"192.177.174.0/24",
			"192.177.175.0/24",
			"192.177.176.0/24",
			"192.177.177.0/24",
			"192.177.178.0/24",
			"192.177.180.0/24",
			"192.177.184.0/24",
			"192.177.187.0/24",
			"192.177.33.0/24",
			"192.177.40.0/24",
			"192.177.56.0/24",
			"192.177.69.0/24",
			"192.177.82.0/24",
			"192.177.98.0/24",
			"193.109.195.0/24",
			"193.142.18.0/24",
			"193.142.4.0/24",
			"193.161.245.0/24",
			"193.228.90.0/24",
			"193.43.143.0/24",
			"194.233.148.0/24",
			"194.233.149.0/24",
			"194.35.225.0/24",
			"194.35.226.0/24",
			"195.180.137.0/24",
			"195.180.149.0/24",
			"199.182.115.0/24",
			"199.250.188.0/23",
			"199.34.83.0/24",
			"199.34.84.0/24",
			"199.34.85.0/24",
			"199.34.86.0/24",
			"199.34.87.0/24",
			"199.34.88.0/24",
			"199.34.89.0/24",
			"199.34.90.0/24",
			"202.43.5.0/24",
			"204.10.18.0/24",
			"204.10.19.0/24",
			"205.164.11.0/24",
			"205.164.28.0/24",
			"205.164.46.0/24",
			"206.198.216.0/22",
			"207.182.24.0/24",
			"207.182.25.0/24",
			"207.182.26.0/24",
			"207.182.27.0/24",
			"207.182.28.0/24",
			"207.182.29.0/24",
			"207.182.30.0/24",
			"207.182.31.0/24",
			"207.229.93.0/24",
			"208.103.166.0/24",
			"209.163.116.0/22",
			"209.251.16.0/24",
			"209.251.17.0/24",
			"209.251.18.0/24",
			"209.251.19.0/24",
			"209.251.20.0/24",
			"209.251.21.0/24",
			"209.251.22.0/24",
			"209.251.23.0/24",
			"209.59.228.0/22",
			"209.73.147.0/24",
			"216.163.199.0/24",
			"216.172.136.0/24",
			"216.180.104.0/24",
			"216.180.105.0/24",
			"216.180.106.0/24",
			"216.180.107.0/24",
			"216.180.108.0/24",
			"216.180.109.0/24",
			"216.180.110.0/24",
			"216.180.111.0/24",
			"216.213.24.0/24",
			"216.213.25.0/24",
			"216.213.26.0/24",
			"216.213.27.0/24",
			"216.213.28.0/24",
			"216.213.29.0/24",
			"216.213.30.0/24",
			"216.213.31.0/24",
			"216.230.30.0/23",
			"216.41.232.0/22",
			"217.19.1.0/24",
			"23.157.192.0/24",
			"23.160.128.0/24",
			"23.230.111.0/24",
			"23.230.12.0/24",
			"23.230.144.0/24",
			"23.230.145.0/24",
			"23.230.167.0/24",
			"23.230.181.0/24",
			"23.230.217.0/24",
			"23.230.219.0/24",
			"23.230.238.0/24",
			"23.230.252.0/24",
			"23.230.39.0/24",
			"23.230.42.0/24",
			"23.230.69.0/24",
			"23.230.70.0/24",
			"23.27.172.0/24",
			"23.27.174.0/24",
			"23.27.186.0/24",
			"23.27.240.0/24",
			"23.27.253.0/24",
			"24.235.12.0/24",
			"24.235.13.0/24",
			"45.141.177.0/24",
			"45.199.139.0/24",
			"45.199.140.0/24",
			"45.199.141.0/24",
			"45.38.158.0/24",
			"45.38.242.0/24",
			"45.38.58.0/24",
			"45.39.212.0/24",
			"45.39.243.0/24",
			"45.39.249.0/24",
			"45.39.72.0/24",
			"50.117.56.0/24",
			"50.118.137.0/24",
			"50.118.138.0/24",
			"50.118.145.0/24",
			"50.118.158.0/24",
			"50.118.189.0/24",
			"50.118.206.0/24",
			"50.118.252.0/24",
			"5.1.40.0/24",
			"52.124.18.0/24",
			"52.124.19.0/24",
			"52.128.31.0/24",
			"64.112.96.0/24",
			"64.112.97.0/24",
			"64.112.98.0/24",
			"66.146.232.0/24",
			"66.146.233.0/24",
			"66.146.234.0/24",
			"66.146.235.0/24",
			"66.146.236.0/24",
			"66.146.237.0/24",
			"66.146.238.0/24",
			"66.146.239.0/24",
			"66.84.88.0/24",
			"66.84.89.0/24",
			"66.84.90.0/24",
			"66.84.91.0/24",
			"66.84.92.0/24",
			"66.84.93.0/24",
			"66.84.94.0/24",
			"66.84.95.0/24",
			"66.97.179.0/24",
			"67.218.4.0/24",
			"67.218.5.0/24",
			"67.226.219.0/24",
			"68.234.40.0/24",
			"68.234.41.0/24",
			"68.234.43.0/24",
			"68.234.44.0/24",
			"68.234.45.0/24",
			"68.234.46.0/24",
			"68.234.47.0/24",
			"68.65.220.0/24",
			"68.65.221.0/24",
			"68.65.222.0/24",
			"68.65.223.0/24",
			"82.115.10.0/24",
			"82.115.11.0/24",
			"82.115.8.0/24",
			"82.115.9.0/24",
			"85.204.37.0/24",
			"85.209.220.0/24",
			"85.209.231.0/24",
			"93.115.155.0/24",
			"98.158.232.0/24",
			"98.158.233.0/24",
			"98.158.234.0/24",
			"98.158.235.0/24",
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
