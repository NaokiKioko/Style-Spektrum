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
                console.log("Received Message:", message.Body);
                if (message.Body.Topic == "Product")
                    // Process message logic here
                    await handleMessage(dal, message.Body);
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
const handleMessage = async (dal, body) => {
    // Add your custom message processing logic here
    console.log("Handling message:", body);
    found = await bal.GetReportsByField(dal, body.Product.URL, "URL")
    if ( found.length > 0) {
        console.log("Product already exists");
        return;
    } else{
        bal.PostCatalog(dal, body.Product);
    }

};

module.exports = { processMessages };