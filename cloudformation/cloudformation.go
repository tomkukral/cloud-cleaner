package cloudformation

import (
	"fmt"
	"regexp"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func checkDelete(stack cloudformation.Stack, exceptions []string) bool {
	ageHours := 24

	// check stack in exceptions
	exceptionMatches := stringInSlice(*stack.StackName, exceptions)

	// check for age
	now := time.Now()
	age := now.Sub(*stack.CreationTime)
	ageMatches := age.Hours() >= float64(ageHours)

	return ageMatches && !exceptionMatches
}

func loadExceptions() []string {
	filename := "stack_exceptions"
	cont, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err)
	} else {
		return strings.Split(string(cont), "\n")
	}

	return []string{}
}

func CleanStacks(dryRun bool) {

	if dryRun {
		fmt.Println("Dry run is active, not stack will be deleted ...")
	}

	profile := "mi"
	regions := [2]string{
		"eu-central-1",
		"us-west-2",
	}

	exceptions := loadExceptions()
	fmt.Println(exceptions)

	for _, region := range regions {
		fmt.Println(region)
		CleanRegion(profile, region, exceptions, dryRun)
	}

}

func CleanRegion(profile string, region string, exceptions []string, dryRun bool) {

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
