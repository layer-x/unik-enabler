package unik_support

import (
	"encoding/json"
	"errors"

	"github.com/cloudfoundry/cli/plugin"
	"github.com/cloudfoundry/cli/plugin/models"
	"github.com/layer-x/layerx-commons/lxerrors"
)

type UnikSupport struct {
	cli plugin.CliConnection
}

type unikError struct {
	Code        int64  `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
	ErrorCode   string `json:"error_code,omitempty"`
}

func NewUnikSupport(cli plugin.CliConnection) *UnikSupport {
	return &UnikSupport{
		cli: cli,
	}
}

func (d *UnikSupport) AddUnikEnv(app plugin_models.GetAppModel, unikIp string) ([]string, error) {
	appEnv := app.EnvironmentVars
	appEnv["UNIK_IP"] = unikIp
	appEnv["UNIKERNEL_NAME"] = app.Name
	envData, err := json.Marshal(appEnv)
	if err != nil {
		return nil, lxerrors.New("could not marshal app env to json", err)
	}
	output, err := d.cli.CliCommandWithoutTerminalOutput("curl", "/v2/apps/"+app.Guid, "-X", "PUT", "-d", `{"environment_json":`+string(envData)+`}`)
	if err != nil {
		return output, err
	}
	if err = checkUnikError(output[0]); err != nil {
		return output, err
	}

	return output, nil
}

type unikVolume struct {
	Name string `json:"Name"`
	Size int `json:"Size"`
	Device string `json:"Device"`
}

func (d *UnikSupport) AddUnikEnvWithVolumes(app plugin_models.GetAppModel, unikIp string, volumeJson string) ([]string, error) {
	appEnv := app.EnvironmentVars
	var unikVol unikVolume
	err := json.Unmarshal([]byte(volumeJson), unikVol)
	if err != nil {
		return nil, lxerrors.New("invalid volumes json!", err)
	}
	appEnv["UNIK_IP"] = unikIp
	appEnv["UNIKERNEL_NAME"] = app.Name
	appEnv["UNIK_VOLUME_DATA"] = volumeJson
	envData, err := json.Marshal(appEnv)
	if err != nil {
		return nil, lxerrors.New("could not marshal app env to json", err)
	}
	output, err := d.cli.CliCommandWithoutTerminalOutput("curl", "/v2/apps/"+app.Guid, "-X", "PUT", "-d", `{"environment_json":`+string(envData)+`}`)
	if err != nil {
		return output, err
	}
	if err = checkUnikError(output[0]); err != nil {
		return output, err
	}

	return output, nil
}

func (d *UnikSupport) RemoveUnikEnv(app plugin_models.GetAppModel) ([]string, error) {
	appEnv := app.EnvironmentVars
	delete(appEnv, "UNIK_IP")
	delete(appEnv, "UNIKERNEL_NAME")
	delete(appEnv, "UNIK_VOLUME_DATA")
	envData, err := json.Marshal(appEnv)
	if err != nil {
		return nil, lxerrors.New("could not marshal app env to json", err)
	}
	output, err := d.cli.CliCommandWithoutTerminalOutput("curl", "/v2/apps/"+app.Guid, "-X", "PUT", "-d", `{"environment_json":`+string(envData)+`}`)
	if err != nil {
		return output, err
	}
	if err = checkUnikError(output[0]); err != nil {
		return output, err
	}

	return output, nil
}

func checkUnikError(jsonRsp string) error {
	b := []byte(jsonRsp)
	unikErr := unikError{}
	err := json.Unmarshal(b, &unikErr)
	if err != nil {
		return err
	}

	if unikErr.ErrorCode != "" || unikErr.Code != 0 {
		return errors.New(unikErr.ErrorCode + " - " + unikErr.Description)
	}

	return nil
}
