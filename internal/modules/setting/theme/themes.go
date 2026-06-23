// Package theme adalah katalog palet tema admin (theme switcher) — PERSIS dengan
// NodeAdmin (@flazhost-nodeadmin/core THEMES). Satu set view didorong oleh 4
// warna (primary/secondary/light/dark) per tema via CSS variable + Tailwind
// config inline di chrome → ganti tema = ganti palet saat render, tanpa rebuild.
package theme

// Theme = satu palet warna (4 nilai, struktur identik antar-tema).
type Theme struct {
	Name      string `json:"name"`
	Primary   string `json:"primary"`
	Secondary string `json:"secondary"`
	Light     string `json:"light"`
	Dark      string `json:"dark"`
}

// Default = tema bawaan bila Setting.theme kosong/invalid.
const Default = "Blue"

// catalog = 9 palet PERSIS NodeAdmin (urut: Blue default lalu alfabetis).
var catalog = []Theme{
	{Name: "Blue", Primary: "#3B82F6", Secondary: "#60A5FA", Light: "#DBEAFE", Dark: "#1E40AF"},
	{Name: "Black", Primary: "#374151", Secondary: "#4B5563", Light: "#6B7280", Dark: "#1F2937"},
	{Name: "Brown", Primary: "#A16207", Secondary: "#D97706", Light: "#FEF3C7", Dark: "#78350F"},
	{Name: "Green", Primary: "#10B981", Secondary: "#34D399", Light: "#D1FAE5", Dark: "#047857"},
	{Name: "Grey", Primary: "#6B7280", Secondary: "#9CA3AF", Light: "#E5E7EB", Dark: "#374151"},
	{Name: "Orange", Primary: "#F59E0B", Secondary: "#FBBF24", Light: "#FEF3C7", Dark: "#D97706"},
	{Name: "Purple", Primary: "#8B5CF6", Secondary: "#A78BFA", Light: "#F3E8FF", Dark: "#6D28D9"},
	{Name: "Red", Primary: "#EF4444", Secondary: "#F87171", Light: "#FECACA", Dark: "#B91C1C"},
	{Name: "Yellow", Primary: "#F59E0B", Secondary: "#FCD34D", Light: "#FEF3C7", Dark: "#D97706"},
}

// All mengembalikan salinan katalog (untuk UI switcher).
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

// IsValid true bila name ada di katalog.
func IsValid(name string) bool {
	for _, t := range catalog {
		if t.Name == name {
			return true
		}
	}
	return false
}

// ByName mengembalikan palet bernama name; fallback ke Default.
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
