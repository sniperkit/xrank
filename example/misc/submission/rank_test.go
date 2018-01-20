package main_test

import (
	"Barracks/data"
	"Barracks/rank"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"reflect"
	"time"
)

func createMockUser(index int, contest *data.Contest) (user data.User) {
	user = data.User{
		ID:        uint(index + 10),
		Name:      "user0" + string(int('0')+index+1),
		StrId:     "user0" + string(int('0')+index+1),
		GroupName: "group01",
		IsAdmin:   false,

		ContestID: contest.ID,
	}

	return
}

func createMockProblem(index int, contest *data.Contest) (problem data.Problem) {
	problem = data.Problem{
		ID:        uint(index),
		Code:      string(int('A') + index),
		ContestID: contest.ID,
	}

	return
}

func createMockSubmission(id uint, result int, problem data.Problem, user data.User, contest *data.Contest, t time.Time, offset time.Duration) (submission data.Submission) {
	submission = data.Submission{
		ID:        id,
		Result:    result,
		ProblemID: problem.ID,
		UserID:    user.ID,
		ContestID: contest.ID,
		CreatedAt: time.Now().Add(offset),
	}

	return
}

var _ = Describe("Rank", func() {
	Context("if contest, problems and users are prepared", func() {
		var (
			contest  *data.Contest
			users    *[]data.User
			problems *[]data.Problem
		)

		BeforeEach(func() {
			By("Creating contest instance")

			contest = &data.Contest{
				ID:    1,
				Name:  "shake17",
				Start: time.Now().Add(-30 * time.Minute),
				End:   time.Now().Add(4*time.Hour + 30*time.Minute),
			}

			By("Creating user instances")

			users = &[]data.User{}
			for i := 0; i < 5; i++ {
				*users = append(*users, createMockUser(i, contest))
			}

			By("Creating problem instances")

			problems = &[]data.Problem{}
			for i := 0; i < 5; i++ {
				*problems = append(*problems, createMockProblem(i, contest))
			}
		})

		Context("when gives contest, problems and users to rank package,", func() {
			It("should init rank data.", func() {
				By("Initiating rank data")

				nowDate := time.Now()
				rank.InitData(contest, users, problems)

				r := rank.MyRankData
				Expect(r.CalcAt).To(BeTemporally("~", nowDate))
				Expect(r.UserRows).To(HaveLen(len(*users)))
				Expect(reflect.DeepEqual(r.ContestInfo, contest)).To(Equal(true))

				u := rank.MyRankData.UserRows[0]
				Expect(u.ProblemStatuses).To(HaveLen(len(*problems)))
			})

			Context("After initiating", func() {
				BeforeEach(func() {
					By("Initiating rank data")

					rank.InitData(contest, users, problems)
				})

				It("should have correct rank data according to given submissions.", func() {
					By("Giving five submissions with four wrong and one correct")

					t := time.Now()

					submissions := &[]data.Submission{
						createMockSubmission(10, data.WRONG_ANSWER, (*problems)[0], (*users)[0], contest, t, 3*time.Second),
						createMockSubmission(11, data.MEMORY_LIMIT_EXCEED, (*problems)[0], (*users)[0], contest, t, 5*time.Second),
						createMockSubmission(12, data.RUNTIME_ERROR, (*problems)[1], (*users)[1], contest, t, 7*time.Second),
						createMockSubmission(13, data.ACCEPTED, (*problems)[0], (*users)[0], contest, t, 10*time.Second),
						createMockSubmission(14, data.TIME_LIMIT_EXCEED, (*problems)[1], (*users)[1], contest, t, 12*time.Second),
					}

					rank.AddSubmissions(submissions)

					userRows := &rank.MyRankData.UserRows
					firstProblemIndex := rank.MyRankData.ProblemMap[(*problems)[0].ID]
					secondProblemIndex := rank.MyRankData.ProblemMap[(*problems)[1].ID]

					firstUserPenaltyEst := 30*time.Minute + 10*time.Second + 2*20*time.Minute

					Expect((*userRows)[0].ProblemStatuses[firstProblemIndex].WrongCount).To(Equal(uint(2)))
					Expect((*userRows)[0].Penalty).To(BeNumerically("~", firstUserPenaltyEst, 10*time.Millisecond))

					Expect((*userRows)[1].ProblemStatuses[secondProblemIndex].WrongCount).To(Equal(uint(2)))
					Expect((*userRows)[1].Penalty).To(Equal(time.Duration(0)))

					eRanks := []uint{1, 2, 2, 2, 2}

					for index, eRank := range eRanks {
						Expect((*userRows)[index].Rank).To(Equal(uint(eRank)))
					}

					By("Third user has accept by one try")
					submissions = &[]data.Submission{
						createMockSubmission(15, data.ACCEPTED, (*problems)[0], (*users)[2], contest, t, 15*time.Second),
					}

					rank.AddSubmissions(submissions)

					Expect((*userRows)[2].ProblemStatuses[firstProblemIndex].WrongCount).To(Equal(uint(0)))
					Expect((*userRows)[2].Penalty).To(BeNumerically("~", 30*time.Minute+15*time.Second, 10*time.Millisecond))

					eRanks = []uint{2, 3, 1, 3, 3}

					for index, eRank := range eRanks {
						Expect((*userRows)[index].Rank).To(Equal(uint(eRank)))
					}

					By("First user has accept second Problem by second try and second user has accept first problem by second try")
					submissions = &[]data.Submission{
						createMockSubmission(16, data.RUNTIME_ERROR, (*problems)[1], (*users)[0], contest, t, 17*time.Second),
						createMockSubmission(17, data.ACCEPTED, (*problems)[1], (*users)[1], contest, t, 19*time.Second),
						createMockSubmission(18, data.ACCEPTED, (*problems)[1], (*users)[0], contest, t, 22*time.Second),
					}

					rank.AddSubmissions(submissions)

					eRanks = []uint{1, 3, 2, 4, 4}

					firstUserPenaltyEst += 30*time.Minute + 22*time.Second + 20*time.Minute
					Expect((*userRows)[0].Penalty).To(BeNumerically("~", firstUserPenaltyEst, 10*time.Millisecond))
					Expect((*userRows)[0].ProblemStatuses[secondProblemIndex].WrongCount).To(Equal(uint(1)))

					for index, eRank := range eRanks {
						Expect((*userRows)[index].Rank).To(Equal(uint(eRank)))
					}
				})
			})
		})
	})
})
