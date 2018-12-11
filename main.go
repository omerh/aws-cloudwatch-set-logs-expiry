package main

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

func main() {
	// Deafult region
	region := "us-east-1"
	// Check if region is overide with environment variable
	awsRegion, ok := os.LookupEnv("AWS_REGION")
	if ok {
		region = awsRegion
	}
	// initialize aws session
	awsSession := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{Region: aws.String(region)},
		// Profile: "profile_name",
	}))
	svc := cloudwatchlogs.New(awsSession)

	// default empty input beofre retriving next tokens
	input := &cloudwatchlogs.DescribeLogGroupsInput{}
	// get all cloudwatch log groups
	cloudwatchGroups, err := svc.DescribeLogGroups(input)

	if err != nil {
		panic(err)
	}

	// create a slice of for having all cloudwatch group names which retention policy is set to never
	var cloudwatchGroupsSlice []string

	// iterate and add all groups with retention is set to never
	for _, group := range cloudwatchGroups.LogGroups {
		if group.RetentionInDays == nil {
			log.Printf("Adding %v to group", *group.LogGroupName)
			cloudwatchGroupsSlice = append(cloudwatchGroupsSlice, *group.LogGroupName)
		}
	}

	// run for loop on all cloudwatch groups according to next token, aws api respone will have 50 entries per response
	for {
		// check if retrival is completed
		if cloudwatchGroups.NextToken == nil {
			log.Println("Finished retriving cloudwatchgroups")
			break
		}

		// set input filter with the next token from aws api response
		nextInput := &cloudwatchlogs.DescribeLogGroupsInput{
			NextToken: cloudwatchGroups.NextToken,
		}

		// retrive next
		cloudwatchGroups, err = svc.DescribeLogGroups(nextInput)
		if err != nil {
			panic(err)
		}

		// iterate on next and add all groups with retention is set to never
		for _, group := range cloudwatchGroups.LogGroups {
			if group.RetentionInDays == nil {
				log.Printf("Adding %v to group", *group.LogGroupName)
				cloudwatchGroupsSlice = append(cloudwatchGroupsSlice, *group.LogGroupName)
			}
		}
	}

	log.Printf("Found %v groups with not retention policy", len(cloudwatchGroupsSlice))

	// iterate and set retention policy to groups
	for _, group := range cloudwatchGroupsSlice {
		// set input filter
		input := &cloudwatchlogs.PutRetentionPolicyInput{
			LogGroupName:    aws.String(group),
			RetentionInDays: aws.Int64(14),
		}

		// put retention policy
		_, err := svc.PutRetentionPolicy(input)
		if err != nil {
			panic(err)
		}
		log.Printf("Retention policy was set to %v", group)
	}
}
