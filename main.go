package main

import (
	"fmt"
	"time"

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
	fmt.Println("fetch messages")
	msgRes, err := sqs.New(sess).ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl:            &queueUrl,
		MaxNumberOfMessages: aws.Int64(10),
		WaitTimeSeconds:     aws.Int64(20),
	})
	fmt.Println("got", len(msgRes.Messages))

	if err != nil {
		return nil, err
	}

	return msgRes, nil
}

func worker(id int, tasks <-chan string, results chan<- string) {
	for task := range tasks {
		fmt.Println("Worker", id, "started")
		if id == 1 {
			time.Sleep(time.Second * 20)
		} else {
		}
		fmt.Println("Worker", id, "done")
		results <- task
	}
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
	fmt.Println(urlRes)

	// Make a worker pool
	const numTasks = 100

	tasks := make(chan string, numTasks)
	results := make(chan string, numTasks)

	// Note that numWorker really depends on machine specs
	const numWorkers = 3
	for i := 1; i <= numWorkers; i++ {
		go worker(i, tasks, results)
	}

	for {
		msgRes, err := getMessages(sess, *urlRes.QueueUrl, 10)
		if err != nil {
			fmt.Println("Unknown err ", err)
		}

		if len(msgRes.Messages) == 0 {
			// fmt.Println("no message")
			continue
		}

		msg := *msgRes.Messages[0].Body
		tasks <- msg

		time.Sleep(time.Second)

		// result := <-results
		// fmt.Println(result)
	}

	// msgRes, err := getMessages(sess, *urlRes.QueueUrl, 1)
	// if err != nil {
	// 	fmt.Println("Fail to get messages: ", err)
	// 	return
	// }

	// fmt.Println(msgRes)
}
