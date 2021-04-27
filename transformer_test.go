package confgotemplate_test

import (
	"fmt"
	"testing"
	"text/template"

	"github.com/stretchr/testify/require"
	"github.com/sv-tools/conf"

	confgotemplate "github.com/sv-tools/conf-transformer-go-template"
)

func TestTransformer(t *testing.T) {
	c := conf.New().WithTransformers(confgotemplate.New(nil, nil))
	c.Set("foo", `{{ Get "bar" }}`)
	c.Set("bar", 42)

	require.Equal(t, 42, c.Get("bar"))
	require.Equal(t, "42", c.Get("foo"))
}

func TestTransformerErrors(t *testing.T) {
	f := func(v int) int {
		return v * v
	}
	c := conf.New().WithTransformers(confgotemplate.New(template.FuncMap{"F": f}, nil))

	c.Set("foo", `{{ Get "bar" `)
	require.Equal(t, `{{ Get "bar" `, c.Get("foo"))

	c.Set("foo", `{{ F "bar" }}`)
	require.Equal(t, `{{ F "bar" }}`, c.Get("foo"))
}

func ExampleNew() {
	c := conf.New().WithTransformers(confgotemplate.New(nil, nil))
	c.Set("foo", `{{ Get "bar" }}`)
	c.Set("bar", 42)

	fmt.Println(c.Get("foo"))
	// Output: 42
}

func ExampleNew_double_func() {
	double := func(v int) int {
		return 2 * v
	}
	c := conf.New().WithTransformers(confgotemplate.New(template.FuncMap{"D": double}, nil))
	c.Set("foo", `{{ Get "bar" | D }}`)
	c.Set("bar", 42)

	fmt.Println(c.Get("foo"))
	// Output: 84
}
