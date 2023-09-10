# Auth Server Modules

All modules (nodejs files) in the modules directory MUST have the following methods:

- `Exec`

The `Exec` method is the main entrypoint to any added module. Exec MUST take at least 1 argument of type `array` which will be an array of arguments for the particular `Exec` function.

For example:

```javascript
// auth_plain.js

// args that are passed into the function
// args[0] == username
// args[1] == encrypted password
args = ["zach", "HDJKALJWNKLJJNL32AKKNKJ"]

const Exec = (args) => {
    // simplified
    const res = authCheckUsernamePassword(args);
    if (!res) {
        throw new Error("Username or password incorrect");
    } 
};
```