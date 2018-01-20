package main

import (
	"Barracks/data"
	"Barracks/httpserver"
	"Barracks/poller"
	"Barracks/rank"
	"Barracks/service"
	"flag"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"os"
)

func main() {
	contestNamePtr := flag.String("contest", "", "contest name")
	portPtr := flag.Int("port", 8080, "port number")
	pushHostPtr := flag.String("pushHost", "", "host domain to push submission info")

	flag.Parse()

	if *contestNamePtr == "" || *pushHostPtr == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	db := data.InitDB()
	defer db.Close()

	contestName := *contestNamePtr

	contest := service.SelectContestByName(db, &contestName)
	users := service.SelectNormalUsersByContest(db, &contest)
	problems := service.SelectProblemsByContest(db, &contest)

	doneChan := make(chan bool)

	rankInfo := rank.NewRankInfo(&contest, &users, &problems)
	rankInfoFreezed := rank.NewRankInfo(&contest, &users, &problems)

	poller.StartPoll(db, rankInfo, rankInfoFreezed, &contest, &doneChan, pushHostPtr)
	httpserver.StartServer(rankInfo, rankInfoFreezed, uint(*portPtr))
}
