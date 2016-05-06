package main

import (
    "fmt"
    "errors"
    "time"
    "encoding/json"
    "math/rand"

    "github.com/parnurzeal/gorequest"
)

type Performance struct {
    IsFirstDayPerformance                bool    `json:"isFirstDayPerformance"`
    PerformanceAcsId                     int     `json:"performanceAcsId"`
    PerformanceDisplayFormat             string  `json:"performanceDisplayFormat"`
    EventId                              int     `json:"eventId"`
    PerformanceId                        int     `json:"performanceId"`
    BookmarkCreateTime                   int     `json:"bookmarkCreateTime"`
    BookmarkStatus                       int     `json:"bookmarkStatus"`
    PerformanceCategoryClass             string  `json:"performanceCategoryClass"`
    TransactionMaxQuota                  int     `json:"transactionMaxQuota"`
    PerformanceDateTime                  int64   `json:"performanceDateTime"`
    IsPurchasable                        bool    `json:"isPurchasable"`
    CounterSalesStartDate                *string `json:"counterSalesStartDate"`
    CounterSalesEndDate                  *string `json:"counterSalesEndDate"`
    DisplayDate                          bool    `json:"displayDate"`
    DisplayTime                          bool    `json:"displayTime"`
    ExternalReferenceKey                 string  `json:"externalReferenceKey"`
    PerformanceDisplayFormatValue        int     `json:"performanceDisplayFormatValue"`
    IsNotAllowedToPurchaseBeforeShowTime bool    `json:"isNotAllowedToPurchaseBeforeShowTime"`
    Note                                 *string `json:"note"`
    PerformanceName                      string  `json:"performanceName"`
}

type PerformanceData struct {
    PerformanceList []Performance `json:"performanceList"`
    StatusList      []string      `json:"performanceQuotaStatusList"`
}

type EventList struct {
    Name    string
    Time    string
    Status  string
}

type Checker struct {
    EventId     string
    Agent       *gorequest.SuperAgent
    PerPage     int
    Events      []EventList
    UserAgent   string
}

func NewChecker(eventId string) *Checker {
    return &Checker{
        EventId  : eventId,
        Agent    : gorequest.New(),
        PerPage  : 5,
        Events   : []EventList{},
        UserAgent: "",
    }
}

func (c Checker) EventList() ([]EventList, error) {
    // Setup user agent first
    c.UserAgent = c.userAgent()

    // Make auth cookie
    if _, err := c.fetchAuth(); err != nil {
        return nil, err
    }

    // Request event list
    if eventList, err := c.fetchEvent(1); err != nil {
        return nil, err
    }else{
        return eventList, nil
    }
}

func (c Checker) fetchAuth() (string, error) {
    resp, body, errs := c.Agent.Get("http://www.urbtix.hk/").
        Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8").
        Set("Accept-Language", "en-US,en;q=0.8").
        Set("Connection", "keep-alive").
        Set("User-Agent", c.UserAgent).
        End()

    if errs != nil {
        return "", errs[0]
    }

    if resp.Status != "200 OK" {
        return "", errors.New("Target page return status code: " + resp.Status)
    }

    if resp.Request.URL.String() != "https://ticket.urbtix.hk/internet/" {
        return "", errors.New("Redirect error: " + resp.Request.URL.String())
    }

    return body, nil
}

func (c *Checker) fetchEvent(pageNo int) ([]EventList, error) {
    timestamp := time.Now().Unix()
    targetUrl := fmt.Sprintf("https://ticket.urbtix.hk/internet/json/event/%s/performance/%d/%d/perf.json?locale=zh_TW&%d", c.EventId, c.PerPage, pageNo, timestamp)

    _, body, errs := c.Agent.Get(targetUrl).
        Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8").
        Set("Accept-Language", "en-US,en;q=0.8").
        Set("Connection", "keep-alive").
        Set("User-Agent", c.UserAgent).
        End()

    if errs != nil {
        return c.Events, errs[0]
    }

    performanceData := &PerformanceData{}

    if err := json.Unmarshal([]byte(body), performanceData); err != nil {
        return c.Events, err
    }else{
        for k, v := range performanceData.PerformanceList {
            timeString := time.Unix(v.PerformanceDateTime/1000, 0).Format(time.RFC3339)

            c.Events = append(c.Events, EventList{
                Name  : v.PerformanceName,
                Time  : timeString,
                Status: performanceData.StatusList[k],
            })
        }

        if len(performanceData.PerformanceList) > 0 {
            return c.fetchEvent(pageNo + 1)
        }else{
            return c.Events, nil
        }
    }
}

func (c Checker) userAgent() string {
    rand.Seed(time.Now().UTC().UnixNano())

    userAgents := []string{
        "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/50.0.2661.94 Safari/537.36",
        "Mozilla/5.0 (Macintosh; U; PPC Mac OS X; fr) AppleWebKit/416.12 (KHTML, like Gecko) Safari/412.5",
        "Mozilla/5.0 (Windows NT 6.1; rv:15.0) Gecko/20120819 Firefox/15.0 PaleMoon/15.0",
        "Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 5.1; Trident/4.0; GTB6; Acoo Browser; .NET CLR 1.1.4322; .NET CLR 2.0.50727)",
        "Mozilla/5.0 (Windows; U; Windows NT 5.1; pt-BR) AppleWebKit/534.12 (KHTML, like Gecko) NavscapeNavigator/Pre-0.1 Safari/534.12",
        "Mozilla/5.0 (Windows; U; WinNT4.0; de-AT; rv:1.7.11) Gecko/20050728",
    }

    return userAgents[rand.Intn(len(userAgents))]
}
