package snapshot

import (
	"encoding/json"
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
