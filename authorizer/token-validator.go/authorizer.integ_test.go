package authorizer

import (
	"fmt"
	"os"
	"testing"

	"github.com/pulpfree/gsales-pdf-reports/config"
	"github.com/stretchr/testify/suite"
)

const (
	defaultsFP   = "../../config/defaults.yml"
	username     = "pulpfree"
	expiredToken = "eyJraWQiOiJmTVlINkJyRHB3T2ZaOUVsNkZkM0N1UCtsME1kTFpLd3gwaUk5VXhOU0RnPSIsImFsZyI6IlJTMjU2In0.eyJzdWIiOiI1ZjlhZTcxZS1iZmRiLTRhYzctOGVhNy1kMTU2ZTdmMmEzODAiLCJkZXZpY2Vfa2V5IjoidXMtZWFzdC0xX2I4YjBhZDYxLTdkMmItNDJkOS1hMTg2LTYyOTY0ODkyOTNhYSIsImV2ZW50X2lkIjoiYTNkMTVhMGEtZGRkYy0xMWU4LTk5NzktZDM5MGVkZmFiZDU0IiwidG9rZW5fdXNlIjoiYWNjZXNzIiwic2NvcGUiOiJhd3MuY29nbml0by5zaWduaW4udXNlci5hZG1pbiIsImF1dGhfdGltZSI6MTU0MTA4MDAzNCwiaXNzIjoiaHR0cHM6XC9cL2NvZ25pdG8taWRwLnVzLWVhc3QtMS5hbWF6b25hd3MuY29tXC91cy1lYXN0LTFfZ3NCNTl3ZnpXIiwiZXhwIjoxNTQxMDgzNjM0LCJpYXQiOjE1NDEwODAwMzQsImp0aSI6IjBjMDMyZjMxLWI0YjgtNDg1OS1hMjUwLWYwYzc2NGY2NGE4MyIsImNsaWVudF9pZCI6IjIwODR1a3Nsc2M4MzFwdDIwMnQyZHVkdDdjIiwidXNlcm5hbWUiOiJwdWxwZnJlZSJ9.M0LneEIY3fAPZng30Shn8Mo440O11XWmNeOfCIrO9lZwYSAa7B7X-norlj1ngV1PCIaCX6WunLNwZDy5Bw5Q8CxoZdr3pbnsaOtbmtwj1Zov1xpxHdusV7dutISPkDKCOlHO9kXS6P159Mfiw8bhzvRtj8tSar-Odk_jGubmaqowLCmw8fo0fYqK3ou1MAMbt4AYRDQFPhxjR9RfUtUZKusyMxJNAEbr2txOHvdkiQdUz0q8qfU_mWp9Cp5VHGIyHR4JEphADKb321R4y7zebf8y_pNse56Cln47HXCW-w-T8jHJPzH04Y7hMMWfJVkmHqehXs_bAOTaE3b_axFaUg"
	validToken   = "eyJraWQiOiJmTVlINkJyRHB3T2ZaOUVsNkZkM0N1UCtsME1kTFpLd3gwaUk5VXhOU0RnPSIsImFsZyI6IlJTMjU2In0.eyJzdWIiOiI1ZjlhZTcxZS1iZmRiLTRhYzctOGVhNy1kMTU2ZTdmMmEzODAiLCJkZXZpY2Vfa2V5IjoidXMtZWFzdC0xX2RkYmZiMWJkLTlhOWQtNDM5Yy1iNDQzLTMzZDE3MzZlY2ZjZCIsInRva2VuX3VzZSI6ImFjY2VzcyIsInNjb3BlIjoiYXdzLmNvZ25pdG8uc2lnbmluLnVzZXIuYWRtaW4iLCJhdXRoX3RpbWUiOjE1Nzg2MDM1OTYsImlzcyI6Imh0dHBzOlwvXC9jb2duaXRvLWlkcC51cy1lYXN0LTEuYW1hem9uYXdzLmNvbVwvdXMtZWFzdC0xX2dzQjU5d2Z6VyIsImV4cCI6MTU3OTEyMjk4MCwiaWF0IjoxNTc5MTE5MzgwLCJqdGkiOiI4NmVkMmMyOS1jZTZhLTQ3NTktYmI3Yy0xODc5NTQxODJkOGEiLCJjbGllbnRfaWQiOiIyMDg0dWtzbHNjODMxcHQyMDJ0MmR1ZHQ3YyIsInVzZXJuYW1lIjoicHVscGZyZWUifQ.fGyUHXgy20RH_EG9PEXAthVeEWWZzfNG-7uUU3DyKH2D3SZ4_WKUyM5Ttp_f6j1-KTaElALiwGhWuURHkb75c6UxnadFmthRRIPBOTg8hc6MCiKLY9JEhv8f5Wh7hlfPHLGUa0fpA8JagNXUQZrJQUj5zCIqm99Ngo5ysKtCmJHYDELs_de7wjdYvgLH7jBoOzujv8Zf2IFBYIWCG1RnrrGB3g4eZmnRZ7KFeVhUqoSXjW5jNlpDFCoE7-FmTDqTDVSV4QI3lTlvpI8HRHCnVxIjAUr9G1h6fHQAxuUUkJLafZb7AxwXBvkRlkLwH4DMUNNBYH_bycwLix2-cn3-Ig"
)

var cfg *config.Config

// IntegSuite struct
type IntegSuite struct {
	suite.Suite
}

// SetupTest method
func (s *IntegSuite) SetupTest() {
	// init config
	os.Setenv("Stage", "test")
	cfg = &config.Config{DefaultsFilePath: defaultsFP}
	err := cfg.Load()
	s.NoError(err)
}

// TestExpiredToken
func (s *IntegSuite) TestExpiredToken() {
	_, err := Validate(cfg.CognitoClientID, expiredToken)
	s.Error(err)
	s.Equal(err.Error(), "Expired token")
}

// TestFetchandTest
func (s *IntegSuite) TestFetchandTest() {
	expectedPrincipalID := fmt.Sprintf("%s|%s", username, cfg.CognitoClientID)
	principalID, err := Validate(cfg.CognitoClientID, validToken)
	s.NoError(err)
	// fmt.Printf("principalID: %s\n", principalID)
	s.Equal(expectedPrincipalID, principalID)
}

// TestIntegrationSuite function
func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegSuite))
}
