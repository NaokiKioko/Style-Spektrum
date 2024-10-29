const security = require('../security/password.js');
const dal = require('../dal/mongoDB.js');
const auth = require('../auth/auth.js');


async function registerUser(email, password) {
    let endcode = null;
    let emailExists = await CheckEmailExists(email);
    if (emailExists == false && ValidateEmail(email)) {
        // register user
        password = security.hashPassword(password);
        try {
            await dal.Interface("post", "Users", { email, password, role: "user", favoriteTags: []}).then((code) => {
                endcode = code;
            });
        }
        catch (err) {
            endcode = 500;
        }
    } else {
        endcode = 400;
    }
    return endcode;
}

async function login(email, password) {
    // login user
    let jwt = null;
    let endcode = null;
    await dal.Interface("get", "Users", { email }).then((users) => {
        if (users.length === 0) {
            endcode = 404;
            return
        }
        user = users[0];
        if (security.comparePassword(password, user.password)) {
            jwt = auth.generateToken(email, user.role);
            endcode = 200;
        }else {
            endcode = 401;
        }
    });
    return {"jwt": jwt, "statuscode": endcode}
}

async function deleteAccount(email, password) {
    let endcode = null;
    await dal.Interface("get", "Users", { email }).then(async (users) => {
        if (users.length === 0) {
            endcode = 404;
        }
        user = users[0];
        if (security.comparePassword(password, user.password)) {
            await dal.Interface("delete", "Users", { email }).then((code) => {
                endcode = code;
            });
        } else {
            endcode = 401;
        }
    });
    return endcode;
}

async function CheckEmailExists(email) {
    let emailExists = null;
    await dal.Interface("get", "Users", { email }).then((users) => {
        if (users.length === 0) {
            emailExists = false;
        }
        else {emailExists = true;}
    });
    return emailExists;
}

function ValidateEmail(email) {
    var regex = /^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9-]+(?:\.[a-zA-Z0-9-]+)*$/;
    const isValid = email && typeof email === 'string' && email.match(regex);

    return isValid ? true : false;
}

function GetUserByEmail(email) {
    return dal.Interface("get", "Users", { email });
}

module.exports = {
    registerUser,
    login,
    deleteAccount,
    GetUserByEmail
}