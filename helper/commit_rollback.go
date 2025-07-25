package helper

import (
	"database/sql"
	"log"
)

func CommitOrRollback(tx *sql.Tx) {
	err := recover()
	if err != nil {
		rollBackErr := tx.Rollback()
		if rollBackErr != nil {
			log.Printf("failed to rollback transaction: %v", rollBackErr)
		}
	} else {
		commitErr := tx.Commit()
		if commitErr != nil {
			log.Printf("failed to commit transaction: %v", commitErr)
		} else {
			log.Printf("transaction committed successfully")
		}
	}
}
