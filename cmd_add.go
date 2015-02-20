package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gonuts/commander"
	"io/ioutil"
	"net/http"
	"os"
)

func make_cmd_add() *commander.Command {
	cmd_add := func(cmd *commander.Command, args []string) error {
		if len(args) != 1 {
			cmd.Usage()
			os.Exit(1)
		}
		token, err := auth()
		if err != nil {
			return err
		}

		for _, arg := range args {
			item := struct {
				Content string `json:"Content"`
			}{arg}

			var buf bytes.Buffer
			err = json.NewEncoder(&buf).Encode(&item)
			if err != nil {
				return err
			}
			req, err := http.NewRequest("POST", "https://todo.ly/api/items.json", &buf)
			if err != nil {
				return err
			}
			req.Header.Add("Token", token)
			res, err := http.DefaultClient.Do(req)
			if err != nil {
				return err
			}
			defer res.Body.Close()

			b, err := ioutil.ReadAll(res.Body)
			if err != nil {
				return err
			}

			var aerr ErrorObject
			err = json.Unmarshal(b, &aerr)
			if err != nil {
				return err
			}
			if aerr.ErrorCode != 0 {
				return fmt.Errorf("%d: %s", aerr.ErrorCode, aerr.ErrorMessage)
			}
		}
		return nil
	}

	return &commander.Command{
		Run:       cmd_add,
		UsageLine: "add [options] [name]",
		Short:     "add task",
	}
}
