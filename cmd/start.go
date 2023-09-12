package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/robyparr/wh/model"
	"github.com/robyparr/wh/repository"
	"github.com/robyparr/wh/util"
	"github.com/spf13/cobra"
)

type startCmdArgs struct {
	timeStr   string
	lengthStr string
	note      string
	dayNote   string
}

var startCmd = &cobra.Command{
	Use:   "start [time]",
	Short: "Start tracking work hours",
	Run: func(cmd *cobra.Command, args []string) {
		repo, err := repository.NewRepo(repository.DefaultDatabasePath)
		if err != nil {
			log.Fatalln(err)
		}

		var cmdArgs startCmdArgs
		if len(args) > 0 {
			cmdArgs.timeStr = args[0]
		}

		cmdArgs.lengthStr = mustGetStringFlag(cmd, "length")
		cmdArgs.note = mustGetStringFlag(cmd, "note")
		cmdArgs.dayNote = mustGetStringFlag(cmd, "day-note")

		if err := runStartCmd(os.Stdout, repo, cmdArgs); err != nil {
			log.Fatalln(err)
		}
	},
}

func init() {
	startCmd.Flags().StringP("length", "l", "", "work day length (e.g. 4h30m)")
	startCmd.Flags().StringP("note", "n", "", "work period note")
	startCmd.Flags().StringP("day-note", "d", "", "work day note")

	rootCmd.AddCommand(startCmd)
}

func runStartCmd(out io.Writer, repo *repository.Repo, args startCmdArgs) error {
	midnight := util.TodayAtMidnight()
	startAt, err := util.ParseTimeString(args.timeStr)
	if err != nil {
		return fmt.Errorf("error parsing time string:, %v", err)
	}

	workDay, err := repo.GetWorkDayByDate(midnight)
	if err != nil {
		return fmt.Errorf("error loading work day: %v", err)
	}

	workPeriods, err := repo.GetWorkPeriods(workDay)
	if err != nil {
		return fmt.Errorf("error loading work periods: %v", err)
	}
	for _, wp := range workPeriods {
		if wp.EndAt.Time.IsZero() {
			fmt.Fprintln(out, "This work day already has an open work period.")
			return nil
		}
	}

	outFormatString := "Started tracking time on work day #%d (%s).\n"
	if workDay.Id == 0 {
		workDay = model.NewWorkDay(midnight)
		if args.lengthStr != "" {
			duration, err := time.ParseDuration(args.lengthStr)
			if err != nil {
				return fmt.Errorf("error parsing length string: %v", err)
			}

			workDay.LengthMins = int(duration.Minutes())
		}

		if args.dayNote != "" {
			workDay.SetNote(args.dayNote)
		}

		workDay, err = repo.CreateWorkDay(workDay)
		if err != nil {
			return fmt.Errorf("error creating new work day: %v", err)
		}

		outFormatString = "Started tracking time on NEW work day #%d (%s).\n"
	}

	period := model.WorkPeriod{WorkDayId: workDay.Id, StartAt: startAt}
	if args.note != "" {
		period.SetNote(args.note)
	}

	_, err = repo.CreateWorkPeriod(period)
	if err != nil {
		return fmt.Errorf("error creating work period: %v", err)
	}

	fmt.Fprintf(out, outFormatString, workDay.Id, util.FormatDate(workDay.Date))
	return nil
}
