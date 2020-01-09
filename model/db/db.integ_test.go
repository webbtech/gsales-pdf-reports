package db

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/pulpfree/gsales-pdf-reports/config"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	date         = "2019-12-21"
	defaultsFP   = "../../config/defaults.yml"
	timeForm     = "2006-01-02"
	recordNum    = "2019-12-21-2"
	stationIDStr = "56cf1815982d82b0f3000001"
)

// IntegSuite struct
type IntegSuite struct {
	cfg       *config.Config
	db        *MDB
	stationID primitive.ObjectID
	suite.Suite
}

// SetupTest method
func (s *IntegSuite) SetupTest() {
	// setup config
	os.Setenv("Stage", "test")
	s.cfg = &config.Config{DefaultsFilePath: defaultsFP}
	err := s.cfg.Load()
	s.NoError(err)

	// Set client options
	clientOptions := options.Client().ApplyURI(s.cfg.GetMongoConnectURL())

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	s.NoError(err)

	s.db = &MDB{
		client: client,
		dbName: s.cfg.DBName,
		db:     client.Database(s.cfg.DBName),
	}

	s.stationID, _ = primitive.ObjectIDFromHex(stationIDStr)
}

// ===================== Exported Functions ================================================ //

// TestNewDB method
func (s *IntegSuite) TestNewDB() {
	_, err := NewDB(s.cfg.GetMongoConnectURL(), s.cfg.DBName)
	s.NoError(err)
}

// TestfetchDay method
func (s *IntegSuite) TestfetchDay() {
	dte, _ := time.Parse(timeForm, date)
	day, err := s.db.fetchDay(dte, s.stationID)
	s.NoError(err)
	s.NotNil(day)
	// fmt.Printf("day %+v\n", day)
}

// TestfetchShift method
func (s *IntegSuite) TestfetchShift() {
	shift, err := s.db.fetchShift(recordNum, s.stationID)
	s.NoError(err)
	s.NotNil(shift)
}

// TestfetchEmployee method
func (s *IntegSuite) TestfetchEmployee() {
	shift, err := s.db.fetchShift(recordNum, s.stationID)
	s.NoError(err)
	employee, err := s.db.fetchEmployee(shift.Attendant.ID)
	s.Equal(employee.ID, shift.Attendant.ID)
}

// TestfetchJournals method
func (s *IntegSuite) TestfetchJournals() {
	recordNum := "2019-10-31-1"
	journals, err := s.db.fetchJournals(recordNum, s.stationID)
	s.NoError(err)
	s.True(len(journals) > 0)
}

// TestGetDay method
func (s *IntegSuite) TestGetDay() {
	dte, _ := time.Parse(timeForm, date)
	day, err := s.db.GetDay(dte, s.stationID)
	s.NoError(err)
	s.NotNil(day)

	futureDate := "2202-02-02"
	dte, _ = time.Parse(timeForm, futureDate)
	day, err = s.db.GetDay(dte, s.stationID)
	s.Error(err)
}

// TestGetShift method
func (s *IntegSuite) TestGetShift() {
	shift, err := s.db.GetShift(recordNum, s.stationID)
	s.NoError(err)
	s.True(shift != nil)

	futureDatedRecordNum := "2202-02-02-1"
	shift, err = s.db.GetShift(futureDatedRecordNum, s.stationID)
	s.Error(err)
}

// TestIntegrationSuite function
func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegSuite))
}
