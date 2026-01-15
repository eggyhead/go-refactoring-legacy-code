package report_fixed

import (
	"fmt"
	"time"
)

// --- 1. Interface Substitution ---
// We define what we NEED, not what we HAVE.
type RevenueReader interface {
	GetMonthlyRevenue(month time.Month) (float64, error)
}

type Uploader interface {
	Upload(key string, data []byte) error
}

// --- 2. Function Field Type ---
// We define a type for the email function to make the struct cleaner.
type EmailSender func(to string, body string) error

type ReportManager struct {
	reader   RevenueReader
	uploader Uploader
	mailer   EmailSender
	now      func() time.Time // Function field to control time
}

// --- 3. Parameterized Constructor ---
// We "inject" all dependencies. The struct no longer creates its own tools.
func NewReportManager(r RevenueReader, u Uploader, m EmailSender, n func() time.Time) *ReportManager {
	return &ReportManager{
		reader:   r,
		uploader: u,
		mailer:   m,
		now:      n,
	}
}

// --- 4. Primitivized Logic ---
// This is a pure function. It doesn't know about databases or S3.
func IsReportDay(t time.Time) bool {
	return t.Day() == 1
}

func (rm *ReportManager) GenerateReport(adminEmail string) error {
	// Use the controlled "now" and the pure logic check
	currentTime := rm.now()
	if !IsReportDay(currentTime) {
		return fmt.Errorf("reports can only be generated on the 1st")
	}

	// Use the interface instead of SQL
	total, err := rm.reader.GetMonthlyRevenue(currentTime.Month())
	if err != nil {
		return err
	}

	// Use the interface instead of AWS S3
	reportContent := []byte(fmt.Sprintf("Revenue: %f", total))
	err = rm.uploader.Upload("monthly-report.txt", reportContent)
	if err != nil {
		return err
	}

	// Use the function field instead of smtp.SendMail
	return rm.mailer(adminEmail, string(reportContent))
}
