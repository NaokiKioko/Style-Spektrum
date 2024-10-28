const { ObjectId } = require('mongodb');
const dal = require('../dal/mongoDB.js');

async function GetCatalogs() {
    return await dal.Interface("get", "catalog", {});
}

async function GetCatalog(id) {
    return await dal.Interface("get", "catalog", {_id: new ObjectId(id)});
}

async function GetCatalogbyTags(tags) {
    return await dal.Interface("get", "catalog", {tags: { $all: tags }});
}
async function GetAllTags(tags) {
    let products = await dal.Interface("get", "catalog", {});
    let tagsArray = [];
    products.forEach(product => {
        product.tags.forEach(tag => {
            if (!tagsArray.includes(tag)) {
                tagsArray.push(tag);
            }
        });
    });
    return {"tags":tagsArray};
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

async function GetTags() {
    return await dal.Interface("get", "tag", {});
}

module.exports = {
    GetCatalogs,
    GetCatalog,
    GetCatalogbyTags,
    PostCatalog,
    PatchCatalog,
    DeleteCatalog,
    GetTags,
    GetAllTags
}