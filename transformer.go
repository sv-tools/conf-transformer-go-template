package confgotemplate

import (
	"text/template"

	bufferspool "github.com/sv-tools/buffers-pool"
	"github.com/sv-tools/conf"
)

// New creates Go Template Transformer to parse and apply the stored templates.
// `funcs` parameter can be used to extend the list of default functions supported by Go Templates.
//
// The `Get` function added to `funcs` map by default to call the `conf.Get`.
//
// The `data` parameter can be used to pass additional data to the `Execute` function.
func New(funcs template.FuncMap, data any) conf.Transform {
	if funcs == nil {
		funcs = make(template.FuncMap)
	}

	return func(key string, value any, c conf.Conf) any {
		funcs["Get"] = c.Get

		text, ok := value.(string)
		if !ok {
			return value
		}

		tmpl, err := template.New(key).Funcs(funcs).Parse(text)
		if err != nil {
			return value
		}

		buf := bufferspool.Get()
		defer bufferspool.Put(buf)

		if err := tmpl.Execute(buf, data); err != nil {
			return value
		}

		return buf.String()
	}
}
