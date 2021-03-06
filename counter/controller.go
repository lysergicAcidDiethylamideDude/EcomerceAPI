package counter

import (
	"github.com/mongodb/mongo-go-driver/bson"
	"go-api-ws/config"
	"go-api-ws/helpers"
)

func GetAndIncreaseItemIdCounterInMongo() int {
	var itemCounter ItemCounter

	db := config.Conf.GetMongoDb()

	err := db.Collection(collectionName).FindOneAndUpdate(nil,
		bson.NewDocument(
			bson.EC.String("_id", "itemid")),
		bson.NewDocument(
			bson.EC.SubDocumentFromElements("$inc",
				bson.EC.Interface("value", 1)))).Decode(&itemCounter)
	helpers.PanicErr(err)
	return itemCounter.Value
}

func GetAndIncreaseQuoteCounterInMySQL() (value int64) {
	db, err := config.Conf.GetDb()
	helpers.PanicErr(err)

	//defer db.Close()

	tx, err := db.Begin()
	helpers.PanicErr(err)

	err = tx.QueryRow("SELECT value FROM counters WHERE name = ?", quoteCounter).Scan(&value)
	if err != nil {
		helpers.CheckErr(err)
		err = tx.Rollback()
		helpers.PanicErr(err)
	}

	_, err = tx.Exec("UPDATE counters SET value = value + 1 WHERE name = ?", quoteCounter)
	if err != nil {
		helpers.CheckErr(err)
		err = tx.Rollback()
		helpers.PanicErr(err)
	}

	err = tx.Commit()
	helpers.PanicErr(err)

	return
}
