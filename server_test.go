package gtw

import (
	"net/http"
	"testing"

	"github.com/vedadiyan/gtw/internal/di"
)

type (
	TestAPI struct {
		Metadata `prefix:"api"`

		Test Service[int] `name:"test"`

		Get Handler `route:"/test/:name" method:"GET"`
	}
)

func (t *TestAPI) GetHandler(httpCtx *HttpCtx) (Status, Response) {
	output := map[string]any{
		"Hello": t.Test.Value(),
	}
	headers := http.Header{}
	headers.Add("x-test", "ok")
	return 201, WithHeader(JSON(output), headers)
}

func TestParse(t *testing.T) {
	di.AddSinletonWithName[int]("test", func() (instance *int, err error) {
		i := 0
		return &i, nil
	})
	server := New()
	err := server.Register(new(TestAPI))
	if err != nil {
		t.FailNow()
	}
	server.Cors(CorsAllowAll()).ListenAndServe(&http.Server{
		Addr: ":8080",
	})
}
