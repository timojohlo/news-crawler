package cmd

import (
	"fmt"
	"path"
	"time"

	"github.com/spf13/cobra"
	"github.com/thesoenke/news-crawler/feedreader"
)

var feedsOutDir string
var cmdFeeds = &cobra.Command{
	Use:   "feeds",
	Short: "Download items from a list of feeds",
	RunE: func(cmd *cobra.Command, args []string) error {
		feedInputFile := args[0]
		reader, err := feedreader.New(feedInputFile)
		if err != nil {
			return err
		}

		location, err := time.LoadLocation(timeZone)
		if err != nil {
			return err
		}

		log := fmt.Sprintf("Feedreader started %s\n", time.Now().In(location))
		fmt.Println(log)
		err = writeLog(logsDir, log)

		dayTime := time.Now().In(location)
		reader.Day = &dayTime
		reader.Verbose = verbose
		start := time.Now()
		// TODO Replace this again with reader.Fetch() when fixed possible deadlock
		reader.FetchSerial()
		items := 0
		for _, feed := range reader.Feeds {
			items += len(feed.Items)
		}

		dir := path.Join(feedsOutDir, lang)
		err = reader.LogFailures(dir, &dayTime)
		if err != nil {
			return err
		}

		file, err := reader.Store(dir, &dayTime)
		if err != nil {
			return err
		}

		log = fmt.Sprintf("Feedreader finished %s\nSuccessful: %d\nFailures: %d\nArticles: %d\nTime: %s\nFile: %s\n", time.Now().In(location), len(reader.Feeds), len(reader.FailedFeeds), items, time.Since(start), file)
		fmt.Println(log)
		err = writeLog(logsDir, log)
		return err
	},
}

func init() {
	cmdFeeds.Args = cobra.ExactArgs(1)
	cmdFeeds.PersistentFlags().StringVarP(&lang, "lang", "l", "", "Language of the content")
	cmdFeeds.PersistentFlags().StringVarP(&timeZone, "timezone", "t", "Europe/Berlin", "Timezone for storing the feeds")
	cmdFeeds.PersistentFlags().StringVarP(&feedsOutDir, "dir", "d", "out/feeds", "Directory to store feed items")
	cmdFeeds.PersistentFlags().StringVar(&logsDir, "logs", "out/events.log", "File to store logs")
	cmdFeeds.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Output more detailed logging")
	RootCmd.AddCommand(cmdFeeds)
}
