// Package fetemplate adalah frontend template switcher (landing) — port konsep
// NodeAdmin (katalog opentailwind). Dua sumber template:
//
//   - BUILTIN (default, minimal) → dirender lewat Go view "home/<slug>".
//   - EKSTERNAL (slug pola opentailwind) → HTML diunduh on-demand & di-cache,
//     lalu disajikan sebagai HTML mentah.
//
// File ini bagian MURNI (tanpa jaringan/state): tipe, validasi slug (anti-SSRF),
// derive metadata, builtins, dan katalog kurasi fallback.
package fetemplate

import (
	"regexp"
	"strings"
)

// Template = satu desain landing di katalog.
type Template struct {
	Slug     string `json:"slug"`
	Name     string `json:"name"`
	Category string `json:"category"`
	Builtin  bool   `json:"builtin"`
}

// DefaultSlug = template bawaan (selalu tersedia, dirender Go view).
const DefaultSlug = "default"

// slugRe membatasi slug eksternal ke pola opentailwind `{kategori}-{NNN}-{nama}`
// (anti-SSRF: charset & struktur tetap → tak bisa memaksa fetch URL sembarang).
var slugRe = regexp.MustCompile(`^([a-z]+(?:-[a-z]+)*)-(\d{3})-([a-z0-9-]+)$`)

// builtins = template GoAdmin (Go view "home/<slug>").
var builtins = []Template{
	{Slug: "default", Name: "Klasik", Category: "Bawaan", Builtin: true},
	{Slug: "minimal", Name: "Minimalis", Category: "Bawaan", Builtin: true},
}

// curated = katalog eksternal kurasi (fallback saat remote mati / Remote=false).
var curated = []Template{
	{Slug: "agency-consulting-002-creative-agency", Name: "Creative Agency", Category: "Agency"},
	{Slug: "agency-consulting-001-digital-marketing-agency", Name: "Digital Marketing Agency", Category: "Agency"},
	{Slug: "technology-saas-001-hero-focused-conversion-page", Name: "Saas Hero Focused Conversion Page", Category: "Technology"},
	{Slug: "ecommerce-retail-001-fashion-boutique", Name: "Fashion Boutique", Category: "Ecommerce"},
	{Slug: "portfolio-creative-001-creative-portfolio", Name: "Creative Portfolio", Category: "Portfolio"},
	{Slug: "professional-services-001-law-firm", Name: "Law Firm", Category: "Professional"},
	{Slug: "real-estate-property-001-real-estate-agency", Name: "Real Estate Agency", Category: "Real Estate"},
	{Slug: "food-hospitality-001-fine-dining-restaurant", Name: "Fine Dining Restaurant", Category: "Food"},
	{Slug: "healthcare-wellness-001-family-doctor-clinic", Name: "Family Doctor Clinic", Category: "Healthcare"},
	{Slug: "education-training-001-private-school", Name: "Private School", Category: "Education"},
	{Slug: "fitness-sports-001-fitness-center", Name: "Fitness Center", Category: "Fitness"},
	{Slug: "travel-tourism-001-travel-agency", Name: "Travel Agency", Category: "Travel"},
}

// Builtins mengembalikan salinan template bawaan.
func Builtins() []Template {
	out := make([]Template, len(builtins))
	copy(out, builtins)
	return out
}

// Curated mengembalikan salinan katalog kurasi eksternal (fallback).
func Curated() []Template {
	out := make([]Template, len(curated))
	copy(out, curated)
	return out
}

// IsBuiltin true bila slug = template bawaan GoAdmin.
func IsBuiltin(slug string) bool {
	for _, t := range builtins {
		if t.Slug == slug {
			return true
		}
	}
	return false
}

// IsValidSlug true bila slug = builtin atau cocok pola opentailwind. Inilah
// gerbang ANTI-SSRF: hanya slug valid yang boleh di-fetch/unduh.
func IsValidSlug(slug string) bool {
	return IsBuiltin(slug) || slugRe.MatchString(slug)
}

// ResolveActive mengembalikan slug aktif valid (atau DefaultSlug bila invalid/kosong).
func ResolveActive(slug string) string {
	if IsValidSlug(slug) {
		return slug
	}
	return DefaultSlug
}

// Derive menyusun metadata tampil dari slug opentailwind; bila tak cocok pola,
// pakai slug apa adanya (kategori "Other").
func Derive(slug string) Template {
	m := slugRe.FindStringSubmatch(slug)
	if m == nil {
		return Template{Slug: slug, Name: titleize(slug), Category: "Other"}
	}
	return Template{Slug: slug, Name: titleize(m[3]), Category: titleize(m[1])}
}

// titleize: "digital-marketing" → "Digital Marketing".
func titleize(s string) string {
	parts := strings.Split(s, "-")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p == "" {
			continue
		}
		out = append(out, strings.ToUpper(p[:1])+p[1:])
	}
	return strings.Join(out, " ")
}
