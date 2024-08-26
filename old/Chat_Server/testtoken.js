let message = 'PRIVMSG #General :Hello : world'
// const messageIndex = message.indexOf(":")
// const msgStart = message.substring(0, messageIndex)
// const msg = message.substring(messageIndex, message.length);
// console.log(`start: ${msgStart}`)
// console.log(`start.split(): ${msgStart.trim().split(' ')}`)
// console.log(`msg: ${msg}`)

let parameters = [];
if (message.startsWith("PRIVMSG")) {
  const messageIndex = message.indexOf(":");
  const msgStart = message.substring(0, messageIndex);
  const msg = message.substring(messageIndex, message.length);
  const msgStartSplit = msgStart.trim().split(' ');
  msgStartSplit.push(msg);
  parameters = msgStartSplit;
}
console.log(parameters)
