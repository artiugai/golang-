package data

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"greenlight.alexedwards.net/internal/validator"
	"time"
)

var (
	ErrDuplicateEmail = errors.New("duplicate email")
)

// User struct represents an individual user. Notice the use of the json:"-" struct tag
// to prevent the Password and Version fields from appearing in JSON output.
// The Password field uses the custom password type defined below.
type User struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  password  `json:"-"`
	Activated bool      `json:"activated"`
	Version   int       `json:"-"`
}

// password is a custom type with plaintext and hashed password fields.
// The plaintext field is a pointer to a string, to distinguish between an absent password
// and an empty string password.
type password struct {
	plaintext *string
	hash      []byte
}

// Set calculates the bcrypt hash of a plaintext password, and stores both the hash
// and the plaintext versions in the struct.
func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}

	p.plaintext = &plaintextPassword
	p.hash = hash
	return nil
}

// Matches checks whether the provided plaintext password matches the hashed password
// stored in the struct, returning true if it matches and false otherwise.
func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}
func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}

func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.Name != "", "name", "must be provided")
	v.Check(len(user.Name) <= 500, "name", "must not be more than 500 bytes long")

	// Call the standalone ValidateEmail() helper.
	ValidateEmail(v, user.Email)

	// If the plaintext password is not nil, call the standalone
	// ValidatePasswordPlaintext() helper.
	if user.Password.plaintext != nil {
		ValidatePasswordPlaintext(v, *user.Password.plaintext)
	}

	// If the password hash is ever nil, this will be due to a logic error in our
	// codebase (probably because we forgot to set a password for the user). It's a
	// useful sanity check to include here, but it's not a problem with the data
	// provided by the client. So rather than adding an error to the validation map,
	// we raise a panic instead.
	if user.Password.hash == nil {
		panic("missing password hash for user")
	}
}

type UserModel struct {
	DB *sql.DB
}

// Insert adds a new record to the users table in the database. The id, created_at, and version
// fields are auto-generated, so the RETURNING clause is used to read them into the User struct.
func (m UserModel) Insert(user *User) error {
	query := `
    INSERT INTO users (name, email, password_hash, activated) VALUES ($1, $2, $3, $4)
    RETURNING id, created_at, version`

	args := []interface{}{user.Name, user.Email, user.Password.hash, user.Activated}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Check for violation of the UNIQUE "users_email_key" constraint.
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}

	return nil
}

// GetByEmail retrieves a user by their email address. It returns ErrRecordNotFound if no
// matching record is found due to the UNIQUE constraint on the email column.
func (m UserModel) GetByEmail(email string) (*User, error) {
	query := `
    SELECT id, created_at, name, email, password_hash, activated, version FROM users
    WHERE email = $1`

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.CreatedAt, &user.Name, &user.Email, &user.Password.hash, &user.Activated, &user.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

// Update modifies the details of a specific user. It includes checks against the version
// field to prevent race conditions and a check for the "users_email_key" constraint.
func (m UserModel) Update(user *User) error {
	query := `
    UPDATE users
    SET name = $1, email = $2, password_hash = $3, activated = $4, version = version + 1
    WHERE id = $5 AND version = $6
    RETURNING version`

	args := []interface{}{user.Name, user.Email, user.Password.hash, user.Activated, user.ID, user.Version}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}
func (m UserModel) GetForToken(tokenScope, tokenPlaintext string) (*User, error) {
	// Calculate the SHA-256 hash of the plaintext token provided by the client.
	// Remember that this returns a byte *array* with length 32, not a slice.
	tokenHash := sha256.Sum256([]byte(tokenPlaintext))

	// Set up the SQL query.
	query := `
	SELECT
		users.id,
		users.created_at,
		users.name,
		users.email,
		users.password_hash,
		users.activated,
		users.version
	FROM
		users
	INNER JOIN
		tokens
	ON
		users.id = tokens.user_id
	WHERE
		tokens.hash = $1
		AND tokens.scope = $2
		AND tokens.expiry > $3
	`

	// Create a slice containing the query arguments. Notice how we use the [:] operator
	// to get a slice containing the token hash, rather than passing in the array (which
	// is not supported by the pq driver), and that we pass the current time as the
	// value to check against the token expiry.
	args := []interface{}{tokenHash[:], tokenScope, time.Now()}

	var user User
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the query, scanning the return values into a User struct. If no matching
	// record is found we return an ErrRecordNotFound error.
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Name,
		&user.Email,
		&user.Password.hash,
		&user.Activated,
		&user.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	// Return the matching user.
	return &user, nil
}
