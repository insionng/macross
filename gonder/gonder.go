package gonder

import (
	"github.com/insionng/macross"
	"io"
	"path/filepath"
	"sync"
	"text/template"
)

type (
	Renderer struct {
		PongorOption
		templates map[string]*template.Template
		lock      sync.RWMutex
	}

	PongorOption struct {
		// Directory to load templates. Default is "templates"
		Directory string
		// Reload to reload templates everytime.
		Reload bool
	}
)

func perparOption(options []PongorOption) PongorOption {
	var opt PongorOption
	if len(options) > 0 {
		opt = options[0]
	}
	if len(opt.Directory) == 0 {
		opt.Directory = "templates"
	}
	return opt
}

func Renderor(opt ...PongorOption) *Renderer {
	o := perparOption(opt)
	r := &Renderer{
		PongorOption: o,
		templates:    make(map[string]*template.Template),
	}
	return r
}

func (r *Renderer) buildTemplatesCache(name string) (t *template.Template, err error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	t, err = template.ParseFiles(filepath.Join(r.Directory, name))
	if err != nil {
		return
	}
	r.templates[name] = t
	return
}

func (r *Renderer) getTemplate(name string) (t *template.Template, err error) {
	name = name + ".html"
	if r.Reload {
		return template.ParseFiles(filepath.Join(r.Directory, name))
	}
	r.lock.RLock()
	var ok bool
	if t, ok = r.templates[name]; !ok {
		r.lock.RUnlock()
		t, err = r.buildTemplatesCache(name)
	} else {
		r.lock.RUnlock()
	}
	return
}

// Render 渲染
func (r *Renderer) Render(w io.Writer, name string, c *macross.Context) (err error) {
	template, err := r.getTemplate(name)
	if err != nil {
		return err
	}
	err = template.Execute(w, c.GetStore())
	return err
}
