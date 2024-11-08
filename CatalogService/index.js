const express = require('express');
const app = express();
const bal = require('./bal/catalog.js');
const DatabaseInterface = require('./dal/mongoDB.js'); // Import the class
const port = 3001;

dal = new DatabaseInterface();

app.use(express.json());

app.get('/', (req, res) => {
    res.send('Style Spektrum!');
});

app.get('/catalog', async (req, res) => {
    const catalogs = await bal.GetCatalogs(dal);
    res.status(200).send(catalogs);
});

app.get('/catalog/tags/:tags', async (req, res) => {
    let tags = req.params.tags.split(',');
    const catalogs = await bal.GetCatalogbyTags(dal, tags);
    res.status(200).send(catalogs);
});

app.get('/catalog/:id', async (req, res) => {
    let id = req.params.id;
    const catalogs = await bal.GetCatalog(dal, id);
    if (catalogs.length === 0) {
        res.status(404).send('Catalog not found');
        return;
    }
    if (catalogs === 500) {
        res.status(500).send('Error updating catalog');
        return;
    }
    res.status(200).send(catalogs[0]);
});


app.post('/catalog', async (req, res) => {
    await bal.PostCatalog(dal, req.body);
    res.sendStatus(201);
});

app.patch('/catalog/:id', async (req, res) => {
    let id = req.params.id;
    let result = await bal.PatchCatalog(dal, id, req.body);
    if (result === 404) {
        res.status(404).send('Catalog not found');
        return;
    }
    if (result === 500) {
        res.status(500).send('Error updating catalog');
        return;
    }
    res.sendStatus(202);
});

app.delete('/catalog/:id', async (req, res) => {
    let id = req.params.id;
    req.body.id = id;
    const result = await bal.DeleteCatalog(dal, id, req.body);
    if (result === 404) {
        res.status(404).send('Catalog not found');
        return;
    }
    if (result === 500) {
        res.status(500).send('Error updating catalog');
        return;
    }
    res.sendStatus(202);
});

app.get('/tags', async (req, res) => {
    const catalogs = await bal.GetTags(dal, []);
    for (let i = 0; i < catalogs.length; i++) {
        delete catalogs[i]._id;
    }
    res.status(200).send(catalogs);
});

app.get('/report/:id', async (req, res) => {
    let id = req.params.id;
    const reports = await bal.GetReports(dal, id);
    if (reports.length === 0) {
        res.status(404).send('Report not found');
        return;
    }
    if (reports === 500) {
        res.status(500).send('Error getting report');
        return;
    }
    res.status(200).send(reports[0]);
});

app.get('/report/:id/field/:field', async (req, res) => {
    let id = req.params.id;
    let field = req.params.field;
    const reports = await bal.GetReportField(dal, id, field);
    if (reports.length === 0) {
        res.status(404).send('Report not found');
        return;
    }
    if (reports === 500) {
        res.status(500).send('Error getting report');
        return;
    }
    res.status(200).send(reports[0]);
});

app.post('/report/:id', async (req, res) => {
    req.body.id = req.params.id;
    await bal.PostReport(dal, req.body);
    res.sendStatus(201);
});




app.post ('/report/:itemID/tag/:tagName', async (req, res) => {
    req.body.id = req.params.itemID;
    req.body.tagName = req.params.tagName;
    await bal.PostReportItemTag(dal, req.body);
    res.sendStatus(201);
});

app.listen(port, () => {
    console.log(`Server listening on port ${port}`);
});