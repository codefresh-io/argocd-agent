package prompt

import (
	"fmt"
	"github.com/manifoldco/promptui"
)

func InputWithDefault(target *string, label string, defaultValue string) error {
	if *target != "" {
		return nil
	}

	prompt := promptui.Prompt{
		Label: fmt.Sprintf("%s, (default: %s)", label, defaultValue),
	}

	result, err := prompt.Run()
	if err != nil {
		return err
	}

	if result == "" {
		result = defaultValue
	}

	*target = result

	return nil
}
