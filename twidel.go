package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/ChimeraCoder/anaconda"
)

const (
	consumerKey    string = "TIIdWkTGq9KtE2lEJZlngu45i"
	consumerSecret string = "wSPphN2UPziOOJKwqfqHiXAKAYXR4usZunTbnYolqH7xyQZoxz"
)

var (
	limit   int
	minfav  int
	minrt   int
	dbgflag bool
)

type Configuration struct {
	AccessToken       string `json:"access_token"`
	AccessTokenSecret string `json:"access_token_secret"`
}

func getSettingFileName(target string) string {
	dir, err := os.Getwd() // current directory
	if err != nil {
		panic(err)
	}
	dir = filepath.Join(dir, "setting")
	if err := os.MkdirAll(dir, 0700); err != nil {
		panic(err)
	}
	return filepath.Join(dir, target+"_setting.json")
}

func getApi(target string) *anaconda.TwitterApi {
	filename := getSettingFileName(target)
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	var conf Configuration
	if err := json.Unmarshal(b, &conf); err != nil {
		panic(err)
	}
	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)
	return anaconda.NewTwitterApi(conf.AccessToken, conf.AccessTokenSecret)
}

// Write access token and secret to json file.
func authenticate(target string) {
	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)
	url, token, err := anaconda.AuthorizationURL("")
	if err != nil {
		panic(err)
	}
	fmt.Println("Access url and log in to get PIN code.")
	fmt.Println(url)

	var verifier string
	fmt.Print("Input PIN Code: ")
	fmt.Scan(&verifier)
	t, _, err := anaconda.GetCredentials(token, verifier)
	if err != nil {
		panic(err)
	}

	conf := Configuration{
		AccessToken:       t.Token,
		AccessTokenSecret: t.Secret,
	}
	outjson, err := json.MarshalIndent(conf, "", "\t")
	if err != nil {
		panic(err)
	}
	filename := getSettingFileName(target)
	err = ioutil.WriteFile(filename, outjson, 0644)
	if err != nil {
		panic(err)
	}
}

func isAuthenticated(target string) bool {
	// Return true if (target)_setting.json exists.
	filename := getSettingFileName(target)
	_, err := os.Stat(filename)
	return err == nil
}

func init() {
	// Set flags
	INF := 114514
	flag.IntVar(&limit, "limit", 3200, "Limit of number to delete tweets")
	flag.IntVar(&minfav, "minfav", INF, "Delete tweet less than minfav")
	flag.IntVar(&minrt, "minrt", INF, "Delete tweet less than minrt")
	flag.BoolVar(&dbgflag, "dbg", false, "Debug mode on if dbg=true")
	flag.Parse()
}

func main() {
	if dbgflag {
		fmt.Println("Debug mode ON!")
	}

	var twitterAccount [2]string
	fmt.Println("In order to avoid mis-execution, 2 times input your twitter account that you want to target.")
	for i := 0; i < 2; i++ {
		fmt.Printf("(%v/2): ", i+1)
		_, err := fmt.Scan(&twitterAccount[i])
		if err != nil {
			panic(err)
		}
	}
	if twitterAccount[0] != twitterAccount[1] {
		fmt.Println("Error: Inputted id is mistaken.")
		os.Exit(0)
	}

	targetAccount := twitterAccount[0]
	if !isAuthenticated(targetAccount) {
		authenticate(targetAccount)
	}

	api := getApi(targetAccount)
	for {
		fmt.Print("Are you sure you want to delete your tweets? (y/n): ")
		var yn string
		fmt.Scan(&yn)
		if yn == "n" || yn == "N" {
			os.Exit(0)
		} else if yn == "y" || yn == "Y" {
			fmt.Println("Start to delete tweets...")
			break
		}
	}

	// Delete user's tweets (up to 3200)
	deleted := 0
	LIM := limit / 200
	v := url.Values{}
	v.Set("count", "200")
	for page := 1; page <= LIM; page++ {
		if page == LIM {
			v.Set("count", strconv.Itoa(limit%200))
		}
		v.Set("page", strconv.Itoa(page))
		timeline, err := api.GetUserTimeline(v)
		if err != nil {
			panic(err)
		}
		if dbgflag {
			fmt.Println("Len: ", len(timeline))
		}
		if len(timeline) == 0 {
			break
		}

		for _, tweet := range timeline {
			if tweet.FavoriteCount >= minfav || tweet.RetweetCount >= minrt {
				continue
			}
			t, err := api.DeleteTweet(tweet.Id, false)
			if err != nil {
				panic(err)
			}
			if dbgflag {
				fmt.Println(t.Text)
			}
			deleted++
			if deleted%100 == 0 {
				fmt.Printf("Deleted %v tweets\n", deleted)
				time.Sleep(10 * time.Second)
			}
		}
	}

	fmt.Printf("Deleted %v tweets\n", deleted)
	fmt.Println("Done")
}
