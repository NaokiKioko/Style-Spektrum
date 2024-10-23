const dal = require('../dal/mongoDB.js');

function GetProducts() {
    return dal.Interface("get", "product", {});
}

function GetProduct(id) {
    return dal.Interface("get", "product", {"_id":id});
}

function PostProducts(product) {
    return dal.Interface("post", "product", product);
}

function PatchProducts(id, product) {
    return dal.Interface("patch", "product", [{"_id":id}, product]);
}

function DeleteProducts(id) {
    return dal.Interface("delete", "product", {"_id":id});
}

module.exports = {
}