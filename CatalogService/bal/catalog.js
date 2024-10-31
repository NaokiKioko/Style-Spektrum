const Databaseinterface = require('../dal/mongoDB.js'); // Import the class

async function GetCatalogs(dal) {
    return await dal.interface("get", "Catalog", {});
}

async function GetCatalog(dal, id) {
    return await dal.interface("get", "Catalog", {_id: new ObjectId(id)});
}

async function GetCatalogbyTags(dal, tags) {
    return await dal.interface("get", "Catalog", {"Tags": { $in: tags }});
}


async function PostCatalog(dal, catalog) {
    let code = await dal.interface("post", "Catalog", catalog);
    if (code === 500) {
        return 500;
    } else {
        await PostTags(catalog.Tags);
    }
    return code;
}

async function PatchCatalog(dal, id, catalog) {
    if (catalog._id) {
        delete catalog._id;
    }
    code = await dal.interface("patch", "Catalog", [{_id: new ObjectId(id)}, catalog]);
    if (code === 500) {
        return 500;
    } else {
        await PostTags(catalog.tags);
    }
    return code;
}

async function DeleteCatalog(dal, id) {
    return await dal.interface("delete", "Catalog", {_id: new ObjectId(id)});
}

async function GetTags(dal) {
    return await dal.interface("get", "Tags", {});
}

async function PostTags(dal, tagNames) {
    for (let i = 0; i < tagNames.length; i++) {
        let tag = await dal.interface("get", "Tags", {"Name": tagNames[i]});
        if (tag === 500) {
            return 500;
        }
        if (tag.length === 0) {
            await dal.interface("post", "Tags", {"Name": tagNames[i], "Favoritecount": 0});
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