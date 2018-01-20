package service

import (
  "github.com/jinzhu/gorm"
  "Barracks/data"
  "time"
)

func SelectContestByName (db *gorm.DB, contestName *string) (contest data.Contest) {
  db.Where(&data.Contest{Name: *contestName}).First(&contest)

  if contest.ID == 0 { // the case not to find contest by name
    panic("해당하는 콘테스트가 없습니다!")
  }

  return
}

func SelectContestFreezeById (db *gorm.DB, contestId uint) (freezeAt time.Time) {
  var contest data.Contest

  db.Select("freezeAt").Where(contestId).First(&contest)
  freezeAt = contest.FreezeAt

  return
}