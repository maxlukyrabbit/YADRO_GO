package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	. "my-app/src/cmd/internal/model"
)

func addTime(time1 string, time2 string) string {
	layout := "15:04"
	time1Parsed, err := time.Parse(layout, time1)
	if err != nil {
		return "Error in time 1 format"
	}
	time2Parsed, err := time.Parse(layout, time2)
	if err != nil {
		return "Error in time 2 format"
	}
	sum := time1Parsed.Add(time2Parsed.Sub(time1Parsed))
	return sum.Format("15:04")
}

func roundUpTime(inputTime string) int {
	parts := strings.Split(inputTime, ":")
	hours, _ := strconv.Atoi(parts[0])
	minutes, _ := strconv.Atoi(parts[1])

	if minutes > 0 {
		hours++
	}

	return hours
}

func checkClientName(name string) bool {
	match, _ := regexp.MatchString("^[a-zA-Z0-9_-]+$", name)
	return match
}

func calculateTimeElapsed(start string, end string) string {
	layout := "15:04"
	startTime, err := time.Parse(layout, start)
	if err != nil {
		return "Ошибка в формате времени начала"
	}
	endTime, err := time.Parse(layout, end)
	if err != nil {
		return "Ошибка в формате времени конца"
	}
	elapsed := endTime.Sub(startTime)
	elapsedHours := int(elapsed.Hours())
	elapsedMinutes := int(elapsed.Minutes()) % 60
	return fmt.Sprintf("%02d:%02d", elapsedHours, elapsedMinutes)
}

func contains(element string, array []string) bool {
	for _, v := range array {
		if v == element {
			return true
		}
	}
	return false
}

func removeElement(element string, array []string) []string {
	index := indexSearch(element, array)
	if index != -1 {
		return append(array[:index], array[index+1:]...)
	}
	return array
}

func indexSearch(element string, array []string) int {
	for i, v := range array {
		if v == element {
			return i
		}
	}
	return -1
}

func readLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

func main() {
	if len(os.Args) == 1 {
		fmt.Println("at least one arg must be passed")
		os.Exit(1)
	}

	var txt string
	txt = os.Args[1]

	lines, err := readLines(txt)
	if err != nil {
		log.Fatalf("readLines: %s", err)
	}

	client_PC := map[string]string{}

	count, err_count := strconv.ParseInt(lines[0], 10, 32)
	time_start, err_time_start := time.Parse("15:04", strings.Split(lines[1], " ")[0])
	time_end, err_time_end := time.Parse("15:04", strings.Split(lines[1], " ")[1])
	price, err_price := strconv.ParseInt(lines[2], 10, 32)

	if err_count != nil {
		fmt.Println(lines[0])
		return
	} else if err_time_start != nil || err_time_end != nil {
		fmt.Println(lines[1])
		return
	} else if err_price != nil {
		fmt.Println(lines[2])
		return
	}

	var events []Event_in

	for _, e := range lines[3:] {
		line := strings.Fields(e)
		event := Event_in{
			Time:   line[0],
			ID:     line[1],
			Client: line[2],
		}
		if len(line) == 4 {
			event.Table = line[3]
		}
		events = append(events, event)
	}

	var profit []Profit
	for i := 1; int64(i) <= count; i++ {
		profitt := Profit{
			PC:         strconv.Itoa(i),
			Full_time:  "00:00",
			Time_start: "00:00",
			Time_end:   "00:00",
		}
		profit = append(profit, profitt)
	}

	var events_out []Event_out

	client_in_club := make([]string, 0)
	busy_table := make([]string, 0)
	queue := make([]string, 0)

	for _, e := range events {
		err_table := error(nil)
		eventTime, err_time := time.Parse("15:04", e.Time)
		if e.Table != "" {
			_, err_table = strconv.ParseInt(e.Table, 10, 32)
		}

		eventClient := checkClientName(e.Client)
		_, err_ID := strconv.ParseInt(e.ID, 10, 32)
		if err_time != nil || err_table != nil || err_ID != nil || !eventClient {
			fmt.Printf("%s %s %s %s\n", e.Time, e.ID, e.Client, e.Table)
			return
		}

		if e.ID == "1" {

			if eventTime.Before(time_start) || eventTime.After(time_end) {
				event_out := Event_out{
					Time:   e.Time,
					ID:     e.ID,
					Client: e.Client,
					Table:  e.Table,
				}

				event_out_ex := Event_out{
					Time:   e.Time,
					ID:     "13",
					Client: "NotOpenYet",
				}
				events_out = append(events_out, event_out)
				events_out = append(events_out, event_out_ex)
			} else if contains(e.Client, client_in_club) {
				event_out := Event_out{
					Time:   e.Time,
					ID:     e.ID,
					Client: e.Client,
					Table:  e.Table,
				}

				event_out_ex := Event_out{
					Time:   e.Time,
					ID:     "13",
					Client: "YouShallNotPass",
				}
				events_out = append(events_out, event_out)
				events_out = append(events_out, event_out_ex)

			} else {
				event_out := Event_out{
					Time:   e.Time,
					ID:     e.ID,
					Client: e.Client,
					Table:  e.Table,
				}
				client_in_club = append(client_in_club, e.Client)
				events_out = append(events_out, event_out)
			}
		} else if e.ID == "2" {
			if !contains(e.Client, client_in_club) {
				event_out := Event_out{
					Time:   e.Time,
					ID:     e.ID,
					Client: e.Client,
					Table:  e.Table,
				}

				event_out_ex := Event_out{
					Time:   e.Time,
					ID:     "13",
					Client: "ClientUnknown",
				}
				events_out = append(events_out, event_out)
				events_out = append(events_out, event_out_ex)
			} else if contains(e.Table, busy_table) {
				event_out := Event_out{
					Time:   e.Time,
					ID:     e.ID,
					Client: e.Client,
					Table:  e.Table,
				}

				event_out_ex := Event_out{
					Time:   e.Time,
					ID:     "13",
					Client: "PlaceIsBusy",
				}
				events_out = append(events_out, event_out)
				events_out = append(events_out, event_out_ex)
			} else {
				event_out := Event_out{
					Time:   e.Time,
					ID:     e.ID,
					Client: e.Client,
					Table:  e.Table,
				}
				events_out = append(events_out, event_out)
				busy_table = append(busy_table, e.Table)
				PC, flag := client_PC[e.Client]
				if flag {
					delete(client_PC, e.Client)
					client_PC[e.Client] = e.Table
					for j, h := range profit {

						if h.PC == PC {
							time := calculateTimeElapsed(h.Time_start, e.Time)
							profit[j].Full_time = addTime(h.Full_time, time)
							profit[j].Time_end = ""
							profit[j].Time_start = ""
						}
						if string(h.PC) == e.Table {
							profit[j].Time_start = e.Time
						}
					}
				} else {
					client_PC[e.Client] = e.Table
					for j, h := range profit {
						if string(h.PC) == e.Table {
							profit[j].Time_start = e.Time
						}
					}
				}

			}
		} else if e.ID == "3" {
			if int64(len((busy_table))) < count {
				event_out := Event_out{
					Time:   e.Time,
					ID:     e.ID,
					Client: e.Client,
					Table:  e.Table,
				}

				event_out_ex := Event_out{
					Time:   e.Time,
					ID:     "13",
					Client: "ICanWaitNoLonger!",
				}
				events_out = append(events_out, event_out)
				events_out = append(events_out, event_out_ex)
			} else if int64(len(queue)) > count {
				event_out := Event_out{
					Time:   e.Time,
					ID:     "11",
					Client: e.Client,
					Table:  e.Table,
				}
				events_out = append(events_out, event_out)
			} else {
				event_out := Event_out{
					Time:   e.Time,
					ID:     e.ID,
					Client: e.Client,
					Table:  e.Table,
				}
				events_out = append(events_out, event_out)
				queue = append(queue, e.Client)
			}
		} else if e.ID == "4" {
			if !contains(e.Client, client_in_club) {
				event_out := Event_out{
					Time:   e.Time,
					ID:     e.ID,
					Client: e.Client,
					Table:  e.Table,
				}
				event_out_ex := Event_out{
					Time:   e.Time,
					ID:     "13",
					Client: "ClientUnknown",
				}
				events_out = append(events_out, event_out)
				events_out = append(events_out, event_out_ex)
			} else {
				PC := client_PC[e.Client]
				delete(client_PC, e.Client)

				event_out := Event_out{
					Time:   e.Time,
					ID:     e.ID,
					Client: e.Client,
					Table:  e.Table,
				}
				events_out = append(events_out, event_out)

				if len(queue) != 0 {

					event_out_sit := Event_out{
						Time:   e.Time,
						ID:     "12",
						Client: queue[0],
						Table:  PC,
					}
					events_out = append(events_out, event_out_sit)
					client_PC[queue[0]] = PC
					queue = removeElement(queue[0], queue)
				} else {
					for j, h := range profit {
						if h.PC == PC {
							time := calculateTimeElapsed(h.Time_start, e.Time)
							profit[j].Full_time = addTime(h.Full_time, time)
							profit[j].Time_end = ""
							profit[j].Time_start = ""
						}
					}
				}

				client_in_club = removeElement(e.Client, client_in_club)
			}
		}
	}

	if len(client_in_club) > 0 {
		sort.Strings(client_in_club)
		for _, e := range client_in_club {
			for j, h := range profit {
				if h.PC == client_PC[e] {
					time := calculateTimeElapsed(h.Time_start, strings.Split(lines[1], " ")[1])
					profit[j].Full_time = addTime(h.Full_time, time)
					profit[j].Time_end = ""
					profit[j].Time_start = ""
				}
			}

			event_out := Event_out{
				Time:   strings.Split(lines[1], " ")[1],
				ID:     "11",
				Client: e,
			}
			events_out = append(events_out, event_out)

		}
	}

	fmt.Println(strings.Split(lines[1], " ")[0])
	for _, e := range events_out {
		fmt.Println(e)
	}
	fmt.Println(strings.Split(lines[1], " ")[1])
	for _, e := range profit {
		fmt.Printf("%s %d %s\n", e.PC, roundUpTime(e.Full_time)*int(price), e.Full_time)
	}

}
