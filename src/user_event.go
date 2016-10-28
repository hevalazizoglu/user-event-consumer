package main

import (
    "encoding/json"
    "log"
    "math/rand"
    "net/http"
    "time"
)

func validateAPIKey(key string) bool {
    for _, value := range KnownClients {
      	if value == key {
      		  return true
      	}
    }
    return false
}

func checkRequiredFields(user_event map[string]interface{}) bool {
    _, o1 := user_event["APIKey"]
    _, o2 := user_event["UserID"]
    _, o3 := user_event["Timestamp"]

    if o1 && o2 && o3 {
        return true
    }
    return false
}

func timeTrack(start time.Time) {
    responseTime := time.Since(start)
    /* here is a little bit tricky,
    responseTime is time.Duration type and during Mongodb conversion it is directly
    being coverted to NumberLong. This is causing us to lose significance. Hence
    coverting to int here.
    */
    responseTimeInMs := responseTime.Nanoseconds()/time.Millisecond.Nanoseconds()
    SaveResponseTime(responseTimeInMs)
}

func sleep() {
    lowerLimit, upperLimit := 1, 100
    sleepTime := rand.Intn(upperLimit-lowerLimit) + lowerLimit

    time.Sleep(time.Duration(sleepTime) * time.Millisecond)
}

func HandleUserEvent(responseWriter http.ResponseWriter, request *http.Request) {
    defer timeTrack(time.Now())

    decoder := json.NewDecoder(request.Body)
    var data interface{}
    err := decoder.Decode(&data)
    if err != nil {
      	log.Panic(err)
    }
    user_event := data.(map[string]interface{})

    if !checkRequiredFields(user_event) {
        http.Error(responseWriter, "Missing Required Fields", http.StatusBadRequest)
        return
    }

    if !validateAPIKey(user_event["APIKey"].(string)) {
        http.Error(responseWriter, "Invalid APIKey", http.StatusBadRequest)
        return
    }

    SaveUserEvent(user_event)
    sleep()

    responseWriter.Header().Set("Content-Type", "text/plain; charset=utf-8")
    responseWriter.WriteHeader(http.StatusOK)
}
