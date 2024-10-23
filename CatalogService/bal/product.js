const dal = require('../dal/mongoDB.js');

async function GetProducts() {
    return await dal.Interface("get", "product", {});
}

async function GetProduct(id) {
    return await dal.Interface("get", "product", {"_id":id});
}

async function PostProducts(product) {
    return await dal.Interface("post", "product", product);
}

async function PatchProducts(id, product) {
    return await dal.Interface("patch", "product", [{"_id":id}, product]);
}

async function DeleteProducts(id) {
    return await dal.Interface("delete", "product", {"_id":id});
}

module.exports = {
    GetProducts,
    GetProduct,
    PostProducts,
    PatchProducts,
    DeleteProducts
}