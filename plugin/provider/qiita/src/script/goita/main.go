package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
	"github.com/mattn/go-runewidth"
	"github.com/urfave/cli"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

type entry struct {
	Title     string
	URL       string
	LikeCount string
}

type QueryResult struct {
	Title   string
	Entries []*entry
}

func main() {
	app := cli.NewApp()
	app.Name = "goita"
	app.Usage = "Command Line Client for Qiita ranking"
	app.Version = "0.0.1"
	app.ArgsUsage = "[tag]"
	app.HideHelp = true
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "number, n",
			Value: 10,
			Usage: "number of output lines",
		},
		cli.IntFlag{
			Name:  "days, d",
			Value: 1,
			Usage: "number of report days",
		},
		cli.BoolFlag{
			Name:  "hatebu, hatena",
			Usage: "ranking sorted by hatebu bookmark number",
		},
	}
	app.Action = func(c *cli.Context) error {
		tag := c.Args().First()
		number := c.Int("number")
		days := c.Int("days")
		ishatebu := c.Bool("hatebu")
		url := buildUrl(tag, days, ishatebu)
		result, err := crawl(url, number)
		if err != nil {
			return err
		}
		showResult(result, url)
		return nil
	}
	app.Run(os.Args)
}

func buildUrl(tag string, days int, ishatebu bool) string {
	var tag_parameter, order string

	if tag != "" {
		tag_parameter = "&tag=" + tag
	}
	if ishatebu {
		order = "&orderby=hatebu_count"
	} else {
		order = "&orderby=like_count"
	}
	d := strconv.Itoa(days)
	return fmt.Sprintf("http://qrank.wbsrv.net/entries?days=%s%s%s", d, tag_parameter, order)
}

func crawl(url string, number int) (QueryResult, error) {
	entries := []*entry{}
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return QueryResult{
			Title:   "",
			Entries: entries,
		}, err
	}
	doc.Find("tr").Each(func(_ int, line *goquery.Selection) {
		cells := [3]interface{}{"", "", ""}
		cells[0] = line.Find("a").Text()
		cells[1], _ = line.Find("a").Attr("href")
		cells[2] = line.Find("td").Eq(2).Text()
		entry := entry{
			Title:     strings.TrimSpace(cells[0].(string)),
			URL:       cells[1].(string),
			LikeCount: cells[2].(string),
		}
		if entry.URL != "" {
			entries = append(entries, &entry)
		}
	})
	if number > len(entries) {
		number = len(entries)
	}
	return QueryResult{
		Title:   doc.Find(".subtitle1").Text(),
		Entries: entries[:number],
	}, nil
}

func maxTitleWidth(entries []*entry) int {
	width := 0
	for _, e := range entries {
		count := runewidth.StringWidth(e.Title)
		if count > width {
			width = count
		}
	}
	return width
}

func maxURLWidth(entries []*entry) int {
	width := 0
	for _, e := range entries {
		count := utf8.RuneCountInString(e.URL)
		if count > width {
			width = count
		}
	}
	return width
}

func showResult(result QueryResult, url string) {
	entries := result.Entries
	if len(entries) == 0 {
		fmt.Println("Rankingｹﾞﾄできなかたyo!")
		fmt.Printf("  url: %s \n\n", url)
		return
	}
	fmt.Printf("%s : %d 件\n",
		result.Title,
		len(entries),
	)
	titleWidth := maxTitleWidth(entries)
	titleFmt := fmt.Sprintf("%%-%ds", titleWidth)

	urlWidth := maxURLWidth(entries)
	urlFmt := fmt.Sprintf("%%-%ds", urlWidth)

	fmt.Fprintf(color.Output, " %s | %s | %s \n",
		color.BlueString(titleFmt, "Title"),
		fmt.Sprintf(urlFmt, "Url"),
		color.CyanString("Like"),
	)
	fmt.Println(strings.Repeat("-", titleWidth+urlWidth+16))
	for _, e := range entries {
		fmt.Fprintf(color.Output, " %s | %s | %s \n",
			color.BlueString(runewidth.FillRight(e.Title, titleWidth)),
			fmt.Sprintf(urlFmt, e.URL),
			color.CyanString(e.LikeCount),
		)
	}
}
