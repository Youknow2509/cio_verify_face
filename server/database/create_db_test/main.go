package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gocql/gocql"
	"github.com/jackc/pgx/v5/pgxpool"
)

// must panics on error
func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	ctx := context.Background()

	// TODO: thay bằng DSN thực tế
	pgpool, err := pgxpool.New(ctx, "postgres://postgres:root1234@localhost:5433/cio_verify_face?sslmode=disable")
	must(err)
	defer pgpool.Close()

	// TODO: thay bằng thông số Scylla thực tế
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "cio_verify_face"
	cluster.Consistency = gocql.Quorum
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: "cassandra",
		Password: "root1234",
	}
	sess, err := cluster.CreateSession()
	must(err)
	defer sess.Close()

	// 1) Lấy company_id FPT
	var companyID string
	err = pgpool.QueryRow(ctx, `SELECT company_id FROM companies WHERE name = 'FPT Software'`).Scan(&companyID)
	must(err)

	// 2) Đảm bảo có các shift (Morning, Afternoon, Flexible)
	ensureShifts(ctx, pgpool, companyID)

	// 3) Lấy user_id của hai nhân viên
	employees := []string{"employee1.fpt@example.com", "employee2.fpt@example.com"}
	var empIDs []string
	for _, email := range employees {
		var uid string
		err := pgpool.QueryRow(ctx, `SELECT user_id FROM users WHERE email = $1`, email).Scan(&uid)
		must(err)
		empIDs = append(empIDs, uid)
	}

	// 4) Gán ca: dùng "Morning Shift" + "Afternoon Shift" cho cả hai, và "Flexible Engineering" cho employee1
	morningID, afternoonID, flexID := assignShifts(ctx, pgpool, companyID, empIDs)
	_ = afternoonID // hiện tại chỉ dùng morning + flexible cho seed demo
	shiftByEmployee := map[string]string{}
	for i, emp := range empIDs {
		shiftByEmployee[emp] = morningID
		if i == 0 {
			shiftByEmployee[emp] = flexID
		}
	}

	// 5) Sinh chấm công từ 2025-10-01 đến 2025-12-08 (bỏ T7, CN)
	start := time.Date(2025, 10, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2025, 12, 8, 23, 59, 59, 0, time.UTC)

	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		if d.Weekday() == time.Saturday || d.Weekday() == time.Sunday {
			continue
		}
		for _, emp := range empIDs {
			// Demo: dùng Morning + Afternoon, có dao động phút nhỏ
			checkIn := time.Date(d.Year(), d.Month(), d.Day(), 8, 0, 0, 0, time.UTC).Add(time.Duration(randMinute(0, 8)) * time.Minute)
			checkOut := time.Date(d.Year(), d.Month(), d.Day(), 17, 30, 0, 0, time.UTC).Add(time.Duration(randMinute(-10, 5)) * time.Minute)

			writeAttendance(ctx, sess, companyID, emp, checkIn, 0)  // record_type 0 = check-in
			writeAttendance(ctx, sess, companyID, emp, checkOut, 1) // record_type 1 = check-out
			writeDailySummary(ctx, sess, companyID, emp, shiftByEmployee[emp], d, checkIn, checkOut)
		}
	}

	fmt.Println("Done seeding shifts and attendance.")
}

// ensureShifts tạo nếu chưa có các ca chuẩn cho FPT
func ensureShifts(ctx context.Context, db *pgxpool.Pool, companyID string) {
	type shiftDef struct {
		Name     string
		Start    string
		End      string
		WorkDays []int
		Flexible bool
	}
	shifts := []shiftDef{
		{"Morning Shift", "08:00:00", "12:00:00", []int{1, 2, 3, 4, 5}, false},
		{"Afternoon Shift", "13:30:00", "17:30:00", []int{1, 2, 3, 4, 5}, false},
		{"Flexible Engineering", "09:00:00", "18:00:00", []int{1, 2, 3, 4, 5}, true},
	}

	for _, s := range shifts {
		_, err := db.Exec(ctx, `
			INSERT INTO work_shifts (company_id, name, start_time, end_time, work_days, is_flexible)
			VALUES ($1, $2, $3, $4, $5, $6)
			ON CONFLICT DO NOTHING`,
			companyID, s.Name, s.Start, s.End, s.WorkDays, s.Flexible)
		must(err)
	}
}

// assignShifts gán Morning + Afternoon cho cả hai; thêm Flexible cho employee1
func assignShifts(ctx context.Context, db *pgxpool.Pool, companyID string, empIDs []string) (morningID, afternoonID, flexID string) {
	// Lấy shift_id theo tên
	must(db.QueryRow(ctx, `SELECT shift_id FROM work_shifts WHERE company_id=$1 AND name='Morning Shift'`, companyID).Scan(&morningID))
	must(db.QueryRow(ctx, `SELECT shift_id FROM work_shifts WHERE company_id=$1 AND name='Afternoon Shift'`, companyID).Scan(&afternoonID))
	must(db.QueryRow(ctx, `SELECT shift_id FROM work_shifts WHERE company_id=$1 AND name='Flexible Engineering'`, companyID).Scan(&flexID))

	effective := "2025-10-01"

	for i, emp := range empIDs {
		// Morning
		_, err := db.Exec(ctx, `
			INSERT INTO employee_shifts (employee_id, shift_id, effective_from)
			VALUES ($1, $2, $3)
			ON CONFLICT DO NOTHING`,
			emp, morningID, effective)
		must(err)
		// Afternoon
		_, err = db.Exec(ctx, `
			INSERT INTO employee_shifts (employee_id, shift_id, effective_from)
			VALUES ($1, $2, $3)
			ON CONFLICT DO NOTHING`,
			emp, afternoonID, effective)
		must(err)
		// Flexible chỉ cho employee1 (index 0)
		if i == 0 {
			_, err = db.Exec(ctx, `
				INSERT INTO employee_shifts (employee_id, shift_id, effective_from)
				VALUES ($1, $2, $3)
				ON CONFLICT DO NOTHING`,
				emp, flexID, effective)
			must(err)
		}
	}

	return
}

// writeAttendance ghi 1 bản ghi vào 2 bảng Scylla: attendance_records & attendance_records_by_user
func writeAttendance(ctx context.Context, sess *gocql.Session, companyID, employeeID string, ts time.Time, recordType int) {
	ym := ts.Format("2006-01")
	deviceID := gocql.UUID{} // để trống demo; thay bằng device thực tế nếu cần
	score := float32(0.97)   // Cassandra float là 32-bit
	method := "face"
	status := "synced"

	// attendance_records
	if err := sess.Query(`
		INSERT INTO attendance_records (
			company_id, year_month, record_time, employee_id, device_id,
			record_type, verification_method, verification_score, face_image_url,
			location_coordinates, metadata, sync_status, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		parseUUID(companyID), ym, ts, parseUUID(employeeID), deviceID,
		recordType, method, score, "", "", map[string]string{}, status, time.Now(),
	).WithContext(ctx).Exec(); err != nil {
		must(err)
	}

	// attendance_records_by_user
	if err := sess.Query(`
		INSERT INTO attendance_records_by_user (
			company_id, employee_id, year_month, record_time, device_id,
			record_type, verification_method, verification_score, face_image_url,
			location_coordinates, metadata, sync_status, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		parseUUID(companyID), parseUUID(employeeID), ym, ts, deviceID,
		recordType, method, score, "", "", map[string]string{}, status, time.Now(),
	).WithContext(ctx).Exec(); err != nil {
		must(err)
	}
}

// writeDailySummary ghi tổng hợp ngày cho cả bảng company và user
func writeDailySummary(ctx context.Context, sess *gocql.Session, companyID, employeeID, shiftID string, workDate time.Time, checkIn, checkOut time.Time) {
	summaryMonth := workDate.Format("2006-01")
	workDate = time.Date(workDate.Year(), workDate.Month(), workDate.Day(), 0, 0, 0, 0, time.UTC)

	expectedStart := time.Date(workDate.Year(), workDate.Month(), workDate.Day(), 8, 0, 0, 0, time.UTC)
	expectedEnd := time.Date(workDate.Year(), workDate.Month(), workDate.Day(), 17, 30, 0, 0, time.UTC)

	lateMinutes := 0
	if checkIn.After(expectedStart) {
		lateMinutes = int(checkIn.Sub(expectedStart).Minutes())
	}
	earlyLeaveMinutes := 0
	if checkOut.Before(expectedEnd) {
		earlyLeaveMinutes = int(expectedEnd.Sub(checkOut).Minutes())
	}
	totalWorkMinutes := int(checkOut.Sub(checkIn).Minutes())
	attendanceStatus := 1 // 1 = đi làm bình thường
	updatedAt := time.Now()

	shiftUUID := parseUUID(shiftID)
	companyUUID := parseUUID(companyID)
	employeeUUID := parseUUID(employeeID)

	// daily_summaries (by company)
	if err := sess.Query(`
		INSERT INTO daily_summaries (
			company_id, summary_month, work_date, employee_id, shift_id,
			actual_check_in, actual_check_out, attendance_status,
			late_minutes, early_leave_minutes, total_work_minutes, notes, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		companyUUID, summaryMonth, workDate, employeeUUID, shiftUUID,
		checkIn, checkOut, attendanceStatus, lateMinutes, earlyLeaveMinutes, totalWorkMinutes, "", updatedAt,
	).WithContext(ctx).Exec(); err != nil {
		must(err)
	}

	// daily_summaries_by_user
	if err := sess.Query(`
		INSERT INTO daily_summaries_by_user (
			company_id, employee_id, summary_month, work_date, shift_id,
			actual_check_in, actual_check_out, attendance_status,
			late_minutes, early_leave_minutes, total_work_minutes, notes, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		companyUUID, employeeUUID, summaryMonth, workDate, shiftUUID,
		checkIn, checkOut, attendanceStatus, lateMinutes, earlyLeaveMinutes, totalWorkMinutes, "", updatedAt,
	).WithContext(ctx).Exec(); err != nil {
		must(err)
	}
}

// parseUUID helper: panic nếu sai
func parseUUID(s string) gocql.UUID {
	u, err := gocql.ParseUUID(s)
	must(err)
	return u
}

// randMinute trả về phút ngẫu nhiên trong [min,max]
func randMinute(min, max int) int {
	if max == min {
		return min
	}
	return min + int(time.Now().UnixNano()%int64(max-min+1))
}
