package repositories

import (
	"time"
	"user-management/config"
)

// VerifyOTP checks if the OTP is valid and not expired for the given email.
func VerifyOTP(email, otp string) (bool, error) {
	var count int
	err := config.DB.QueryRow(
		`SELECT COUNT(*) FROM otp_verifications
		 WHERE email = @p1 AND otp = @p2 AND expires_at > @p3`,
		email, otp, time.Now(),
	).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// DeleteOTP removes all OTP records for the given email after successful verification.
func DeleteOTP(email string) error {
	_, err := config.DB.Exec(
		"DELETE FROM otp_verifications WHERE email = @p1", email,
	)
	return err
}

// UpdatePasswordByEmail updates the user's password for the given email.
func UpdatePasswordByEmail(email, newPassword string) error {
	_, err := config.DB.Exec(
		"UPDATE users SET password = @p1 WHERE email = @p2",
		newPassword, email,
	)
	return err
}