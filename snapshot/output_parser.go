package snapshot

import (
	"encoding/json"
	"fmt"
	"time"
)

type GroupData struct {
	GroupKey  GroupKey   `json:"group_key"`
	Snapshots []Snapshot `json:"snapshots"`
}

type GroupKey struct {
	Hostname string   `json:"hostname"`
	Tags     []string `json:"tags"`
}

type Snapshot struct {
	Time     time.Time `json:"time"`
	Hostname string    `json:"hostname"`
	Tags     []string  `json:"tags"`
	Summary  Summary   `json:"summary"`
}

type Summary struct {
	BackupStart         time.Time `json:"backup_start"`
	BackupEnd           time.Time `json:"backup_end"`
	FilesNew            int       `json:"files_new"`
	FilesChanged        int       `json:"files_changed"`
	FilesUnmodified     int       `json:"files_unmodified"`
	DirsNew             int       `json:"dirs_new"`
	DirsChanged         int       `json:"dirs_changed"`
	DirsUnmodified      int       `json:"dirs_unmodified"`
	DataBlobs           int       `json:"data_blobs"`
	TreeBlobs           int       `json:"tree_blobs"`
	DataAdded           int       `json:"data_added"`
	DataAddedPacked     int       `json:"data_added_packed"`
	TotalFilesProcessed int       `json:"total_files_processed"`
	TotalBytesProcessed int       `json:"total_bytes_processed"`
}

type SnapshotMetrics struct {
	Time                time.Time
	BackupStart         time.Time
	BackupEnd           time.Time
	FilesNew            int
	FilesChanged        int
	FilesUnmodified     int
	DirsNew             int
	DirsChanged         int
	DirsUnmodified      int
	DataBlobs           int
	TreeBlobs           int
	DataAdded           int
	DataAddedPacked     int
	TotalFilesProcessed int
	TotalBytesProcessed int
}

func readJson(data []byte) ([]GroupData, error) {
	var o []GroupData
	err := json.Unmarshal(data, &o)
	if err != nil {
		return o, err
	}
	return o, nil
}

func getTotalSnapshotCount(data []GroupData) int {
	sum := 0

	for _, group := range data {
		sum += len(group.Snapshots)
	}

	return sum
}

func getSnapshotCountByGroup(group GroupData) (GroupKey, int) {
	return group.GroupKey, len(group.Snapshots)
}

func getLastSnapshotByGroup(group GroupData) (GroupKey, Snapshot, error) {
	if len(group.Snapshots) == 0 {
		return GroupKey{}, Snapshot{}, fmt.Errorf("no snapshots found")
	}

	lastSnapshot := group.Snapshots[len(group.Snapshots)-1]

	for _, snapshot := range group.Snapshots {
		if snapshot.Time.After(lastSnapshot.Time) {
			lastSnapshot = snapshot
		}
	}

	return group.GroupKey, lastSnapshot, nil
}

func getSnapshotMetricsByGroup(group GroupData) (GroupKey, SnapshotMetrics, error) {
	_, snapshot, err := getLastSnapshotByGroup(group)
	if err != nil {
		return GroupKey{}, SnapshotMetrics{}, err
	}

	return group.GroupKey, SnapshotMetrics{
		Time:                snapshot.Time,
		BackupStart:         snapshot.Summary.BackupStart,
		BackupEnd:           snapshot.Summary.BackupEnd,
		FilesNew:            snapshot.Summary.FilesNew,
		FilesChanged:        snapshot.Summary.FilesChanged,
		FilesUnmodified:     snapshot.Summary.FilesUnmodified,
		DirsNew:             snapshot.Summary.DirsNew,
		DirsChanged:         snapshot.Summary.DirsChanged,
		DirsUnmodified:      snapshot.Summary.DirsUnmodified,
		DataBlobs:           snapshot.Summary.DataBlobs,
		TreeBlobs:           snapshot.Summary.TreeBlobs,
		DataAdded:           snapshot.Summary.DataAdded,
		DataAddedPacked:     snapshot.Summary.DataAddedPacked,
		TotalFilesProcessed: snapshot.Summary.TotalFilesProcessed,
		TotalBytesProcessed: snapshot.Summary.TotalBytesProcessed,
	}, nil
}
