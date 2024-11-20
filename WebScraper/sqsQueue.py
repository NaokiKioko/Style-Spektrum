import os
import logging
import boto3
from botocore.exceptions import ClientError

logger = logging.getLogger(__name__)
session = boto3.Session(profile_name="default")  
sqs = session.client('sqs')

def get_queue_url(name):
    try:
        response = sqs.get_queue_url(QueueName=name)
        queue_url = response['QueueUrl']
        logger.info("Got queue '%s' with URL=%s", name, queue_url)
    except ClientError as error:
        logger.exception("Couldn't get queue named %s.", name)
        raise error
    else:
        return queue_url


def send_message(name, message_body, message_attributes=None):
    queue_url = get_queue_url(name)
    if not message_attributes:
        message_attributes = {}

    try:
        response = sqs.send_message(
            QueueUrl=queue_url,
            MessageBody=message_body,
            MessageAttributes=message_attributes
        )
        logger.info("Sent message to queue '%s'. Message ID: %s", name, response['MessageId'])
    except ClientError as error:
        logger.exception("Send message failed: %s", message_body)
        raise error
    else:
        return response

__all__ = ["send_message"]
