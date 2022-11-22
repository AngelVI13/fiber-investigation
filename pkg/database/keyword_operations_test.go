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

	f.Add(6, 3, "business")
	f.Add(3, 3, "business")
	f.Add(2, 2, "business")
	f.Add(1, 1, "business")

	f.Add(6, 3, "technical")
	f.Add(5, 2, "technical")
	f.Add(4, 1, "technical")
	f.Add(1, 0, "technical")

	f.Fuzz(func(t *testing.T, version int, count int, kwType string) {
		keywordsForVersion, _ := GetAllKeywordsForVersion(testDb, version, kwType)

		if len(keywordsForVersion) != count {
			t.Errorf(
				"Wrong number of business kwds found: expected %d but got %d",
				count,
				len(keywordsForVersion),
			)
		}
		var found bool
		var expId int
		// for each count of keywords make sure ids are correct.
		// business ids 1->3, technical 4->6
		for idx := 1; idx <= count; idx++ {
			found = false
			for _, kwd := range keywordsForVersion {
				expId = idx
				if kwd.KwType == "technical" {
					expId = expId + 3
				}
				if int(kwd.ID) == expId {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Failed to found keyword with id: %v", keywordsForVersion)
			}
		}
	})

}

func populateMockKeywords(db *gorm.DB) error {
	var kwType string
	for i := 1; i <= 6; i++ {
		if i <= 3 {
			kwType = "business"
		} else {
			kwType = "technical"
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
