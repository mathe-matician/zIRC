package chat

// REGISTER <account> <password> <email>
// <account>: Desired account name.
// <password>: Chosen password for the account.
// <email>: Email address associated with the account.
// If the server supports SASL, a user may first authenticate with their chosen SASL mechanism and then use the REGISTER command without a password.
// The server can validate the user based on the SASL credentials provided.
// REGISTER <account> * <email>
func register(params map[string]interface{}) Response {

	// TODO
	// if client has already AUTHENTICATEd, then if they want to register
	// they don't need a password
	// e.g. REGISTER <account> * <email>
	// we would need to check whether they have authenticated

	// var _username, _password string
	// _, err := register_stmt.QueryOne(pg.Scan(&_username, &_password), _username, _password)

	// if err != nil {
	// 	log.Error().Msgf("Error executing PLAIN auth query: %s", err)
	// 	return ERR_SASLFAIL("SASL Failed")
	// }

	return &Reply{}
}

func verify(params map[string]interface{}) Response {
	return &Reply{}
}
