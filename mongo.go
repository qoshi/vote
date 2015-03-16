package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fatih/structs"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

type Sb struct {
	Iq int
	Eq int
}

const (
	mongoServer = "localhost"
	mongoPort   = "27017"
	DBName      = "test"
	tableName   = "vote"
)

func getResult(id string) ([]byte, error) {
	log.Println("getting vote ID =", id)
	mongoStr := mongoServer + ":" + mongoPort
	session, err := mgo.Dial(mongoStr)
	if err != nil {
		return nil, err
	}
	defer session.Close()
	collection := session.DB(DBName).C(tableName)

	oid := bson.ObjectIdHex(id)
	query := bson.M{"_id": oid}
	var result = bson.M{}
	err = collection.Find(query).One(&result)
	if err != nil {
		return nil, err
	}
	ret, err3 := json.Marshal(result)
	if err3 != nil {
		ret = []byte("")
	}
	log.Println("result", string(ret), err3)
	return ret, err3
}

func newVote(v Vote) (string, error) {
	log.Println("new vote", v)
	oid := bson.NewObjectId()
	mongoStr := mongoServer + ":" + mongoPort
	session, err := mgo.Dial(mongoStr)
	if err != nil {
		return "", err
	}
	defer session.Close()
	collection := session.DB(DBName).C(tableName)

	var insertMap = bson.M{}
	var tmap = structs.Map(&v)
	insertMap["_id"] = oid
	for key, val := range tmap {
		insertMap[key] = val
	}
	err = collection.Insert(insertMap)
	if err != nil {
		return "", err
	}
	var ret string
	_, err = fmt.Sscanf(oid.String(), "ObjectIdHex(%q)", &ret)
	if err != nil {
		return "", err
	}
	log.Println("new vote created", ret)
	return ret, nil
}

func vote(id string, v Vote) error {
	log.Println("vote for", id, v)
	mongoStr := mongoServer + ":" + mongoPort
	oid := bson.ObjectIdHex(id)
	session, err := mgo.Dial(mongoStr)
	if err != nil {
		return err
	}
	defer session.Close()
	collection := session.DB(DBName).C(tableName)

	query := bson.M{"_id": oid}
	t := bson.M{}
	err = collection.Find(query).One(&t)
	if err != nil {
		return err
	}
	tDetail, ok := t["Detail"]
	if !ok {
		return errors.New("no such vote")
	}
	tDetailMap := tDetail.(bson.M)
	if err != nil {
		return err
	}
	for key1, val1 := range tDetailMap {
		if valt, ok := v.Detail[key1]; ok != false && valt != 0 {
			v.Detail[key1] = 1 + val1.(int)
		} else {
			v.Detail[key1] = val1.(int)
		}
	}
	updateT := bson.M{}
	for key, val := range structs.Map(&v) {
		updateT[key] = val
	}
	collection.Update(query, &updateT)
	log.Println("vote success", updateT)
	return nil
}

func test() {
	m := make(map[string]int)
	m["AAA"] = 1
	m["BBB"] = 1
	m["CCC"] = 1
	v := Vote{"test", 0, m, "adf"}
	id, err := newVote(v)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("insert id", id)
	r, er1 := getResult(id)
	if er1 != nil {
		fmt.Println(er1)
	}
	fmt.Println("result", string(r))
	m = make(map[string]int)
	m["AAA"] = 1
	v2 := Vote{"test", 0, m, "adf"}
	err = vote(id, v2)
	m = make(map[string]int)
	m["BBB"] = 1
	v2 = Vote{"test", 0, m, "adf"}
	err = vote(id, v2)
	m = make(map[string]int)
	m["CCC"] = 1
	v2 = Vote{"test", 0, m, "adf"}
	err = vote(id, v2)
	if err != nil {
		fmt.Println(err)
	}
	r2, err2 := getResult(id)
	if err != nil {
		fmt.Println(err2)
	}
	fmt.Println(string(r2))
}

// func main() {
// 	test()
// }
