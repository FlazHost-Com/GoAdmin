package integration

import (
	"context"
	"testing"

	"goadmin/internal/config"
	accessmig "goadmin/internal/modules/access/migration"
	"goadmin/internal/modules/dashboard/service"
	"goadmin/tests/testutil"
)

// DB kosong (skema access ada, tanpa seed) → semua hitungan nol.
func TestDashboardService_StatsEmpty(t *testing.T) {
	c := testutil.NewContainer(t, config.ModeFull)
	svc := service.NewDashboardService(c.DB)

	st, err := svc.Stats(context.Background())
	if err != nil {
		t.Fatalf("stats: %v", err)
	}
	if st.Users != 0 || st.Roles != 0 || st.Permissions != 0 {
		t.Fatalf("harusnya nol, dapat: %+v", st)
	}
}

// Setelah seed admin (1 user, 1 role Administrator, semua permission inti).
func TestDashboardService_StatsAfterSeed(t *testing.T) {
	c := testutil.NewContainer(t, config.ModeFull)
	testutil.SeedAdmin(t, c)
	svc := service.NewDashboardService(c.DB)

	st, err := svc.Stats(context.Background())
	if err != nil {
		t.Fatalf("stats: %v", err)
	}
	if st.Users != 1 {
		t.Fatalf("users: harap 1, dapat %d", st.Users)
	}
	if st.Roles != 1 {
		t.Fatalf("roles: harap 1, dapat %d", st.Roles)
	}
	if want := int64(len(accessmig.CorePermissions)); st.Permissions != want {
		t.Fatalf("permissions: harap %d, dapat %d", want, st.Permissions)
	}
}
