package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// User has many CreditCards, UserID is the foreign key
type User struct {
	gorm.Model
	CreditCards []CreditCard
}

type CreditCard struct {
	gorm.Model
	Number string
	UserID uint
}

type Hashes struct {
	ID    uint
	Hash  string `gorm:"unique"`
	Files []Files
}

type Files struct {
	ID       uint
	FilePath string
	HashesID uint
}

// type FileHashes struct {
// 	ID      uint
// 	Hash    string `gorm:"unique"`
// 	Counter int
// }

// type Files struct {
// 	gorm.Model
// 	FilePath string
// 	// HashId   int
// 	HashId FileHashes
// }

func dbInteract(hash, path string, db *gorm.DB) {
	// db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	// if err != nil {
	// 	panic("failed to connect database")
	// }

	// // Migrate the schema
	// // db.AutoMigrate(&Product{})
	// db.AutoMigrate(&Hashes{}, &Files{})

	// // Start Association Mode
	// var hashes Hashes
	// db.Model(&hashes).Association("Hash")

	// user := User{Name: "Jinzhu", Age: 18, Birthday: time.Now()}

	// result := db.Create(&user) // pass pointer of data to Create

	// user.ID             // returns inserted data's primary key
	// result.Error        // returns error
	// result.RowsAffected // returns inserted records count

	// Create
	// db.Create(&Product{Code: "D42", Price: 100})

	hashInstance := Hashes{Hash: hash}
	if result := db.Create(&hashInstance); result.Error != nil {
		if result.Error.Error() == "UNIQUE constraint failed: hashes.hash" {
			fmt.Println("not a unique value")
		} else {
			// log.Fatal(result.Error)
			fmt.Println(result.Error)
			fmt.Println("some error")
			fmt.Printf("%+v\n", result)
		}
	}

	// result := db.Where("hash = ?", hash).First(&hashes)

	fmt.Printf("%+v\n", hashInstance)
	db.Create(&Files{FilePath: path, HashesID: hashInstance.ID})

	// result2 := db.Create(&hashInstance)
	// fmt.Printf("%+v\n", result2)

	// db.Create(&Files{FilePath: path})
	// db.Create(&Files{FilePath: "/Users/morggoth/wk/p/dfinder/tests/MD5SUMS", HashId: 1})
	// db.Create(&Files{FilePath: path})

	// // Start Association Mode
	// var user User
	// db.Model(&user).Association("Languages")
	// // `user` is the source model, it must contains primary key
	// // `Languages` is a relationship's field name
	// // If the above two requirements matched, the AssociationMode should be started successfully, or it should return error
	// db.Model(&user).Association("Languages").Error

	// 3e90dcae65ccd49e160b02a5dd83f649c6be27bc9765cd038a9d1d2547a179fb        /Users/morggoth/wk/p/dfinder/tests/MD5SUMS

	// // Read
	// var product Product
	// db.First(&product, 1)                 // find product with integer primary key
	// db.First(&product, "code = ?", "D42") // find product with code D42

	// // Update - update product's price to 200
	// db.Model(&product).Update("Price", 200)
	// // Update - update multiple fields
	// db.Model(&product).Updates(Product{Price: 200, Code: "F42"}) // non-zero fields
	// db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

	// // Delete - delete product
	// db.Delete(&product, 1)
}

func check_error(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func hashWalk(path string, info os.FileInfo, err error) error {
	if info.IsDir() {

		fmt.Printf("%s\n", info.Name())
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

		sum := sha256Hash.Sum(nil)

		absPath, err := filepath.Abs(path)
		if err != nil {
			return err
		}

		fmt.Printf("%x	%s\n", sum, absPath)
	}

	return nil
}

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	// db.AutoMigrate(&Product{})
	db.AutoMigrate(&Hashes{}, &Files{})

	// err := filepath.Walk(".", hashWalk)
	// check_error(err)

	hash := "3e90dcae65ccd49e160b02a5dd83f649c6be27bc9765cd038a9d1d2547a179fb"
	path := "/Users/morggoth/wk/p/dfinder/tests/MD5SUMS"

	dbInteract(hash, path, db)

	hash1 := "another_hash"
	path1 := "another_path"
	dbInteract(hash1, path1, db)

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
