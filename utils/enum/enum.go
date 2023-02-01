package enum

type Admin string

const (
	Superuser Admin = "SUPERUSER"
	Verifier  Admin = "VERIFIER"
	Analytics Admin = "ANALYTICS"
)