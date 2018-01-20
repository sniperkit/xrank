package rank

import (
  "time"
  "Barracks/data"
)

type problemStatusSummary struct {
  ProblemId uint `json:"problemId"`
  ProblemCode string `json:"problemCode"`
  Accepted  bool `json:"accepted"`
  Wrong     bool `json:"wrong"`
}

type UserRankSummary struct {
  LastSubId uint                        `json:"lastSubId"`
  StrId string                          `json:"strId"`
  Penalty         time.Duration         `json:"penalty"`
  UserId uint                           `json:"userId"`
  Rank uint                             `json:"rank"`
  AcceptedCnt uint                      `json:"acceptedCnt"`
  ProblemStatus []problemStatusSummary  `json:"problemStatus,omitempty"`
}

type problemStatus struct {
  WrongCount      uint              `json:"wrongCnt"`
  Status          string            `json:"status"`
  Accepted        bool              `json:"accepted"`
}

type userRow struct {
  Rank            uint              `json:"rank"`
  StrId           string            `json:"strId"`
  ID              uint              `json:"id"`
  AcceptedCnt     uint              `json:"acceptedCnt"`
  Penalty         time.Duration     `json:"penalty"`
  ProblemStatuses []problemStatus
}

type rankData struct {
  CalcAt      time.Time
  UserRows    []userRow
  ContestInfo *data.Contest
  UserMap     map[uint]uint
  ProblemMap  map[uint]uint
  ProblemCodeMap map[uint]string
}

type rankNode struct {
  UserIndex     uint
  Penalty       time.Duration
  AcceptedCnt   uint
}

type rankHeap []rankNode

func (h rankHeap) Len() int {
  return len(h)
}

func (h rankHeap) Less(i, j int) bool {
  if h[i].AcceptedCnt == h[j].AcceptedCnt {
    return h[i].Penalty < h[j].Penalty
  }

  return h[i].AcceptedCnt > h[j].AcceptedCnt
}

func (h rankHeap) Swap(i, j int) {
  h[i], h[j] = h[j], h[i]
}

func (h *rankHeap) Push(element interface{}) {
  *h = append(*h, element.(rankNode))
}

func (h *rankHeap) Pop() interface{} {
  old := *h
  n := len(old)
  element := old[n-1]
  *h = old[0 : n-1]
  return element
}