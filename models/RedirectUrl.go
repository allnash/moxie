package models

import (
    "net/url"
)

// ProxyTarget defines the upstream target.
type RedirectUrl struct {
    URL  *url.URL
}