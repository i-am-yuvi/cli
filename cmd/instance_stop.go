package cmd

import (
	"fmt"
	"github.com/briandowns/spinner"
	"os"
	"time"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

var waitStop bool
var instanceStopCmd = &cobra.Command{
	Use:     "stop",
	Short:   "Stop an instance",
	Aliases: []string{"shutdown"},
	Long: `Pull the power from the specified instance by part of the ID or name.
If you wish to use a custom format, the available fields are:

	* ID
	* Hostname

Example: civo instance stop ID/NAME`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			fmt.Printf("Unable to create a Civo API Client: %s\n", aurora.Red(err))
			os.Exit(1)
		}

		instance, err := client.FindInstance(args[0])
		if err != nil {
			fmt.Printf("Finding instance: %s\n", aurora.Red(err))
			os.Exit(1)
		}

		_, err = client.StopInstance(instance.ID)
		if err != nil {
			fmt.Printf("Stopping instance: %s\n", aurora.Red(err))
			os.Exit(1)
		}

		if waitStop == true {

			stillStopping := true
			s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
			s.Suffix = " Stopping instance..."
			s.Start()

			for stillStopping {
				instanceCheck, _ := client.FindInstance(instance.ID)
				if instanceCheck.Status == "SHUTOFF" {
					stillStopping = false
					s.Stop()
				}
				time.Sleep(5 * time.Second)
			}
		}

		if outputFormat == "human" {
			fmt.Printf("The instance %s (%s) is being stopped\n", aurora.Green(instance.Hostname), instance.ID)
		} else {
			ow := utility.NewOutputWriter()
			ow.StartLine()
			ow.AppendData("ID", instance.ID)
			ow.AppendData("Hostname", instance.Hostname)
			if outputFormat == "json" {
				ow.WriteSingleObjectJSON()
			} else {
				ow.WriteCustomOutput(outputFields)
			}
		}
	},
}
