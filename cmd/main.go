package main

import (
    "crypto/tls"
    "github.com/allnash/moxie/models"
    "github.com/joho/godotenv"
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
    "golang.org/x/crypto/acme"
    "golang.org/x/crypto/acme/autocert"
    "net/http"
    "os"
    "time"
)

const AppEnvFilename = "app.env"

func load() error {
    return godotenv.Load(AppEnvFilename)
}

func main() {

    // Create Echo server
    e := echo.New()

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
    api.Use(middleware.Logger())
    api.Use(middleware.Recover())

    hosts[os.Getenv("API")] = &models.Host{Echo: api}

    api.GET("/", func(c echo.Context) error {
        return c.String(http.StatusOK, "API")
    })

    //------
    // Asset
    //------

    blog := echo.New()
    blog.Use(middleware.Logger())
    blog.Use(middleware.Recover())

    hosts[os.Getenv("ASSET_DOMAIN")] = &models.Host{Echo: blog}

    blog.GET("/", func(c echo.Context) error {
        return c.String(http.StatusOK, "ASSET")
    })

    // Server
    e.Pre(middleware.HTTPSRedirect())
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

    //// Set up Proxy
    //proxy, err := url.Parse("http://localhost:8000")
    //if err != nil {
    //    e.Logger.Fatal(err)
    //}
    //apiTargets := []*middleware.ProxyTarget{
    //    {
    //        URL: proxy,
    //    },
    //}
    //e.Use(middleware.Proxy(middleware.NewRoundRobinBalancer(apiTargets)))

    //e.GET("/", func(c echo.Context) error {
    //    return c.HTML(http.StatusOK, `
	//		<h1>Welcome to Echo!</h1>
	//		<h3>TLS certificates automatically installed from Let's Encrypt :)</h3>
	//	`)
    //})

    autoTLSManager := autocert.Manager{
        Prompt: autocert.AcceptTOS,
        // Cache certificates to avoid issues with rate limits (https://letsencrypt.org/docs/rate-limits)
        Cache: autocert.DirCache("/var/www/.cache"),
        HostPolicy: autocert.HostWhitelist(os.Getenv("API_DOMAIN"), os.Getenv("ASSET_DOMAIN")),
    }
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
    }
}
