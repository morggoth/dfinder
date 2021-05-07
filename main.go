package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Hashes struct {
	ID    uint
	Hash  string `gorm:"unique"`
	Files []Files
}

type Files struct {
	ID       uint
	FilePath string `gorm:"unique"`
	HashesID uint
}

func dbInteract(hash, path string, db *gorm.DB) {
	hashInstance := Hashes{Hash: hash}
	log.Printf("hashInstance struct now is: %v", hashInstance)

	if result := db.Create(&hashInstance); result.Error != nil {
		if result.Error.Error() == "UNIQUE constraint failed: hashes.hash" {
			log.Printf("Hash %s already exists", hashInstance.Hash)

			db.Where("hash = ?", hash).First(&hashInstance)
			log.Printf("hashInstance struct now is: %v", hashInstance)
		} else {
			log.Fatal(result.Error)
		}
	}

	fileInstance := Files{FilePath: path, HashesID: hashInstance.ID}
	log.Printf("fileInstance struct now is: %v", fileInstance)

	if result2 := db.Create(&fileInstance); result2.Error != nil {
		if result2.Error.Error() == "UNIQUE constraint failed: files.file_path" {
			log.Printf("File %s already exists", fileInstance.FilePath)
		} else {
			log.Fatal(result2.Error)
		}
	}
}

func check_error(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

// func hashWalk(path string, info os.FileInfo, err error) error {
// 	if info.IsDir() {
// 		log.Printf("%s\n", info.Name())
// 	} else {
// 		dat, err := ioutil.ReadFile(path)
// 		if err != nil {
// 			return err
// 		}

// 		sha256Hash := sha256.New()

// 		_, err = io.Copy(sha256Hash, bytes.NewReader(dat))
// 		if err != nil {
// 			return err
// 		}

// 		sum := sha256Hash.Sum(nil)

// 		absPath, err := filepath.Abs(path)
// 		if err != nil {
// 			return err
// 		}

// 		log.Printf("%x	%s\n", sum, absPath)
// 	}

// 	return nil
// }

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&Hashes{}, &Files{})

	var walk = func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			log.Printf("%s\n", info.Name())
		} else {
			dat, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			sha256Hash := sha256.New()

			_, err = io.Copy(sha256Hash, bytes.NewReader(dat))
			if err != nil {
				return err
			}

			sum := hex.EncodeToString(sha256Hash.Sum(nil))

			absPath, err := filepath.Abs(path)
			if err != nil {
				return err
			}

			dbInteract(sum, absPath, db)

			log.Printf("%s	%s\n", sum, absPath)
		}

		return nil
	}

	err = filepath.Walk(".", walk)
	check_error(err)

	// hash := "3e90dcae65ccd49e160b02a5dd83f649c6be27bc9765cd038a9d1d2547a179fb"
	// path := "/Users/morggoth/wk/p/dfinder/tests/MD5SUMS"

	// dbInteract(hash, path, db)

	// hash1 := "another_hash"
	// path1 := "another_path"
	// dbInteract(hash1, path1, db)

	// dat, err := ioutil.ReadFile("tests/groovy-server-cloudimg-armhf-root.tar.xz")
	// check_error(err)

	// // md5Hash := md5.New()
	// sha256Hash := sha256.New()

	// _, err = io.Copy(sha256Hash, bytes.NewReader(dat))
	// check_error(err)

	// sum := sha256Hash.Sum(nil)

	// fmt.Printf("%x\n", sum)
	// // sha1Hash := sha1.New()
	// // sha256Hash := sha256.New()

}
