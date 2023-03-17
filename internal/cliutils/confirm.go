package cliutils

import (
	"fmt"
	"log"
	"os"
	"strings"

	"gitlab.com/zephinzer/go-devops"
)

func Confirm(text string) error {
	questionText := strings.Trim(
		fmt.Sprintf(`
%s

confirm?`,
			text,
		), "\n ")
	if confirmed, err := devops.Confirm(devops.ConfirmOpts{
		MatchExact: "yes",
		Output:     os.Stderr,
		Question:   questionText,
	}); err != nil {
		log.Fatalf("failed to request for confirmation: %s", err)
		return err
	} else if !confirmed {
		log.Fatalf("refusing to continue")
		return fmt.Errorf("user said no")
	}

	return nil
}
