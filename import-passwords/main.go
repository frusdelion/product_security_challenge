package main

import (
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/klauspost/password"
	"github.com/klauspost/password/drivers/boltpw"
	"github.com/klauspost/password/tokenizer"
	"github.com/snwfdhmp/errlog"
	"os"
)

func Import(mem password.DbWriter) {
	r, err := os.Open("crackstation-human-only.txt.gz")
	if err != nil {
		panic(err)
	}
	defer r.Close()

	in, err := tokenizer.NewGzLine(r)
	if err != nil {
		panic(err)
	}

	err = password.Import(in, mem, nil)
	if err != nil {
		panic(err)
	}
}

var (
	CommonPasswordsDbFile = "../common-passwords.db"
)

func main() {
	file, err := os.Open(CommonPasswordsDbFile)
	if errlog.Debug(err) {
		if !os.IsNotExist(err) {
			panic(err)
		}
	}

	if os.IsExist(err) {
		panic(errors.New(fmt.Sprintf("%s already exists. quitting import", CommonPasswordsDbFile)))
		defer file.Close()
	}

	// Open the database using the Bolt driver
	// You probably have this elsewhere if you already use Bolt
	db, err := bolt.Open(CommonPasswordsDbFile, 0666, nil)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Use the driver to read/write to the bucket "commonpwd"
	chk, err := boltpw.New(db, "commonpwd")
	if err != nil {
		panic(err)
	}

	Import(chk)
}
