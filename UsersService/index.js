const express = require('express');
const app = express();

const port = 3000;

app.get('/', (req, res) => {
  res.send('Style Spektrum!');
});

app.get('/register', (req, res) => {
    res.send('Register!');
    });

app.get('/login', (req, res) => {
    res.send('Login!');
    }
);

app.get('/logout', (req, res) => {
    res.send('Logout!');
    }
);

app.listen(port, () => {
  console.log(`Server listening on port ${port}`);
});