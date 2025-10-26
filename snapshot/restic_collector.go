package snapshot

import (
	"os/exec"

	"github.com/prometheus/client_golang/prometheus"
)

type Collector struct{}

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
	cmd := exec.Command("restic", "snapshots", "--json", "--no-lock")

	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	groupData, err := readJson(out)
	if err != nil {
		panic(err)
	}

	totalSnapshotCount := getTotalSnapshotCount(groupData)
	ch <- prometheus.MustNewConstMetric(snapshotCountTotalDesc, prometheus.GaugeValue, float64(totalSnapshotCount))
}
