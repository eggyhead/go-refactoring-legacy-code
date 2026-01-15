package report_fixed

import (
	"testing"
	"time"
)

// --- Fakes for Testing ---
type FakeDB struct{ revenue float64 }

func (f *FakeDB) GetMonthlyRevenue(m time.Month) (float64, error) { return f.revenue, nil }

type FakeS3 struct{}

func (f *FakeS3) Upload(k string, d []byte) error { return nil }

func TestGenerateReport_Success(t *testing.T) {
	// 1. We freeze time to the 1st of the month
	frozenTime := time.Date(2026, time.February, 1, 10, 0, 0, 0, time.UTC)

	// 2. We provide a fake email function that just returns nil
	fakeEmail := func(to string, body string) error { return nil }

	// 3. Assemble the manager with fakes
	mgr := NewReportManager(
		&FakeDB{revenue: 5000.0},
		&FakeS3{},
		fakeEmail,
		func() time.Time { return frozenTime },
	)

	// 4. Run the test
	err := mgr.GenerateReport("admin@test.com")
	if err != nil {
		t.Errorf("Expected success, got error: %v", err)
	}
}

func TestIsReportDay(t *testing.T) {
	// Testing the primitivized logic is now trivial!
	if !IsReportDay(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)) {
		t.Error("Expected Jan 1st to be a report day")
	}
	if IsReportDay(time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC)) {
		t.Error("Expected Jan 2nd NOT to be a report day")
	}
}
