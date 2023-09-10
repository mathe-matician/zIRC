const capabilities = {
    "draft/account-registration": "before-connect,email-required,custom-account-name", 
    "sasl": "PLAIN",
    "chilkat": "ENCRYPT",
    "cap-notify": null
};

module.exports = { capabilities };