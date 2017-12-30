package main

import (
	"fmt"
	"log"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func getAllObjectsFromTable(session *mgo.Session, table string, out interface{}) {
	cSession := session.Copy()

	defer cSession.Close()

	fmt.Printf("Getting Objects From Database: %v, and table: %v", DATABASE, table)

	c := cSession.DB(DATABASE).C(table)

	c.Find(bson.M{}).All(out)
	fmt.Printf("Objects: %v\n", &out)
}

func insertObjectIntoTable(session *mgo.Session, table string, obj interface{}) {
	cSession := session.Copy()

	defer cSession.Close()

	fmt.Printf("Inserting %v into table %s", obj, table)

	c := cSession.DB(DATABASE).C(table)

	err := c.Insert(obj)

	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}
