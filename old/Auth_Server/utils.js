const isJson = (str) => {
    try {
      JSON.parse(str);
    } catch (e) {
      return false;
    }  
    return true;
};

const isTokenExpired = (token_expr) => {
  const now = new Date();
  const expiration = new Date(token_expr);
  if (now > expiration) {
    return true;
  }
  return false;
};

module.exports = {
    isJson,
    isTokenExpired
}