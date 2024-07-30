package gtw

import (
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
	return 200, JSON(map[string]any{
		"Hello": t.Test.Value(),
	})
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
