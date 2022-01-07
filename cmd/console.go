/*
Copyright © 2021 Jean-Marc Meessen, ON4KJM  <on4kjm@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"errors"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"

	"github.com/on4kjm/morserino_display/morserino"
)

// TODO(jly):
// Many things are missing here, few on top on my mind right now:
// - Signal handling: how does the user exits this CLI => let's build a context that is cancelled when a signal is received.

// consoleCmd represents the console command
var consoleCmd = &cobra.Command{
	Use:   "console",
	Short: "Displays the Morserino output to the console",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info().Msg("Starting morserino console")

		// First of all let's try to detect the ports based on the input.
		ports, err := morserino.NewPortDetector(morserinoPortName).Detect()

		if err != nil {
			log.Err(err).Msg("Failed to detect morserino devices")
			return
		}

		switch len(ports) {
		case 0:
			// If we don't have any device, then exit.
			log.Error().Msg("No morserino devices found, exiting...")
			return
		case 1:
			// If we have one, then we're good to go, let's continue.
		default:
			// if we have more than one, we can't take any action, so let's tell the user to be more specific.
			log.Error().Int("Found", len(ports)).Msg("Found an unexpected amount of devices, please use a device name")
			return
		}

		log.Info().Str("port", ports[0]).Msg("Selected port")

		// Let's attempt to open our morserino device.
		dev, err := morserino.Open(ports[0])

		if err != nil {
			log.Err(err).Msg("Unable to open morserino device")
			return
		}

		defer dev.Close()

		var (
			// An errgroup is an abstraction that allows to wait for the completion of a group of goroutines.
			// see https://pkg.go.dev/golang.org/x/sync/errgroup.
			gr, gctx = errgroup.WithContext(cmd.Context())
			console  = morserino.NewConsole()
			emitter  = morserino.NewEmitter(dev, console)
		)

		// Let'snow start our listener goroutine, that will emit events as soon as it gets a meaningful message.
		gr.Go(func() error {
			return emitter.Listen(gctx)
		})

		// To finish, let's start our console goroutine that will handle events coming from the user and from the listener.
		gr.Go(func() error {
			return console.Run(gctx)
		})

		// Let's now wait for completion of all of those goroutines, and complain in case of error.
		if err := gr.Wait(); err != nil && !errors.Is(err, morserino.ErrEmitterDone) {
			log.Err(err).Msg("Received an unexpected error, exiting")
			return
		}

		log.Info().Msg("Successfully terminated the console.")
	},
}

func init() {
	rootCmd.AddCommand(consoleCmd)

	//FIXME: the port parameter should be moved here

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// consoleCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// consoleCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
