package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/daviddengcn/go-colortext"
	"github.com/gonuts/commander"
	"io/ioutil"
	"net/http"
)

var priorityColor = map[int64]ct.Color{
	1: ct.Red,
	2: ct.Blue,
	3: ct.Green,
	4: ct.White,
}

func make_cmd_list() *commander.Command {
	cmd_list := func(cmd *commander.Command, args []string) error {
		token, err := auth()
		if err != nil {
			return err
		}

		req, err := http.NewRequest("GET", "https://todo.ly/api/items.json", nil)
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

		var items []ItemObject
		err = json.Unmarshal(b, &items)
		if err != nil {
			return err
		}
		for _, item := range items {
			ct.ChangeColor(ct.Magenta, false, ct.None, false)
			fmt.Printf("%8d", item.Id)
			ct.ResetColor()
			fmt.Print(" ")
			if item.Checked {
				ct.ChangeColor(ct.Red, false, ct.None, false)
				fmt.Print("✕")
				ct.ResetColor()
			} else {
				ct.ChangeColor(ct.Green, false, ct.None, false)
				fmt.Print("✓")
				ct.ResetColor()
			}
			fmt.Print(" ")
			if pc, ok := priorityColor[item.Priority]; ok {
				ct.ChangeColor(pc, false, ct.None, false)
			}
			fmt.Print(item.Content)
			ct.ResetColor()
			fmt.Print(" ")
			ct.ChangeColor(ct.Black, true, ct.None, false)
			fmt.Print(item.CreatedDate.Format("2006/01/02 15:04:05"))
			ct.ResetColor()
			fmt.Println()
		}
		return nil
	}

	return &commander.Command{
		Run:       cmd_list,
		UsageLine: "list [options]",
		Short:     "show list index",
	}
}
