package bufinsert

import (
	"gorm.io/gorm"
)

func ExampleNewInserter() {
	var db = &gorm.DB{} // fake database
	var inserter = NewInserter(db, 2)

	for i := 0; i < 11; i++ {
		go inserter.Insert(i)
	}

	inserter.Flush() // flush and empty
}

func ExampleBufferedInserter_Insert() {
	var db = &gorm.DB{} // fake database
	var inserter = NewInserter(db, 2)

	// way 1
	inserter.Insert("1") // not flush
	inserter.Insert("2") // flush and empty
	inserter.Insert("3") // not flush
	// way 2
	inserter.Insert([]interface{}{"4", "5", "6", "7"})

	inserter.Flush() // flush and empty
}
