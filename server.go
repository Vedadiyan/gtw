package gtw

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"unsafe"
)

type (
	Server struct {
		mux        *http.ServeMux
		routeTable *RouteTable
	}
)

var (
	_defaultServer *Server
)

func init() {
	_defaultServer = New()
}

func DefaultServer() *Server {
	return _defaultServer
}

func New() *Server {
	mux := new(http.ServeMux)
	server := new(Server)
	server.mux = mux
	server.routeTable = NewRouteTable()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		route, err := server.routeTable.Find(r.URL, r.Method)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		route.ServeHTTP(w, r)
	})
	return server
}

func (srv *Server) Handle(route string, method string, handlerFunc Handler) error {
	url, err := url.Parse(route)
	if err != nil {
		return err
	}
	srv.routeTable.Register(url, method, handlerFunc)
	return nil
}

func (srv *Server) ListenAndServe(server *http.Server) error {
	server.Handler = srv.mux
	return server.ListenAndServe()
}

func (srv *Server) Register(v any) error {
	t := reflect.TypeOf(v)
	if t.Kind() != reflect.Pointer {
		return fmt.Errorf("expected pointer buy found value")
	}
	if t.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("expected struct but found %T", v)
	}
	val := reflect.ValueOf(v)
	handlerType := reflect.TypeOf(func(*HttpCtx) (Status, Response) { return 0, nil })
	metadataType := reflect.TypeOf(Metadata(0))
	lenOfFields := t.Elem().NumField()
	prefix := t.Elem().Name()
	for i := 0; i < lenOfFields; i++ {
		field := t.Elem().Field(i)
		if field.Name == "Metadata" && field.Type.AssignableTo(metadataType) {
			prefix = strings.TrimPrefix(field.Tag.Get("prefix"), "/")
			continue
		}
		if field.Type.AssignableTo(handlerType) {
			route, ok := field.Tag.Lookup("route")
			if !ok {
				route = field.Name
			}
			httpMethod, ok := field.Tag.Lookup("method")
			if !ok {
				httpMethod = "GET"
			}
			methodName := fmt.Sprintf("%sHandler", field.Name)
			_, ok = t.MethodByName(methodName)
			if !ok {
				continue
			}
			method := val.MethodByName(methodName).Interface().(func(*HttpCtx) (Status, Response))
			srv.Handle(fmt.Sprintf("/%s/%s", strings.TrimSuffix(prefix, "/"), strings.TrimPrefix(route, "/")), httpMethod, method)
			continue
		}
		if strings.HasPrefix(field.Type.Name(), "Service[") && field.Type.PkgPath() == "github.com/vedadiyan/gtw" {
			rf := val.Elem().Field(i)
			name, ok := field.Tag.Lookup("name")
			if ok {
				f := rf.FieldByName("name")
				f = reflect.NewAt(f.Type(), unsafe.Pointer(f.UnsafeAddr())).Elem()
				f.Set(reflect.ValueOf(name))
			}
		}
	}
	return nil
}
