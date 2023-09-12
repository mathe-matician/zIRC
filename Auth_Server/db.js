const { Client, Pool } = require('pg')
const os = require('os');

require('dotenv').config();

const DB = () => {
    const postgresConfig = {
        user: process.env.POSTGRES_USER,
        password: process.env.POSTGRES_PASSWORD,
        host: process.env.POSTGRES_HOST,
        database: process.env.POSTGRES_DATABASE,
        port: process.env.POSTGRES_PORT,
    };

    /**
     * 
     * @param {string} statementName Name of the prepared statement
     * @param {string} query The pg query
     * @param {[*]} values This can be an array of values used positionally in the query or a single value
     * @returns The prepared object
     */
    const Prepare = (statementName, query, values=null) => {
        if (!values) {
            return {name: statementName, text: query};
        }
        return {
            name: statementName,
            text: query,
            values: values
        }
    };

    const Exec = async (preparedStatement) => {
        console.log(`Exec start: preparedStatement: ${JSON.stringify(preparedStatement)}`)
        const client = await new Client(postgresConfig);
        try {
            await client.connect();
            const res = await client.query(preparedStatement);
            if (!res) {
                throw new Error("query error");
            }
            console.log(`Postgres Exec res = ${JSON.stringify(res)}`);
            console.log(`Postgres Exec res.rows[0] = ${JSON.stringify(res.rows[0])}`);
            return res.rows[0];
        } catch (error) {
            console.log(`Exec error: ${error}`);
            return null;
        } finally {
            client.end();
        }
    };

    return { Prepare, Exec };
};

module.exports = {
    DB
}

// const pool = new Pool({
//     host: '127.0.0.1',
//     user: process.env.POSTGRES_USER,
//     max: 20,
//     idleTimeoutMillis: 30000,
//     connectionTimeoutMillis: 2000,
//   })

// const doit = async () => {
//     const client = await pool.connect()
//     // await client.query('SELECT $1::text as message', ['Hello world!']);
//     await client.query('SELECT NOW()')
//     client.release()
// };

// doit();