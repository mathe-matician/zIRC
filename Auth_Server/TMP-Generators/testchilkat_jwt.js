const chilkatManager = require('./Chilkat/chilkat_manager');
require('dotenv').config();

chilkatManager.unlock_bundle();

const jwt = new chilkatManager.chilkat.Jwt();

var jose = new chilkatManager.chilkat.JsonObject();
var success = jose.AppendString("alg", "HS256");
success = jose.AppendString("typ", "JWT");

// create claims also known as jwt payload
var claims = new chilkatManager.chilkat.JsonObject();
success = claims.AppendString("iss", "http://example.org");
success = claims.AppendString("sub", "Zach Mathe");
success = claims.AppendString("aud", "http://example.com");

var curDateTime = jwt.GenNumericDate(0);
// set timestamp 
success = claims.AddIntAt(-1, "iat", curDateTime);
// set the "not process before" timestamp to now
success = claims.AddIntAt(-1, "nbf", curDateTime);
// set expiration time
success = claims.AddIntAt(-1, "exp", curDateTime+3600);

// produce smallest JWT
jwt.AutoCompact = true;

var token = jwt.CreateJwt(jose.Emit(), claims.Emit(), "secret");

console.log(`JWT token = ${token}`);

var sigVerified = jwt.VerifyJwt(token, "secret");
console.log(`with correct password: ${sigVerified}`)

var sigVerified = jwt.VerifyJwt(token, "secret2");
console.log(`with incorrect password: ${sigVerified}`)

var leeway = 60;
var bTimeValid = jwt.IsTimeValid(token, leeway);
console.log(`time constraints valid: ${bTimeValid}`);

var payload = jwt.GetPayload(token);
console.log(payload)

var json = new chilkatManager.chilkat.JsonObject();
var success = json.Load(payload);
json.EmitCompact = false;
console.log(json.Emit());

var joseHeader = jwt.GetHeader(token);
console.log(joseHeader)

success = json.Load(joseHeader);
json.EmitCompact = false
console.log(json.Emit());