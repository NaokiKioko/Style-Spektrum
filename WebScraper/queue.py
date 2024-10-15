# file for all things mongoDB
import os
import json
import boto3


sqs = boto3.client('sqs')
queue_url = os.environ['QUEUE_URL']

def SendMessage(topic: str, message: json):
    response = sqs.send_message(
        QueueUrl=queue_url,
        MessageBody=json.dumps({
            'topic': topic,
            'message': message
        })
    )
    return response

__all__ = ["SendMessage"]