import Client from "../client";
import assert from 'node:assert/strict';

/**
 * Test valid nickname
 */
const TestNickNameValid = () => {
    const expected = Error;
    let nicknames = [];

    // Nickname length cannot be longer than 64 characters
    let c = Client("JeffJeffJeffJeffJeffJeffJeffJeffJeffJeffJeffJeffJeffJeffJeffJeffJeff");
    assert.throws(c, expected);

    // Nickname cannot be prefixed with specific characters $:#&
    nicknames = [
        ":Jeff",
        "$Jeff",
        "#Jeff",
        "&Jeff",
    ];

    for (let i = 0; i < nicknames.length; i++) {
        assert.throws(Client(nicknames[i]), expected);
    }

    // Nickname cannot be null or empty
    c = Client();
    assert.throws(c, expected);

    c = Client("");
    assert.throws(c, expected);

    // Nickname cannot have specific characters \s,*?!@.
    nicknames = [
        " ",
        " Jeff", "J eff", "Je ff", "Jef f", "Jeff ",
        ",Jeff", "J,eff", "Je,ff", "Jef,f", "Jeff,",
        "*Jeff", "J*eff", "Je*ff", "Jef*f", "Jeff*",
        "?Jeff", "J?eff", "Je?ff", "Jef?f", "Jeff?",
        "!Jeff", "J!eff", "Je!ff", "Jef!f", "Jeff!",
        "@Jeff", "J@eff", "Je@ff", "Jef@f", "Jeff@",
    ];

    for (let i = 0; i < nicknames.length; i++) {
        assert.throws(Client(nicknames[i]), expected);
    }
}