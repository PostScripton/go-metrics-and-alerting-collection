package storage

import (
	"time"

	log "github.com/rs/zerolog/log"
)

type Restorer struct {
	backupStorage Storager
	storage       Storager
}

func NewRestorer(backupStorage Storager, storage Storager) *Restorer {
	return &Restorer{
		backupStorage: backupStorage,
		storage:       storage,
	}
}

func (r *Restorer) Run(shouldRestore bool, storeInterval time.Duration) {
	if shouldRestore {
		if err := r.restore(); err != nil {
			log.Warn().Err(err).Msg("Error on restoring from backup storage")
		}
	}

	if storeInterval == 0 {
		log.Info().Msg("Synchronously save to disk")
		// todo не знаю как сделать, чтобы сохраняло синхронно
	} else {
		log.Info().Dur("interval", storeInterval).Msg("Asynchronous save to disk")
		go func() {
			if err := r.runStoring(storeInterval); err != nil {
				log.Warn().Err(err).Msg("Storing backup")
			}
		}()
	}
}

func (r *Restorer) Store() error {
	collection, err := r.storage.GetCollection()
	if err != nil {
		return err
	}

	if err = r.storage.CleanUp(); err != nil {
		return err
	}

	if err = r.backupStorage.StoreCollection(collection); err != nil {
		return err
	}

	log.Print("Backup stored")
	return nil
}

func (r *Restorer) restore() error {
	collection, err := r.backupStorage.GetCollection()
	if err != nil {
		return err
	}

	if err = r.storage.CleanUp(); err != nil {
		return err
	}

	if err = r.storage.StoreCollection(collection); err != nil {
		return err
	}

	return nil
}

func (r *Restorer) runStoring(interval time.Duration) error {
	storeInterval := time.NewTicker(interval)
	for {
		<-storeInterval.C

		if err := r.Store(); err != nil {
			return err
		}
	}
}
