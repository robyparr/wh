package cmd

import (
	"fmt"
	"io"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/robyparr/wh/model"
	"github.com/robyparr/wh/repository"
	"github.com/robyparr/wh/util"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a new work day",
	Run: func(cmd *cobra.Command, args []string) {
		repo, err := repository.NewRepo("./db.sqlite")
		if err != nil {
			log.Fatalln(err)
		}

		if err := runAddCmd(os.Stdout, repo); err != nil {
			log.Fatalln(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}

func runAddCmd(w io.Writer, repo *repository.Repo) error {
	midnight := util.TodayAtMidnight()
	workDay, err := repo.GetWorkDayByDate(midnight)
	if err != nil {
		return err
	}

	if workDay.Id != 0 {
		fmt.Fprintf(w, "Work day on %s already exists.\n", util.FormatDate(midnight))
		return nil
	}

	workDay, err = repo.CreateWorkDay(model.NewWorkDay())
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "Added work day #%d on %s\n", workDay.Id, util.FormatDate(workDay.Date))
	return nil
}