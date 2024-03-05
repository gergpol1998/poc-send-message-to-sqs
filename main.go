package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/joho/godotenv"
)

type MessageRecord struct {
	API_Key string `json:"api_key"`
	Reason  string `json:"reason,omitempty"`
	Status  string `json:"status"`
}

func GetQueueURL(sess *session.Session, queue *string) (*sqs.GetQueueUrlOutput, error) {
	// Create an SQS service client
	svc := sqs.New(sess)

	result, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: queue,
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func SendMsg(sess *session.Session, queueURL *string, data *MessageRecord) error {
	// Create an SQS service client
	// snippet-start:[sqs.go.send_message.call]
	svc := sqs.New(sess)

	// Marshal the struct into a JSON string
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = svc.SendMessage(&sqs.SendMessageInput{
		MessageBody: aws.String(string(jsonData)),
		QueueUrl:    queueURL,
	})
	// snippet-end:[sqs.go.send_message.call]
	if err != nil {
		return err
	}

	return nil
}

func MessageToSQS(data *MessageRecord) error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get queueName
	queueName := os.Getenv("QUEUE_NAME")
	accessKey := os.Getenv("ACCESS_KEY")
	secretKey := os.Getenv("SECRET_KEY")

	// Check if the queueName is empty
	if queueName == "" {
		fmt.Printf("You must supply the name of a queue (-q QUEUE)")
		return nil
	}

	// Specify your AWS credentials directly
	awsConfig := aws.Config{
		Region:      aws.String("ap-southeast-1"),
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
	}

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config:            awsConfig,
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Get URL of queue
	result, err := GetQueueURL(sess, &queueName)
	if err != nil {
		fmt.Printf("Error getting the queue URL: %v", err)
		return nil
	}

	queueURL := result.QueueUrl

	// Send message to SQS
	err = SendMsg(sess, queueURL, data)
	if err != nil {
		fmt.Printf("Error sending the message: %v", err)
		return nil
	}

	fmt.Println("Sent message to queue")
	return nil
}

func main() {
	// Send message to sqs
	messageData := &MessageRecord{API_Key: "test3", Reason: "good", Status: "PASS"}
	MessageToSQS(messageData)

}
