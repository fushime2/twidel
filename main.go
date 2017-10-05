package main

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"strconv"

	"github.com/ChimeraCoder/anaconda"
)

func getLines(fname string) ([]string, error) {
	fp, err := os.Open(fname)
	if err != nil {
		panic(err)
	}
	defer fp.Close()
	scanner := bufio.NewScanner(fp)
	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	return lines, nil
}

func getApi() *anaconda.TwitterApi {
	info, err := getLines("oauth.txt")
	if err != nil {
		panic(err)
	}
	anaconda.SetConsumerKey(info[0])
	anaconda.SetConsumerSecret(info[1])
	return anaconda.NewTwitterApi(info[2], info[3])
}

func authenticate() {
	info, err := getLines("oauth.txt")
	if err != nil {
		panic(err)
	}
	anaconda.SetConsumerKey(info[0])
	anaconda.SetConsumerSecret(info[1])
	url, token, err := anaconda.AuthorizationURL("")
	if err != nil {
		panic(err)
	}
	fmt.Println(url)
	fmt.Println(token.Token, token.Secret)

	var verifier string
	fmt.Print("Input PIN Code: ")
	fmt.Scan(&verifier)
	token, val, err := anaconda.GetCredentials(token, verifier)
	fmt.Println(token)
	fmt.Println(val)
	if err != nil {
		panic(err)
	}
	//api := anaconda.NewTwitterApi(token.Token, token.Secret)
}

func test() {
	fmt.Println("0")
}

func main() {
	dbgflag := false
	if dbgflag {
		fmt.Println("Debug mode ON!")
		test()
		return
	}

	info, err := getLines("oauth.txt")
	if err != nil {
		panic(err)
	}

	anaconda.SetConsumerKey(info[0])
	anaconda.SetConsumerSecret(info[1])
	api := anaconda.NewTwitterApi(info[2], info[3])

	// Delete user's tweets (up to 3200)
	LIM := 16
	v := url.Values{}
	v.Set("count", "200")
	for page := 1; page <= LIM; page++ {
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
			t, err := api.DeleteTweet(tweet.Id, false)
			if err != nil {
				panic(err)
			}
			fmt.Println(t.Text)
		}
	}

	fmt.Println("Done")
}
