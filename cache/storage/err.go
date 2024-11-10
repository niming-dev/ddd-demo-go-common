package storage

const (
	NotFound      = Error("not found")
	TypeMismatch  = Error("type mismatch")
	DestUnsetable = Error("dest unsetable")
	DestCantBeNil = Error("dest can't be nil")
	DestMustBePtr = Error("dest must be ptr")
)

type Error string

func (e Error) Error() string {
	return string(e)
}
