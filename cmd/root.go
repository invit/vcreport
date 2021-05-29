package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/invit/vcreport/internal/lib/version"
	"github.com/olekukonko/tablewriter"
	"github.com/prometheus/common/expfmt"
	"github.com/spf13/cobra"
)

type Pod struct {
	Namespace string `json:"namespace"`
	Pod       string `json:"pod"`
	Container string `json:"container"`
}

type Image struct {
	Image          string `json:"image"`
	CurrentVersion string `json:"current_version"`
	LatestVersion  string `json:"latest_version"`
	IsLatest       bool   `json:"is_latest"`
	Pods           []Pod  `json:"pods"`
}

func (i *Image) Version() string {
	if i.IsLatest {
		return fmt.Sprintf("%s (Up to date)", i.CurrentVersion)
	}

	return fmt.Sprintf("%s > %s", i.CurrentVersion, i.LatestVersion)
}

func (i *Image) AddPod(p Pod) {
	i.Pods = append(i.Pods, p)
}

func init() {
	rootCmd.Flags().BoolP("all", "a", false, "Show all images, not just outdated ones")
	rootCmd.Flags().BoolP("brief", "b", false, "Just show images, but no pods")
}

var rootCmd = &cobra.Command{
	Use:           "vcreport metrics-url",
	Short:         "Displays a human-readable report from version-checker metrics",
	SilenceErrors: true,
	Version:       fmt.Sprintf("%s-%s", version.Version, version.Commit),
	Args:          cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		displayAll, _ := cmd.Flags().GetBool("all")
		displayBrief, _ := cmd.Flags().GetBool("brief")

		c := http.Client{}
		resp, err := c.Get(args[0])
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			p := expfmt.TextParser{}
			mf, err := p.TextToMetricFamilies(resp.Body)

			if err != nil {
				return err
			}

			images := map[string]*Image{}

			for _, v := range mf {
				for _, m := range v.Metric {
					labels := map[string]string{}
					latest := *m.Gauge.Value > 0

					for _, l := range m.Label {
						labels[l.GetName()] = l.GetValue()
					}

					if !displayAll && latest {
						continue
					}

					key := fmt.Sprintf("%s:%s", labels["image"], labels["current_version"])

					if _, ok := images[key]; !ok {
						images[key] = &Image{
							Image:          labels["image"],
							CurrentVersion: labels["current_version"],
							LatestVersion:  labels["latest_version"],
							IsLatest:       latest,
							Pods: []Pod{{
								Namespace: labels["namespace"],
								Pod:       labels["pod"],
								Container: labels["container"],
							}},
						}
					} else {
						images[key].AddPod(
							Pod{
								Namespace: labels["namespace"],
								Pod:       labels["pod"],
								Container: labels["container"],
							},
						)
					}
				}
			}

			var keys = []string{}

			for k, _ := range images {
				keys = append(keys, k)
			}

			sort.Strings(keys)

			table := tablewriter.NewWriter(os.Stdout)
			table.SetAlignment(tablewriter.ALIGN_LEFT)
			table.SetRowLine(true)

			if displayBrief {
				table.SetHeader([]string{"Image", "Version"})

				for _, k := range keys {
					table.Append([]string{
						images[k].Image,
						images[k].Version(),
					})
				}

			} else {
				table.SetHeader([]string{"Image", "Version", "Pods"})

				for _, k := range keys {
					pods := []string{}

					for _, p := range images[k].Pods {
						pods = append(
							pods,
							fmt.Sprintf("%s/%s/%s", p.Namespace, p.Pod, p.Container),
						)
					}

					table.Append([]string{
						images[k].Image,
						images[k].Version(),
						strings.Join(pods, "\n"),
					})
				}
			}

			table.Render()
		} else {
			b, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("Error from %s (Code %d): %s", args[0], resp.StatusCode, b)
		}

		return nil
	},
}

// Execute runs root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
