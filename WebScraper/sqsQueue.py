# file for all things mongoDB
import os
import json
import boto3


sqs = boto3.client('sqs')
queue_url = os.environ.get('QUEUE_URL')

def SendMessage(topic: str, message: json):
    if (message == None):
        return None
    response = sqs.send_message(
        QueueUrl=queue_url,
        MessageBody=json.dumps({
            'topic': topic,
            'message': message
        })
    )

__all__ = ["SendMessage"]