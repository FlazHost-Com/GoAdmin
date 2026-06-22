package integration

import (
	"bytes"
	"html/template"
	"strings"
	"testing"
)

// View showcase Components = inti modul. Test memastikan template valid secara
// sintaks, ber-nama "components/index" (sesuai RenderView), dan render tanpa error.
func TestComponentsView_ParsesAndRenders(t *testing.T) {
	const path = "../../internal/modules/components/view/index.html"

	tmpl, err := template.New("").ParseFiles(path)
	if err != nil {
		t.Fatalf("parse template: %v", err)
	}
	if tmpl.Lookup("components/index") == nil {
		t.Fatal(`define name salah — harus {{define "components/index"}}`)
	}

	data := map[string]interface{}{
		"title":  "Komponen UI",
		"stats":  []map[string]interface{}{{"label": "Pengguna", "value": 128, "cls": "from-blue-500 to-blue-600"}},
		"badges": []map[string]interface{}{{"text": "Active", "cls": "bg-emerald-100 text-emerald-700"}},
		"alerts": []map[string]interface{}{{"msg": "ok", "cls": "x"}},
		"rows":   []map[string]interface{}{{"code": "U-001", "name": "Budi", "email": "b@x.com", "status": "Active"}},
	}

	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "components/index", data); err != nil {
		t.Fatalf("render: %v", err)
	}

	out := buf.String()
	for _, want := range []string{"Komponen UI", "Stat Card", "Pengguna", "Active", "U-001"} {
		if !strings.Contains(out, want) {
			t.Fatalf("output tak memuat %q", want)
		}
	}
}
