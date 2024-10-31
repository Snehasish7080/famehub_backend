package user

import (
	"context"
	"errors"
	"time"

	"github.com/gocql/gocql"
	"github.com/snehasish7080/famehub/pkg/hash"
	"github.com/snehasish7080/famehub/pkg/jwtclaim"
	"github.com/snehasish7080/famehub/pkg/otp"
)

type UserStorage struct {
	session *gocql.Session
}

func NewUserStorage(session *gocql.Session) *UserStorage {
	return &UserStorage{
		session: session,
	}
}

func (u *UserStorage) createUserTables() error {
	userTable := `CREATE TABLE IF NOT EXISTS users (
        uuid UUID,
        email TEXT,
        password TEXT,
        otp TEXT,
        created_at TIMESTAMP,
        updated_at TIMESTAMP,
        PRIMARY KEY (uuid, email)
    );`

	usersByUsernameTable := `CREATE TABLE IF NOT EXISTS users_by_username (
        uuid UUID,
        email TEXT,
        username TEXT,
        fullName TEXT,
        bio TEXT,
        isBrand BOOLEAN,
        created_at TIMESTAMP,
        updated_at TIMESTAMP,
        PRIMARY KEY (uuid, username, fullName, email)
    );`

	if err := u.session.Query(userTable).Exec(); err != nil {
		return err
	}

	if err := u.session.Query(usersByUsernameTable).Exec(); err != nil {
		return err
	}

	return nil
}

func (u *UserStorage) signUp(email string, password string, ctx context.Context) (string, error) {

	// Create User Table if don't exists.
	if err := u.createUserTables(); err != nil {
		return "", err
	}

	// Check if email already exists
	exists, err := u.emailExists(email)
	if err != nil {
		return "", err
	}

	if exists {
		return "", errors.New("Email already exists")
	}

	// Generate Otp
	generatedOtp := otp.EncodeToString(6)
	// Generate a new UUID for the user

	UUID := gocql.TimeUUID().String()

	// Generate timestamp for the user
	currentTime := time.Now()
	createdAt := currentTime
	updatedAt := currentTime

	// Generate hashed password for the user
	hashedPassword, err := hash.HashPassword(password)
	if err != nil {
		return "", err
	}

	// Insert the new user into the main users table
	query := `INSERT INTO users (uuid, email, password, otp, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`
	if err := u.session.Query(query, UUID, email, hashedPassword, generatedOtp, createdAt, updatedAt).Exec(); err != nil {
		return "", err
	}

	query = `INSERT INTO users_by_username (uuid, email, username, fullName, bio, isBrand, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	if err := u.session.Query(query, UUID, email, "", "", "", false, createdAt, updatedAt).Exec(); err != nil {
		return "", err
	}

	verifyToken, err := jwtclaim.CreateJwtToken(UUID, false)
	if err != nil {
		return "", err
	}

	return verifyToken, nil
}

func (u *UserStorage) emailExists(email string) (bool, error) {
	var exists bool
	query := "SELECT COUNT(*) FROM users WHERE email = ?"

	var count int
	if err := u.session.Query(query, email).Scan(&count); err != nil {
		return false, err
	}
	// If the count is greater than 0, the email exists
	exists = count > 0
	return exists, nil
}

func (u *UserStorage) verifyOtp(userId string, otp string, ctx context.Context) (string, error) {

	var storedOTP string

	// getting otp from db
	query := "SELECT otp FROM users WHERE uuid = ?"
	err := u.session.Query(query, userId).Scan(&storedOTP)

	if err != nil {
		return "", err
	}

	if storedOTP != otp {
		return "", errors.New("Invalid OTP")
	}

	// generate verified token
	verifyToken, err := jwtclaim.CreateJwtToken(userId, true)

	if err != nil {
		return "", err
	}

	return verifyToken, nil

}

type User struct {
	UUID     string `json:"uuid"`
	Password string `json:"password"`
}

func (u *UserStorage) login(email string, password string, ctx context.Context) (string, error) {

	// Check if email exists or not
	exists, err := u.emailExists(email)
	if err != nil {
		return "", err
	}

	if !exists {
		return "", errors.New("Email does not exists. Please signup.")
	}

	// getting user id and password
	var user User
	query := "SELECT uuid, password FROM users WHERE email = ?"
	err = u.session.Query(query, email).Scan(&user.UUID, &user.Password)

	if err != nil {
		return "", err
	}

	// checking password
	if !hash.CheckPasswordHash(password, user.Password) {
		return "", errors.New("incorrect email or password")
	}

	// updating the otp
	generatedOtp := otp.EncodeToString(6)

	// Generate timestamp for the user
	currentTime := time.Now()
	updatedAt := currentTime
	query = "UPDATE users SET otp = ?, updated_at = ? WHERE uuid = ? AND email = ?"
	err = u.session.Query(query, generatedOtp, updatedAt, user.UUID, email).Exec()
	if err != nil {
		return "", err
	}

	// generating jwt for otp
	verifyToken, err := jwtclaim.CreateJwtToken(user.UUID, false)
	if err != nil {
		return "", err
	}

	return verifyToken, nil

}
