const bcrypt = require ("bcryptjs")

const hashPassword = (password) => {
    let salt = bcrypt.genSaltSync()
    let hash = bcrypt.hashSync(password, salt)
    return hash;
}

    const comparePassword = (password1, password2) => {
        let isValid =  bcrypt.compareSync(password1, password2)
        return isValid
    }

    module.exports = {
        hashPassword,
        comparePassword
    }