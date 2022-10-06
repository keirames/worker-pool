package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func getQueueUrl(sess *session.Session, queue string) (*sqs.GetQueueUrlOutput, error) {
	result, err := sqs.New(sess).GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: &queue,
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func getMessages(sess *session.Session, queueUrl string, maxMessages int) (*sqs.ReceiveMessageOutput, error) {
	msgRes, err := sqs.New(sess).ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl:            &queueUrl,
		MaxNumberOfMessages: aws.Int64(1),
	})

	if err != nil {
		return nil, err
	}

	return msgRes, nil
}

func main() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Profile: "default",
		Config: aws.Config{
			Region: aws.String("ap-southeast-1"),
		},
	}))

	queueName := "tasks"
	urlRes, err := getQueueUrl(sess, queueName)
	if err != nil {
		fmt.Println("Fail to get url: ", err)
		return
	}

	msgRes, err := getMessages(sess, *urlRes.QueueUrl, 1)
	if err != nil {
		fmt.Println("Fail to get messages: ", err)
		return
	}

	fmt.Println(msgRes)
}
