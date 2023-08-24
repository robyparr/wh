package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/robyparr/wh/model"
	"github.com/robyparr/wh/repository"
	"github.com/robyparr/wh/util"
	"github.com/spf13/cobra"
)

var exactTimeRegex = regexp.MustCompile(`^\d{2}:\d{2}$`)
var relativeTimeRegex = regexp.MustCompile(`^(-?\d+h(\d+m)?)|(-?\d+m)$`)

var startCmd = &cobra.Command{
	Use:   "start [time]",
	Short: "Start tracking work hours",
	Run: func(cmd *cobra.Command, args []string) {
		repo, err := repository.NewRepo(repository.DefaultDatabasePath)
		if err != nil {
			log.Fatalln(err)
		}

		var timeStr string
		if len(args) > 0 {
			timeStr = args[0]
		}

		if err := runStartCmd(os.Stdout, repo, timeStr); err != nil {
			log.Fatalln(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}

func runStartCmd(out io.Writer, repo *repository.Repo, timeStr string) error {
	workDay, err := repo.GetWorkDayByDate(util.TodayAtMidnight())
	if err != nil {
		return fmt.Errorf("error loading work day: %v", err)
	}

	startAt := time.Now()
	switch {
	case exactTimeRegex.MatchString(timeStr):
		duration, err := parseExactTimeString(timeStr)
		if err != nil {
			return fmt.Errorf("error parsing time string:, %v", err)
		}

		startAt = util.TodayAtMidnight().Add(duration)

	case relativeTimeRegex.MatchString(timeStr):
		duration, err := time.ParseDuration(timeStr)
		if err != nil {
			return fmt.Errorf("error parsing time string: %v", err)
		}

		startAt = startAt.Add(duration)
	}

	period := model.WorkPeriod{WorkDayId: workDay.Id, StartAt: startAt}
	_, err = repo.CreateWorkPeriod(period)
	if err != nil {
		return err
	}

	fmt.Fprintf(out, "Started tracking time on work day #%d (%s).\n", workDay.Id, util.FormatDate(workDay.Date))
	return nil
}

func parseExactTimeString(str string) (time.Duration, error) {
	timeStrParts := strings.Split(str, ":")
	hour, err := strconv.Atoi(timeStrParts[0])
	if err != nil {
		return 0, nil
	}
	min, err := strconv.Atoi(timeStrParts[1])
	if err != nil {
		return 0, err
	}

	totalMinutes := (hour * 60) + min
	return time.Duration(totalMinutes) * time.Minute, nil
}
