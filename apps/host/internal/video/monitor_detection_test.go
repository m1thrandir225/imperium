package video

import (
	"testing"
)

func TestGetPrimaryMonitorInfo(t *testing.T) {
	monitor, err := GetPrimaryMonitorInfo()
	if err != nil {
		t.Fatalf("GetPrimaryMonitorInfo() failed: %v", err)
	}

	if monitor == nil {
		t.Fatal("GetPrimaryMonitorInfo() returned nil monitor")
	}

	if monitor.Width <= 0 {
		t.Errorf("Expected positive width, got %d", monitor.Width)
	}

	if monitor.Height <= 0 {
		t.Errorf("Expected positive height, got %d", monitor.Height)
	}

	if !monitor.IsPrimary {
		t.Error("Expected primary monitor to have IsPrimary=true")
	}

	t.Logf("Primary monitor: %dx%d at (%d,%d)",
		monitor.Width, monitor.Height, monitor.OffsetX, monitor.OffsetY)
}

func TestGetMonitorCount(t *testing.T) {
	count, err := GetMonitorCount()
	if err != nil {
		t.Fatalf("GetMonitorCount() failed: %v", err)
	}

	if count <= 0 {
		t.Errorf("Expected at least 1 monitor, got %d", count)
	}

	t.Logf("Monitor count: %d", count)
}

func TestGetAllMonitorsInfo(t *testing.T) {
	monitors, err := GetAllMonitorsInfo()
	if err != nil {
		t.Fatalf("GetAllMonitorsInfo() failed: %v", err)
	}

	if len(monitors) == 0 {
		t.Fatal("GetAllMonitorsInfo() returned no monitors")
	}

	// Check that at least one monitor is marked as primary
	foundPrimary := false
	for i, monitor := range monitors {
		if monitor == nil {
			t.Errorf("Monitor %d is nil", i)
			continue
		}

		if monitor.Width <= 0 {
			t.Errorf("Monitor %d has invalid width: %d", i, monitor.Width)
		}

		if monitor.Height <= 0 {
			t.Errorf("Monitor %d has invalid height: %d", i, monitor.Height)
		}

		if monitor.IsPrimary {
			if foundPrimary {
				t.Error("Multiple monitors marked as primary")
			}
			foundPrimary = true
		}

		t.Logf("Monitor %d: %dx%d at (%d,%d), primary: %v",
			i, monitor.Width, monitor.Height, monitor.OffsetX, monitor.OffsetY, monitor.IsPrimary)
	}

	if !foundPrimary {
		t.Error("No monitor marked as primary")
	}
}

func TestCaching(t *testing.T) {
	cachedPrimaryMonitor = nil
	monitor1, err := GetPrimaryMonitorInfo()
	if err != nil {
		t.Fatalf("First GetPrimaryMonitorInfo() call failed: %v", err)
	}

	monitor2, err := GetPrimaryMonitorInfo()
	if err != nil {
		t.Fatalf("Second GetPrimaryMonitorInfo() call failed: %v", err)
	}

	if monitor1 != monitor2 {
		t.Error("Expected cached monitor to be the same instance")
	}

	if cachedPrimaryMonitor == nil {
		t.Error("Expected cache to be populated")
	}

	if lastPrimaryFetch.IsZero() {
		t.Error("Expected lastPrimaryFetch to be set")
	}
}

func TestMonitorInfoConsistency(t *testing.T) {
	primaryMonitor, err := GetPrimaryMonitorInfo()
	if err != nil {
		t.Fatalf("GetPrimaryMonitorInfo() failed: %v", err)
	}

	allMonitors, err := GetAllMonitorsInfo()
	if err != nil {
		t.Fatalf("GetAllMonitorsInfo() failed: %v", err)
	}

	var foundPrimary *MonitorInfo
	for _, monitor := range allMonitors {
		if monitor.IsPrimary {
			foundPrimary = monitor
			break
		}
	}

	if foundPrimary == nil {
		t.Fatal("No primary monitor found in GetAllMonitorsInfo() result")
	}

	if primaryMonitor.Width != foundPrimary.Width {
		t.Errorf("Width mismatch: GetPrimaryMonitorInfo()=%d, GetAllMonitorsInfo()=%d",
			primaryMonitor.Width, foundPrimary.Width)
	}

	if primaryMonitor.Height != foundPrimary.Height {
		t.Errorf("Height mismatch: GetPrimaryMonitorInfo()=%d, GetAllMonitorsInfo()=%d",
			primaryMonitor.Height, foundPrimary.Height)
	}
}

func BenchmarkGetPrimaryMonitorInfo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := GetPrimaryMonitorInfo()
		if err != nil {
			b.Fatalf("GetPrimaryMonitorInfo() failed: %v", err)
		}
	}
}

func BenchmarkGetMonitorCount(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := GetMonitorCount()
		if err != nil {
			b.Fatalf("GetMonitorCount() failed: %v", err)
		}
	}
}

func BenchmarkGetAllMonitorsInfo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := GetAllMonitorsInfo()
		if err != nil {
			b.Fatalf("GetAllMonitorsInfo() failed: %v", err)
		}
	}
}
