const { MongoClient, ServerApiVersion } = require('mongodb');

class DatabaseInterface {
    constructor(connectionString = process.env.MONGO_CONNECTION_STRING) {
        this.connectionString = connectionString || null;
        this.databaseName = "StyleSpektrum";
        this.client = new MongoClient(this.connectionString, {
            serverApi: {
                version: ServerApiVersion.v1,
                strict: true,
                deprecationErrors: true,
            }
        });
        this.isConnected = false; // Track connection state
    }

    async connect() {
        if (!this.isConnected) {
            try {
                await this.client.connect();
                await this.client.db("admin").command({ ping: 1 });
                this.isConnected = true;
                console.log("Connected to the database");
            } catch (err) {
                console.error("Failed to connect to the database:", err);
                throw err;
            }
        }
    }

    async disconnect() {
        if (this.isConnected) {
            await this.client.close();
            this.isConnected = false;
            console.log("Disconnected from the database");
        }
    }

    async checkConnectionString() {
        return this.connectionString === null ? 500 : 200;
    }

    async interface(request, type, object) {
        try {
            await this.connect();
            type = type + "";
            switch (request.toLowerCase()) {
                case "post":
                    return await this.postAny(type, object);
                case "get":
                    return await this.getAny(type, object);
                case "patch":
                    return await this.patchAny(type, object);
                case "delete":
                    return await this.deleteAny(type, object);
                default:
                    return 404;
            }
        } catch (err) {
            console.error("Request handling error:", err);
            return 500;
        }
    }

    async getAny(type, object) {
        if (await this.checkConnectionString() === 500) return;
        try {
            const cursor = await this.client.db(this.databaseName).collection(type).find(object);
            return await cursor.toArray();
        } catch (err) {
            console.error("Get operation error:", err);
            return 500;
        }
    }

    async postAny(type, object) {
        if (await this.checkConnectionString() === 500) return;
        try {
            await this.client.db(this.databaseName).collection(type).insertOne(object);
            return 200;
        } catch (err) {
            console.error("Post operation error:", err);
            return 500;
        }
    }

    async patchAny(type, object) {
        if (await this.checkConnectionString() === 500) return;
        try {
            const options = { upsert: true };
            const set = { $set: object[1] };
            await this.client.db(this.databaseName).collection(type).updateOne(object[0], set, options);
            return 200;
        } catch (err) {
            console.error("Patch operation error:", err);
            return 500;
        }
    }

    async deleteAny(type, object) {
        if (await this.checkConnectionString() === 500) return;
        try {
            await this.client.db(this.databaseName).collection(type).deleteOne(object);
            return 200;
        } catch (err) {
            console.error("Delete operation error:", err);
            return 500;
        }
    }
}

module.exports = DatabaseInterface;
