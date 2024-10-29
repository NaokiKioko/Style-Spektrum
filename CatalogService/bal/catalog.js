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
    let tags = await dal.Interface("get", "tag", {});
    if (tags === 500) {
        return 500;
    }
    if (tags.length === 0) {
        return {"tags":[]};
    }
    return tags;
}

async function PostCatalog(catalog) {
    let code = await dal.Interface("post", "catalog", catalog);
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
    code = await dal.Interface("patch", "catalog", [{_id: new ObjectId(id)}, catalog]);
    if (code === 500) {
        return 500;
    } else {
        await PostTags(catalog.tags);
    }
    return code;
}

async function DeleteCatalog(id) {
    return await dal.Interface("delete", "catalog", {_id: new ObjectId(id)});
}

async function GetTags() {
    return await dal.Interface("get", "tag", {});
}

async function PostTags(tagNames) {
    for (let i = 0; i < tagNames.length; i++) {
        let tag = await dal.Interface("get", "tag", {"name": tagNames[i]});
        if (tag === 500) {
            return 500;
        }
        if (tag.length === 0) {
            await dal.Interface("post", "tag", {"name": tagNames[i], "favoritecount": 0});
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
    GetAllTags
}