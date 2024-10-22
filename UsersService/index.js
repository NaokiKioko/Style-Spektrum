const express = require('express');
const app = express();
const auth = require('./auth/auth.js');
const bal = require('./bal/users.js');
const port = 3000;

app.get('/', (req, res) => {
    res.send('Style Spektrum!');
});

app.get('/register', (req, res) => {
    code = bal.registerUser(req.body.email, req.body.password);
    if (code === 200) {
        res.sendStatus(200).send('User registered!');
    } else if (code === 400) {
        res.sendStatus(400).send('Invalid email or Email in use!');
    } else {
        res.sendStatus(500).send('Error registering user!');
    }});

app.get('/login', (req, res) => {
    let token = bal.login(req.body.email, req.body.password);
    if (token === 401) {
        res.sendStatus(401).send('Invalid email or password!');
    } else if (token === 404) {
        res.sendStatus(404).send('User not found!');
    } else {
        res.sendStatus(200).send('Token: %s', token);
    }
});

app.get('/logout', (req, res) => {
    res.sendStatus(200).send('just delete the jwt <3');
});

app.get("/authenticate", auth.authenticateToken, (req, res) => {
    res.send('Authenticated user: %s', req.user.email);
});

app.listen(port, () => {
    console.log(`Server listening on port ${port}`);
});