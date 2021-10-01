package main

import (
    "crypto/tls"
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
    "golang.org/x/crypto/acme"
    "golang.org/x/crypto/acme/autocert"
    "net/http"
)

func main() {
    e := echo.New()
    e.Use(middleware.Recover())
    e.Use(middleware.Logger())
    e.GET("/", func(c echo.Context) error {
        return c.HTML(http.StatusOK, `
			<h1>Welcome to Echo!</h1>
			<h3>TLS certificates automatically installed from Let's Encrypt :)</h3>
		`)
    })

    autoTLSManager := autocert.Manager{
        Prompt: autocert.AcceptTOS,
        // Cache certificates to avoid issues with rate limits (https://letsencrypt.org/docs/rate-limits)
        Cache: autocert.DirCache("/var/www/.cache"),
        //HostPolicy: autocert.HostWhitelist("<DOMAIN>"),
    }
    s := http.Server{
        Addr:    ":443",
        Handler: e, // set Echo as handler
        TLSConfig: &tls.Config{
            //Certificates: nil, // <-- s.ListenAndServeTLS will populate this field
            GetCertificate: autoTLSManager.GetCertificate,
            NextProtos:     []string{acme.ALPNProto},
        },
        //ReadTimeout: 30 * time.Second, // use custom timeouts
    }
    if err := s.ListenAndServeTLS("", ""); err != http.ErrServerClosed {
        e.Logger.Fatal(err)
    }
}
