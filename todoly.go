package main

import (
	"encoding/json"
	"fmt"
	"github.com/daviddengcn/go-colortext"
	"github.com/fhs/go-netrc/netrc"
	"io"
	"io/ioutil"
	"log"
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

var priorityColor = map[int64]ct.Color{
	1: ct.Red,
	2: ct.Blue,
	3: ct.Green,
	4: ct.White,
}

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

func main() {
	home := os.Getenv("HOME")
	if runtime.GOOS == "windows" && home == "" {
		home = os.Getenv("USERPROFILE")
	}
	m, _ := netrc.FindMachine(filepath.Join(home, ".netrc"), "todo.ly")

	req, err := http.NewRequest("GET", "https://todo.ly/api/authentication/token.json", nil)
	if err != nil {
		log.Fatal(err)
	}
	if m != nil {
		req.SetBasicAuth(m.Login, m.Password)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var aerr ErrorObject
	err = json.Unmarshal(b, &aerr)
	if err != nil {
		log.Fatal(err)
	}
	if aerr.ErrorCode != 0 {
		log.Fatalf("%d: %s", aerr.ErrorCode, aerr.ErrorMessage)
	}

	var ares TokenObject
	err = json.Unmarshal(b, &ares)
	if err != nil {
		log.Fatal(err)
	}

	req, err = http.NewRequest("GET", "https://todo.ly/api/items.json", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Token", ares.TokenString)
	res, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	var items []ItemObject
	err = json.NewDecoder(res.Body).Decode(&items)
	if err != nil {
		log.Fatal(err)
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
}
