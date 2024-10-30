async function GetCatalogs(dal) {
    return await dal.Interface("get", "Catalog", {});
}

async function GetCatalog(dal, id) {
    return await dal.Interface("get", "Catalog", {_id: new ObjectId(id)});
}

async function GetCatalogbyTags(dal, tags) {
    return await dal.Interface("get", "Catalog", {"Tags": { $in: tags }});
}


async function PostCatalog(dal, catalog) {
    let code = await dal.Interface("post", "Catalog", catalog);
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
    code = await dal.Interface("patch", "Catalog", [{_id: new ObjectId(id)}, catalog]);
    if (code === 500) {
        return 500;
    } else {
        await PostTags(catalog.tags);
    }
    return code;
}

async function DeleteCatalog(dal, id) {
    return await dal.Interface("delete", "Catalog", {_id: new ObjectId(id)});
}

async function GetTags(dal) {
    return await dal.Interface("get", "Tags", {});
}

async function PostTags(dal, tagNames) {
    for (let i = 0; i < tagNames.length; i++) {
        let tag = await dal.Interface("get", "Tags", {"Name": tagNames[i]});
        if (tag === 500) {
            return 500;
        }
        if (tag.length === 0) {
            await dal.Interface("post", "Tags", {"Name": tagNames[i], "Favoritecount": 0});
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