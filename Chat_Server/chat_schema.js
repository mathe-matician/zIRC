const ChatMessage = (athleteID, server, channel, message, timestamp = Date.now(), tags = {}) => {
    return {athleteID, server, channel, message, timestamp, tags, supportedVersions};
}
  
module.exports = ChatMessage;