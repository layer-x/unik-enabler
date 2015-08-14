package diego_support

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/cloudfoundry/cli/plugin"
)

type DiegoSupport struct {
	cli plugin.CliConnection
}

type diegoError struct {
	Code        int64  `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
	ErrorCode   string `json:"error_code,omitempty"`
}

func NewDiegoSupport(cli plugin.CliConnection) *DiegoSupport {
	return &DiegoSupport{
		cli: cli,
	}
}

func (d *DiegoSupport) SetDiegoFlag(appGuid string, enable bool) ([]string, error) {
	output, err := d.cli.CliCommandWithoutTerminalOutput("curl", "/v2/apps/"+appGuid, "-X", "PUT", "-d", `{"diego":`+strconv.FormatBool(enable)+`}`)
	if err != nil {
		return output, err
	}

	if err = checkDiegoError(output[0]); err != nil {
		return output, err
	}

	return output, nil
}

func checkDiegoError(jsonRsp string) error {
	b := []byte(jsonRsp)
	diegoErr := diegoError{}
	err := json.Unmarshal(b, &diegoErr)
	if err != nil {
		return err
	}

	if diegoErr.ErrorCode != "" || diegoErr.Code != 0 {
		return errors.New(diegoErr.ErrorCode + " - " + diegoErr.Description)
	}

	return nil
}