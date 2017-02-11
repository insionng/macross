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
		Option
		templates map[string]*template.Template
		lock      sync.RWMutex
	}

	Option struct {
		// Directory to load templates. Default is "templates"
		Directory string
		// Reload to reload templates everytime.
		Reload bool
		// DelimLeft "{{"
		DelimLeft string
		// DelimRight "}}"
		DelimRight string
	}
)

func perparOption(options []Option) Option {
	var opt Option
	if len(options) > 0 {
		opt = options[0]
	}
	if len(opt.Directory) == 0 {
		opt.Directory = "templates"
	}
	if len(opt.DelimLeft) == 0 {
		opt.DelimLeft = "{{"
	}
	if len(opt.DelimRight) == 0 {
		opt.DelimRight = "}}"
	}
	return opt
}

func Renderor(opt ...Option) *Renderer {
	o := perparOption(opt)
	r := &Renderer{
		Option:    o,
		templates: make(map[string]*template.Template),
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
	var okay bool
	if t, okay = r.templates[name]; !okay {
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
	template.Delims(r.DelimLeft, r.DelimRight)
	err = template.Execute(w, c.GetStore())
	return err
}
