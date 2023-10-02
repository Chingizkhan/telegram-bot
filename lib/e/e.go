package e

import "fmt"

// todo: in rabbitmq realize the same logic with defer and log

func Wrap(msg string, err error) error {
	if err == nil {
		return nil
	}

	return fmt.Errorf("%s: %w", msg, err)
}
