package nickserv

// SHOULD SUPPORT THESE COMMANDS
// REGISTER: Allows users to register their current nickname with a password and email address. This secures the nickname, preventing others from using it without authentication.
// IDENTIFY: Used to authenticate the user to their registered nickname by providing the password. This is necessary to prove ownership of the nickname.
// DROP: Unregisters a nickname, making it available for anyone to use. Only the owner of the nickname or network staff can execute this command.
// RECOVER: Reclaims a nickname that is currently in use by another user (possibly someone who failed to identify in time). It forces the user with that nickname to disconnect or changes their nickname, allowing the rightful owner to reclaim it.
// RELEASE: Releases a nickname that has been held or locked due to a RECOVER command. Once released, the rightful owner can re-use or claim it.

// SET PASSWORD: Allows users to change the password associated with their nickname. This is important for maintaining security if the user believes their account may have been compromised.
// SET EMAIL: Updates the email address associated with the nickname. The email is often used for recovery and verification purposes.
// SET ENFORCE: Enables or disables enforcement of nickname ownership. When enabled, if someone tries to use a registered nickname without identifying, they will be automatically disconnected or forced to change their nickname.
// SET KILL: Determines how NickServ handles users who attempt to use a registered nickname without identifying. It can be set to:
// ON: Enforces nickname ownership strictly, disconnecting users who fail to authenticate within a certain timeframe.
// QUICK: Enforces immediately but with shorter time limits.
// OFF: Disables enforcement, allowing anyone to use the nickname.
// SET PRIVATE: Hides the user’s nickname information from public view in /whois or other commands, enhancing privacy.
// SET SECURITY: Enforces stricter security policies for the nickname, such as requiring identification before joining channels or sending messages.

// GROUP: Groups additional nicknames under the same user account. This allows a user to manage multiple nicknames with a single login, making it convenient to use different nicknames while maintaining the same registration.
// UNGROUP: Removes a nickname from the group, separating it from the user’s account and making it independent again.
// ACC: Checks the account status of a user’s nickname, showing whether they are identified as the owner of that nickname.

// SENDPASS: Sends a password reset email to the registered email address, allowing users to recover access if they forget their password. This often requires administrator approval or verification for security.
// VERIFY: Confirms email changes or new nickname registrations. The user may need to provide a verification code sent to their email to complete the process.

// ACCESS: Manages a list of trusted hosts or IP addresses that can automatically use the registered nickname without requiring explicit identification. This can be helpful for users who frequently connect from the same location or want to bypass identification from known devices.
// SET NOOP: Prevents IRC operators (IRCops) from overriding or taking control of the user’s nickname, giving additional security to the user.

// INFO: Provides detailed information about a registered nickname, such as its registration date, last login time, email address (if visible), and other related settings.
// GHOST: Forces a disconnect of a “ghost” user who is occupying the nickname but is no longer connected to the network. This is useful when the user’s connection was dropped unexpectedly, and the nickname is still registered as being in use.
// LOGOUT: Logs out the user from their nickname, removing their authentication status. This can be used if the user wishes to switch accounts without reconnecting.

// Administrative and Oper Features (for IRC Operators)
// FORCE IDENTIFY: Allows network staff to authenticate a user manually, usually for troubleshooting or assisting users who are unable to identify on their own.
// SUSPEND: Temporarily suspends a nickname, preventing its use until the suspension is lifted. This might be used for disciplinary reasons or to enforce network rules.
// UNSUSPEND: Reinstates a suspended nickname, allowing the owner to resume using it.
// DROP (forcible): An operator can forcibly drop a nickname registration if it is deemed necessary, for example, if the nickname violates network policies.

// Help and Documentation
// HELP: Provides users with a list of available commands and detailed instructions on how to use them. Users can type:

var command_map = map[string]func(){
	"REGISTER":       register,
	"IDENTIFY":       identify,
	"DROP":           drop,
	"RECOVER":        recover,
	"RELEASE":        release,
	"SET PASSWORD":   set_password,
	"SET EMAIL":      set_email,
	"SET ENFORCE":    set_enforce,
	"SET KILL":       set_kill,
	"ON":             on,
	"QUICK":          quick,
	"OFF":            off,
	"SET PRIVATE":    set_private,
	"SET SECURITY":   set_security,
	"GROUP":          group,
	"UNGROUP":        ungroup,
	"ACC":            acc,
	"SENDPASS":       sendpass,
	"VERIFY":         verify,
	"ACCESS":         access,
	"SET NOOP":       set_noop,
	"INFO":           info,
	"GHOST":          ghost,
	"LOGOUT":         logout,
	"FORCE IDENTIFY": force_identify,
	"SUSPEND":        suspend,
	"UNSUSPEND":      unsuspend,
	"HELP":           help,
	// "DROP":           *NewCommand(quit, make(map[string]string), false),
}
