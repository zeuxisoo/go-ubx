package main

import (
    "fmt"
    "errors"

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
    IsNotAllowedToPurchaseBeforeShowTime string  `json:"isNotAllowedToPurchaseBeforeShowTime"`
    Note                                 *string `json:"note"`
    PerformanceName                      string  `json:"performanceName"`
}

type PerformanceData struct {
    PerformanceList []Performance `json:"performanceList"`
    StatusList      []string      `json:"performanceQuotaStatusList"`
}

type Checker struct {
    EventId string
    Agent   *gorequest.SuperAgent
}

func NewChecker(eventId string) *Checker {
    return &Checker{
        EventId: eventId,
        Agent  : gorequest.New(),
    }
}

func (c Checker) Check() (string, error) {
    _, err := c.fetchAuth()

    if err != nil {
        return "", err
    }



    return "", nil
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
