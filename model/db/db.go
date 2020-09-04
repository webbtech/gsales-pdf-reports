package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	log "github.com/sirupsen/logrus"

	pkgerrors "github.com/pulpfree/go-errors"
	"github.com/pulpfree/gsales-pdf-reports/model"
)

// MDB struct
type MDB struct {
	client *mongo.Client
	dbName string
	db     *mongo.Database
}

// DB and Table constants
const (
	colConfig       = "config"
	colEmployees    = "employees"
	colJournals     = "journals"
	colNonFuelSales = "non-fuel-sales"
	colProducts     = "products"
	colSales        = "sales"
	colStations     = "stations"
)

const noRecordsMsg = "No records found matching criteria"

// ======================== Exported Functions ================================================= //

// NewDB sets up new MDB struct
func NewDB(connection string, dbNm string) (*MDB, error) {

	clientOptions := options.Client().ApplyURI(connection)
	err := clientOptions.Validate()
	if err != nil {
		log.Fatal(err)
	}

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Infoln("Connected to MongoDB!")

	return &MDB{
		client: client,
		dbName: dbNm,
		db:     client.Database(dbNm),
	}, err
}

// ======================== Exported Methods =================================================== //

// Close method
func (db *MDB) Close() {
	err := db.client.Disconnect(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	log.Infoln("Connection to MongoDB closed.")
}

// GetDay method
func (db *MDB) GetDay(date time.Time, stationID primitive.ObjectID) (day bson.M, err error) {

	day, err = db.fetchDay(date, stationID)
	if err != nil {
		return nil, err
	}
	if day == nil {
		return nil, &pkgerrors.StdError{Err: "", Caller: "db.GetDay", Msg: noRecordsMsg}
	}

	return day, err
}

// GetEmployee method
func (db *MDB) GetEmployee(attendantID primitive.ObjectID) (employee *model.Employee, err error) {

	employee, err = db.fetchEmployee(attendantID)
	if err != nil {
		errStr := fmt.Sprintf("Failed to fetch employee record with id:%s", attendantID)
		return nil, &pkgerrors.StdError{Err: errStr, Caller: "db.GetShift", Msg: "Failed to fetch employee"}
	}

	return employee, err
}

// GetJournals method
func (db *MDB) GetJournals(recordNum string, stationID primitive.ObjectID) (journals []*model.Journal, err error) {

	journals, err = db.fetchJournals(recordNum, stationID)
	if err != nil {
		errStr := fmt.Sprintf("Failed to fetch journal records with recordNum:%s and stationID:%v", recordNum, stationID)
		return nil, &pkgerrors.StdError{Err: errStr, Caller: "db.GetJournals", Msg: "Failed to fetch journal entries"}
	}

	return journals, err
}

// GetShift method
func (db *MDB) GetShift(recordNum string, stationID primitive.ObjectID) (shift *model.Sales, err error) {

	shift, err = db.fetchShift(recordNum, stationID)
	if err != nil {
		return nil, err
	}
	if shift == nil {
		return nil, &pkgerrors.StdError{Err: "", Caller: "db.GetShift", Msg: noRecordsMsg}
	}

	return shift, err
}

// GetStation method
func (db *MDB) GetStation(stationID primitive.ObjectID) (station *model.Station, err error) {

	station, err = db.fetchStation(stationID)
	if err != nil {
		errStr := fmt.Sprintf("Failed to fetch station record with id:%s", stationID)
		return nil, &pkgerrors.StdError{Err: errStr, Caller: "db.GetShift", Msg: "Failed to fetch station"}
	}

	return station, err
}

// ======================== Un-exported Methods ================================================ //

// fetchDay method
func (db *MDB) fetchDay(date time.Time, stationID primitive.ObjectID) (day bson.M, err error) {

	col := db.db.Collection(colSales)
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		{
			primitive.E{
				Key: "$match",
				Value: bson.D{
					primitive.E{
						Key:   "recordDate",
						Value: date,
					},
					primitive.E{
						Key:   "stationID",
						Value: stationID,
					},
				},
			},
		},
		{
			primitive.E{
				Key: "$group",
				Value: bson.D{
					primitive.E{
						Key:   "_id",
						Value: "$stationID",
					},
					primitive.E{
						Key: "cash_bills",
						Value: bson.D{
							primitive.E{
								Key:   "$sum",
								Value: "$cash.bills",
							},
						},
					},
					primitive.E{
						Key: "cash_debit",
						Value: bson.D{
							primitive.E{
								Key:   "$sum",
								Value: "$cash.debit",
							},
						},
					},
					primitive.E{
						Key: "cash_dieselDiscount",
						Value: bson.D{
							primitive.E{
								Key:   "$sum",
								Value: "$cash.dieselDiscount",
							},
						},
					},
					primitive.E{
						Key: "cash_other",
						Value: bson.D{
							primitive.E{
								Key:   "$sum",
								Value: "$cash.other",
							},
						},
					},
					primitive.E{
						Key: "cash_payout",
						Value: bson.D{
							primitive.E{
								Key:   "$sum",
								Value: "$cash.payout",
							},
						},
					},
					primitive.E{
						Key: "cash_driveOffNSF",
						Value: bson.D{
							primitive.E{
								Key:   "$sum",
								Value: "$cash.driveOffNSF",
							},
						},
					},
					primitive.E{
						Key: "cash_galesLoyaltyRedeem",
						Value: bson.D{
							primitive.E{
								Key:   "$sum",
								Value: "$cash.galesLoyaltyRedeem",
							},
						},
					},
					primitive.E{
						Key: "cash_giftCertRedeem",
						Value: bson.D{
							primitive.E{
								Key:   "$sum",
								Value: "$cash.giftCertRedeem",
							},
						},
					},
					primitive.E{
						Key: "cash_lotteryPayout",
						Value: bson.D{
							primitive.E{
								Key:   "$sum",
								Value: "$cash.lotteryPayout",
							},
						},
					},
					primitive.E{
						Key: "cash_osAdjusted",
						Value: bson.D{
							primitive.E{
								Key:   "$sum",
								Value: "$cash.osAdjusted",
							},
						},
					},
					primitive.E{
						Key: "cash_writeOff",
						Value: bson.D{
							primitive.E{
								Key:   "$sum",
								Value: "$cash.writeOff",
							},
						},
					},
					primitive.E{
						Key: "cc_amex",
						Value: bson.D{
							primitive.E{
								Key:   "$sum",
								Value: "$creditCard.amex",
							},
						},
					},
					primitive.E{
						Key: "cc_discover",
						Value: bson.D{
							primitive.E{
								Key:   "$sum",
								Value: "$creditCard.discover",
							},
						},
					},
					primitive.E{
						Key: "cc_gales",
						Value: bson.D{
							primitive.E{
								Key:   "$sum",
								Value: "$creditCard.gales",
							},
						},
					},
					primitive.E{
						Key: "cc_mastercard",
						Value: bson.D{
							primitive.E{
								Key:   "$sum",
								Value: "$creditCard.mc",
							},
						},
					},
					primitive.E{
						Key: "cc_visa",
						Value: bson.D{
							primitive.E{
								Key:   "$sum",
								Value: "$creditCard.visa",
							},
						},
					},
					primitive.E{
						Key: "fuel_1_dollar",
						Value: bson.D{
							primitive.E{
								Key:   "$sum",
								Value: "$salesSummary.fuel.fuel_1.dollar",
							},
						},
					},
					primitive.E{
						Key: "fuel_1_litre",
						Value: bson.D{
							primitive.E{
								Key:   "$sum",
								Value: "$salesSummary.fuel.fuel_1.litre",
							},
						},
					},
					primitive.E{
						Key: "fuel_2_dollar",
						Value: bson.D{
							primitive.E{
								Key:   "$sum",
								Value: "$salesSummary.fuel.fuel_2.dollar",
							},
						},
					},
					primitive.E{
						Key: "fuel_2_litre",
						Value: bson.D{
							primitive.E{
								Key:   "$sum",
								Value: "$salesSummary.fuel.fuel_2.litre",
							},
						},
					},
					primitive.E{
						Key: "fuel_3_dollar",
						Value: bson.D{
							primitive.E{
								Key:   "$sum",
								Value: "$salesSummary.fuel.fuel_3.dollar",
							},
						},
					},
					primitive.E{
						Key: "fuel_3_litre",
						Value: bson.D{
							primitive.E{
								Key:   "$sum",
								Value: "$salesSummary.fuel.fuel_3.litre",
							},
						},
					},
					primitive.E{
						Key: "fuel_4_dollar",
						Value: bson.D{
							primitive.E{
								Key:   "$sum",
								Value: "$salesSummary.fuel.fuel_4.dollar",
							},
						},
					},
					primitive.E{
						Key: "fuel_4_litre",
						Value: bson.D{
							primitive.E{
								Key:   "$sum",
								Value: "$salesSummary.fuel.fuel_4.litre",
							},
						},
					},
					primitive.E{
						Key: "fuel_5_dollar",
						Value: bson.D{
							primitive.E{
								Key:   "$sum",
								Value: "$salesSummary.fuel.fuel_5.dollar",
							},
						},
					},
					primitive.E{
						Key: "fuel_5_litre",
						Value: bson.D{
							primitive.E{
								Key:   "$sum",
								Value: "$salesSummary.fuel.fuel_5.litre",
							},
						},
					},
					primitive.E{
						Key: "fuel_6_dollar",
						Value: bson.D{
							primitive.E{
								Key:   "$sum",
								Value: "$salesSummary.fuel.fuel_6.dollar",
							},
						},
					},
					primitive.E{
						Key: "fuel_6_litre",
						Value: bson.D{
							primitive.E{
								Key:   "$sum",
								Value: "$salesSummary.fuel.fuel_6.litre",
							},
						},
					},
					primitive.E{
						Key: "total_fuelDollar",
						Value: bson.D{
							primitive.E{
								Key:   "$sum",
								Value: "$salesSummary.fuelDollar",
							},
						},
					},
					primitive.E{
						Key: "total_fuelLitre",
						Value: bson.D{
							primitive.E{
								Key:   "$sum",
								Value: "$salesSummary.fuelLitre",
							},
						},
					},
					primitive.E{
						Key: "total_nonFuel",
						Value: bson.D{
							primitive.E{
								Key:   "$sum",
								Value: "$salesSummary.totalNonFuel",
							},
						},
					},
					primitive.E{
						Key: "total_sales",
						Value: bson.D{
							primitive.E{
								Key:   "$sum",
								Value: "$salesSummary.totalSales",
							},
						},
					},
					primitive.E{
						Key: "total_cash",
						Value: bson.D{
							primitive.E{
								Key:   "$sum",
								Value: "$salesSummary.cashTotal",
							},
						},
					},
					primitive.E{
						Key: "total_cashAndCC",
						Value: bson.D{
							primitive.E{
								Key:   "$sum",
								Value: "$salesSummary.cashCCTotal",
							},
						},
					},
					primitive.E{
						Key: "total_creditCard",
						Value: bson.D{
							primitive.E{
								Key:   "$sum",
								Value: "$salesSummary.creditCardTotal",
							},
						},
					},
					primitive.E{
						Key: "overshort",
						Value: bson.D{
							primitive.E{
								Key:   "$sum",
								Value: "$overshort.amount",
							},
						},
					},
				},
			},
		},
	}

	cur, err := col.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var results []bson.M
	if err = cur.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}

	// as there is only one result, we need to extract
	if len(results) > 0 {
		day = results[0]
	}

	return day, err
}

// fetchEmployee method
func (db *MDB) fetchEmployee(attendantID primitive.ObjectID) (employee *model.Employee, err error) {

	col := db.db.Collection(colEmployees)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.D{primitive.E{Key: "_id", Value: attendantID}}
	err = col.FindOne(ctx, filter).Decode(&employee)

	return employee, err
}

// fetchJournals method
func (db *MDB) fetchJournals(recordNum string, stationID primitive.ObjectID) (journals []*model.Journal, err error) {

	col := db.db.Collection(colJournals)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	findOptions := options.Find()
	findOptions.SetSort(bson.D{primitive.E{Key: "recordNum", Value: 1}})
	filter := bson.D{
		primitive.E{Key: "recordNum", Value: recordNum},
		primitive.E{Key: "stationID", Value: stationID},
		primitive.E{Key: "type", Value: "nonFuelSaleAdjust"},
	}
	cur, err := col.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}

	if err := cur.All(ctx, &journals); err != nil {
		return nil, err
	}

	return journals, err
}

// fetchShift method
func (db *MDB) fetchShift(recordNum string, stationID primitive.ObjectID) (shift *model.Sales, err error) {

	col := db.db.Collection(colSales)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.D{primitive.E{Key: "recordNum", Value: recordNum}, primitive.E{Key: "stationID", Value: stationID}}
	err = col.FindOne(ctx, filter).Decode(&shift)

	return shift, err
}

// fetchStation method
func (db *MDB) fetchStation(stationID primitive.ObjectID) (station *model.Station, err error) {

	col := db.db.Collection(colStations)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.D{primitive.E{Key: "_id", Value: stationID}}
	err = col.FindOne(ctx, filter).Decode(&station)

	return station, err
}
