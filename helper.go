/********************************
*** Multiplexer for Go        ***
*** Bone is under MIT license ***
*** Code by CodingFerret      ***
*** github.com/go-zoo         ***
*********************************/

package bone

import (
	"net/http"
	"strings"

	"github.com/go-zoo/duck"
)

func (m *Mux) parse(rw http.ResponseWriter, req *http.Request) bool {
	for _, r := range m.Routes[req.Method] {
		ok := r.parse(rw, req)
		if ok {
			return true
		}
	}
	// If no HEAD method, default to GET
	if req.Method == "HEAD" {
		for _, r := range m.Routes["GET"] {
			ok := r.parse(rw, req)
			if ok {
				return true
			}
		}
	}
	return false
}

// StaticRoute check if the request path is for Static route
func (m *Mux) staticRoute(rw http.ResponseWriter, req *http.Request) bool {
	for _, s := range m.Routes[static] {
		if len(req.URL.Path) >= s.Size {
			if req.URL.Path[:s.Size] == s.Path {
				s.Handler.ServeHTTP(rw, req)
				return true
			}
		}
	}
	return false
}

// HandleNotFound handle when a request does not match a registered handler.
func (m *Mux) HandleNotFound(rw http.ResponseWriter, req *http.Request) {
	if m.notFound != nil {
		m.notFound.ServeHTTP(rw, req)
	} else {
		http.NotFound(rw, req)
	}
}

// Check if the path don't end with a /
func (m *Mux) validate(rw http.ResponseWriter, req *http.Request) bool {
	if len(req.URL.Path) > 1 && req.URL.Path[len(req.URL.Path)-1:] == "/" {
		cleanURL(&req.URL.Path)
		rw.Header().Set("Location", req.URL.Path)
		rw.WriteHeader(http.StatusFound)
	}
	// Retry to find a route that match
	return m.parse(rw, req)
}

func valid(path string) bool {
	if len(path) > 1 && path[len(path)-1:] == "/" {
		return false
	}
	return true
}

// Clean url path
func cleanURL(url *string) {
	if len((*url)) > 1 {
		if (*url)[len((*url))-1:] == "/" {
			*url = (*url)[:len((*url))-1]
			cleanURL(url)
		}
	}
}

// GetValue return the key value, of the current *http.Request
func GetValue(req *http.Request, key string) string {
	if ok, value := extractParams(req); ok {
		return value[key]
	}
	return ""
}

// GetAllValues return the req PARAMs
func GetAllValues(req *http.Request) map[string]string {
	if ok, values := extractParams(req); ok {
		return values
	}
	return nil
}

func extractParams(req *http.Request) (bool, map[string]string) {
	var ss = strings.Split(req.URL.Path, "/")
	var params = make(map[string]string)
	var r = GetRequestRoute(req)
	if r != nil {
		if r.Atts&REGEX != 0 {
			for k := range r.Compile {
				params[r.Tag[k]] = ss[k]
			}
		}

		for k, v := range r.Pattern {
			params[v] = ss[k]
		}

		return true, params
	}
	return false, nil
}

// GetRequestRoute returns the route of given Request
func GetRequestRoute(req *http.Request) *Route {
	rt := duck.GetContext(req, "route")
	if rt != nil {
		return rt.(*Route)
	}
	return nil
}
