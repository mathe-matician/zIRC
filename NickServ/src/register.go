package nickserv

import (
	"context"
	"fmt"
	"strings"

	"github.com/phuslu/log"
	"golang.org/x/crypto/bcrypt"
)

// REGISTER <account> <password> <email>
// <account>: Desired account name.
// <password>: Chosen password for the account.
// <email>: Email address associated with the account.
// If the server supports SASL, a user may first authenticate with their chosen SASL mechanism and then use the REGISTER command without a password.
// The server can validate the user based on the SASL credentials provided.
// REGISTER <account> * <email>

func register(params string) string {
	// func register(account, auth, email string) string {
	// check if nick is already registered
	params_split := strings.Split(params, " ")
	if len(params_split) < 3 {
		return ""
	}

	account := params_split[0]
	password := params_split[1]
	email := params_split[2]

	// TODO
	// check password requirements

	// password_hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	_, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		// TODO
		// if this fails then the user will need to use the SET command to set their password? idk
		return ""
	}

	var res string
	query := "with new_user as (insert into users (id, nick, email) values (default, $1, $2) on conflict (email) do nothing returning id) insert into auth (user_id, auth_type, credentials) select id, 'PLAIN', 'pw' from new_user"
	err = g_DB.QueryRow(context.Background(), query, account, email).Scan(&res)
	if err != nil {
		if strings.Contains(err.Error(), "there is no unique or exclusion constraint") {
			// return
			return fmt.Sprintf("*** Error: Nickname '%s' is currently in use by another user.", account)
		}
		log.Error().Msgf("Error executing PLAIN auth query: %s", err)
		panic(err)
	}

	log.Debug().Msgf("REGISTER db res: %s", res)

	var hash_res string
	err = g_DB.QueryRow(context.Background(), "", account).Scan(&hash_res)
	if err != nil {

	}

	log.Debug().Msgf("REGISTER hash res: %s", hash_res)

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
