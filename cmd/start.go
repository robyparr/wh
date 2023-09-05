package cmd

import (
	"database/sql"
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

		lengthStr := mustGetStringFlag(cmd, "length")
		note := mustGetStringFlag(cmd, "note")
		if err := runStartCmd(os.Stdout, repo, timeStr, lengthStr, note); err != nil {
			log.Fatalln(err)
		}
	},
}

func init() {
	startCmd.Flags().StringP("length", "l", "", "work day length (e.g. 4h30m)")
	startCmd.Flags().StringP("note", "n", "", "work period note")

	rootCmd.AddCommand(startCmd)
}

func runStartCmd(out io.Writer, repo *repository.Repo, timeStr string, lengthStr string, note string) error {
	midnight := util.TodayAtMidnight()
	startAt := time.Now()

	switch {
	case exactTimeRegex.MatchString(timeStr):
		duration, err := parseExactTimeString(timeStr)
		if err != nil {
			return fmt.Errorf("error parsing time string:, %v", err)
		}

		startAt = midnight.Add(duration)

	case relativeTimeRegex.MatchString(timeStr):
		duration, err := time.ParseDuration(timeStr)
		if err != nil {
			return fmt.Errorf("error parsing time string: %v", err)
		}

		startAt = startAt.Add(duration)
	}

	workDay, err := repo.GetWorkDayByDate(midnight)
	if err != nil {
		return fmt.Errorf("error loading work day: %v", err)
	}

	outFormatString := "Started tracking time on work day #%d (%s).\n"
	if workDay.Id == 0 {
		workDay = model.NewWorkDay(midnight)
		if lengthStr != "" {
			duration, err := time.ParseDuration(lengthStr)
			if err != nil {
				return fmt.Errorf("error parsing length string: %v", err)
			}

			workDay.LengthMins = int(duration.Minutes())
		}

		workDay, err = repo.CreateWorkDay(workDay)
		if err != nil {
			return fmt.Errorf("error creating new work day: %v", err)
		}

		outFormatString = "Started tracking time on NEW work day #%d (%s).\n"
	}

	period := model.WorkPeriod{WorkDayId: workDay.Id, StartAt: startAt}
	if note != "" {
		period.Note = sql.NullString{Valid: true, String: note}
	}

	_, err = repo.CreateWorkPeriod(period)
	if err != nil {
		return fmt.Errorf("error creating work period: %v", err)
	}

	fmt.Fprintf(out, outFormatString, workDay.Id, util.FormatDate(workDay.Date))
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
