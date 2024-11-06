const security = require('../security/password.js');
const DatabaseInterface = require('../dal/mongoDB.js'); ;
const auth = require('../auth/auth.js');


async function registerUser(dal, email, password) {
    let endcode = null;
    let emailExists = await CheckEmailExists(dal, email);
    if (emailExists == false && ValidateEmail(email)) {
        // register user
        password = security.hashPassword(password);
        try {
            email = email.toLowerCase();
            await dal.interface("post", "Users", { email, password, role: "user", favoriteTags: []}).then((code) => {
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

async function login(dal, email, password) {
    // login user
    let jwt = null;
    let endcode = null;
    email = email.toLowerCase();
    await dal.interface("get", "Users", { email }).then((users) => {
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

async function DeleteAccount(dal, email, password) {
    let endcode = null;
    await dal.interface("get", "Users", { email }).then(async (users) => {
        if (users.length === 0) {
            endcode = 404;
        }
        user = users[0];
        if (security.comparePassword(password, user.password)) {
            await dal.interface("delete", "Users", { email }).then((code) => {
                endcode = code;
            });
        } else {
            endcode = 401;
        }
    });
    return endcode;
}

async function CheckEmailExists(dal, email) {
    let emailExists = null;
    await dal.interface("get", "Users", { email }).then((users) => {
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

async function GetUserByEmail(dal, email) {
    let user
    await dal.interface("get", "Users", { "email":email }).then((users) => {
        if (users.length === 0) {
            return 404;
        }
        user = users[0];
    });
    let tags = [];
    if (user.favoriteTags.length > 0) {
        await dal.interface("get", "Tags", { Name: { $in: user.favoriteTags } }).then((tag) => {
            tags = tag;
        });
    }
    for (let i = 0; i < tags.length; i++) {
        delete tags[i]._id;
    }
    return {
        email: user.email,
        role: user.role,
        favoriteTags: tags
    };
    
}

async function AddFavoriteTag(dal, email, tag) {
    let user;
    let returncode = 200;
    await dal.interface("get", "Users", { "email":email }).then((users) => {
        if (users.length === 0) {
            return 404;
        }
        user = users[0];
    });
    if (user.favoriteTags.includes(tag)) {
        returncode = 400;
    }
    if (returncode === 400) {
        return returncode;
    }
    user.favoriteTags.push(tag);
    await dal.interface("patch", "Users", [{ "email":email }, { favoriteTags: user.favoriteTags }]).then((code) => {
        returncode = code;
    });
    AlterTagFavoriteCount(dal, tag, 1);
    return returncode;
}

async function RemoveFavoriteTag(dal, email, tag) {
    let user;
    let returncode = 200;
    await dal.interface("get", "Users", { "email":email }).then((users) => {
        if (users.length === 0) {
            return 404;
        }
        user = users[0];
    });
    if (!user.favoriteTags.includes(tag)) {
        returncode = 400;
    }
    if (returncode === 400) {
        return returncode;
    }
    user.favoriteTags = user.favoriteTags.filter(e => e !== tag);
    await dal.interface("patch", "Users", [{ "email":email }, { favoriteTags: user.favoriteTags }]).then((code) => {
        returncode = code;
    });
    AlterTagFavoriteCount(dal, tag, -1);
    return returncode;
}

async function AlterTagFavoriteCount(dal, tag, ammount) {
    dal.interface("get", "Tags", { "name":tag }).then((tags) => {
        if (tags.length === 0) {
            return 404;
        }
        let tagData = tags[0];
        tagData.favoritecount = tagData.favoritecount + ammount;
        dal.interface("patch", "Tags", [{ "name":tag }, { "favoritecount": tagData.favoritecount }])
    });
}


module.exports = {
    registerUser,
    login,
    DeleteAccount,
    GetUserByEmail,
    AddFavoriteTag,
    RemoveFavoriteTag
}