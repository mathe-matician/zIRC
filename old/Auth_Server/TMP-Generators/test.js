const util = require('util');
const fs = require('node:fs');

const statAsync = util.promisify(fs.stat);

const doit = async () => {
    try {
        const stats = await statAsync("./modules/auth_plain.js");
        return stats.size
    } catch (error) {
        console.log(error);
    }
}

const main = async () => {
    const res = await doit();
    console.log(`res = ${res}`);
}

main()