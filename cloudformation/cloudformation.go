package cloudformation

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
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

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	filepath := path.Join(dir, filename)
	fmt.Printf("Reading exceptions from %s\n", filepath)

	cont, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}

	raw := strings.Split(string(cont), "\n")
	output := make([]string, 0)

	for _, value := range raw {
		if value != "" {
			output = append(output, value)
		}
	}

	return output
}

func CleanStacks(dryRun bool) {

	if dryRun {
		fmt.Println("Dry run is active, not stack will be deleted ...")
	}

	profile := "mi"
	regions := [...]string{
		"eu-central-1",
		"us-west-2",
	}

	exceptions := loadExceptions()
	fmt.Printf("Stack exceptions:\n")
	for _, stackName := range exceptions {
		fmt.Printf(" * %s\n", stackName)
	}

	for _, region := range regions {
		fmt.Printf("Cleaning stacks in region: %s\n", region)
		CleanRegion(profile, region, exceptions, dryRun)
	}

}

func CleanRegion(profile string, region string, exceptions []string, dryRun bool) {

	sess, err := session.NewSessionWithOptions(session.Options{
		Profile: profile,
		Config:  aws.Config{Region: aws.String(region)},
	})

	if err != nil {
		fmt.Println("failed")
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

		if checkDelete(*stack, exceptions) {
			fmt.Println(" * deleting", *stack.StackName)

			if !dryRun {
				time.Sleep(3000 * time.Millisecond)

				deleteInput := cloudformation.DeleteStackInput{
					StackName: stack.StackName,
				}

				_, err := cfn.DeleteStack(&deleteInput)
				if err != nil {
					fmt.Println("  FAILED to delete stack", *stack.StackName, err)
				}
			}

		}
	}

}
