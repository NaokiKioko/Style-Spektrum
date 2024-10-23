const express = require('express');
const app = express();
const bal = require('./bal/product.js');
const port = 3001;

app.use(express.json());

app.get('/', (req, res) => {
    res.send('Style Spektrum!');
});

app.get('/products', async (req, res) => {
    const products = await bal.GetProducts();
    res.send(products);
});

app.get('/product', async (req, res) => {
    const products = await bal.GetProducts();
    res.send(products[0]);
});

app.post('/product', async (req, res) => {
    const result = await bal.PostProducts(req.body);
    res.send(result);
});

app.patch('/product', async (req, res) => {
    const result = await bal.PatchProducts(req.body);
    res.send(result);
});

app.delete('/product', async (req, res) => {
    const result = await bal.DeleteProducts(req.body);
    res.send(result);
});

app.listen(port, () => {
    console.log(`Server listening on port ${port}`);
});