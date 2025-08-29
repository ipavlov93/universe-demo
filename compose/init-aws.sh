#!/bin/bash
set -euxo pipefail

echo "Initializing SQS DLQ and applying policy..."

# Create the dead letter queue if not exist
awslocal sqs get-queue-url --queue-name DeadLetterQueue >/dev/null 2>&1 || \
awslocal sqs create-queue --queue-name DeadLetterQueue

DLQ_URL=$(awslocal sqs get-queue-url --queue-name DeadLetterQueue --query 'QueueUrl' --output text)
DLQ_ARN=$(awslocal sqs get-queue-attributes --queue-url "$DLQ_URL" --attribute-names QueueArn --query 'Attributes.QueueArn' --output text)

# Create ProductSV queue
awslocal sqs get-queue-url --queue-name ProductSV >/dev/null 2>&1 || \
awslocal sqs create-queue --queue-name ProductSV \
--attributes "{\"RedrivePolicy\":\"{\\\"deadLetterTargetArn\\\":\\\"$DLQ_ARN\\\",\\\"maxReceiveCount\\\":\\\"100\\\"}\",\"ReceiveMessageWaitTimeSeconds\":\"20\"}"

echo "SQS queues created and policy applied."
