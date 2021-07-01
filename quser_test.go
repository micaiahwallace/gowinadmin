package gowinadmin_test

import (
	"testing"

	"github.com/micaiahwallace/gowinadmin"
	"github.com/stretchr/testify/assert"
)

const QuserResponseFull = `
USERNAME              SESSIONNAME        ID  STATE   IDLE TIME  LOGON TIME
user1                 console             2  Active       4:33  6/28/2021 8:35 AM
`

const QuserResponseNoSessionName = `
USERNAME              SESSIONNAME        ID  STATE   IDLE TIME  LOGON TIME
testuser2                                 3  Disc      3+02:14  6/28/2021 10:00 AM
`

func TestParseQuserOutput(t *testing.T) {
	testCases := []struct {
		name           string
		inputData      string
		expectedResult []*gowinadmin.QuserSession
	}{
		{
			name:      "ParseQuserOutput parses sessions when all columns are present",
			inputData: QuserResponseFull,
			expectedResult: []*gowinadmin.QuserSession{
				{
					UserName:    "user1",
					SessionName: "console",
					Id:          "2",
					State:       "Active",
					IdleTime:    "4:33",
					LogonTime:   "6/28/2021 8:35 AM",
				},
			},
		},
		{
			name:      "ParseQuserOutput parses sessions when session name is missing",
			inputData: QuserResponseNoSessionName,
			expectedResult: []*gowinadmin.QuserSession{
				{
					UserName:    "testuser2",
					SessionName: "",
					Id:          "3",
					State:       "Disc",
					IdleTime:    "3+02:14",
					LogonTime:   "6/28/2021 10:00 AM",
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			actual, err := gowinadmin.ParseQuserOutput(testCase.inputData)
			assert.Nil(t, err)
			assert.Len(t, actual, len(testCase.expectedResult), "number of resulting sessions should equal expected session count")

			for i, expected := range testCase.expectedResult {
				assert.Equal(t, expected.Id, actual[i].Id)
				assert.Equal(t, expected.IdleTime, actual[i].IdleTime)
				assert.Equal(t, expected.LogonTime, actual[i].LogonTime)
				assert.Equal(t, expected.SessionName, actual[i].SessionName)
				assert.Equal(t, expected.State, actual[i].State)
				assert.Equal(t, expected.UserName, actual[i].UserName)
			}
		})
	}
}
