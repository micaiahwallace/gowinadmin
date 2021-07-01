package gowinadmin

import (
	"fmt"
	"regexp"
	"strings"
)

type QuserSession struct {
	UserName    string
	SessionName string
	Id          string
	State       string
	IdleTime    string
	LogonTime   string
}

type QuserRequest struct {
	Session string
	Server  string
}

// Take the cli arguments and execute command returning stdout and stderr strings
func (r *QuserRequest) RunQuser() ([]*QuserSession, error) {

	// Get arguments and run quser command
	args := r.getArgs()
	stdout, stderr, execErr := executeCommand("quser", args)
	if execErr != nil {
		return nil, execErr
	}

	// Check for any errors
	if parsedError := ParseQuserError(stderr); parsedError != nil {
		return nil, parsedError
	}

	return ParseQuserOutput(stdout)
}

// Get the request arguments with defaults applied to pass to quser
func (r *QuserRequest) getArgs() []string {

	// Set the return defaults
	server := "127.0.0.1"
	session := "*"

	// Set the specified remote host
	if strings.TrimSpace(r.Server) != "" {
		server = r.Server
	}

	// Set a specified session filter
	if strings.TrimSpace(r.Session) != "" {
		session = r.Session
	}

	// Return the ordered arguments to pass to quser
	return []string{session, fmt.Sprintf("/server:%s", server)}
}

// Parse the resulting stdout and stderr data into UserSessions
func ParseQuserOutput(stdout string) ([]*QuserSession, error) {

	// Store captured sessions
	sessions := make([]*QuserSession, 0)

	// Create regex to separate columns in session rows
	reg, regErr := regexp.Compile(`\s{2,}`)
	if regErr != nil {
		return nil, regErr
	}

	// Get the response rows indicating each session
	sessionRows := strings.Split(strings.TrimSpace(stdout), "\n")

	// Parse each session, skipping the header row
	for _, sessionRow := range sessionRows[1:] {

		// Parse each column of tab separated row
		csvRow := reg.ReplaceAllLiteralString(sessionRow, ",")
		sessionData := strings.Split(csvRow, ",")

		// Session name is missing, shift elements to add a blank space for session name
		if len(sessionData) == 5 {
			sessionData = append(sessionData, "")
			copy(sessionData[2:], sessionData[1:])
			sessionData[1] = ""
		}

		// Append the parsed data to the resulting session list
		sessions = append(sessions, &QuserSession{
			UserName:    sessionData[0],
			SessionName: sessionData[1],
			Id:          sessionData[2],
			State:       sessionData[3],
			IdleTime:    sessionData[4],
			LogonTime:   sessionData[5],
		})
	}
	return sessions, nil
}

// Parse a quser error message
func ParseQuserError(stderr string) error {

	// Check for empty errors
	if len(strings.TrimSpace(stderr)) == 0 {
		return nil
	}

	// Check for a response where no users came back
	if strings.Contains(stderr, "No User exists") {
		return ErrSessionNotExist
	}

	// Otherwise it is likely that the host didn't respond to the request
	if strings.Contains(stderr, "The RPC server is unavailable") {
		return ErrNoResponse
	}

	// Unknown case
	return fmt.Errorf("unknown error occurred: %s", stderr)
}
