// Package view menyediakan helper render template HTML terpusat (padanan
// renderView() di NodeAdmin). Controller web WAJIB lewat RenderView — bukan
// memanggil c.HTML dengan path mentah (di-enforce checker).
//
// Template di-load sekali saat start (ParseGlob), memetakan layout + partial
// (head/sidebar/topbar/foot) + view modul. Helper Route() & hasAccess di-inject
// sebagai FuncMap agar named-routes & sidebar dinamis tersedia di template.
package view

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"

	"goadmin/internal/router"
)

// Engine membungkus template ter-parse + dipasang ke gin.
type Engine struct {
	tmpl *template.Template
}

// Funcs default tersedia di semua template.
func funcMap() template.FuncMap {
	return template.FuncMap{
		// route("nama", "id", "7") → URL bernama (named-routes di template).
		"route": func(name string, pairs ...string) string {
			params := map[string]string{}
			for i := 0; i+1 < len(pairs); i += 2 {
				params[pairs[i]] = pairs[i+1]
			}
			return router.Route(name, params)
		},
	}
}

// Load mem-parse seluruh template (layout + modul) dari glob pattern.
// Mengembalikan nil-aman bila tak ada file (mode api / belum ada view).
func Load(patterns ...string) (*Engine, error) {
	tmpl := template.New("").Funcs(funcMap())
	var loaded bool
	for _, p := range patterns {
		t, err := tmpl.ParseGlob(p)
		if err != nil {
			// Glob tanpa match bukan error fatal (modul mungkin tak punya view).
			continue
		}
		tmpl = t
		loaded = true
	}
	if !loaded {
		return &Engine{tmpl: template.New("").Funcs(funcMap())}, nil
	}
	return &Engine{tmpl: tmpl}, nil
}

// Attach memasang template engine ke gin.
func (e *Engine) Attach(r *gin.Engine) {
	r.SetHTMLTemplate(e.tmpl)
}

// RenderView merender satu view dengan locals (data) + status 200.
// Locals otomatis diperkaya dengan user terautentikasi bila ada di context.
func RenderView(c *gin.Context, name string, locals gin.H) {
	if locals == nil {
		locals = gin.H{}
	}
	if u, ok := c.Get("auth_user"); ok {
		locals["currentUser"] = u
	}
	// Token CSRF (diset middleware.CSRF) → form menyertakan <input name="_csrf">.
	if tok, ok := c.Get("csrf_token"); ok {
		locals["_csrf"] = tok
	}
	// Flash one-shot (diset middleware.Flash) → banner sukses/error.
	if v, ok := c.Get("flash_success"); ok {
		locals["flash_success"] = v
	}
	if v, ok := c.Get("flash_error"); ok {
		locals["flash_error"] = v
	}
	c.HTML(http.StatusOK, name, locals)
}
