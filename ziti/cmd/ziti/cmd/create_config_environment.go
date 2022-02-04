/*
	Copyright NetFoundry, Inc.

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

	https://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

package cmd

import (
	_ "embed"
	"fmt"
	cmdhelper "github.com/openziti/ziti/ziti/cmd/ziti/cmd/helpers"
	"github.com/openziti/ziti/ziti/cmd/ziti/cmd/templates"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"
)

var (
	createConfigEnvironmentLong = templates.LongDesc(`
		Displays available environment variable manual overrides
`)

	createConfigEnvironmentExample = templates.Examples(`
		# Display environment variables and their values 
		ziti create config environment

		# Print an environment file to the console
		ziti create config environment --output stdout
	`)
)

//go:embed config_templates/environment.yml
var environmentConfigTemplate string

// CreateConfigEnvironmentOptions the options for the create environment command
type CreateConfigEnvironmentOptions struct {
	CreateConfigOptions
	ConfigTemplateValues
	cmdhelper.EnvVariables
	output string
}

// NewCmdCreateConfigEnvironment creates a command object for the "environment" command
func NewCmdCreateConfigEnvironment(data *ConfigTemplateValues) *cobra.Command {
	environmentOptions := &CreateConfigEnvironmentOptions{
		ConfigTemplateValues: *data,
		EnvVariables:         cmdhelper.EnvVariableNames,
	}

	cmd := &cobra.Command{
		Use:     "environment",
		Short:   "Display config environment variables",
		Aliases: []string{"env"},
		Long:    createConfigEnvironmentLong,
		Example: createConfigEnvironmentExample,
		PreRun: func(cmd *cobra.Command, args []string) {
			// Setup logging
			var logOut *os.File
			environmentOptions.OS = runtime.GOOS
			if environmentOptions.Verbose {
				logrus.SetLevel(logrus.DebugLevel)
				// Only print log to stdout if not printing config to stdout
				if strings.ToLower(environmentOptions.Output) != "stdout" {
					logOut = os.Stdout
				} else {
					logOut = os.Stderr
				}
				logrus.SetOutput(logOut)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			environmentOptions.Cmd = cmd
			environmentOptions.Args = args
			err := environmentOptions.run()
			cmdhelper.CheckErr(err)
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			// Reset log output after run completes
			logrus.SetOutput(os.Stdout)
		},
	}

	createConfigLong := fmt.Sprintf("Creates a config file for specified Ziti component using environment variables which have default values but can be manually set to override the config output.\n\n"+
		"The following environment variables can be set to override config values (current value is displayed):\n"+
		"%-36s %s\n"+
		"%-36s %s\n"+
		"%-36s %s\n"+
		"%-36s %s\n"+
		"%-36s %s\n"+
		"%-36s %s\n"+
		"%-36s %s\n"+
		"%-36s %s\n"+
		"%-36s %s\n"+
		"%-36s %s\n"+
		"%-36s %s\n"+
		"%-36s %s\n"+
		"%-36s %s\n"+
		"%-36s %s\n"+
		"%-36s %s\n"+
		"%-36s %s\n"+
		"%-36s %s\n"+
		"%-36s %s",
		cmdhelper.EnvVariableNames.ZitiHomeVarName, data.ZitiHome,
		cmdhelper.EnvVariableNames.ZitiCtrlListenerHostPortVarName, data.Controller.ListenerHostPort,
		cmdhelper.EnvVariableNames.ZitiCtrlMgmtListenerHostPortVarName, data.Controller.MgmtListenerHostPort,
		cmdhelper.EnvVariableNames.ZitiCtrlNameVarName, data.Controller.Name,
		cmdhelper.EnvVariableNames.ZitiEdgeCtrlListenerHostPortVarName, data.Controller.Edge.ListenerHostPort,
		cmdhelper.EnvVariableNames.ZitiEdgeCtrlAdvertisedHostPortVarName, data.Controller.Edge.AdvertisedHostPort,
		cmdhelper.EnvVariableNames.ZitiEdgeRouterHostnameVarName, data.Router.Edge.Hostname,
		cmdhelper.EnvVariableNames.ZitiEdgeRouterPortVarName, data.Router.Edge.Port,
		cmdhelper.EnvVariableNames.ZitiCtrlIdentityCertVarName, data.Controller.IdentityCert,
		cmdhelper.EnvVariableNames.ZitiCtrlIdentityServerCertVarName, data.Controller.IdentityServerCert,
		cmdhelper.EnvVariableNames.ZitiCtrlIdentityKeyVarName, data.Controller.IdentityKey,
		cmdhelper.EnvVariableNames.ZitiCtrlIdentityCAVarName, data.Controller.IdentityCA,
		cmdhelper.EnvVariableNames.ZitiEdgeCtrlIdentityCertVarName, data.Controller.Edge.IdentityCert,
		cmdhelper.EnvVariableNames.ZitiEdgeCtrlIdentityServerCertVarName, data.Controller.Edge.IdentityServerCert,
		cmdhelper.EnvVariableNames.ZitiEdgeCtrlIdentityKeyVarName, data.Controller.Edge.IdentityKey,
		cmdhelper.EnvVariableNames.ZitiEdgeCtrlIdentityCAVarName, data.Controller.Edge.IdentityCA,
		cmdhelper.EnvVariableNames.ZitiSigningCertVarName, data.ZitiSigningCert,
		cmdhelper.EnvVariableNames.ZitiSigningKeyVarName, data.ZitiSigningKey)

	cmd.Long = createConfigLong

	environmentOptions.addCreateFlags(cmd)
	environmentOptions.addFlags(cmd)

	return cmd
}

func (options *CreateConfigEnvironmentOptions) addFlags(cmd *cobra.Command) {
	// cmd.Flags().StringVar(&options.CtrlListener, optionCtrlListener, constants.DefaultZitiControllerListenerHostPort, "address and port of the config controller listener")
	// cmd.Flags().StringVar(&options.DatabaseFile, optionDatabaseFile, "ctrl.db", "location of the database file")
	// cmd.Flags().StringVar(&options.MgmtListener, optionMgmtListener, constants.DefaultZitiMgmtControllerListenerHostPort, "address and port of the config management listener")
}

// run implements the command
func (options *CreateConfigEnvironmentOptions) run() error {

	tmpl, err := template.New("environment-config").Parse(environmentConfigTemplate)
	if err != nil {
		return err
	}

	var f *os.File
	if strings.ToLower(options.Output) != "stdout" {
		// Check if the path exists, fail if it doesn't
		basePath := filepath.Dir(options.Output) + "/"
		if _, err := os.Stat(filepath.Dir(basePath)); os.IsNotExist(err) {
			logrus.Fatalf("Provided path: [%s] does not exist\n", basePath)
			return err
		}

		f, err = os.Create(options.Output)
		logrus.Debugf("Created output file: %s", options.Output)
		if err != nil {
			return errors.Wrapf(err, "unable to create config file: %s", options.Output)
		}
	} else {
		f = os.Stdout
	}
	defer func() { _ = f.Close() }()

	if err := tmpl.Execute(f, options); err != nil {
		return errors.Wrap(err, "unable to execute template")
	}

	logrus.Debugf("Environment configuration file generated successfully and written to: %s", options.Output)

	return nil
}
