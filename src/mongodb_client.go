package main

import (
  	"log"
  	"time"

  	"gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
)

var connection, err = mgo.Dial("localhost")
var database = connection.DB("user_event")

func save(data map[string]interface{}, collection string) bool {
  	if err != nil {
    		log.Fatal(err)
  	}

  	cursor := database.C(collection)
  	err = cursor.Insert(data)
  	if err != nil {
    		log.Panic(err)
    		return false
  	}
  	return true
}

func SaveUserEvent(userEvent map[string]interface{}) {
  	save(userEvent, userEvent["APIKey"].(string))
    return
}

func SaveResponseTime(responseTime int64) {
  	var responseTimeData = map[string]interface{}{
    		"response_time": responseTime,
    		"timestamp": time.Now(),
  	}

  	save(responseTimeData, "time_stats")
    return
}

func GetResponseTimeStats() ([]bson.M, interface{}) {
    pipeline := []bson.M{
        bson.M{
            "$project": bson.M{
                "range": bson.M{
                    "$concat": []bson.M{
                        bson.M{"$cond": []interface{}{ bson.M{ "$and": []interface{}{ bson.M{ "$gt": []interface{}{"$response_time", 0,},}, bson.M{ "$lte": []interface{}{"$response_time", 1, },},},},"<1ms", "",},},
                        bson.M{"$cond": []interface{}{ bson.M{ "$and": []interface{}{ bson.M{ "$gt": []interface{}{"$response_time", 1,},}, bson.M{ "$lte": []interface{}{"$response_time", 5, },},},},"<5ms", "",},},
                        bson.M{"$cond": []interface{}{ bson.M{ "$and": []interface{}{ bson.M{ "$gt": []interface{}{"$response_time", 5,},}, bson.M{ "$lte": []interface{}{"$response_time", 10, },},},},"<10ms", "",},},
                        bson.M{"$cond": []interface{}{ bson.M{ "$and": []interface{}{ bson.M{ "$gt": []interface{}{"$response_time", 10,},}, bson.M{ "$lte": []interface{}{"$response_time", 20, },},},},"<20ms", "",},},
                        bson.M{"$cond": []interface{}{ bson.M{ "$and": []interface{}{ bson.M{ "$gt": []interface{}{"$response_time", 20,},}, bson.M{ "$lte": []interface{}{"$response_time", 50, },},},},"<50ms", "",},},
                        bson.M{"$cond": []interface{}{ bson.M{ "$and": []interface{}{ bson.M{ "$gt": []interface{}{"$response_time", 50,},}, bson.M{ "$lte": []interface{}{"$response_time", 100, },},},},"<100ms", "",},},
                    },
                },
            },
        },
        bson.M{
            "$group": bson.M{
                "_id": "$range",
                "count": bson.M{
                    "$sum": 1.0,
                },
            },
        },
    }
    pipe := database.C("time_stats").Pipe(pipeline)
    results := []bson.M{}
    err := pipe.All(&results)
    return results, err
}
