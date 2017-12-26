package cloudformation

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

func CleanStacks() {

	profile := "mi"

	sess, err := session.NewSessionWithOptions(session.Options{
		Profile: profile,
	})

	if err != nil {
		fmt.Println("fuck")
	}

	cfn := cloudformation.New(sess)

	stacks, err = cfn.ListStacks(cloudformation.ListStackInput{})
	if err != nil {
		fmt.Println("fuck")
	}

	fmt.Println(stacks)

}
