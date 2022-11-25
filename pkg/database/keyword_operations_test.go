package database

import (
	"fmt"
	"os"
	"testing"

	"gorm.io/gorm"
)

const TestDBFileName = "db_for_test.db"

func FuzzGetKwdsForVersion(f *testing.F) {
	testDb, err := Create(TestDBFileName)
	if err != nil {
		f.Errorf("Failed to create test db object")
	}
	// teardown: try to remove temp db file created for testing
	f.Cleanup(func() {
		dbInstance, _ := testDb.DB()
		_ = dbInstance.Close()
		_ = os.Remove(TestDBFileName)
	})

	err = populateMockKeywords(testDb)
	if err != nil {
		f.Errorf("Failed to populate test db with mock data")
	}

	f.Add(6, 3, Business)
	f.Add(3, 3, Business)
	f.Add(2, 2, Business)
	f.Add(1, 1, Business)

	f.Add(6, 3, Technical)
	f.Add(5, 2, Technical)
	f.Add(4, 1, Technical)
	f.Add(1, 0, Technical)

	f.Add(6, 6, All)
	f.Add(3, 3, All)
	f.Add(1, 1, All)
	f.Add(0, 0, All)

	f.Fuzz(func(t *testing.T, version int, count int, kwType string) {
		keywordsForVersion, _ := KeywordsForVersion(testDb, version, kwType)

		if len(keywordsForVersion) != count {
			t.Errorf(
				"Wrong number of business kwds found: expected %d but got %d",
				count,
				len(keywordsForVersion),
			)
		}
	})

}

func populateMockKeywords(db *gorm.DB) error {
	var kwType string
	for i := 1; i <= 6; i++ {
		if i <= 3 {
			kwType = Business
		} else {
			kwType = Technical
		}
		// TODO: InsertNewKeyword is not tested yet
		err := InsertNewKeyword(
			db,
			fmt.Sprintf("Name%d", i),
			fmt.Sprintf("Args%d", i),
			fmt.Sprintf("Docs%d", i),
			kwType,
		)
		if err != nil {
			return err
		}
	}
	return nil
}
