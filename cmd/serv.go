/*
Copyright © 2023 Chris Berthiaume, University of Washington <chrisbee@uw.edu>

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
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"

	"ingestwebhook/serve"
)

var (
	addr string
)

// servCmd represents the serv
var servCmd = &cobra.Command{
	Use:   "serv",
	Short: "Oceanographic cruise data ingest minio webhook server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(os.Stderr, "ingestwebhook version %v\n", version)
		fmt.Fprintf(os.Stderr, "starting server at %v\n", addr)
		err := serve.Start(addr)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintln(os.Stderr, "closing server")
	},
}

func init() {
	rootCmd.AddCommand(servCmd)
	servCmd.PersistentFlags().StringVarP(&addr, "address", "a", ":9010", "server bind address")
}
