package httpserver

import (
  "github.com/labstack/echo"
  "strconv"
  "net/http"
  "Barracks/rank"
  //"github.com/labstack/echo/middleware"
)

func StartServer(rankInfo *rank.RankInfo, rankInfoFreeze *rank.RankInfo, port uint) {
  e := echo.New()

  //e.Use(middleware.Logger())

  e.GET("/api/:contestName/ranking", func(ctx echo.Context) error {
    if rankInfoFreeze.RankData.ContestInfo.Name != ctx.Param("contestName") {
      return ctx.NoContent(http.StatusNotFound)
    }

    //r := rankInfoFreeze.GetRanking()
    //if r == nil {
    //  return ctx.NoContent(http.StatusNotFound)
    //}

    var r []rank.UserRankSummary

    for k := range rankInfoFreeze.RankData.UserMap {
      r = append(r, *rankInfoFreeze.GetUserSummary(uint(k), 0))
    }

    return ctx.JSON(http.StatusOK, r)
  })

  e.GET("/api/acceptedCnts/:userId", func(ctx echo.Context) error {
    userId, err := strconv.Atoi(ctx.Param("userId"))

    if err != nil {
      return ctx.NoContent(http.StatusNotFound)
    }

    r := rankInfo.GetUserSummary(uint(userId), 0)

    if r == nil {
      return ctx.NoContent(http.StatusNotFound)
    }

    r.Rank = rankInfoFreeze.GetUserSummary(uint(userId), 0).Rank
    return ctx.JSON(http.StatusOK, r)
  })

  e.GET("/api/problemStatuses/:userId", func(ctx echo.Context) error {
    userId, err := strconv.Atoi(ctx.Param("userId"))
    if err != nil {
      return ctx.NoContent(http.StatusNotFound)
    }

    r := rankInfo.GetUserProblemStatusSummary(uint(userId))
    if r == nil {
      return ctx.NoContent(http.StatusNotFound)
    }
    return ctx.JSON(http.StatusOK, r)
  })

  e.Logger.Fatal(e.Start(":"+strconv.Itoa(int(port))))
}
