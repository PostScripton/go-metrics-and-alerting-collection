package repository

import (
	"fmt"
	"time"
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
			fmt.Printf("Error on restoring from backup storage: %s\n", err)
		}
	}

	if storeInterval == 0 {
		fmt.Println("Synchronously save to disk")
		// todo не знаю как сделать, чтобы сохраняло синхронно
	} else {
		fmt.Printf("Asynchronous save to disk with [%s] interval\n", storeInterval)
		go func() {
			if err := r.runStoring(storeInterval); err != nil {
				fmt.Printf("Error on storing backup: %s\n", err)
			}
		}()
	}
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

		collection, err := r.storage.GetCollection()
		if err != nil {
			return err
		}

		if err = r.storage.CleanUp(); err != nil {
			return err
		}

		if err := r.backupStorage.StoreCollection(collection); err != nil {
			return err
		}

		fmt.Printf("Backup stored\n")
	}
}
