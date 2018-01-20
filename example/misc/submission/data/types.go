package data

import "time"

type Contest struct {
  ID          uint        `gorm:"column:id;primary_key"`
  Name        string

  FreezeAt    time.Time   `gorm:"column:freezeAt"`
  Start       time.Time
  End         time.Time
}

func (Contest) TableName() string {
  return "Contests"
}

type User struct {
  ID          uint        `gorm:"column:id;primary_key"`
  Name        string
  StrId       string      `gorm:"column:strId"`
  GroupName   string      `gorm:"column:groupName"`
  IsAdmin     bool        `gorm:"column:isAdmin"`

  ContestID   uint        `gorm:"column:ContestId"`
}

func (User) TableName() string {
  return "Users"
}

type Problem struct {
  ID          uint        `gorm:"column:id;primary_key"`
  Code        string      `gorm:"type:text"`

  ContestID   uint        `gorm:"column:ContestId"`
}

func (Problem) TableName() string {
  return "Problems"
}

type Submission struct {
  ID            uint      `gorm:"column:id;primary_key"`
  Result        int
  ProblemCode   string

  CreatedAt     time.Time `gorm:"column:createdAt"`

  ContestID     uint       `gorm:"column:ContestId"`
  ProblemID     uint       `gorm:"column:ProblemId"`
  UserID        uint       `gorm:"column:UserId"`
}

func (Submission) TableName() string {
  return "Submissions"
}

const(
  ACCEPTED = 4
  PRESENTATION_ERROR = 5
  WRONG_ANSWER = 6
  TIME_LIMIT_EXCEED = 7
  MEMORY_LIMIT_EXCEED = 8
  OUTPUT_LIMIT_EXCEED = 9
  RUNTIME_ERROR = 10
  COMPILE_ERROR = 11
)

func IsAccepted(status int) bool {
  return status == ACCEPTED
}

func IsError(status int) bool {
  return status != ACCEPTED
}