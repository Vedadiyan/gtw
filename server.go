package gtw

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"
	"unsafe"

	"github.com/vedadiyan/gtw/internal/di"
)

type (
	Service[T any] struct {
		name     string
		hasScope bool
		scopeId  uint64
		ttl      time.Duration
	}
	Server struct {
		mux        *http.ServeMux
		routeTable *RouteTable
	}
)

func AddSingleton[T any](fn func() (instance *T, err error)) {
	di.AddSinleton[T](fn)
}

func AddSingletonWithName[T any](name string, fn func() (instance *T, err error)) {
	di.AddSinletonWithName[T](name, fn)
}

func AddTransient[T any](fn func() (instance *T, err error)) {
	di.AddTransient[T](fn)
}

func AddTransientWithName[T any](name string, fn func() (instance *T, err error)) {
	di.AddTransientWithName[T](name, fn)
}

func AddScoped[T any](fn func() (instance *T, err error)) {
	di.AddScoped[T](fn)
}

func AddScopedWithName[T any](name string, fn func() (instance *T, err error)) {
	di.AddScopedWithName[T](name, fn)
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

func (srv *Server) ListenAndServe(addr string) error {
	server := http.Server{
		Addr:    addr,
		Handler: srv.mux,
	}
	return server.ListenAndServe()
}

func (srv *Server) Register(v any) error {
	t := reflect.TypeOf(v)
	if t.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("expected struct but found %T", v)
	}
	val := reflect.ValueOf(v)
	handlerType := reflect.TypeOf(func(*HttpCtx) (Status, Response) { return 0, nil })
	metadataType := reflect.TypeOf(Metadata(0))
	lenOfFields := t.Elem().NumField()
	prefix := ""
	for i := 0; i < lenOfFields; i++ {
		field := t.Elem().Field(i)
		if field.Name == "Metadata" && field.Type.AssignableTo(metadataType) {
			prefix = strings.TrimPrefix(field.Tag.Get("prefix"), "/")
			continue
		}
		if field.Type.AssignableTo(handlerType) {
			route := field.Tag.Get("route")
			httpMethod := field.Tag.Get("method")
			methodName := fmt.Sprintf("%sHandler", field.Name)
			_, ok := t.MethodByName(methodName)
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

func (i *Service[T]) Value() *T {
	var options *di.Options
	if i.hasScope {
		options = di.NewOptions(i.scopeId, i.ttl)
	}
	if len(i.name) == 0 {
		return di.ResolveOrPanic[T](options)
	}
	return di.ResolveWithNameOrPanic[T](i.name, options)
}

func (i *Service[T]) Scope(scopeId uint64, ttl time.Duration) *Service[T] {
	copy := *i
	copy.hasScope = true
	copy.scopeId = scopeId
	copy.ttl = ttl
	return &copy
}
