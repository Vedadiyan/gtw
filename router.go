package gtw

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/vedadiyan/gtw/internal/structutil"
)

type (
	Metadata    byte
	Status      = int
	Response    func(int, http.ResponseWriter)
	RouteValues map[string]any
	RouterError string
	Reader      http.Request
	HttpCtx     struct {
		Response http.ResponseWriter
		Request  struct {
			*Reader
			RouteValues RouteValues
		}
	}
	HttpError struct {
		Status  int
		Message any
	}
	Handler    func(*HttpCtx) (Status, Response)
	RouteTable struct {
		routes  map[int][]*Route
		configs map[string]Handler
	}
	Route struct {
		path        string
		method      string
		routeValues map[int]string
		routeParams map[int]string
		hash        string
	}
)

const (
	NO_MATCH_FOUND    RouterError = "no match found"
	NO_URL_REGISTERED RouterError = "no url registered"
)

var (
	_upgrader websocket.Upgrader
)

func NewRouteTable() *RouteTable {
	routeTable := RouteTable{
		routes:  map[int][]*Route{},
		configs: make(map[string]Handler),
	}
	return &routeTable
}

func (routerError RouterError) Error() string {
	return string(routerError)
}

func (r *Route) Bind(route *Route) map[string]any {
	rank := RouteCompare(r, route)
	if rank == 0 {
		return nil
	}
	routeValues := make(map[string]any)
	for key, value := range r.routeValues {
		k := r.routeValues[key]
		if value == "?" {
			k = r.routeParams[key]
		}
		routeValues[k] = route.routeValues[key]
	}
	return routeValues
}

func (route *Route) GetHash() string {
	return route.hash
}

func ParseRoute(url *url.URL, method string) *Route {
	routeValues := make(map[int]string)
	routeParams := make(map[int]string)
	for index, segment := range strings.Split(url.Path, "/") {
		if len(segment) == 0 {
			continue
		}
		if strings.HasPrefix(segment, ":") {
			routeValues[index] = "?"
			routeParams[index] = segment[1:]
			continue
		}
		routeValues[index] = segment
		routeValues[index] = segment
	}

	hash := CreateHash(url, method)
	route := Route{
		path:        url.Path,
		routeValues: routeValues,
		routeParams: routeParams,
		method:      strings.ToUpper(method),
		hash:        hash,
	}
	return &route
}

func RouteCompare(preferredRoute *Route, route *Route) int {
	if len(preferredRoute.routeValues) != len(route.routeValues) {
		return 0
	}
	rank := 1
	for key, value := range preferredRoute.routeValues {
		if value == "?" {
			rank += 1
			continue
		}
		if value != route.routeValues[key] {
			rank = 0
			break
		}
		rank += 2
	}
	return rank
}

func CreateHash(url *url.URL, method string) string {
	buffer := bytes.NewBufferString(strings.ToUpper(method))
	buffer.WriteString(":")
	buffer.WriteString(url.Path)
	sha256 := sha256.New()
	sha256.Write(buffer.Bytes())
	hash := hex.EncodeToString(sha256.Sum(nil))
	return hash
}

func (rt *RouteTable) Register(url *url.URL, method string, handlerFunc Handler) {
	route := ParseRoute(url, method)
	len := len(route.routeValues)
	if _, ok := rt.configs[route.hash]; ok {
		return
	}
	rt.configs[route.hash] = handlerFunc
	_, ok := rt.routes[len]
	if !ok {
		rt.routes[len] = make([]*Route, 0)
	}
	rt.routes[len] = append(rt.routes[len], route)
}

func (rt RouteTable) Find(url *url.URL, method string) (http.HandlerFunc, error) {
	if len(rt.routes) == 0 {
		return nil, NO_URL_REGISTERED
	}
	prt := ParseRoute(url, method)
	routes, ok := rt.routes[len(prt.routeValues)]
	if !ok {
		return nil, NO_MATCH_FOUND
	}
	lrnk := 0
	var lrt *Route
	for _, url := range routes {
		if url.method != strings.ToUpper(method) {
			continue
		}
		rnk := RouteCompare(url, prt)
		if rnk != 0 {
			if rnk > lrnk {
				lrnk = rnk
				lrt = url
			}
		}
	}
	if lrnk == 0 {
		return nil, NO_MATCH_FOUND
	}
	return func(w http.ResponseWriter, r *http.Request) {
		httpCtx := &HttpCtx{
			Response: w,
			Request: struct {
				*Reader
				RouteValues RouteValues
			}{
				Reader:      (*Reader)(r),
				RouteValues: lrt.Bind(prt),
			},
		}
		status, value := rt.GetHandlerFunc(lrt.hash)(httpCtx)
		value(status, w)

	}, nil
}

func (rt RouteTable) GetHandlerFunc(hash string) Handler {
	return rt.configs[hash]
}

func (rv RouteValues) Unmarshal(v any) error {
	return structutil.Unmarshal(rv, v)
}

func (r *Reader) Unmarshal(v any) error {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

func (httpCtx *HttpCtx) Upgrade(headers http.Header) (*websocket.Conn, error) {
	return _upgrader.Upgrade(httpCtx.Response, (*http.Request)(httpCtx.Request.Reader), headers)
}

func WithHeader(r func(status int, w http.ResponseWriter), h http.Header) func(status int, w http.ResponseWriter) {
	return func(status int, w http.ResponseWriter) {
		Header(h)(0, w)
		r(status, w)
	}
}

func Header(headers http.Header) func(_ int, w http.ResponseWriter) {
	return func(_ int, w http.ResponseWriter) {
		for key := range headers {
			w.Header().Add(key, headers.Get(key))
		}
	}
}

func JSON(v any) func(status int, w http.ResponseWriter) {
	return func(status int, w http.ResponseWriter) {
		json, _err := json.Marshal(v)
		if _err != nil {
			http.Error(w, _err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(status)
		w.Write(json)
	}
}

func Raw(data []byte) func(status int, w http.ResponseWriter) {
	return func(status int, w http.ResponseWriter) {
		w.WriteHeader(status)
		w.Write(data)
	}
}

func Empty() func(status int, w http.ResponseWriter) {
	return func(status int, w http.ResponseWriter) {
		w.WriteHeader(status)
	}
}
