package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/gonuts/commander"
)

func init() {
	makecmd := func(check bool) func(cmd *commander.Command, args []string) error {
		return func(cmd *commander.Command, args []string) error {
			if len(args) != 1 {
				cmd.Usage()
				os.Exit(1)
			}
			token, err := auth()
			if err != nil {
				return err
			}

			for _, arg := range args {
				i, err := strconv.Atoi(arg)
				if err != nil {
					return err
				}
				item := struct {
					Checked bool `json:"Checked"`
				}{check}

				var buf bytes.Buffer
				err = json.NewEncoder(&buf).Encode(&item)
				if err != nil {
					return err
				}
				req, err := http.NewRequest("POST", fmt.Sprintf("https://todo.ly/api/items/%d.json", i), &buf)
				if err != nil {
					return err
				}
				req.Header.Add("Token", token)
				res, err := http.DefaultClient.Do(req)
				if err != nil {
					return err
				}
				defer res.Body.Close()

				if res.StatusCode != 200 {
					return errors.New(res.Status)
				}
				b, err := ioutil.ReadAll(res.Body)
				if err != nil {
					return err
				}

				var aerr ErrorObject
				err = json.Unmarshal(b, &aerr)
				if err == nil {
					if aerr.ErrorCode != 0 {
						return fmt.Errorf("%d: %s", aerr.ErrorCode, aerr.ErrorMessage)
					}
				}
			}
			return nil
		}
	}

	commander.Defaults.Subcommands = append(commander.Defaults.Subcommands, &commander.Command{
		Run:       makecmd(true),
		UsageLine: "check [options] [id]",
		Short:     "check the todo",
	})

	commander.Defaults.Subcommands = append(commander.Defaults.Subcommands, &commander.Command{
		Run:       makecmd(false),
		UsageLine: "uncheck [options] [id]",
		Short:     "uncheck the todo",
	})
}
