package main

import (
  "io"
  "fmt"
  "strings"
  "strconv"
  "time"
  "net/http"
)

func main() {
  failedAttempts := 0
  handleFailedAttempt := func() {
    failedAttempts += 1

    if failedAttempts < 3 {
      return
    }

    failedAttempts = 0

    fmt.Println("Unable to fetch server statistic")
  }

  for {
    // fetch statistics
    response, err := http.Get("http://srv.msk01.gigacorp.local/_stats")
    if err != nil {
      handleFailedAttempt()
      continue
    }

    // only 200 status code is allowed
    if response.StatusCode != 200 {
      handleFailedAttempt()
      continue
    }

    // read response body
    body, err := io.ReadAll(response.Body)
    response.Body.Close()
    if err != nil {
      handleFailedAttempt()
      continue
    }

    // parse raw response data
    rawDataArray := strings.Split(string(body), ",")
    if (len(rawDataArray) != 7) {
      handleFailedAttempt()
      continue
    }

    // convert raw data into integer array
    var data = [7]int{}
    var isDataConvertFailed = false
    for index, element := range rawDataArray {
      i, err := strconv.Atoi(element)
      if err != nil {
        isDataConvertFailed = true
        break;
      }

      data[index] = i
    }
    if isDataConvertFailed == true {
      handleFailedAttempt()
      continue
    }

    // check server stats
    if data[0] > 30 {
      fmt.Println("Load Average is too high: %d", data[0])
    }

    time.Sleep(time.Second)
  }
}
