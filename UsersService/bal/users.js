const security = require('../security/password.js');
const dal = require('../dal/mongoDB.js');
const auth = require('../auth/auth.js');


function registerUser(email, password) {
    if (CheckEmailValidAndExists(email) === 404) {
        // register user
        password = security.hashPassword(password);
        try {
            dal.Interface("post", "Users", { email, password, role: "user" });
        } catch (error) {
            return 500;
        }
        return 200;
    }
    return 400;
}

function login(email, password) {
    // login user
    let user = dal.Interface("get", "Users", { email });
    if (user.length === 0) {
        return 404;
    }
    user = user[0];
    if (security.comparePassword(password, user.password)) {
        return auth.generateToken(email, user.role);
    }
    return 401;
}

function logout(jwt) {
    auth.invlaidateToken(jwt);
    return 200;
}

function CheckEmailValidAndExists(email) {
    // check if email exists
    if (!ValidateEmail(email)) {
        return 400;
    }
    let user = dal.Interface("get", "Users", { email });
    if (user.length === 0) {
        return 404;
    }
    return 200;
}

function ValidateEmail(input) {
    var validRegex = /^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9-]+(?:\.[a-zA-Z0-9-]+)*$/;
    if (match(validRegex)) {
        return true;
    } else {
        return false;
    }
}

module.exports = {
    registerUser,
    login,
    logout
}