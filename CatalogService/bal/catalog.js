const { ObjectId } = require('mongodb');
const dal = require('../dal/mongoDB.js');

async function GetCatalogs() {
    return await dal.Interface("get", "Catalog", {});
}

async function GetCatalog(id) {
    return await dal.Interface("get", "Catalog", {_id: new ObjectId(id)});
}

async function GetCatalogbyTags(tags) {
    return await dal.Interface("get", "Catalog", {tags: { $all: tags }});
}


async function PostCatalog(catalog) {
    let code = await dal.Interface("post", "Catalog", catalog);
    if (code === 500) {
        return 500;
    } else {
        await PostTags(catalog.tags);
    }
    return code;
}

async function PatchCatalog(id, catalog) {
    if (catalog._id) {
        delete catalog._id;
    }
    code = await dal.Interface("patch", "Catalog", [{_id: new ObjectId(id)}, catalog]);
    if (code === 500) {
        return 500;
    } else {
        await PostTags(catalog.tags);
    }
    return code;
}

async function DeleteCatalog(id) {
    return await dal.Interface("delete", "Catalog", {_id: new ObjectId(id)});
}

async function GetTags() {
    return await dal.Interface("get", "Tags", {});
}

async function PostTags(tagNames) {
    for (let i = 0; i < tagNames.length; i++) {
        let tag = await dal.Interface("get", "Tags", {"name": tagNames[i]});
        if (tag === 500) {
            return 500;
        }
        if (tag.length === 0) {
            await dal.Interface("post", "Tags", {"name": tagNames[i], "favoritecount": 0});
        }
    }
}

module.exports = {
    GetCatalogs,
    GetCatalog,
    GetCatalogbyTags,
    PostCatalog,
    PatchCatalog,
    DeleteCatalog,
    GetTags,
}