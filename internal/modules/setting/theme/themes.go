// Package theme adalah katalog palet warna admin (theme switcher). Palet
// disimpan di kode (nama + warna), nama palet AKTIF disimpan di Setting.theme
// (DB). Layout memetakan warna ke CSS variable → ganti tema tanpa rebuild
// (padanan src/config/themes.ts di NodeAdmin).
package theme

// Theme = satu palet warna.
type Theme struct {
	Name    string `json:"name"`
	Primary string `json:"primary"`
	Accent  string `json:"accent"`
}

// Default = palet bawaan bila Setting.theme kosong/invalid.
const Default = "Blue"

// catalog = daftar palet tersedia (urut tampil).
var catalog = []Theme{
	{Name: "Blue", Primary: "#2563eb", Accent: "#1d4ed8"},
	{Name: "Green", Primary: "#16a34a", Accent: "#15803d"},
	{Name: "Purple", Primary: "#7c3aed", Accent: "#6d28d9"},
	{Name: "Rose", Primary: "#e11d48", Accent: "#be123c"},
	{Name: "Slate", Primary: "#475569", Accent: "#334155"},
}

// All mengembalikan salinan katalog (hindari mutasi dari luar).
func All() []Theme {
	out := make([]Theme, len(catalog))
	copy(out, catalog)
	return out
}

// Names mengembalikan nama-nama palet.
func Names() []string {
	out := make([]string, 0, len(catalog))
	for _, t := range catalog {
		out = append(out, t.Name)
	}
	return out
}

// ByName mengembalikan palet bernama name; fallback ke Default bila tak ada.
func ByName(name string) Theme {
	for _, t := range catalog {
		if t.Name == name {
			return t
		}
	}
	for _, t := range catalog {
		if t.Name == Default {
			return t
		}
	}
	return catalog[0]
}

// IsValid true bila name ada di katalog.
func IsValid(name string) bool {
	for _, t := range catalog {
		if t.Name == name {
			return true
		}
	}
	return false
}
