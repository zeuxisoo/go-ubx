package main

import (
    "fmt"
    "errors"
    "time"
    "encoding/json"

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
    EventId string
    Agent   *gorequest.SuperAgent
    PerPage int
    PageNo  int
}

func NewChecker(eventId string) *Checker {
    return &Checker{
        EventId: eventId,
        Agent  : gorequest.New(),
        PerPage: 5,
        PageNo : 1,
    }
}

func (c Checker) EventList() ([]EventList, error) {
    if _, err := c.fetchAuth(); err != nil {
        return nil, err
    }

    if eventList, err := c.fetchEvent(); err != nil {
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
        Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/50.0.2661.94 Safari/537.36").
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

func (c Checker) fetchEvent() ([]EventList, error) {
    var eventList []EventList

    timestamp := time.Now().Unix()
    targetUrl := fmt.Sprintf("https://ticket.urbtix.hk/internet/json/event/%s/performance/%d/%d/perf.json?locale=zh_TW&%d", c.EventId, c.PerPage, c.PageNo, timestamp)

    _, body, errs := c.Agent.Get(targetUrl).
        Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8").
        Set("Accept-Language", "en-US,en;q=0.8").
        Set("Connection", "keep-alive").
        Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/50.0.2661.94 Safari/537.36").
        End()

    if errs != nil {
        return eventList, errs[0]
    }

    performanceData := &PerformanceData{}

    if err := json.Unmarshal([]byte(body), performanceData); err != nil {
        return eventList, err
    }else{
        for k, v := range performanceData.PerformanceList {
            timeString := time.Unix(v.PerformanceDateTime/1000, 0).Format(time.RFC3339)

            eventList = append(eventList, EventList{
                Name  : v.PerformanceName,
                Time  : timeString,
                Status: performanceData.StatusList[k],
            })
        }

        return eventList, nil
    }
}
