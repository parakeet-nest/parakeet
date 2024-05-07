package bbolt

import (
	"bytes"
	"log"

	bolt "go.etcd.io/bbolt"
)


func Initialize(dbPath, bucketName string) (*bolt.DB, error) {
	//db, err := bolt.Open("my.db", 0600, &bolt.Options{Timeout: 1 * time.Second})

	db, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	//defer db.Close()

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return err
		}
		return nil
	})
	return db, err
}

func Save(db *bolt.DB, bucketName, key, value string) error {
	// Save data
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		err := b.Put([]byte(key), []byte(value))
		return err
	})
	return err
}

func Get(db *bolt.DB, bucketName, key string) string {

	var strReturnValue string

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		bufferValue := b.Get([]byte(key))

		if bufferValue == nil {
			strReturnValue = ""
		} else {
			strReturnValue = string(bufferValue)
		}
		return nil
	})
	return strReturnValue
}

func GetAll(db *bolt.DB, bucketName string) map[string]string {
	records := map[string]string{}

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		c := b.Cursor()

		for bufferKey, bufferValue := c.First(); bufferKey != nil; bufferKey, bufferValue = c.Next() {
			records[string(bufferKey)] = string(bufferValue)
		}

		return nil
	})
	return records
}

func GetAllWitPrefix(db *bolt.DB, bucketName, prefix string) map[string]string {
	records := map[string]string{}

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		c := b.Cursor()
		bufferPrefix := []byte(prefix)

		for bufferKey, bufferValue := c.Seek(bufferPrefix); bufferKey != nil && bytes.HasPrefix(bufferKey, bufferPrefix); bufferKey, bufferValue = c.Next() {
			records[string(bufferKey)] = string(bufferValue)
		}

		return nil
	})
	return records
}

func Delete(db *bolt.DB,bucketName, key string) error {
	// Save data
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		err := b.Delete([]byte(key))
		return err
	})
	return err
}
