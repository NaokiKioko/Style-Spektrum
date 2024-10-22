const { MongoClient } = require("mongodb");
const { type } = require("os");
const { env } = require("process");
const { object } = require("webidl-conversions");

const connectionString = process.env.MONGO_CONNECTION_STRING || null
const databaseName = "GameExchange"
const user = "Users"


const Interface = async (request, type, object) => {
    // console.log(request, type, object);
    type = type + "";
    switch (request.toLowerCase()) {
        case "post":
            return await PostAny(type, object);
        case "get":
            return await GetAny(type, object);
        case "patch":
            return await PatchAny(type, object);
        case "delete":
            return await DeleteAny(type, object);
        default:
            return 404;
    }
}

const GetAny = async (type, object) => {
    if (CheckConnectionString() === 500) { return; }
    const client = new MongoClient(connectionString)
    const database = client.db(databaseName);
    const Employees = database.collection(type);

    const cursor = Employees.find(object)

    return await cursor.toArray()

}

const PostAny = async (type, object) => {
    if (CheckConnectionString() === 500) { return; }
    const client = new MongoClient(connectionString)
    const database = client.db(databaseName);
    const Employees = database.collection(type);

    await Employees.insertOne(object);
    return 200;
}

// object[0] = query
// object[1] = set
const PatchAny = async (type, object) => {
    if (CheckConnectionString() === 500) { return; }
    const client = new MongoClient(connectionString)
    const database = client.db(databaseName);
    const Employees = database.collection(type);

    const options = { upsert: true };
    const set = { $set: object[1] }

    await Employees.updateOne(object[0], set, options);
    return 200
}

const DeleteAny = async (type, object) => {
    if (CheckConnectionString() === 500) { return; }
    const client = new MongoClient(connectionString)
    const database = client.db(databaseName);
    const Employees = database.collection(type);

    await Employees.deleteOne(object);
    return 200;
}

const CheckConnectionString = async () => {
    if (connectionString === null) {
        return 500;
    }
    return 200;
}

module.exports = {
    interface: Interface,
}