const express = require('express');
const app = express();
const auth = require('./auth/auth.js');
const bal = require('./bal/users.js');
const DatabaseInterface = require('./dal/mongoDB.js'); // Import the class
const port = 3000;

dal = new DatabaseInterface();

app.use(express.json());

app.get('/', (req, res) => {
    res.send('Style Spektrum!');
});

app.post('/register', async (req, res) => {
    const email = req.body.email;
    const password = req.body.password;
    const code = await bal.registerUser(dal, email, password);
    if (code === 200) {
        return res.status(200).send('User registered!');  // return to prevent further execution
    } else if (code === 400) {
        return res.status(400).send('Invalid email or Email in use!');
    } else {
        return res.status(500).send('Error registering user!');
    }
});

app.post('/login', async (req, res) => {
    let response = await bal.login(dal, req.body.email, req.body.password);
    if (response.statuscode === 401) {
        res.status(401).send('Invalid email or password!');
    } else if (response.statuscode === 404) {
        res.status(404).send('User not found!');
    } else {
        res.status(response.statuscode).send({token: response.jwt});
    }
});

app.get('/me', auth.authenticateToken, async (req, res) => {
    let user = await bal.GetUserByEmail(dal, req.user.email);
    res.status(200).send(user);
});

app.get('/logout', (req, res) => {
    res.status(200).send('just delete the jwt <3');
});

app.get("/authenticate", auth.authenticateToken, (req, res) => {
    res.send('Authenticated user: ' + req.user.email);
});

app.delete("/delete", auth.authenticateToken, async (req, res) => {
    let email = req.user.email
    let password = req.body.password
    let code = await bal.deleteAccount(dal, email, password)
    if (code === 200) {
        res.status(code).send('User deleted!')
    } else {
        res.status(code).send('Error deleting user!')
    }
});

app.post('/favorite/tag/:tag', auth.authenticateToken, async (req, res) => {
    let email = req.user.email
    let tag = req.params.tag
    let code = await bal.AddFavoriteTag(dal, email, tag)
    if (code === 200) {
        res.status(code).send('Tag added to favorites!')
    } else {
        res.status(code).send('Error adding tag to favorites!')
    }
});

app.delete('/favorite/tag/:tag', auth.authenticateToken, async (req, res) => {
    let email = req.user.email
    let tag = req.params.tag
    let code = await bal.RemoveFavoriteTag(dal, email, tag)
    if (code === 200) {
        res.status(code).send('Tag removed from favorites!')
    } else {
        res.status(code).send('Error removing tag from favorites!')
    }
});

app.listen(port, () => {
    console.log(`Server listening on port ${port}`);
});