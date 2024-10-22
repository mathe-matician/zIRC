package nickserv

import (
	"context"
	"fmt"
	"strings"

	"github.com/phuslu/log"
)

// REGISTER <account> <password> <email>
// <account>: Desired account name.
// <password>: Chosen password for the account.
// <email>: Email address associated with the account.
// If the server supports SASL, a user may first authenticate with their chosen SASL mechanism and then use the REGISTER command without a password.
// The server can validate the user based on the SASL credentials provided.
// REGISTER <account> * <email>

func register(params map[string]string) string {
	// func register(account, auth, email string) string {
	// check if nick is already registered
	account := params["account"]
	email := params["email"]

	var res string
	err := g_DB.QueryRow(context.Background(), "INSERT INTO users VALUES (default, $1, $2) ON CONFLICT (nick, email) DO NOTHING", account, email).Scan(&res)
	if err != nil {
		if strings.Contains(err.Error(), "there is no unique or exclusion constraint") {
			// return
			return fmt.Sprintf("*** Error: Nickname '%s' is currently in use by another user.", account)
		}
		log.Error().Msgf("Error executing PLAIN auth query: %s", err)
		panic(err)
	}

	// if not then register it
	// if there is an existing user with that nick, kick them!
	// send email for email verification
	return ""
}

// responses
// Incorrect Password:
// Message: *** Error: Password for nickname 'username' is incorrect.
// Nickname in Use:
// Message: *** Error: Nickname 'username' is currently in use by another user.
// Unregistered Nickname:
// Message: *** Error: Nickname 'username' is not registered. Please register it first.
// Account Not Verified:
// Message: *** Error: Your account must be verified before you can use this nickname.
// Password Too Short:
// Message: *** Error: The password you provided is too short.
// Registration Failed:
// Message: *** Error: Unable to register nickname 'username'. Please try again later.
// Nickname Already Registered:
// Message: *** Error: Nickname 'username' is already registered. Please choose a different nickname.
