package user_input

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"tadpoles-backup/config"
	"tadpoles-backup/internal/utils"
	"tadpoles-backup/pkg/headings"

	"golang.org/x/crypto/ssh/terminal"
)

func GetResetCode() (string, string) {
	var resetCode string
	var newPassword string

	if config.HasEnvCreds() {
		resetCode = config.EnvResetCode
		newPassword = config.EnvNewPassword
	} else if config.IsInteractive() {
		resetCode, newPassword = cliResetCodeNewPassword()
	} else {
		utils.CmdFailed(
			errors.New("credentials must be supplied from the environment if running in non-interactive mode"),
		)
	}

	return resetCode, newPassword
}

// get username and password from user user_input
func cliResetCodeNewPassword() (string, string) {
	utils.WriteInfo(
		"Input",
		fmt.Sprintf("%s email sent, enter reset-code and new-password", config.Provider.String()),
	)
	reader := bufio.NewReader(os.Stdin)

	utils.WriteInfo("Reset Code", "", headings.NoNewLine)
	resetCode, _ := reader.ReadString('\n')

	utils.WriteInfo("Password", "", headings.NoNewLine)
	bytePassword, _ := terminal.ReadPassword(int(os.Stdin.Fd()))
	newPassword := string(bytePassword)
	fmt.Println()

	return strings.TrimSpace(resetCode), strings.TrimSpace(newPassword)
}
