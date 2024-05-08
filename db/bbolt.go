package bbolt

import (
	"bytes"
	"encoding/json"
	"log"

	bolt "go.etcd.io/bbolt"
)

// Initialize initializes a new BoltDB database with the given database path and bucket name.
//
// Parameters:
// - dbPath: The path to the BoltDB database file.
// - bucketName: The name of the bucket to be created in the database.
//
// Returns:
// - *bolt.DB: The initialized BoltDB database connection.
// - error: An error if the initialization fails.
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

// Save saves the provided value in the specified bucket using the given key.
//
// Parameters:
// - db: The BoltDB database connection.
// - bucketName: The name of the bucket where the data will be stored.
// - key: The key under which the data will be stored.
// - value: The value to be saved.
// Return type: error.
func Save(db *bolt.DB, bucketName, key, value string) error {
	// Save data
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		err := b.Put([]byte(key), []byte(value))
		return err
	})
	return err
}

// SaveAsJson saves a map[string]interface{} value as JSON in a BoltDB database.
//
// Parameters:
// - db: The BoltDB database connection.
// - bucketName: The name of the bucket where the JSON data will be stored.
// - key: The key under which the JSON data will be stored.
// - value: The map[string]interface{} value to be saved as JSON.
//
// Returns:
// - map[string]interface{}: The saved value.
// - error: An error if the JSON marshaling or saving operation fails.
func SaveAsJson(db *bolt.DB, bucketName, key string, value map[string]interface{}) (map[string]interface{},error) {
	jsonStr, err := json.Marshal(&value)
	if err != nil {
		return nil, err
	}
	err = Save(db, bucketName, key, string(jsonStr))
	if err != nil {
		return nil, err
	}
	return value, nil
}

// Get retrieves the value associated with the given key from the specified bucket in the BoltDB database.
//
// Parameters:
// - db: The BoltDB database connection.
// - bucketName: The name of the bucket where the key-value pair is stored.
// - key: The key for which the value is retrieved.
//
// Returns:
// - string: The value associated with the key, or an empty string if the key is not found.
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
	//? should I return error here?
}

// GetFromJson retrieves a map[string]interface{} value from a BoltDB database using the specified bucket name and key.
//
// Parameters:
// - db: The BoltDB database connection.
// - bucketName: The name of the bucket where the key-value pair is stored.
// - key: The key for which the value is retrieved.
//
// Returns:
// - map[string]interface{}: The retrieved value as a map, or nil if an error occurs.
// - error: An error if the JSON unmarshaling operation fails.
func GetFromJson(db *bolt.DB, bucketName, key string) (map[string]interface{},error) {
	jsonStr := Get(db, bucketName, key)
	value := map[string]interface{}{}
	err := json.Unmarshal([]byte(jsonStr), &value)
	if err != nil {
		return nil, err
	}
	return value, nil
}

// GetAll retrieves all key-value pairs from the specified bucket in the BoltDB database.
//
// Parameters:
// - db: The BoltDB database connection.
// - bucketName: The name of the bucket from which the key-value pairs are retrieved.
//
// Returns:
// - map[string]string: A map containing the key-value pairs from the specified bucket.
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

// GetAllWitPrefix retrieves all key-value pairs from the specified bucket in the BoltDB database that have a given prefix.
//
// Parameters:
// - db: The BoltDB database connection.
// - bucketName: The name of the bucket from which the key-value pairs are retrieved.
// - prefix: The prefix used to filter the key-value pairs.
//
// Returns:
// - map[string]string: A map containing the key-value pairs from the specified bucket that have the given prefix.
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

// Delete deletes a key-value pair from the specified bucket in the BoltDB database.
//
// Parameters:
// - db: The BoltDB database connection.
// - bucketName: The name of the bucket from which the key-value pair is deleted.
// - key: The key of the key-value pair to be deleted.
//
// Returns:
// - error: An error if the deletion operation fails.
func Delete(db *bolt.DB, bucketName, key string) error {
	// Save data
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		err := b.Delete([]byte(key))
		return err
	})
	return err
}
