package snapshot

import (
	"os/exec"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

type Collector struct {
	resticExecutablePath string
}

func NewSnapshotCollector(resticExecutablePath string) *Collector {
	return &Collector{
		resticExecutablePath: resticExecutablePath,
	}
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- snapshotCountTotalDesc
	ch <- snapshotCountDesc
	ch <- lastSnapshotTimeDesc
	ch <- lastSnapshotBackupStartDesc
	ch <- lastSnapshotBackupEndDesc
	ch <- lastSnapshotFilesNewDesc
	ch <- lastSnapshotFilesChangedDesc
	ch <- lastSnapshotFilesUnmodifiedDesc
	ch <- lastSnapshotDirsNewDesc
	ch <- lastSnapshotDirsChangedDesc
	ch <- lastSnapshotDirsUnmodifiedDesc
	ch <- lastSnapshotDataBlobsDesc
	ch <- lastSnapshotTreeBlobsDesc
	ch <- lastSnapshotDataAddedDesc
	ch <- lastSnapshotDataAddedPackedDesc
	ch <- lastSnapshotTotalFilesProcessedDesc
	ch <- lastSnapshotTotalBytesProcessedDesc
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	c.collectWithExecutor(ch, execCommandExecutor)
}

func (c *Collector) collectWithExecutor(ch chan<- prometheus.Metric, commandExecutor CommandExecutor) {
	out, err, exitCode := commandExecutor(c.resticExecutablePath, "snapshots", "--json", "--no-lock", "--group-by", "host,tags")

	ch <- prometheus.MustNewConstMetric(snapshotExitCode, prometheus.GaugeValue, float64(exitCode))
	if err != nil {
		return
	}

	groupData, err := readJson(out)
	if err != nil {
		panic(err)
	}

	totalSnapshotCount := getTotalSnapshotCount(groupData)
	ch <- prometheus.MustNewConstMetric(snapshotCountTotalDesc, prometheus.GaugeValue, float64(totalSnapshotCount))

	for _, group := range groupData {
		hostname := group.GroupKey.Hostname
		tags := strings.Join(group.GroupKey.Tags, ",")
		_, snapshotCount := getSnapshotCountByGroup(group)
		ch <- prometheus.MustNewConstMetric(
			snapshotCountDesc,
			prometheus.GaugeValue,
			float64(snapshotCount),
			hostname,
			tags,
		)

		if snapshotCount > 0 {
			_, metrics, err := getSnapshotMetricsByGroup(group)
			if err != nil {
				panic(err)
			}

			ch <- prometheus.MustNewConstMetric(lastSnapshotTimeDesc, prometheus.GaugeValue, float64(metrics.Time.Unix()), hostname, tags)
			ch <- prometheus.MustNewConstMetric(lastSnapshotBackupStartDesc, prometheus.GaugeValue, float64(metrics.BackupStart.Unix()), hostname, tags)
			ch <- prometheus.MustNewConstMetric(lastSnapshotBackupEndDesc, prometheus.GaugeValue, float64(metrics.BackupEnd.Unix()), hostname, tags)
			ch <- prometheus.MustNewConstMetric(lastSnapshotFilesNewDesc, prometheus.GaugeValue, float64(metrics.FilesNew), hostname, tags)
			ch <- prometheus.MustNewConstMetric(lastSnapshotFilesChangedDesc, prometheus.GaugeValue, float64(metrics.FilesChanged), hostname, tags)
			ch <- prometheus.MustNewConstMetric(lastSnapshotFilesUnmodifiedDesc, prometheus.GaugeValue, float64(metrics.FilesUnmodified), hostname, tags)
			ch <- prometheus.MustNewConstMetric(lastSnapshotDirsNewDesc, prometheus.GaugeValue, float64(metrics.DirsNew), hostname, tags)
			ch <- prometheus.MustNewConstMetric(lastSnapshotDirsChangedDesc, prometheus.GaugeValue, float64(metrics.DirsChanged), hostname, tags)
			ch <- prometheus.MustNewConstMetric(lastSnapshotDirsUnmodifiedDesc, prometheus.GaugeValue, float64(metrics.DirsUnmodified), hostname, tags)
			ch <- prometheus.MustNewConstMetric(lastSnapshotDataBlobsDesc, prometheus.GaugeValue, float64(metrics.DataBlobs), hostname, tags)
			ch <- prometheus.MustNewConstMetric(lastSnapshotTreeBlobsDesc, prometheus.GaugeValue, float64(metrics.TreeBlobs), hostname, tags)
			ch <- prometheus.MustNewConstMetric(lastSnapshotDataAddedDesc, prometheus.GaugeValue, float64(metrics.DataAdded), hostname, tags)
			ch <- prometheus.MustNewConstMetric(lastSnapshotDataAddedPackedDesc, prometheus.GaugeValue, float64(metrics.DataAddedPacked), hostname, tags)
			ch <- prometheus.MustNewConstMetric(lastSnapshotTotalFilesProcessedDesc, prometheus.GaugeValue, float64(metrics.TotalFilesProcessed), hostname, tags)
			ch <- prometheus.MustNewConstMetric(lastSnapshotTotalBytesProcessedDesc, prometheus.GaugeValue, float64(metrics.TotalBytesProcessed), hostname, tags)
		}
	}
}

// CommandExecutor is a function that executes a command and returns the output, error and exit code.
type CommandExecutor func(name string, arg ...string) ([]byte, error, int)

// execCommandExecutor executes a command with exec.Command and returns the output, error and exit code.
var execCommandExecutor CommandExecutor = func(name string, arg ...string) ([]byte, error, int) {
	cmd := exec.Command(name, arg...)
	output, err := cmd.Output()
	exitCode := cmd.ProcessState.ExitCode()
	return output, err, exitCode
}
