package main

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

func main() {
	awsSession := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{Region: aws.String("us-east-1")},
		// Profile: "profile_name",
	}))
	svc := cloudwatchlogs.New(awsSession)

	input := &cloudwatchlogs.DescribeLogGroupsInput{}
	cloudwatchGroups, err := svc.DescribeLogGroups(input)

	if err != nil {
		panic(err)
	}

	var cloudwatchGroupsSlice []string

	for _, group := range cloudwatchGroups.LogGroups {
		if group.RetentionInDays == nil {
			log.Printf("Adding %v to group", *group.LogGroupName)
			cloudwatchGroupsSlice = append(cloudwatchGroupsSlice, *group.LogGroupName)
		}
	}

	for {
		if cloudwatchGroups.NextToken == nil {
			log.Println("Finished retriving cloudwatchgroups")
			break
		}

		nextInput := &cloudwatchlogs.DescribeLogGroupsInput{
			NextToken: cloudwatchGroups.NextToken,
		}

		cloudwatchGroups, err = svc.DescribeLogGroups(nextInput)
		if err != nil {
			panic(err)
		}

		for _, group := range cloudwatchGroups.LogGroups {
			if group.RetentionInDays == nil {
				cloudwatchGroupsSlice = append(cloudwatchGroupsSlice, *group.LogGroupName)
			}
		}
	}

	log.Printf("Found %v groups with not retention policy", len(cloudwatchGroupsSlice))

	for _, group := range cloudwatchGroupsSlice {
		input := &cloudwatchlogs.PutRetentionPolicyInput{
			LogGroupName:    aws.String(group),
			RetentionInDays: aws.Int64(14),
		}

		_, err := svc.PutRetentionPolicy(input)
		if err != nil {
			panic(err)
		}
		log.Printf("Retention policy was set to %v", group)
	}
}
