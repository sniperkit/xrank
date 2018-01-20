package rank

import (
  "time"
  "Barracks/data"
  "container/heap"
)

type RankInfo struct {
  RankData *rankData
  RankHeap *rankHeap
}

func newUserRow (user *data.User, problems *[]data.Problem) (u userRow) {
  u = userRow{
    Rank: 1,
    StrId: (*user).StrId,
    ID: (*user).ID,
    ProblemStatuses: make([]problemStatus, len(*problems)),
  }

  return
}

func newRankData (contest *data.Contest, users *[]data.User, problems *[]data.Problem) (r *rankData) {

  r = &rankData{
    CalcAt: time.Now(),
    ContestInfo: contest,
    UserRows: make([]userRow, len(*users)),
    UserMap: make(map[uint]uint),
    ProblemMap: make(map[uint]uint),
    ProblemCodeMap: make(map[uint]string),
  }

  for index, problem := range *problems {
    r.ProblemMap[problem.ID] = uint(index)
    r.ProblemCodeMap[problem.ID] = problem.Code
  }

  for index, user := range *users {
    r.UserMap[user.ID] = uint(index)
    r.UserRows[index] = newUserRow(&user, problems)
  }

  return
}

func NewRankInfo (contest *data.Contest, users *[]data.User, problems *[]data.Problem) (r *RankInfo){
  r = &RankInfo{}
  r.RankData = newRankData(contest, users, problems)

  r.RankHeap = &rankHeap{}
  heap.Init(r.RankHeap)

  return
}

func (r RankInfo) calcRanks() {
  for index, userRow := range r.RankData.UserRows {
    heap.Push(r.RankHeap, rankNode{Penalty: userRow.Penalty, AcceptedCnt: userRow.AcceptedCnt, UserIndex: uint(index)})
  }

  rankValue := uint(1)
  var beforeRankNode *rankNode

  for r.RankHeap.Len() > 0 {
    popRankNode := rankNode(heap.Pop(r.RankHeap).(rankNode))

    if beforeRankNode != nil && (beforeRankNode.AcceptedCnt != popRankNode.AcceptedCnt || beforeRankNode.Penalty != popRankNode.Penalty) {
      rankValue++
    }

    r.RankData.UserRows[popRankNode.UserIndex].Rank = rankValue
    beforeRankNode = &popRankNode
  }
}

func (r RankInfo) analyzeSubmissions(submissions []data.Submission) {
  contestInfo := r.RankData.ContestInfo

  // 각각의 제출에 대하여
  for _, submission := range submissions {
    // userRow와 problemStatus를 구한다.

    if _, ok := r.RankData.UserMap[submission.UserID]; !ok {
      continue
    }

    userRow := &r.RankData.UserRows[r.RankData.UserMap[submission.UserID]]
    problemStatus := &userRow.ProblemStatuses[r.RankData.ProblemMap[submission.ProblemID]]

    // 만약 제출이 정답 소스코드라면
    if data.IsAccepted(submission.Result) {
      // 문제가 맞지 않은 상황이라면
      if !problemStatus.Accepted {
        penalty := submission.CreatedAt.Sub(contestInfo.Start) + time.Duration(problemStatus.WrongCount) * 20 * time.Minute

        (*userRow).Penalty += penalty
        (*problemStatus).Accepted = true
        (*userRow).AcceptedCnt++
      }

    } else { // 제출이 틀린 소스코드라면
      if !problemStatus.Accepted {
        (*problemStatus).WrongCount++
      }
    }
  }
}

func (r RankInfo) AddSubmissions (submissions []data.Submission) {
  r.analyzeSubmissions(submissions)
  r.calcRanks()
}

func (r RankInfo) GetUserProblemStatusSummary (userId uint) (summary []problemStatusSummary) {
  mappedId, ok := r.RankData.UserMap[userId]
  if !ok {
    summary = nil
    return
  }

  userRowRef := &r.RankData.UserRows[mappedId]
  summary = make([]problemStatusSummary, len(r.RankData.ProblemMap))
  idx := 0

  for key, val := range r.RankData.ProblemMap {
    summary[idx] = problemStatusSummary{
      ProblemCode: r.RankData.ProblemCodeMap[key],
      ProblemId: key,
      Accepted: userRowRef.ProblemStatuses[val].Accepted,
      Wrong: !summary[idx].Accepted && userRowRef.ProblemStatuses[val].WrongCount > 0,
    }

    idx++
  }

  return
}

func (r RankInfo) GetUserSummary(userId uint, subId uint) (summary *UserRankSummary) {
  mappedId, ok := r.RankData.UserMap[userId]
  if !ok {
    summary = nil
    return
  }

  userRowRef := &r.RankData.UserRows[mappedId]
  summary = &UserRankSummary{
    LastSubId: subId,
    Penalty: userRowRef.Penalty,
    StrId: userRowRef.StrId,
    UserId: userId,
    AcceptedCnt: userRowRef.AcceptedCnt,
    Rank: userRowRef.Rank,
    ProblemStatus: r.GetUserProblemStatusSummary(userId),
  }

  return
}

func (r RankInfo) GetRanking () (summary []userRow) {
  return r.RankData.UserRows
}