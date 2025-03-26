package model

import (
	"net/http"
	"regexp"
	"slender/internal/ip"
	"strings"
)

// HomeFolderListItem defines folder list item used by the homepage.
type HomeFolderListItem struct {
	ID    MyInt64  `db:"id"`
	Name  MyString `db:"name"`
	Des   MyString `db:"description"`
	Large MyBool   `db:"large"`
}

// HomeBookmarkListItem defines bookmark list item used by the homepage.
type HomeBookmarkListItem struct {
	ID       MyInt64  `db:"id"`
	URL      MyString `db:"url"`
	Name     MyString `db:"name"`
	Des      MyString `db:"description"`
	Icon     MyString `db:"icon"`
	FolderID MyInt64  `db:"folder_id"`

	HideInOther MyBool `db:"hide_in_other"`
}

// HomeSearchEngineListItem defines search engine list item used by the homepage.
type HomeSearchEngineListItem struct {
	ID     MyInt64  `db:"id"`
	Name   MyString `db:"name"`
	Method MyString `db:"method"`
	URL    MyString `db:"url"`
	Body   MyString `db:"body"`
	Icon   MyString `db:"icon"`
}

// PageDynamicURL defines dynamic URL.
type PageDynamicURL struct {
	Parsed bool

	Host     string
	Hostname string
	Href     string
	Origin   string
	Pathname string
	Port     string
	Protocol string
}

// Parse parses network request.
func (d *PageDynamicURL) Parse(r *http.Request) {
	d.Parsed = true

	scheme := ip.GetProtocol(r) + ":"
	defaultPort := "80"
	if scheme == "https:" {
		defaultPort = "443"
	}
	host := r.Host

	hostname := host
	port := defaultPort
	reg, _ := regexp.Compile(`([\w+\.-]+):(\d+)$`)
	portMatch := reg.FindStringSubmatch(host)
	if portMatch != nil {
		hostname = portMatch[1]
		port = portMatch[2]
	}

	d.Host = host
	d.Hostname = hostname
	d.Href = strings.Join([]string{scheme, "//", host, r.RequestURI}, "")
	d.Origin = strings.Join([]string{scheme, "//", host}, "")
	d.Pathname = r.URL.Path
	d.Port = port
	d.Protocol = scheme
}

// Convert converts URL.
func (d *PageDynamicURL) Convert(url string) (res string) {
	res = url

	if d.Parsed {
		replacements := map[string]string{
			"{host}":     d.Host,
			"{hostname}": d.Hostname,
			"{href}":     d.Href,
			"{origin}":   d.Origin,
			"{pathname}": d.Pathname,
			"{port}":     d.Port,
			"{protocol}": d.Protocol,
		}

		for key, value := range replacements {
			res = strings.ReplaceAll(res, key, value)
		}
	}

	return
}

// IsInSameNetwork checks if URL is in the same network.
func (d *PageDynamicURL) IsInSameNetwork(url string) bool {
	if d.Parsed {
		if !strings.HasPrefix(url, "http") {
			return true
		} else if strings.Contains(d.Convert(url), d.Hostname) {
			return true
		}
	}
	return false
}
