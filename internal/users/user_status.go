package users

type UserStatus string

const (
	// user registered and confirmed email
	StatusActive UserStatus = "Active"

	// user registered but hasn't confirmed email
	StatusPending UserStatus = "Pending"

	// user was blocked due to rules violation
	StatusBlocked UserStatus = "Blocked"
)
