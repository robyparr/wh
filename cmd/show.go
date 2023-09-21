package cmd

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/robyparr/wh/repository"
	"github.com/robyparr/wh/template"
	"github.com/robyparr/wh/util"
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show [date]",
	Short: "Shows details about a work day",
	Run: func(cmd *cobra.Command, args []string) {
		repo, err := repository.NewRepo(repository.DefaultDatabasePath)
		if err != nil {
			log.Fatalln(err)
		}

		var dateStr string
		if len(args) > 0 {
			dateStr = args[0]
		}
		if err := runShowCmd(os.Stdout, repo, dateStr); err != nil {
			log.Fatalln(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}

func runShowCmd(out io.Writer, repo *repository.Repo, dateStr string) error {
	date, err := util.ParseDateString(dateStr)
	if err != nil {
		return fmt.Errorf("error parsing date: %v", err)
	}

	workDay, err := repo.GetWorkDayByDate(date)
	if err != nil {
		return fmt.Errorf("error loading work day: %v", err)
	}

	if workDay.Id == 0 {
		fmt.Fprintf(out, "No work day for %s yet.\n", dateStr)
		return nil
	}

	workPeriods, err := repo.GetWorkPeriods(workDay)
	if err != nil {
		return fmt.Errorf("error loading work periods: %v", err)
	}

	workDay.SetWorkPeriods(workPeriods)
	vm := showViewModel{
		Title:           util.Underline(workDay.Date.Format("January 02, 2006 (Mon)")),
		DayLength:       util.FormatDuration(workDay.Length()),
		TimeWorked:      util.FormatDuration(workDay.TimeWorked()),
		TimeRemaining:   util.FormatDuration(workDay.TimeRemaining()),
		EstimatedFinish: util.FormatDateTime(workDay.EstimatedFinish()),
		Note:            workDay.Note.String,
	}
	if err := template.Render(out, "work_day_show.txt", vm); err != nil {
		return err
	}

	return nil
}

type showViewModel struct {
	Title           string
	DayLength       string
	TimeWorked      string
	TimeRemaining   string
	EstimatedFinish string
	Note            string
}
