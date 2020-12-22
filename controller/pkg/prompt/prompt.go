package prompt

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"os"
)

func InputWithDefault(target *string, label string, defaultValue string) error {
	if *target != "" {
		return nil
	}

	prompt := &survey.Input{
		Message: label,
		Default: defaultValue,
	}

	err := survey.AskOne(prompt, target)
	if err == terminal.InterruptErr {
		os.Exit(0)
	}
	if err != nil {
		return err
	}

	return nil
}

func InputPassword(target *string, label string) error {
	if *target != "" {
		return nil
	}

	prompt := &survey.Password{
		Message: label,
	}

	err := survey.AskOne(prompt, target)
	if err == terminal.InterruptErr {
		os.Exit(0)
	}
	if err != nil {
		return err
	}

	return nil
}

func Input(target *string, label string) error {
	if *target != "" {
		return nil
	}

	prompt := &survey.Input{
		Message: label,
	}

	err := survey.AskOne(prompt, target)
	if err == terminal.InterruptErr {
		os.Exit(0)
	}
	if err != nil {
		return err
	}

	return nil
}

func Confirm(label string) (error, bool) {
	result := false

	prompt := &survey.Confirm{
		Message: label,
	}

	err := survey.AskOne(prompt, &result)
	if err == terminal.InterruptErr {
		os.Exit(0)
	}
	if err != nil {
		return err, false
	}

	return nil, result
}

func Multiselect(items []string, label string) (error, []string) {
	result := make([]string, 0)

	var multiQs = []*survey.Question{
		{
			Name: "letter",
			Prompt: &survey.MultiSelect{
				Message: "Choose one or more words :",
				Options: items,
			},
		},
	}

	err := survey.Ask(multiQs, &result)
	if err == terminal.InterruptErr {
		os.Exit(0)
	}
	if err != nil {
		return err, nil
	}

	return nil, result
}

func Select(items []string, label string) (error, string) {
	result := ""

	prompt := &survey.Select{
		Options: items,
		Message: label,
	}

	err := survey.AskOne(prompt, &result)
	if err == terminal.InterruptErr {
		os.Exit(0)
	}
	if err != nil {
		return err, ""
	}

	return nil, result
}
