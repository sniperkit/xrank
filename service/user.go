package service

import (
  "Barracks/data"
  "github.com/jinzhu/gorm"
)

func SelectNormalUsersByContest(db *gorm.DB, contest *data.Contest) (users []data.User) {
  db.Where(map[string]interface{}{"contestId": contest.ID, "isAdmin": false}).Find(&users)

  return
}