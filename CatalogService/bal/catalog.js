const dal = require('../dal/mongoDB.js');

async function GetCatalogs() {
    return await dal.Interface("get", "catalog", {});
}

async function GetCatalog(id) {
    return await dal.Interface("get", "catalog", {"id":id});
}

async function PostCatalog(catalog) {
    return await dal.Interface("post", "catalog", catalog);
}

async function PatchCatalog(id, catalog) {
    return await dal.Interface("patch", "catalog", [{"_id":id}, catalog]);
}

async function DeleteCatalog(id) {
    return await dal.Interface("delete", "catalog", {"_id":id});
}

module.exports = {
    GetCatalogs,
    GetCatalog,
    PostCatalog,
    PatchCatalog,
    DeleteCatalog
}