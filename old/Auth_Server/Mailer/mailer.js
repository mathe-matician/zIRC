const nodemailer = require("nodemailer");
const cryptoManager = require("../Crypto/crypto_manager");

const Mailer = () => {

    const CreateTransport = () => {
        const transporter = nodemailer.createTransport({
          host: "smtp.forwardemail.net",
          port: 465,
          secure: true,
          auth: {
            // TODO: replace `user` and `pass` values from <https://forwardemail.net>
            user: "REPLACE-WITH-YOUR-ALIAS@YOURDOMAIN.COM",
            pass: "REPLACE-WITH-YOUR-GENERATED-PASSWORD",
          },
        });
    };

    const CreateTransportOAUTH = () => {
        const transporter = nodemailer.createTransport({
          service: 'gmail',
          auth: {
            type: 'OAuth2',
            user: process.env.MAIL_USERNAME,
            pass: process.env.MAIL_PASSWORD,
            clientId: process.env.OAUTH_CLIENTID,
            clientSecret: process.env.OAUTH_CLIENT_SECRET,
            refreshToken: process.env.OAUTH_REFRESH_TOKEN
          }
        });
    };

    const SendMail = async (token) => {
        const info = await transporter.sendMail({
          from: '"Fred Foo 👻" <foo@example.com>', // sender address
          to: "bar@example.com, baz@example.com", // list of receivers
          subject: "Hello ✔", // Subject line
          text: "Hello world?", // plain text body
          html: "<b>Hello world?</b>", // html body
        });

        console.log("Message sent: %s", info.messageId);
        // Message sent: <b658f8ca-6296-ccf4-8306-87d57a0b4321@example.com>

        //
        // NOTE: You can go to https://forwardemail.net/my-account/emails to see your email delivery status and preview
        //       Or you can use the "preview-email" npm package to preview emails locally in browsers and iOS Simulator
        //       <https://github.com/forwardemail/preview-email>
        //
    };

    const SendMailTest = async (token) => {
        const transporter = nodemailer.createTransport({
          port: 25, // Postfix uses port 25
          host: 'localhost',
          tls: {
            rejectUnauthorized: false
          },
        });
          
        var message = {
          from: 'noreply@domain.com',
          to: 'user@example.com',
          subject: 'Confirm Email',
          text: 'Please confirm your email',
          html: '<p>Please confirm your email</p>'
        };
          
        transporter.sendMail(message, (error, info) => {
          if (error) {
              return console.log(error);
          }
          console.log('Message sent: %s', info.messageId);
        });
    }; 

    return { 
        CreateTransport, 
        CreateTransportOAUTH, 
        SendMail, 
        SendMailTest 
    };
}

// SendMail().catch(console.error);

module.exports = {
    Mailer
}