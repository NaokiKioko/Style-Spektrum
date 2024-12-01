// Load the AWS SDK for Node.js
var AWS = require("aws-sdk");
const bal = require('./catalog.js');

// Set the region
AWS.config.update({ region: process.env.REGION });


// Create an SQS service object
var sqs = new AWS.SQS({ apiVersion: "2012-11-05" });

var queueURL = process.env.AWS_SQS_URL;

const processMessages = async (dal) => {
    const params = {
        AttributeNames: ["SentTimestamp"],
        MaxNumberOfMessages: 10,
        MessageAttributeNames: ["All"],
        QueueUrl: queueURL,
        VisibilityTimeout: 20,
        WaitTimeSeconds: 10, // Long polling for efficiency
    };

    try {
        const data = await sqs.receiveMessage(params).promise();
        if (data.Messages && data.Messages.length > 0) {
            for (const message of data.Messages) {
                // console.log("Received Message:", message.Body);
                message.Body = JSON.parse(message.Body);
                if (message.Body.Topic == "Product")
                    // Process message logic here
                    await handleMessage(dal, message.Body.Product);
                    // Delete message from queue
                    const deleteParams = {
                        QueueUrl: queueURL,
                        ReceiptHandle: message.ReceiptHandle,
                    };
                    await sqs.deleteMessage(deleteParams).promise();
                    console.log("Message Deleted:", message.MessageId);
            }
        }
    } catch (err) {
        console.error("Error processing messages:", err);
    }
};

// Example message handler function
const handleMessage = async (dal, Product) => {
    // Add your custom message processing logic here
    // console.log("Handling message:", Product);
    found = await bal.GetCatalogsByFeild(dal, "URL", Product.URL)
    if ( found.length > 0) {
        console.log("Product already exists");
        return;
    }
    found = await bal.GetCatalogsByFeild(dal, "Name", Product.Name)
    // check if the found product comes from the same website
    if (!found || found.length === 0) {
        console.log("No matching products found");
        bal.PostCatalog(dal, Product);
        return;
    } else{
        ProductBase = Product.URL.split("/")[2].toLowerCase();
        for (let i = 0; i < found.length; i++) {
            FoundBase = found[i].URL.split("/")[2].toLowerCase();
            if (ProductBase == FoundBase && found[i].Price == Product.Price && CompareStringBeginings(found[i].Description, Product.Description)) {
                console.log("Product already exists");
                return;
            }
        }
        bal.PostCatalog(dal, Product);
    }
};

const CompareStringBeginings = (str1, str2) => {
    minString = str1.length > str2.length ? str2.length : str1.length;
    if (minString > 15) {
        minString = 15;
    }
    string1 = str1.substring(0, minString).toLowerCase();
    string2 = str2.substring(0, minString).toLowerCase();
    return string1.includes(string2) || string2.includes(string1);
};

module.exports = { processMessages };