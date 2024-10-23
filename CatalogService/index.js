const express = require('express');
const app = express();
const bal = require('./bal/catalog.js');
const port = 3001;

app.use(express.json());

app.get('/', (req, res) => {
    res.send('Style Spektrum!');
});

app.get('/catalogs', async (req, res) => {
    const catalogs = await bal.GetCatalogs();
    res.status(200).send(catalogs);
});

app.get('/catalog/:id', async (req, res) => {
    let id = req.params.id;
    const catalogs = await bal.GetCatalog(id);
    if (catalogs.length === 0) {
        res.status(404).send('Catalog not found');
    }
    res.status(200).send(catalogs[0]);
});

app.post('/catalog', async (req, res) => {
    const result = await bal.PostCatalog(req.body);
    res.status(201).send(result);
});

app.patch('/catalog/:id', async (req, res) => {
    let id = req.params.id;
    req.body.id = id;
    let result = await bal.PatchCatalog(id, req.body);
    res.status(202).send(result);
});

app.delete('/catalog/:id', async (req, res) => {
    let id = req.params.id;
    req.body.id = id;
    const result = await bal.DeleteCatalog(id, req.body);
    res.status(202).send(result);
});

app.listen(port, () => {
    console.log(`Server listening on port ${port}`);
});