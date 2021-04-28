package confgotemplate

import (
	"text/template"

	bufferspool "github.com/sv-tools/buffers-pool"
	"github.com/sv-tools/conf"
)

// New creates Go Template Trasformer to parse and apply the stored templates.
// `funcs` paramater can be used to extend the list of default functions supported by Go Templates.
//         The `Get` function added by default to call the `conf.Get`.
// `data` parameter can be used to pass additional data to the `Execute` function.
func New(funcs template.FuncMap, data interface{}) conf.Transform {
	if funcs == nil {
		funcs = make(template.FuncMap)
	}

	return func(key string, value interface{}, c conf.Conf) interface{} {
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
