package constant

type UserStatus string

const (
	UserStatusActive   UserStatus = "Active"
	UserStatusInactive UserStatus = "Inactive"
)

func (s UserStatus) String() string {
	return string(s)
}
