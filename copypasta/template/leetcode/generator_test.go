package leetcode

import (
	"encoding/json"
	"fmt"
	"github.com/levigross/grequests"
	"golang.org/x/net/html"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
)

const (
	leetCodeZH = "leetcode-cn.com"
	leetCodeEN = "leetcode.com"

	ua = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.108 Safari/537.36"
)

const (
	host      = leetCodeZH
	contestID = 165
)

var cookies []*http.Cookie

func init() {
	return

	csrftokenUrl := fmt.Sprintf("https://%s/graphql/", host)
	resp, err := grequests.Post(csrftokenUrl, &grequests.RequestOptions{
		UserAgent: ua,
		JSON:      map[string]string{"operationName": "globalData", "query": "query globalData {\n  feature {\n    questionTranslation\n    subscription\n    signUp\n    discuss\n    mockInterview\n    contest\n    store\n    book\n    chinaProblemDiscuss\n    socialProviders\n    studentFooter\n    cnJobs\n    __typename\n  }\n  userStatus {\n    isSignedIn\n    isAdmin\n    isStaff\n    isSuperuser\n    isTranslator\n    isPremium\n    isVerified\n    isPhoneVerified\n    isWechatVerified\n    checkedInToday\n    username\n    realName\n    userSlug\n    groups\n    jobsCompany {\n      nameSlug\n      logo\n      description\n      name\n      legalName\n      isVerified\n      permissions {\n        canInviteUsers\n        canInviteAllSite\n        leftInviteTimes\n        maxVisibleExploredUser\n        __typename\n      }\n      __typename\n    }\n    avatar\n    optedIn\n    requestRegion\n    region\n    activeSessionId\n    permissions\n    notificationStatus {\n      lastModified\n      numUnread\n      __typename\n    }\n    completedFeatureGuides\n    useTranslation\n    __typename\n  }\n  siteRegion\n  chinaHost\n  websocketUrl\n}\n"},
	})
	if err != nil {
		panic(err)
	}
	if !resp.Ok {
		panic(resp.StatusCode)
	}
	var csrfToken string
	for _, c := range resp.RawResponse.Cookies() {
		if c.Name == "csrftoken" {
			csrfToken = c.Value
			break
		}
	}
	if csrfToken == "" {
		panic("csrftoken not found")
	}

	loginUrl := fmt.Sprintf("https://%s/accounts/login/", host)
	resp, err = grequests.Post(loginUrl, &grequests.RequestOptions{
		UserAgent: ua,
		Data: map[string]string{
			"csrfmiddlewaretoken": csrfToken, // csrfToken,
			"login":               os.Getenv("USERNAME"),
			"password":            os.Getenv("PASSWORD"),
			"next":                "/",
		},
		Headers: map[string]string{
			"origin":  "https://leetcode-cn.com",
			"referer": "https://leetcode-cn.com/",
		},
	})
	if err != nil {
		panic(err)
	}
	if !resp.Ok {
		panic(resp.StatusCode)
	}
	for _, c := range resp.RawResponse.Cookies() {
		if c.Name == "csrftoken" || c.Name == "LEETCODE_SESSION" {
			cookies = append(cookies, c)
		}
	}
	if len(cookies) != 2 {
		panic(cookies)
	}
}

var contestDir = fmt.Sprintf("../../../leetcode/%d/", contestID)

func createDir(problemID string) error {
	dirPath := contestDir + problemID
	return os.MkdirAll(dirPath, os.ModePerm)
}

func writeMainFile(problemID, defaultCode string) error {
	defaultCode = strings.TrimSpace(defaultCode)
	mainStr := fmt.Sprintf(`package main

%s
`, defaultCode)
	filePath := contestDir + fmt.Sprintf("%[1]s/%[1]s.go", problemID)
	return ioutil.WriteFile(filePath, []byte(mainStr), 0644)
}

func writeTestFile(problemID, funcName string, sampleIns, sampleOuts [][]string) error {
	if funcName == "" {
		return fmt.Errorf("funcName is empty")
	}
	if len(sampleIns) != len(sampleOuts) {
		return fmt.Errorf("len(sampleIns) != len(sampleOuts) : %d != %d", len(sampleIns), len(sampleOuts))
	}

	funcName = strings.TrimSpace(funcName)
	sampleToStr := func(samples [][]string) (s string) {
		for i, args := range samples {
			if i > 0 {
				s += ", "
			}
			s += "{"
			for j, arg := range args {
				if j > 0 {
					s += ", "
				}
				s += "`" + arg + "`"
			}
			s += "}"
		}
		return
	}
	testStr := fmt.Sprintf(`package main

import (
	"github.com/EndlessCheng/codeforces-go/leetcode/testutil"
	"testing"
)

func Test(t *testing.T) {
	t.Log("Current test is [%s]")
	sampleIns := [][]string{%s}
	sampleOuts := [][]string{%s}
	if err := testutil.RunLeetCodeFunc(t, %s, sampleIns, sampleOuts); err != nil {
		t.Fatal(err)
	}
}
`, problemID, sampleToStr(sampleIns), sampleToStr(sampleOuts), funcName)
	filePath := contestDir + fmt.Sprintf("%[1]s/%[1]s_test.go", problemID)
	return ioutil.WriteFile(filePath, []byte(testStr), 0644)
}

func parseHTML(fileName string, htmlURL string) error {
	resp, err := grequests.Get(htmlURL, &grequests.RequestOptions{
		UserAgent: ua,
		Cookies: []*http.Cookie{
			{
				Name:  "LEETCODE_SESSION",
				Value: os.Getenv("LEETCODE_SESSION"),
			},
			{
				Name:  "csrftoken",
				Value: os.Getenv("CSRF_TOKEN"),
			},
		},
	})
	if err != nil {
		return err
	}
	if !resp.Ok {
		return fmt.Errorf("GET %s return code %d", htmlURL, resp.StatusCode)
	}

	root, err := html.Parse(resp)
	if err != nil {
		return err
	}

	htmlNode := root.FirstChild.NextSibling
	var bodyNode *html.Node
	for o := htmlNode.FirstChild; o != nil; o = o.NextSibling {
		if o.Type == html.ElementNode && o.Data == "body" {
			bodyNode = o
			break
		}
	}

	var funcName string
	var genTestFile bool
	for o := bodyNode.FirstChild; o != nil; o = o.NextSibling {
		if o.Type == html.ElementNode && o.Data == "script" && o.FirstChild != nil {
			jsText := o.FirstChild.Data
			if start := strings.Index(jsText, "codeDefinition:"); start != -1 {
				end := strings.Index(jsText, "enableTestMode")
				jsonText := jsText[start+len("codeDefinition:") : end]
				jsonText = strings.TrimSpace(jsonText)
				jsonText = jsonText[:len(jsonText)-3] + "]" // remove , at end
				jsonText = strings.Replace(jsonText, `'`, `"`, -1)

				d := []struct {
					Value       string `json:"value"`
					DefaultCode string `json:"defaultCode"`
				}{}
				if err := json.Unmarshal([]byte(jsonText), &d); err != nil {
					return err
				}

				for _, e := range d {
					if e.Value == "golang" {
						funcName, genTestFile = parseFuncName(e.DefaultCode)
						if err := writeMainFile(fileName, e.DefaultCode); err != nil {
							return err
						}
						break
					}
				}
				break
			}
		}
	}

	if !genTestFile {
		return nil
	}

	const (
		tokenInputZH  = "输入："
		tokenOutputZH = "输出："

		tokenInputEN  = "Input:"
		tokenOutputEN = "Output:"
	)

	var sampleIns, sampleOuts [][]string
	parseSampleText := func(text string, parseArgs bool) (sample []string) {
		text = strings.TrimSpace(text)
		text = strings.Replace(text, "\n", "", -1)
		if !parseArgs {
			return []string{text}
		}
		splits := strings.Split(text, "=")
		sample = make([]string, 0, len(splits)-1)
		for _, s := range splits[1 : len(splits)-1] {
			end := strings.LastIndexByte(s, ',')
			sample = append(sample, strings.TrimSpace(s[:end]))
		}
		sample = append(sample, strings.TrimSpace(splits[len(splits)-1]))
		return
	}
	var f func(*html.Node)
	f = func(o *html.Node) {
		if o.Type == html.TextNode {
			if o.Data == tokenInputZH {
				sample := parseSampleText(o.Parent.NextSibling.Data, true)
				sampleIns = append(sampleIns, sample)
			} else if o.Data == tokenOutputZH {
				sample := parseSampleText(o.Parent.NextSibling.Data, false)
				sampleOuts = append(sampleOuts, sample)
			}
			return
		}
		for c := o.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(bodyNode)

	if err := writeTestFile(fileName, funcName, sampleIns, sampleOuts); err != nil {
		return err
	}

	return nil
}

func TestGenLeetCodeTests(t *testing.T) {
	apiInfoContest := fmt.Sprintf("https://%s/contest/api/info/weekly-contest-%d/", host, contestID)
	resp, err := grequests.Get(apiInfoContest, &grequests.RequestOptions{UserAgent: ua})
	if err != nil {
		t.Fatal(err)
	}
	if !resp.Ok {
		t.Fatalf("GET %s return code %d", apiInfoContest, resp.StatusCode)
	}
	d := struct {
		Questions []struct {
			TitleSlug string `json:"title_slug"`
		} `json:"questions"`
	}{}
	if err := resp.JSON(&d); err != nil {
		t.Fatal(err)
	}

	problemURLs := make([]string, len(d.Questions))
	for i, q := range d.Questions {
		problemURLs[i] = fmt.Sprintf("https://%s/contest/weekly-contest-%d/problems/%s/", host, contestID, q.TitleSlug)
	}
	for i, pUrl := range problemURLs {
		problemID := string('a' + i)
		fmt.Println(problemID, pUrl)
		if err := createDir(problemID); err != nil {
			t.Fatal(err)
		}
		if err := parseHTML(problemID, pUrl); err != nil {
			t.Fatal(err, problemID, pUrl)
		}
	}
}
