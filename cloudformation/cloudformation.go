package cloudformation

import (
	"fmt"
	"regexp"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

func checkDelete(stack cloudformation.Stack) bool {
	nameRegexp := "^jenkins-"
	ageHours := 24

	// check name
	nameMatches, err := regexp.MatchString(nameRegexp, *stack.StackName)
	if err != nil {
		fmt.Println("Failed to match stack name:", err)
		return false
	}

	// check for age
	now := time.Now()
	age := now.Sub(*stack.CreationTime)
	ageMatches := age.Hours() >= float64(ageHours)

	return nameMatches && ageMatches
}

func CleanStacks() {

	profile := "mi"
	region := "eu-central-1"

	sess, err := session.NewSessionWithOptions(session.Options{
		Profile: profile,
		Config:  aws.Config{Region: aws.String(region)},
	})

	if err != nil {
		fmt.Println("fuck")
		return
	}

	cfn := cloudformation.New(sess)

	res, err := cfn.DescribeStacks(&cloudformation.DescribeStacksInput{})
	if err != nil {
		fmt.Println("Error", err)
		return
	}

	// delete stacks matching: jenkins-*, older than 1 day

	for i := 0; i < len(res.Stacks); i++ {
		stack := res.Stacks[i]

		if checkDelete(*stack) {
			fmt.Println("Deleting", *stack.StackName)

			deleteInput := cloudformation.DeleteStackInput{
				StackName: stack.StackName,
			}

			_, err := cfn.DeleteStack(&deleteInput)
			if err != nil {
				fmt.Println("Failed to delete stack", *stack.StackName, err)
			}

		}
	}

}
