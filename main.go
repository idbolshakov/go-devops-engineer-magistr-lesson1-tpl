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
    if isDataConvertFailed {
      handleFailedAttempt()
      continue
    }

    // check Load Average
    if data[0] > 30 {
      fmt.Printf("Load Average is too high: %d\n", data[0])
    }

    // check Memory usage
    memoryUsage := data[2] * 100 / data[1]
    if memoryUsage > 80 {
      fmt.Printf("Memory usage too high: %d%%\n", memoryUsage)
    }

    // check free disk space
    if data[4] * 100 / data[3] > 90 {
      freeDiskSpace := (data[3] - data[4]) / (1024*1024)
      fmt.Printf("Free disk space is too low: %d Mb left\n", freeDiskSpace)
    }

    // check network bandwidth
    if data[6] * 100 / data[5] > 90 {
      freeBandwidth := (data[5] - data[6]) / 1000000
      fmt.Printf("Network bandwidth usage high: %d Mbit/s available\n", freeBandwidth)
    }

    time.Sleep(time.Second)
  }
}
