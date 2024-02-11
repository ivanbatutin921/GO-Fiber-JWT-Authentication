package initializers

import (
	//"errors"
	"fmt"
)


type Table interface {
	Migration(data interface{}) error
}

func (d *Data) MigrateTable(table Table) error {
	if err := d.DB.AutoMigrate(table); err != nil {
		return fmt.Errorf("failed to migrate table: %w", err)
	}
	return nil
}

