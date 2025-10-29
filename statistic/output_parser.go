package statistic

import (
	"encoding/json"
)

type RawDataMetrics struct {
	TotalSize              int     `json:"total_size"`
	TotalUncompressedSize  int     `json:"total_uncompressed_size"`
	CompressionRatio       float64 `json:"compression_ratio"`
	CompressionProgress    int     `json:"compression_progress"`
	CompressionSpaceSaving float64 `json:"compression_space_saving"`
	TotalBlobCount         int     `json:"total_blob_count"`
	SnapshotCount          int     `json:"snapshots_count"`
}

func readJson(data []byte) (RawDataMetrics, error) {
	var o RawDataMetrics
	err := json.Unmarshal(data, &o)
	if err != nil {
		return o, err
	}
	return o, nil
}
