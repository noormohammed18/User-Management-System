package repositories

import (
	"database/sql"
	"fmt"
	"time"
	"user-management/config"
	"user-management/models"
	"math/rand"
)

func GetUsers(page, limit int) ([]models.User, error) {
	offset := (page - 1) * limit

	query := `
		SELECT id, name, email, age, is_admin
		FROM users
		ORDER BY id
		OFFSET @p1 ROWS
		FETCH NEXT @p2 ROWS ONLY
	`

	rows, err := config.DB.Query(query, offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var u models.User

		err := rows.Scan(
			&u.ID,
			&u.Name,
			&u.Email,
			&u.Age,
			&u.IsAdmin,
		)

		if err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	return users, nil
}

func GetUserCount() (int, error) {
	var count int

	err := config.DB.QueryRow(
		"SELECT COUNT(*) FROM users",
	).Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil
}

func GetUserByID(id string) (*models.User, error) {
	row := config.DB.QueryRow(
		"SELECT id, name, email, age, is_admin FROM users WHERE id=@p1", id,
	)
	var u models.User
	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Age, &u.IsAdmin)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func GetUserByEmailAndPassword(email, password string) (*models.User, error) {
	row := config.DB.QueryRow(
		"SELECT id, name, email, age, is_admin FROM users WHERE email=@p1 AND password=@p2",
		email, password,
	)
	var u models.User
	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Age, &u.IsAdmin)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func CreateUser(name, email, password, age string) error {
	query := `INSERT INTO users(name, email, password, age, is_admin) VALUES(@p1, @p2, @p3, @p4, 0)`
	_, err := config.DB.Exec(query, name, email, password, age)
	return err
}

func UpdateUser(id, name, email, age string) error {
	query := `UPDATE users SET name=@p1, email=@p2, age=@p3 WHERE id=@p4`
	_, err := config.DB.Exec(query, name, email, age, id)
	return err
}

func DeleteUser(id string) error {
	_, err := config.DB.Exec("DELETE FROM users WHERE id=@p1", id)
	return err
}
func ToggleAdmin(id, value string) error {
	_, err := config.DB.Exec(
		"UPDATE users SET is_admin = @p1 WHERE id = @p2", value, id,
	)
	return err
}

func FilterUsers(name, email, age, role string) ([]models.User, error) {
	query := `SELECT id, name, email, age, is_admin FROM users WHERE 1=1`
	args := []interface{}{}
	i := 1

	if name != "" {
		query += fmt.Sprintf(" AND name LIKE @p%d", i)
		args = append(args, "%"+name+"%")
		i++
	}
	if email != "" {
		query += fmt.Sprintf(" AND email LIKE @p%d", i)
		args = append(args, "%"+email+"%")
		i++
	}
	if age != "" {
		query += fmt.Sprintf(" AND age = @p%d", i)
		args = append(args, age)
		i++
	}
	if role == "admin" {
		query += fmt.Sprintf(" AND is_admin = @p%d", i)
		args = append(args, 1)
		i++
	} else if role == "user" {
		query += fmt.Sprintf(" AND is_admin = @p%d", i)
		args = append(args, 0)
		i++
	}

	rows, err := config.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		rows.Scan(&u.ID, &u.Name, &u.Email, &u.Age, &u.IsAdmin)
		users = append(users, u)
	}
	return users, nil
}

func GenerateOTP() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func SaveOTP(email, otp string) error {

	_, err := config.DB.Exec(
		`
		INSERT INTO otp_verifications
		(email, otp, expires_at)
		VALUES(@p1,@p2,@p3)
		`,
		email,
		otp,
		time.Now().Add(5*time.Minute),
	)

	return err
}