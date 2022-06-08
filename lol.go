package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type LeagueAPI struct {
	region string
	key    string
}

func (l *LeagueAPI) Config(region string, key string) {
	l.region = region
	l.key = key
}

func (l *LeagueAPI) GetSummonerByName(summoner string) (map[string]interface{}, error) {

	resSumm, err := http.Get("https://" + l.region + ".api.riotgames.com/lol/summoner/v4/summoners/by-name/" + summoner + "?api_key=" + l.key)
	if err != nil {
		log.Fatalln(err)
	}

	defer resSumm.Body.Close()
	body, err := io.ReadAll(resSumm.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (l *LeagueAPI) GetMatchesPerPuuid(puuid string, from string, to string) ([]string, error) {

	resMatch, err := http.Get("https://europe.api.riotgames.com/lol/match/v5/matches/by-puuid/" + puuid + "/ids?startTime=" + from + "&endTime=" + to + "&start=0&count=100&api_key=" + l.key)
	if err != nil {
		return nil, err
	}

	defer resMatch.Body.Close()

	body, err := io.ReadAll(resMatch.Body)

	var resultMatch []string
	err = json.Unmarshal(body, &resultMatch)
	if err != nil {
		return nil, err
	}

	return resultMatch, nil
}

func (l *LeagueAPI) GetMatchesPerSummonerName(summoner string, from string, to string) ([]string, error) {

	resSumm, err := l.GetSummonerByName(summoner)
	if err != nil {
		return nil, err
	}

	puuid := resSumm["puuid"].(string)

	resMatches, err := l.GetMatchesPerPuuid(puuid, from, to)
	if err != nil {
		return nil, err
	}

	return resMatches, nil
}

func (l *LeagueAPI) GetMatchDetail(match string) (map[string]interface{}, error) {

	resMatch, err := http.Get("https://europe.api.riotgames.com/lol/match/v5/matches/" + match + "?api_key=" + l.key)
	if err != nil {
		return nil, err
	}

	defer resMatch.Body.Close()

	body, err := io.ReadAll(resMatch.Body)
	if err != nil {
		log.Fatalln("GetMatchDetail")
		return nil, err
	}

	var resultMatch map[string]interface{}

	err = json.Unmarshal(body, &resultMatch)
	if err != nil {
		log.Fatalln("GetMatchDetail")
		return nil, err
	}

	return resultMatch, nil
}

func (l *LeagueAPI) GetMatchTimeline(match string) (map[string]interface{}, error) {

	resMatch, err := http.Get("https://europe.api.riotgames.com/lol/match/v5/matches/" + match + "/timeline?api_key=" + l.key)
	if err != nil {
		return nil, err
	}

	defer resMatch.Body.Close()

	body, err := io.ReadAll(resMatch.Body)
	if err != nil {
		return nil, err
	}

	var resultMatch map[string]interface{}

	err = json.Unmarshal(body, &resultMatch)
	if err != nil {
		return nil, err
	}

	return resultMatch, nil
}
