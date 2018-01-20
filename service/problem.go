package service

import (
  "Barracks/data"
  "github.com/jinzhu/gorm"
)

func SelectProblemsByContest(db *gorm.DB, contest *data.Contest) (problems []data.Problem) {
  db.Select("id, code, ContestId").Where(map[string]interface{}{"contestId": contest.ID}).Find(&problems)

  return
}
