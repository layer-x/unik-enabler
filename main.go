package main

import (
	"fmt"
	"os"

	"github.com/layer-x/Unik-Enabler/unik_support"
	"github.com/cloudfoundry/cli/plugin"
)

type UnikEnabler struct{}

func (c *UnikEnabler) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "Unik-Enabler",
		Version: plugin.VersionType{
			Major: 0,
			Minor: 0,
			Build: 0,
		},
		Commands: []plugin.Command{
			{
				Name:     "enable-unik",
				HelpText: "enable Unik support for an app",
				UsageDetails: plugin.Usage{
					Usage: "cf enable-unik APP_NAME UNIK_IP VOLUME_DATA",
				},
			},
			{
				Name:     "disable-unik",
				HelpText: "disable Unik support for an app",
				UsageDetails: plugin.Usage{
					Usage: "cf disable-unik APP_NAME",
				},
			},
			{
				Name:     "has-unik-enabled",
				HelpText: "Check if Unik support is enabled for an app",
				UsageDetails: plugin.Usage{
					Usage: "cf has-unik-enabled APP_NAME",
				},
			},
		},
	}
}

func main() {
	plugin.Start(new(UnikEnabler))
}

func (c *UnikEnabler) Run(cliConnection plugin.CliConnection, args []string) {
	if args[0] == "enable-unik" && len(args) == 3 {
		c.enableUnikSupport(cliConnection, args[1], args[2])
	} else if args[0] == "enable-unik" && len(args) == 4 {
		c.toggleUnikSupportWithVolumes(true, cliConnection, args[1], args[2], args[3])
	} else if args[0] == "disable-unik" && len(args) == 2 {
		c.enableUnikSupport(false, cliConnection, args[1])
	} else if args[0] == "has-unik-enabled" && len(args) == 2 {
		c.isUnikEnabled(cliConnection, args[1])
	} else {
		c.showUsage(args)
	}
}

func (c *UnikEnabler) showUsage(args []string) {
	for _, cmd := range c.GetMetadata().Commands {
		if cmd.Name == args[0] {
			fmt.Println("Invalid Usage: \n", cmd.UsageDetails.Usage)
		}
	}
}

func (c *UnikEnabler) enableUnikSupport(cliConnection plugin.CliConnection, appName, unikIp string) {
	d := unik_support.NewUnikSupport(cliConnection)

	fmt.Printf("Enabling Unik support for app %s using Unik Backend at %s, with no volumes\n", appName, unikIp)
	app, err := cliConnection.GetApp(appName)
	if err != nil {
		exitWithError(err, []string{})
	}

	if output, err := d.AddUnikEnv(app, unikIp); err != nil {
		fmt.Println("err 1", err, output)
		exitWithError(err, output)
	}
	sayOk()

	fmt.Printf("Verifying %s Unik support is enabled\n", appName)
	app, err = cliConnection.GetApp(appName)
	if err != nil {
		exitWithError(err, []string{})
	}

	if _, ok := app.EnvironmentVars["UNIK_IP"]; ok {
		sayOk()
	} else {
		sayFailed()
		fmt.Printf("Unik support for %s is NOT enabled\n\n", appName)
		os.Exit(1)
	}
}func (c *UnikEnabler) disableUnikSupport(cliConnection plugin.CliConnection, appName string) {
	d := unik_support.NewUnikSupport(cliConnection)

	fmt.Printf("Disabling Unik support for app %s\n", appName)
	app, err := cliConnection.GetApp(appName)
	if err != nil {
		exitWithError(err, []string{})
	}

	if output, err := d.RemoveUnikEnv(app); err != nil {
		fmt.Println("err 1", err, output)
		exitWithError(err, output)
	}
	sayOk()

	fmt.Printf("Verifying %s Unik support is disabled\n", appName)
	app, err = cliConnection.GetApp(appName)
	if err != nil {
		exitWithError(err, []string{})
	}

	if _, ok := app.EnvironmentVars["UNIK_IP"]; !ok {
		sayOk()
	} else {
		sayFailed()
		fmt.Printf("Unik support for %s is STILL enabled\n\n", appName)
		os.Exit(1)
	}
}

func (c *UnikEnabler) isUnikEnabled(cliConnection plugin.CliConnection, appName string) {
	app, err := cliConnection.GetApp(appName)
	if err != nil {
		exitWithError(err, []string{})
	}

	if app.Guid == "" {
		sayFailed()
		fmt.Printf("App %s not found\n\n", appName)
		os.Exit(1)
	}

	if _, ok := app.EnvironmentVars["UNIK_IP"]; ok {
		fmt.Println("true")
	} else {
		fmt.Println("false")
	}
}

func exitWithError(err error, output []string) {
	sayFailed()
	fmt.Println("Error: ", err)
	for _, str := range output {
		fmt.Println(str)
	}
	os.Exit(1)
}

func say(message string, color uint, bold int) string {
	return fmt.Sprintf("\033[%d;%dm%s\033[0m", bold, color, message)
}

func sayOk() {
	fmt.Println(say("Ok\n", 32, 1))
}

func sayFailed() {
	fmt.Println(say("FAILED", 31, 1))
}
