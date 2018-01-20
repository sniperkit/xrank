package data

import "github.com/jinzhu/gorm"

func InitDB () (db *gorm.DB){
  var err error

  db, err = gorm.Open("mysql", dbConfig)

  if err != nil {
    panic(err)
  }

  db.LogMode(true)

  return
}
