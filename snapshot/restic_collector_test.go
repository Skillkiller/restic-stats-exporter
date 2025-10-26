package snapshot

import (
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func TestCollector_Describe(t *testing.T) {
	c := &Collector{
		resticExecutablePath: "",
	}

	expectedDesc := map[string]bool{
		snapshotCountTotalDesc.String():              true,
		snapshotCountDesc.String():                   true,
		lastSnapshotTimeDesc.String():                true,
		lastSnapshotBackupStartDesc.String():         true,
		lastSnapshotBackupEndDesc.String():           true,
		lastSnapshotFilesNewDesc.String():            true,
		lastSnapshotFilesChangedDesc.String():        true,
		lastSnapshotFilesUnmodifiedDesc.String():     true,
		lastSnapshotDirsNewDesc.String():             true,
		lastSnapshotDirsChangedDesc.String():         true,
		lastSnapshotDirsUnmodifiedDesc.String():      true,
		lastSnapshotDataBlobsDesc.String():           true,
		lastSnapshotTreeBlobsDesc.String():           true,
		lastSnapshotDataAddedDesc.String():           true,
		lastSnapshotDataAddedPackedDesc.String():     true,
		lastSnapshotTotalFilesProcessedDesc.String(): true,
		lastSnapshotTotalBytesProcessedDesc.String(): true,
		snapshotExitCode.String():                    true,
	}

	expectedCount := len(expectedDesc)

	ch := make(chan *prometheus.Desc, expectedCount)
	done := make(chan struct{})

	go func() {
		c.Describe(ch)
		close(done)
	}()

	select {
	case <-done:
		// Describe has finished
	case <-time.After(200 * time.Millisecond):
		t.Fatalf("Describe blocked or did not return within timeout — possible extra sends")
	}

	close(ch)

	got := map[string]bool{}
	for d := range ch {
		if d == nil {
			t.Fatalf("received nil descriptor")
		}
		got[d.String()] = true
	}

	if len(got) != expectedCount {
		t.Fatalf("wrong number of descriptors: got %d, want %d", len(got), expectedCount)
	}

	for want := range expectedDesc {
		if !got[want] {
			t.Fatalf("descriptor not sent: %s", want)
		}
	}
}
