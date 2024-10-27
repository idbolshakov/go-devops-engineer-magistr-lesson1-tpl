package main

import (
  "io"
  "fmt"
  "strings"
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
    response, err := http.Get("http://srv.msk01.gigacorp.local/_stats")
    if err != nil {
      handleFailedAttempt()
      continue
    }

    if response.StatusCode != 200 {
      handleFailedAttempt()
      continue
    }

    body, err := io.ReadAll(response.Body)
    response.Body.Close()
    if err != nil {
      handleFailedAttempt()
      continue
    }

    rawDataArray := strings.Split(string(body), ",")
    if (len(rawDataArray) != 7) {
      handleFailedAttempt()
      continue
    }

    time.Sleep(time.Second)
  }
}
