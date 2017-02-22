// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"bytes"
	"crypto/md5"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var (
	root  string
	state string

	checkCmd = &cobra.Command{
		Use:   "check",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		RunE: check,
	}
)

func init() {
	RootCmd.AddCommand(checkCmd)

	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	checkCmd.PersistentFlags().StringVar(&root, "root", "", "A file glob to check for changes")
	checkCmd.PersistentFlags().StringVar(&state, "state", "", "The state to store")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

func getCurrentState(root string) (b []byte, err error) {

	matches, err := filepath.Glob(root)
	log.Printf("matches: %v", matches)

	sort.Strings(matches)

	h := md5.New()
	for _, match := range matches {
		io.WriteString(h, match)
	}
	return h.Sum(nil), err

}

func saveState(state []byte, stateFile string) (err error) {
	if err := ioutil.WriteFile(stateFile, state, 0644); err != nil {
		log.Printf("stateFile: %v", stateFile)
		log.Print(err)
	}
	return err
}

func check(cmd *cobra.Command, args []string) (err error) {

	var currentState []byte
	var lastState []byte

	if currentState, err = getCurrentState(root); err != nil {
		log.Printf("err: %#v", err)
		if err != os.ErrNotExist {
		}
		return err
	}
	log.Printf("currentState: %x", currentState)

	if lastState, err = ioutil.ReadFile(state); err != nil {
		log.Printf("err: %#v", err)
		return saveState(currentState, state)
	}

	if !bytes.Equal(lastState, currentState) {
		log.Printf("currentState:\t%x", currentState)
		log.Printf("lastState:\t%x", lastState)
		return errors.New("State does not match")
	}

	log.Print("serve called")
	// log.Printf("cmd: %v", cmd)
	log.Printf("args: %v", args)

	return err
}
