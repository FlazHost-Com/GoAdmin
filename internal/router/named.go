// Package router menyediakan named-routes + registry modul (padanan namedRoutes
// & registrasi modul di NodeAdmin). URL dirujuk lewat NAMA (Route("admin.v1.user.index"))
// bukan string hardcode → mudah refactor.
package router

import (
	"regexp"
	"strings"
	"sync"
)

// namedRegistry menyimpan pemetaan nama → pola path.
type namedRegistry struct {
	mu     sync.RWMutex
	routes map[string]string
}

var registry = &namedRegistry{routes: map[string]string{}}

// Register menyimpan nama→path. Dipanggil saat modul mendaftarkan route.
func Register(name, path string) {
	registry.mu.Lock()
	defer registry.mu.Unlock()
	registry.routes[name] = path
}

var paramRe = regexp.MustCompile(`:([A-Za-z0-9_]+)`)

// Route mengembalikan URL untuk nama tertentu, mensubstitusi parameter `:id`
// dengan nilai dari params (urutan map tak relevan — substitusi by-key).
//
// Contoh: Route("admin.v1.user.edit", map[string]string{"id": "7"})
//         → "/admin/v1/users/7/edit"
// Bila nama tak terdaftar, mengembalikan "#" (gagal jelas di UI, bukan panik).
func Route(name string, params ...map[string]string) string {
	registry.mu.RLock()
	path, ok := registry.routes[name]
	registry.mu.RUnlock()
	if !ok {
		return "#"
	}
	if len(params) == 0 || len(params[0]) == 0 {
		return path
	}
	p := params[0]
	return paramRe.ReplaceAllStringFunc(path, func(m string) string {
		key := strings.TrimPrefix(m, ":")
		if v, ok := p[key]; ok {
			return v
		}
		return m
	})
}

// All mengembalikan salinan seluruh named routes (untuk debug/template helper).
func All() map[string]string {
	registry.mu.RLock()
	defer registry.mu.RUnlock()
	out := make(map[string]string, len(registry.routes))
	for k, v := range registry.routes {
		out[k] = v
	}
	return out
}

// reset hanya untuk test.
func reset() {
	registry.mu.Lock()
	defer registry.mu.Unlock()
	registry.routes = map[string]string{}
}
