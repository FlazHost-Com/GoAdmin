package integration

import (
	"bytes"
	"html/template"
	"strings"
	"testing"
)

// Memuat partials chrome + semua view modul dalam SATU set lalu merender tiap
// view — membuktikan cross-reference {{template "admin/head"}} resolve & semua
// view valid + render tanpa error (deteksi dini sebelum runtime).
func TestViews_AllRenderWithChrome(t *testing.T) {
	files := []string{
		"../../web/templates/partials/admin_chrome.html",
		"../../internal/modules/access/view/users.html",
		"../../internal/modules/access/view/auth_login.html",
		"../../internal/modules/access/view/auth_forgot.html",
		"../../internal/modules/access/view/auth_reset.html",
		"../../internal/modules/access/view/users_form.html",
		"../../internal/modules/access/view/roles_index.html",
		"../../internal/modules/access/view/roles_form.html",
		"../../internal/modules/access/view/permissions_index.html",
		"../../internal/modules/access/view/permissions_form.html",
		"../../internal/modules/dashboard/view/index.html",
		"../../internal/modules/setting/view/index.html",
		"../../internal/modules/profile/view/index.html",
		"../../internal/modules/home/view/default.html",
		"../../internal/modules/home/view/minimal.html",
		"../../internal/modules/home/view/appearance.html",
		"../../internal/modules/components/view/index.html",
	}
	set, err := template.New("").ParseFiles(files...)
	if err != nil {
		t.Fatalf("parse views: %v", err)
	}

	cases := []struct {
		name string
		data map[string]interface{}
		want string
	}{
		{"auth/login", map[string]interface{}{"title": "Masuk", "flash_error": "Salah"}, "Masuk GoAdmin"},
		{"auth/forgot", map[string]interface{}{"title": "Lupa Password", "_csrf": "tok"}, "Kirim OTP"},
		{"auth/reset", map[string]interface{}{"title": "Reset Password", "_csrf": "tok", "email": "a@b.com"}, "Kode OTP"},
		{"dashboard/index", map[string]interface{}{
			"title": "Dashboard", "active": "dashboard", "currentUser": nil,
			"stats": map[string]interface{}{"Users": 1, "Roles": 1, "Permissions": 12},
		}, "Pengguna"},
		{"access/users/index", map[string]interface{}{
			"title": "User", "active": "users", "currentUser": nil, "search": "",
			"users": []map[string]interface{}{{"Code": "U-1", "Name": "Budi", "Email": "b@x.com", "Status": "Active"}},
			"meta":  map[string]interface{}{"From": 1, "To": 1, "Total": 1, "CurrentPage": 1, "LastPage": 1},
		}, "Budi"},
		{"setting/index", map[string]interface{}{
			"title": "Pengaturan", "active": "setting", "currentUser": nil,
			"setting": map[string]interface{}{"Name": "Toko", "Initial": "T", "Email": "", "Phone": "", "Address": "", "Description": "", "Copyright": "", "Theme": "Blue"},
			"themes":  []map[string]interface{}{{"Name": "Blue"}, {"Name": "Green"}},
		}, "Toko"},
		{"profile/index", map[string]interface{}{
			"title": "Profil", "active": "profile", "currentUser": nil,
			"profile": map[string]interface{}{"Name": "Budi", "Email": "b@x.com", "Phone": "", "Timezone": "UTC"},
		}, "Konfirmasi Password"},
		{"home/default", map[string]interface{}{
			"landing": map[string]interface{}{"AppName": "Toko Saya", "Description": "Halo", "Logo": "", "Email": "", "Phone": "", "Address": "", "Copyright": "© X", "ThemeName": "Green", "Primary": "#16a34a", "Accent": "#15803d", "Template": "default"},
		}, "Toko Saya"},
		{"home/minimal", map[string]interface{}{
			"landing": map[string]interface{}{"AppName": "Toko Saya", "Description": "Halo", "Email": "", "Phone": "", "Copyright": "© X", "ThemeName": "Green", "Primary": "#16a34a", "Accent": "#15803d", "Template": "minimal"},
		}, "Toko Saya"},
		{"home/appearance", map[string]interface{}{
			"title": "Tampilan", "active": "appearance", "currentUser": nil, "_csrf": "tok",
			"search": "", "category": "", "activeSlug": "default", "page": 1, "lastPage": 1, "total": 2,
			"categories": []string{"Bawaan", "Agency"},
			"templates": []map[string]interface{}{
				{"Slug": "default", "Name": "Klasik", "Category": "Bawaan", "Builtin": true},
				{"Slug": "minimal", "Name": "Minimalis", "Category": "Bawaan", "Builtin": true},
			},
		}, "Minimalis"},
		{"users/form", map[string]interface{}{
			"title": "Tambah Pengguna", "active": "users", "currentUser": nil, "_csrf": "tok",
			"user": nil, "action": "/admin/v1/users", "selected": map[string]bool{},
			"roles": []map[string]interface{}{{"ID": "r1", "Name": "Administrator"}},
		}, "Administrator"},
		{"roles/index", map[string]interface{}{
			"title": "Role", "active": "roles", "currentUser": nil, "search": "", "_csrf": "tok",
			"roles": []map[string]interface{}{{"ID": "r1", "Name": "Editor", "Permissions": []interface{}{}}},
			"meta":  map[string]interface{}{"Total": 1, "CurrentPage": 1, "LastPage": 1},
		}, "Editor"},
		{"roles/form", map[string]interface{}{
			"title": "Tambah Role", "active": "roles", "currentUser": nil, "_csrf": "tok",
			"role": nil, "action": "/admin/v1/roles", "selected": map[string]bool{},
			"permissions": []map[string]interface{}{{"ID": "p1", "Name": "user.view"}},
		}, "user.view"},
		{"permissions/index", map[string]interface{}{
			"title": "Permission", "active": "permissions", "currentUser": nil, "search": "", "_csrf": "tok",
			"permissions": []map[string]interface{}{{"ID": "p1", "Name": "user.view"}},
			"meta":        map[string]interface{}{"Total": 1, "CurrentPage": 1, "LastPage": 1},
		}, "user.view"},
		{"permissions/form", map[string]interface{}{
			"title": "Tambah Permission", "active": "permissions", "currentUser": nil, "_csrf": "tok",
			"permission": nil, "action": "/admin/v1/permissions",
		}, "Nama Permission"},
	}

	for _, tc := range cases {
		var buf bytes.Buffer
		if err := set.ExecuteTemplate(&buf, tc.name, tc.data); err != nil {
			t.Fatalf("render %s: %v", tc.name, err)
		}
		if !strings.Contains(buf.String(), tc.want) {
			t.Fatalf("view %s tak memuat %q", tc.name, tc.want)
		}
	}
}
