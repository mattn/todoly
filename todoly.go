package main

import (
	"encoding/json"
	"fmt"
	"github.com/fhs/go-netrc/netrc"
	"github.com/gonuts/commander"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"time"
)

var _ io.Writer
var _ *os.File

var dateRe = regexp.MustCompile(`^"\\/Date\(([0-9]+)\)\\/"$`)

type ErrorObject struct {
	ErrorMessage string `json:"ErrorMessage"`
	ErrorCode    int64  `json:"ErrorCode"`
}

type ItemType int

const (
	CheckItem ItemType = iota + 1
	ProjectItem
	DoneItem
	FilterItem
	PlaceholderItem
	PlaceholderDoneItem
	DeletedItem
)

type RepeatType int

const (
	Daily RepeatType = iota + 1
	Weekly
	Monthly
	Yearly
)

type TokenObject struct {
	UserEmail      string
	TokenString    string
	ExpirationTime JsonDate
}

type JsonDate struct {
	time.Time
}

func (jd *JsonDate) UnmarshalJSON(b []byte) error {
	token := dateRe.FindSubmatch(b)
	if len(token) == 2 {
		i, err := strconv.ParseInt(string(token[1]), 10, 64)
		if err != nil {
			return err
		}
		jd.Time = time.Unix(i/1000, i%1000)
	}
	return nil
}

type RecurrenceObject struct {
	RepeatType      RepeatType
	SelectDays      int
	SelectWeeks     int
	Weekday0        bool
	Weekday1        bool
	Weekday2        bool
	Weekday3        bool
	Weekday4        bool
	Weekday5        bool
	Weekday6        bool
	SelectMonths    int
	MonthByMonthDay bool
	MonthByDay      bool
	SelectYears     int
	OriginalDate    time.Time
}

type ItemObject struct {
	Checked            bool             `json:"Checked"`
	Children           []ItemObject     `json:"Children"`
	Collapsed          bool             `json:"Collapsed"`
	Content            string           `json:"Content"`
	CreatedDate        JsonDate         `json:"CreatedDate"`
	DateString         string           `json:"DateString"`
	DateStringPriority int64            `json:"DateStringPriority"`
	Deleted            bool             `json:"Deleted"`
	DueDate            string           `json:"DueDate"`
	DueDateTime        string           `json:"DueDateTime"`
	DueTimeSpecified   bool             `json:"DueTimeSpecified"`
	Id                 int64            `json:"Id"`
	InHistory          bool             `json:"InHistory"`
	ItemOrder          int64            `json:"ItemOrder"`
	ItemType           ItemType         `json:"ItemType"`
	LastCheckedDate    string           `json:"LastCheckedDate"`
	LastSyncedDateTime string           `json:"LastSyncedDateTime"`
	LastUpdatedDate    string           `json:"LastUpdatedDate"`
	Notes              string           `json:"Notes"`
	OwnerId            int64            `json:"OwnerId"`
	ParentId           int64            `json:"ParentId"`
	Path               string           `json:"Path"`
	Priority           int64            `json:"Priority"`
	ProjectId          int64            `json:"ProjectId"`
	Recurrence         RecurrenceObject `json:"Recurrence"`
}

func auth() (string, error) {
	home := os.Getenv("HOME")
	if runtime.GOOS == "windows" && home == "" {
		home = os.Getenv("USERPROFILE")
	}
	m, _ := netrc.FindMachine(filepath.Join(home, ".netrc"), "todo.ly")

	req, err := http.NewRequest("GET", "https://todo.ly/api/authentication/token.json", nil)
	if err != nil {
		return "", err
	}
	if m != nil {
		req.SetBasicAuth(m.Login, m.Password)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var aerr ErrorObject
	err = json.Unmarshal(b, &aerr)
	if err != nil {
		return "", err
	}
	if aerr.ErrorCode != 0 {
		return "", fmt.Errorf("%d: %s", aerr.ErrorCode, aerr.ErrorMessage)
	}

	var ares TokenObject
	err = json.Unmarshal(b, &ares)
	if err != nil {
		return "", err
	}
	return ares.TokenString, nil
}

func main() {
	command := &commander.Command{
		UsageLine: os.Args[0],
		Short:     "cli interface for todo.ly",
	}
	command.Subcommands = []*commander.Command{
		make_cmd_list(),
		make_cmd_add(),
		make_cmd_del(),
	}
	err := command.Dispatch(os.Args[1:])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
