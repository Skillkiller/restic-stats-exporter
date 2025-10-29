package statistic

import (
	"reflect"
	"testing"
)

func Test_readJson(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    RawDataMetrics
		wantErr bool
	}{
		{
			name: "empty data",
			args: args{
				data: []byte(""),
			},
			want:    RawDataMetrics{},
			wantErr: true,
		},
		{
			name: "empty repository",
			args: args{
				data: []byte("{\"total_size\":0,\"snapshots_count\":0}"),
			},
			want:    RawDataMetrics{},
			wantErr: false,
		},
		{
			name: "filled repository",
			args: args{
				data: []byte("{\"total_size\":181885552,\"total_uncompressed_size\":203507483,\"compression_ratio\":1.1188765724503504,\"compression_progress\":100,\"compression_space_saving\":10.624636834607204,\"total_blob_count\":979,\"snapshots_count\":5}"),
			},
			want: RawDataMetrics{
				TotalSize:              181885552,
				TotalUncompressedSize:  203507483,
				CompressionRatio:       1.1188765724503504,
				CompressionProgress:    100,
				CompressionSpaceSaving: 10.624636834607204,
				TotalBlobCount:         979,
				SnapshotCount:          5,
			},
			wantErr: false,
		},
		{
			name: "filled repository 2",
			args: args{
				data: []byte("{\"total_size\":591642919300,\"total_uncompressed_size\":624191665868,\"compression_ratio\":1.05501417410101,\"compression_progress\":100,\"compression_space_saving\":5.214543600600264,\"total_blob_count\":1162492,\"snapshots_count\":79}"),
			},
			want: RawDataMetrics{
				TotalSize:              591642919300,
				TotalUncompressedSize:  624191665868,
				CompressionRatio:       1.05501417410101,
				CompressionProgress:    100,
				CompressionSpaceSaving: 5.214543600600264,
				TotalBlobCount:         1162492,
				SnapshotCount:          79,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readJson(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("readJson() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readJson() got = %v, want %v", got, tt.want)
			}
		})
	}
}
