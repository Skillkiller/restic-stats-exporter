package snapshot

import (
	"os/exec"

	"github.com/prometheus/client_golang/prometheus"
)

type Collector struct{}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- snapshotCountTotal
	ch <- snapshotCount
	ch <- lastSnapshotTime
	ch <- lastSnapshotBackupStart
	ch <- lastSnapshotBackupEnd
	ch <- lastSnapshotFilesNew
	ch <- lastSnapshotFilesChanged
	ch <- lastSnapshotFilesUnmodified
	ch <- lastSnapshotDirsNew
	ch <- lastSnapshotDirsChanged
	ch <- lastSnapshotDirsUnmodified
	ch <- lastSnapshotDataBlobs
	ch <- lastSnapshotTreeBlobs
	ch <- lastSnapshotDataAdded
	ch <- lastSnapshotDataAddedPacked
	ch <- lastSnapshotTotalFilesProcessed
	ch <- lastSnapshotTotalBytesProcessed
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	cmd := exec.Command("restic", "snapshots", "--json", "--no-lock")

	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	print(string(out))
}
