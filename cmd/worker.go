package cmd

import (
	"awesomeProject/prometheus"
	"awesomeProject/worker"
	"github.com/spf13/cobra"
	"log"
	"sync"
	"time"
)

func init() {
	cmd := &cobra.Command{
		Use:   "worker",
		Short: "Worker",
		Long:  "Worker",
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("starting..")

			cfg := getConfig(cmd)

			wg := sync.WaitGroup{}

			ch := make(chan bool)

			go func() {
				time.Sleep(time.Minute * time.Duration(cfg.WorkerConfig.Time))
				close(ch)
			}()

			newStatist := prometheus.NewPrometheus(&cfg.StatistConfig)

			newWorker := worker.NewWorker(&cfg.WorkerConfig, &wg, &ch, newStatist)

			for i := 0; i < cfg.WorkerConfig.Threads; i++ {
				newWorker.Start(i)
			}

			wg.Wait()
		},
	}
	rootCmd.AddCommand(cmd)
	cmd.PersistentFlags().String("config", "config.yaml", "config.yaml")
}
