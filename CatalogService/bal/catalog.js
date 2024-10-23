const { ObjectId } = require('mongodb');
const dal = require('../dal/mongoDB.js');

async function GetCatalogs() {
    return await dal.Interface("get", "catalog", {});
}

async function GetCatalog(id) {
    return await dal.Interface("get", "catalog", {_id: new ObjectId(id)});
}

async function PostCatalog(catalog) {
    return await dal.Interface("post", "catalog", catalog);
}

async function PatchCatalog(id, catalog) {
    if (catalog._id) {
        delete catalog._id;
    }
    return await dal.Interface("patch", "catalog", [{_id: new ObjectId(id)}, catalog]);
}

async function DeleteCatalog(id) {
    return await dal.Interface("delete", "catalog", {_id: new ObjectId(id)});
}

module.exports = {
    GetCatalogs,
    GetCatalog,
    PostCatalog,
    PatchCatalog,
    DeleteCatalog
}