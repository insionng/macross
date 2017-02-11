// Copyright 2016 Insionng
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

// Usage
//
// template := "http://{{host}}/?q={{query}}&foo={{bar}}{{bar}}"
// t := femplate.New(template, "{{", "}}")
// s := t.ExecuteString(map[string]interface{}{
//        "host":  "google.com",
//        "query": url.QueryEscape("hello=world"),
//        "bar":   "foobar",
//    })
//    fmt.Printf("%s", s)
//
//    Output:
//    http://google.com/?q=hello%3Dworld&foo=foobarfoobar
//
//
// Advanced usage
//
// template := "Hello, [user]! You won [prize]!!! [foobar]"
// t, err := fasttemplate.NewTemplate(template, "[", "]")
// if err != nil {
// 	log.Fatalf("unexpected error when parsing template: %s", err)
// }
// s := t.ExecuteFuncString(func(w io.Writer, tag string) (int, error) {
// 	switch tag {
// 	case "user":
// 		return w.Write([]byte("John"))
// 	case "prize":
// 		return w.Write([]byte("$100500"))
// 	default:
// 		return w.Write([]byte(fmt.Sprintf("[unknown tag %q]", tag)))
// 	}
// })
// fmt.Printf("%s", s)
//
// Output:
// Hello, John! You won $100500!!! [unknown tag "foobar"]
//

package fempla

import (
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"sync"

	"github.com/insionng/macross"
	"github.com/insionng/macross/libraries/femplate"
)

type (
	Renderer struct {
		Option
		templates map[string]*femplate.Template
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
		templates: make(map[string]*femplate.Template),
	}
	return r
}

func (r *Renderer) fromFile(path string) (t *femplate.Template, err error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	t = femplate.New(string(buf), r.Option.DelimLeft, r.Option.DelimRight)
	return t, nil
}

func (r *Renderer) buildTemplatesCache(name string) (t *femplate.Template, err error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	t, err = r.fromFile(filepath.Join(r.Directory, name))
	if err != nil {
		return
	}
	r.templates[name] = t
	return
}

func (r *Renderer) getTemplate(name string) (t *femplate.Template, err error) {
	name = name + ".html"
	if r.Reload {
		return r.fromFile(filepath.Join(r.Directory, name))
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

/*
func getContext(templateData interface{}) map[string]interface{} {
	if templateData == nil {
		return nil
	}
	contextData, isMap := templateData.(map[string]interface{})
	if isMap {
		return contextData
	}
	return nil
}
*/

func (r *Renderer) Render(w io.Writer, name string, ctx *macross.Context) error {
	template, err := r.getTemplate(name)
	if err != nil {
		return err
	}
	s := template.ExecuteString(ctx.GetStore())
	//_, err = io.Copy(w, bytes.NewReader([]byte(s)))
	_, err = fmt.Fprintf(w, "%s", s)
	return err
}
