package cmd

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/robyparr/wh/repository"
	"github.com/robyparr/wh/util"
	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop [time]",
	Short: "Stop tracking work hours",
	Run: func(cmd *cobra.Command, args []string) {
		repo, err := repository.NewRepo(repository.DefaultDatabasePath)
		if err != nil {
			log.Fatalln(err)
		}

		var timeStr string
		if len(args) != 0 {
			timeStr = args[0]
		}

		note := mustGetStringFlag(cmd, "note")
		if err := runStopCmd(os.Stdout, repo, timeStr, note); err != nil {
			log.Fatalln(err)
		}
	},
}

func init() {
	stopCmd.Flags().StringP("note", "n", "", "work period note")
	rootCmd.AddCommand(stopCmd)
}

func runStopCmd(out io.Writer, repo *repository.Repo, timeStr string, note string) error {
	workDay, err := repo.GetWorkDayByDate(util.TodayAtMidnight())
	if err != nil {
		return err
	}

	period, err := repo.GetOpenWorkPeriod(workDay)
	if err != nil {
		return err
	}

	if period.Id == 0 {
		fmt.Fprintln(out, "Unable to find an ongoing work period.")
		return nil
	}

	endAt, err := util.ParseTimeString(timeStr)
	if err != nil {
		return err
	}

	period.EndAt = sql.NullTime{Valid: true, Time: endAt}
	if note != "" {
		period.SetNote(note)
	}

	if _, err := repo.UpdateWorkPeriod(period); err != nil {
		return err
	}

	return nil
}
