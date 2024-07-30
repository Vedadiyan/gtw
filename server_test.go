package gtw

import (
	"net/http"
	"testing"

	"github.com/vedadiyan/gtw/internal/di"
)

type (
	TestAPI struct {
		Metadata `prefix:"api"`
		Test     Inject[int] `name:"test"`
		Get      HandlerFunc `route:"/test/:name" method:"GET"`
	}
)

func (t *TestAPI) GetHandler(httpCtx *HttpCtx) (Status, Response) {
	output := map[string]any{
		"Hello": t.Test.Value(),
	}
	heders := http.Header{}
	heders.Add("x-test", "ok")
	return 200, WithHeader(JSON(output), Header(heders))
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
	server.ListenAndServe(":8081")
}
