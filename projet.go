package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"
)

var alreadycheck []string
var urlOrigin string
var nbOK int
var nbError int
var okList string
var errorList string
var finalList string

type error interface {
	Error() string
}

func main() {
	start := time.Now()

	urlStart := os.Args[1]
	regOrigin := regexp.MustCompile(`https?://.*?/`)
	urlOrigin = regOrigin.FindString(urlStart)
	nbOK = 0
	nbError = 0
	checkurl(urlOrigin)

	nbErrorString := strconv.Itoa(nbError)
	nbOKString := strconv.Itoa(nbOK)
	nbTotal := strconv.Itoa(len(alreadycheck))

	timeElapsed := time.Since(start)
	fmt.Println(len(alreadycheck), " url find in ", timeElapsed)

	finalList = nbTotal + " url find in " + timeElapsed.String() + "\nlist error (" + nbErrorString + ") :\n\n" + errorList + "\nlist OK (" + nbOKString + "):\n\n" + okList
	out := []byte(finalList)
	result := ioutil.WriteFile("result.txt", out, 0644)
	checkError(result)
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
}

func checkurl(link string) {
	for _, value := range alreadycheck {
		if link == value {
			return
		}
	}
	alreadycheck = append(alreadycheck, link)

	page, err := http.Get(link)
	if err != nil {
		errorList += link + " -> Error : " + err.Error() + "\n"
		nbError = nbError + 1
		return
	}
	defer page.Body.Close()

	if page.StatusCode == 200 {
		okList += link + " -> OK\n"
		nbOK = nbOK + 1
	} else {
		errorList += link + " -> Error : " + strconv.Itoa(page.StatusCode) + "\n"
		nbError = nbError + 1
	}

	regO := regexp.MustCompile(urlOrigin + `.*`)

	if regO.MatchString(link) {
		content, err := ioutil.ReadAll(page.Body)
		checkError(err)
		stringContent := string(content)
		findurl(stringContent)
	}
}

func findurl(body string) {

	regA := regexp.MustCompile(`https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,4}\b([-a-zA-Z0-9@:%_\+.~#?&//=]*)`)
	regR := regexp.MustCompile(`"\/([\w\-\.,@?^=%&amp;:/~\+#]*[\w\-\@?^=%&amp;/~\+#])+"`)
	listA := regA.FindAllString(string(body), -1)
	listR := regR.FindAllString(string(body), -1)

	var newlistR []string
	var allurl []string

	for _, value := range listR {
		newhref := urlOrigin + value[2:len(value)-1]
		newlistR = append(newlistR, newhref)
	}

	allurl = append(listA, newlistR...)

	for _, value := range allurl {
		checkurl(value)
	}
}
