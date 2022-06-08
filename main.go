package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {

	key := flag.String("key", "", "Riot games api key")
	summ := flag.String("summ", "", "Summoner name")
	region := flag.String("region", "", "Region")
	flag.Parse()

	f, _ := os.OpenFile(*summ+".json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()

	lolApi := LeagueAPI{}
	lolApi.Config(*region, *key)

	var matches []string

	resSumm, err := lolApi.GetSummonerByName(*summ)
	if err != nil {
		log.Fatalln(err)
	}

	puuid := resSumm["puuid"].(string)

	start, _ := time.Parse("2006-01-02", "2021-07-31")
	to := start
	from := start
	for {
		to = from.AddDate(0, 0, 15)
		log.Println(from, ":", to)

		b := false
		if to.Unix() > time.Now().Unix() {
			to = time.Now()

			b = true
		}

		resMatches, err := lolApi.GetMatchesPerPuuid(puuid, strconv.FormatInt(from.Unix(), 10), strconv.FormatInt(to.Unix(), 10))
		if err != nil {
			log.Fatalln(err)
		}

		time.Sleep(5 * time.Second)

		matches = append(matches, resMatches...)

		if b {
			break
		}

		from = to
	}

	log.Println(matches)

	for _, m := range matches {

		log.Println("Analyze " + m)
		match, err := lolApi.GetMatchDetail(m)
		if err != nil {
			log.Fatalln(err)
		}

		timeline, err := lolApi.GetMatchTimeline(m)
		if err != nil {
			log.Fatalln(err)
		}

		bm, _ := json.Marshal(match)
		bt, _ := json.Marshal(timeline)

		f.WriteString("{ \"match\" : ")
		f.Write(bm)
		f.WriteString(", \"timeline\": ")
		f.Write(bt)
		f.WriteString("}\n")

		time.Sleep(2 * time.Second)
	}

}
