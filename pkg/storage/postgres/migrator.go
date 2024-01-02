package postgres

import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun/migrate"
)

func (t *storage) runMigrator(migrations *migrate.Migrations) error {
	mgrtr := migrate.NewMigrator(t.db, migrations)
	err := mgrtr.Init(context.Background())
	if err != nil {
		log.Error().Msg(err.Error())
	}
	if err := mgrtr.Lock(context.Background()); err != nil {
		log.Error().Msg(err.Error())
	}
	group, err := mgrtr.Migrate(context.Background())
	if err != nil {
		log.Warn().Msg(err.Error())
		mgrtr.Unlock(context.Background())
		return err
	}
	if group.IsZero() {
		log.Warn().Msg("there are no new migrations to run (database is up to date)")
	}
	log.Info().Msgf("migrated to %s\\n", group)
	mgrtr.Unlock(context.Background())
	return nil
}
