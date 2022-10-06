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
	msgRes, err := sqs.New(sess).ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl:            &queueUrl,
		MaxNumberOfMessages: aws.Int64(1),
	})

	if err != nil {
		return nil, err
	}

	return msgRes, nil
}

func worker(tasks <-chan int, results chan<- int) {
	fmt.Println("Worker initialize...")
	for task := range tasks {
		time.Sleep(time.Second)
		fmt.Println("Done task", task)
		results <- task
	}
}

func fillTasks(tasks chan<- int) {
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

	tasks := make(chan int, numTasks)
	results := make(chan int, numTasks)

	// Note that numWorker really depends on machine specs
	const numWorkers = 1
	for i := 0; i < numWorkers; i++ {
		go worker(tasks, results)
	}

	for {
		result := <-results
		fmt.Println(result)
	}

	// msgRes, err := getMessages(sess, *urlRes.QueueUrl, 1)
	// if err != nil {
	// 	fmt.Println("Fail to get messages: ", err)
	// 	return
	// }

	// fmt.Println(msgRes)
}
