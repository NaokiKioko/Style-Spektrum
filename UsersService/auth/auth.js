// generate jwt token
const e = require('express');
const jwt = require('jsonwebtoken');
const secret = process.env.JWT_SECRET
const timeout = process.env.JWT_TIMEOUT

const generateToken = (email, role) => {
    let token = jwt.sign({ email, role }, secret, { expiresIn: timeout });
    return token;
};

const authenticateToken = (req, res, next) => {
    const authHeader = req.headers["authorization"];
    // Extracting token from authorization header
    const token = authHeader && authHeader.split(" ")[1];
    // Checking if the token is null
    if (!token) {
        return res.status(401).send("Authorization failed. No access token.");
    }

    // Verifying if the token is valid
    jwt.verify(token, secret, (err, user) => {
        if (err) {
            console.log(err);
            return res.status(403).send("Could not verify token");
        }
        req.user = user;
        next();  // Ensure next() is called only after successful verification
    });
};

module.exports = {
    generateToken,
    authenticateToken,
}