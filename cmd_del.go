package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gonuts/commander"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

func init() {
	commander.Defaults.Subcommands = append(commander.Defaults.Subcommands, &commander.Command{
		Run: func(cmd *commander.Command, args []string) error {
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
				req, err := http.NewRequest("DELETE", fmt.Sprintf("https://todo.ly/api/%d.json", i), nil)
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
					return errors.New("todo not found")
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
		},
		UsageLine: "del [options] [id]",
		Short:     "del todo",
	})
}
