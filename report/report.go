package report

import (
	"database/sql"
	"fmt"
	"net/smtp"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

type ReportManager struct {
	db *sql.DB // Hard dependency on SQL
}

func NewReportManager(connString string) *ReportManager {
	db, _ := sql.Open("postgres", connString)
	return &ReportManager{db: db}
}

func (rm *ReportManager) GenerateReport(adminEmail string) error {
	// 1. Logic Blocker: Hard-coded time check
	// This makes it impossible to test on any day except the 1st.
	now := time.Now()
	if now.Day() != 1 {
		return fmt.Errorf("reports can only be generated on the 1st of the month")
	}

	// 2. Database Blocker: Hard-coded SQL call
	var totalRevenue float64
	err := rm.db.QueryRow("SELECT SUM(amount) FROM orders WHERE month = $1", now.Month()).Scan(&totalRevenue)
	if err != nil {
		return err
	}

	// 3. Side Effect Blocker: Hard-coded S3 Upload
	// You can't run this without AWS credentials.
	s3Client := s3.New(nil) // Imagine real config here
	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String("my-reports-bucket"),
		Key:    aws.String("monthly-report.txt"),
		Body:   nil, // Imagine the report content here
	})
	if err != nil {
		return err
	}

	// 4. Global Blocker: Hard-coded SMTP call
	// This will actually try to send an email every time you run a test.
	auth := smtp.PlainAuth("", "user@example.com", "password", "smtp.example.com")
	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: Revenue\r\n\r\nTotal: %f", adminEmail, totalRevenue))
	return smtp.SendMail("smtp.example.com:587", auth, "sender@example.com", []string{adminEmail}, msg)
}
