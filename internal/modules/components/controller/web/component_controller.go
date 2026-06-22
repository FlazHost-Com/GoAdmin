// Package web berisi controller halaman showcase komponen UI — acuan visual +
// markup untuk membuat elemen serupa (stat card, badge, tabel, form, alert,
// button). Statis, tanpa service. Render lewat view.RenderView (di-enforce checker).
package web

import (
	"github.com/gin-gonic/gin"

	"goadmin/internal/view"
)

// ComponentController menyajikan katalog komponen UI.
type ComponentController struct{}

// NewComponentController merakit controller (tanpa dependency).
func NewComponentController() *ComponentController {
	return &ComponentController{}
}

// Index → GET /admin/v1/components (showcase statis dengan data contoh).
func (ctl *ComponentController) Index(c *gin.Context) {
	view.RenderView(c, "components/index", gin.H{
		"title": "Komponen UI",
		"stats": []gin.H{
			{"label": "Pengguna", "value": 128, "cls": "from-blue-500 to-blue-600"},
			{"label": "Peran", "value": 4, "cls": "from-emerald-500 to-emerald-600"},
			{"label": "Izin", "value": 24, "cls": "from-violet-500 to-violet-600"},
		},
		"badges": []gin.H{
			{"text": "Active", "cls": "bg-emerald-100 text-emerald-700"},
			{"text": "Inactive", "cls": "bg-slate-100 text-slate-600"},
			{"text": "Blocked", "cls": "bg-rose-100 text-rose-700"},
		},
		"alerts": []gin.H{
			{"msg": "Data berhasil disimpan.", "cls": "bg-emerald-50 text-emerald-800 border-emerald-200"},
			{"msg": "Terjadi kesalahan saat memproses.", "cls": "bg-rose-50 text-rose-800 border-rose-200"},
			{"msg": "Cache pengaturan baru saja di-refresh.", "cls": "bg-blue-50 text-blue-800 border-blue-200"},
		},
		"rows": []gin.H{
			{"code": "U-001", "name": "Budi", "email": "budi@example.com", "status": "Active"},
			{"code": "U-002", "name": "Ani", "email": "ani@example.com", "status": "Inactive"},
			{"code": "U-003", "name": "Citra", "email": "citra@example.com", "status": "Active"},
		},
	})
}
