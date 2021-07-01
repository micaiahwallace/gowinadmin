package gowinadmin

// Indicates that a session was not found
type SessionNotExistError string

func (e SessionNotExistError) Error() string {
	return string(e)
}

// No session was found
const ErrSessionNotExist SessionNotExistError = "no session found"

// Indicates that the host did not respond
type NoResponseError string

func (e NoResponseError) Error() string {
	return string(e)
}

// Host did not send a response
const ErrNoResponse SessionNotExistError = "no response from host"
