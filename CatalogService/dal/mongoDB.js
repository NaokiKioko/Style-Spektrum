const { MongoClient, ServerApiVersion } = require('mongodb');
const connectionString = process.env.MONGO_CONNECTION_STRING || null;
const databaseName = "StyleSpektrum";
const client = new MongoClient(connectionString,
    {
        serverApi: {
            version: ServerApiVersion.v1,
            strict: true,
            deprecationErrors: true,
        }
    });
    
    async function connectClient() {
        if (!client.isConnected()) {
            await client.connect();
        }
        return client;
    }
    
    const Interface = async (request, type, object) => {
        try {
            const dbClient = await connectClient();
        await dbClient.db("type").command({ ping: 1 });
        type = type + "";
        switch (request.toLowerCase()) {
            case "post":
                return await PostAny(dbClient, type, object);
            case "get":
                return await GetAny(dbClient, type, object);
            case "patch":
                return await PatchAny(dbClient, type, object);
            case "delete":
                return await DeleteAny(dbClient, type, object);
            default:
                return 404;
        }
    } catch (err) {
        console.log(err);
        return 500;
    } finally {
        // Ensure the client will close when you finish/error
        await client.close();
    }
};

const GetAny = async (client, type, object) => {
    if (CheckConnectionString() === 500) { return; }
    try {
        const cursor = await client.db(databaseName).collection(type).find(object)
        return await cursor.toArray()
    } catch (err) {
        console.log(err);
        return 500;
    }
}

const PostAny = async (client, type, object) => {
    if (CheckConnectionString() === 500) { return; }
    try {
        await client.db(databaseName).collection(type).insertOne(object);
        return 200;
    } catch (err) {
        console.log(err);
        return 500;
    }
}

// object[0] = query
// object[1] = set
const PatchAny = async (client, type, object) => {
    if (CheckConnectionString() === 500) { return; }
    try {
        const options = { upsert: true };
        const set = { $set: object[1] }
        await client.db(databaseName).collection(type).updateOne(object[0], set, options);
        return 200
    } catch (err) {
        console.log(err);
        return 500;
    }
}

const DeleteAny = async (client, type, object) => {
    if (CheckConnectionString() === 500) { return; }
    try {
        await client.db(databaseName).collection(type).deleteOne(object);
        return 200;
    } catch (err) {
        console.log(err);
        return 500;
    }
}

const CheckConnectionString = async () => {
    if (connectionString === null) {
        return 500;
    }
    return 200;
}

const CheckConnection = async (client) => {

}



module.exports = {
    Interface,
}