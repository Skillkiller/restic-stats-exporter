package snapshot

import (
	"os"
	"reflect"
	"testing"
	"time"
)

func Test_readJson(t *testing.T) {
	type args struct {
		testFileName string
	}
	tests := []struct {
		name    string
		args    args
		want    []GroupData
		wantErr bool
	}{
		{
			name: "valid array",
			args: args{
				testFileName: "testdata/input1.json",
			},
			want: []GroupData{
				{
					GroupKey: GroupKey{Hostname: "SK12", Tags: []string{"full-server"}},
					Snapshots: []Snapshot{
						{
							Time:     mustParse(t, "2025-10-07T17:56:01.685056163+02:00"),
							Hostname: "SK12",
							Tags:     []string{"full-server"},
							Summary: Summary{
								BackupStart:         mustParse(t, "2025-10-07T17:56:01.685056163+02:00"),
								BackupEnd:           mustParse(t, "2025-10-07T18:02:13.257197421+02:00"),
								FilesNew:            72799,
								FilesChanged:        0,
								FilesUnmodified:     0,
								DirsNew:             10352,
								DirsChanged:         0,
								DirsUnmodified:      0,
								DataBlobs:           10352,
								TreeBlobs:           9919,
								DataAdded:           3155414417,
								DataAddedPacked:     1241206301,
								TotalFilesProcessed: 72799,
								TotalBytesProcessed: 3771315033,
							},
						},
						{
							Time:     mustParse(t, "2025-10-12T00:35:01.347812525+02:00"),
							Hostname: "SK12",
							Tags:     []string{"full-server"},
							Summary: Summary{
								BackupStart:         mustParse(t, "2025-10-12T00:35:01.347812525+02:00"),
								BackupEnd:           mustParse(t, "2025-10-12T00:36:10.861283523+02:00"),
								FilesNew:            1,
								FilesChanged:        69,
								FilesUnmodified:     73040,
								DirsNew:             0,
								DirsChanged:         80,
								DirsUnmodified:      10419,
								DataBlobs:           230,
								TreeBlobs:           73,
								DataAdded:           138116104,
								DataAddedPacked:     44128842,
								TotalFilesProcessed: 73110,
								TotalBytesProcessed: 3927390997,
							},
						},
					},
				},
				{
					GroupKey: GroupKey{Hostname: "SK12", Tags: []string{"kuma"}},
					Snapshots: []Snapshot{
						{
							Time:     mustParse(t, "2025-10-08T05:23:10.031203027+02:00"),
							Hostname: "SK12",
							Tags:     []string{"kuma"},
							Summary: Summary{
								BackupStart:         mustParse(t, "2025-10-08T05:23:10.031203027+02:00"),
								BackupEnd:           mustParse(t, "2025-10-08T05:23:26.623964395+02:00"),
								FilesNew:            7,
								FilesChanged:        0,
								FilesUnmodified:     0,
								DirsNew:             12,
								DirsChanged:         0,
								DirsUnmodified:      0,
								DataBlobs:           6,
								TreeBlobs:           8,
								DataAdded:           8550737,
								DataAddedPacked:     1825021,
								TotalFilesProcessed: 7,
								TotalBytesProcessed: 16325329,
							},
						},
						{
							Time:     mustParse(t, "2025-10-12T05:23:09.346002024+02:00"),
							Hostname: "SK12",
							Tags:     []string{"kuma"},
							Summary: Summary{
								BackupStart:         mustParse(t, "2025-10-12T05:23:09.346002024+02:00"),
								BackupEnd:           mustParse(t, "2025-10-12T05:23:24.842472768+02:00"),
								FilesNew:            0,
								FilesChanged:        1,
								FilesUnmodified:     6,
								DirsNew:             0,
								DirsChanged:         9,
								DirsUnmodified:      3,
								DataBlobs:           7,
								TreeBlobs:           7,
								DataAdded:           9785903,
								DataAddedPacked:     2079457,
								TotalFilesProcessed: 7,
								TotalBytesProcessed: 18025169,
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "empty array",
			args: args{
				testFileName: "testdata/empty.json",
			},
			want:    []GroupData{},
			wantErr: false,
		},
		{
			name: "invalid json",
			args: args{
				testFileName: "testdata/invalid.json",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := os.ReadFile(tt.args.testFileName)
			if err != nil {
				t.Fatalf("readJson test file: %v", err)
			}
			got, err := readJson(data)
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

// mustParse parses an RFC3339Nano string into time.Time and fails the test on error.
// Use t.Helper() to trace errors back to the caller.
func mustParse(t testing.TB, s string) time.Time {
	t.Helper()
	tt, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		t.Fatalf("failed to parse time %q: %v", s, err)
	}
	return tt
}

func Test_getTotalSnapshotCount(t *testing.T) {
	type args struct {
		data []GroupData
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "valid array",
			args: args{
				data: []GroupData{
					{
						GroupKey: GroupKey{Hostname: "SK12", Tags: []string{"full-server"}},
						Snapshots: []Snapshot{
							{
								Time:     mustParse(t, "2025-10-07T17:56:01.685056163+02:00"),
								Hostname: "SK12",
								Tags:     []string{"full-server"},
							},
						},
					},
				},
			},
			want: 1,
		},
		{
			name: "empty array",
			args: args{
				data: []GroupData{},
			},
			want: 0,
		},
		{
			name: "multiple group keys array",
			args: args{
				data: []GroupData{
					{
						GroupKey: GroupKey{Hostname: "SK12", Tags: []string{"full-server"}},
						Snapshots: []Snapshot{
							{
								Time:     mustParse(t, "2025-10-07T17:56:01.685056163+02:00"),
								Hostname: "SK12",
								Tags:     []string{"full-server"},
							},
							{
								Time:     mustParse(t, "2025-10-12T00:35:01.347812525+02:00"),
								Hostname: "SK12",
								Tags:     []string{"full-server"},
							},
						},
					},
					{
						GroupKey: GroupKey{Hostname: "SK12", Tags: []string{"kuma"}},
						Snapshots: []Snapshot{
							{
								Time:     mustParse(t, "2025-10-08T05:23:10.031203027+02:00"),
								Hostname: "SK12",
								Tags:     []string{"kuma"},
							},
							{
								Time:     mustParse(t, "2025-10-12T05:23:09.346002024+02:00"),
								Hostname: "SK12",
								Tags:     []string{"kuma"},
							},
						},
					},
				},
			},
			want: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getTotalSnapshotCount(tt.args.data); got != tt.want {
				t.Errorf("getTotalSnapshotCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getSnapshotCountByGroup(t *testing.T) {
	type args struct {
		group GroupData
	}
	tests := []struct {
		name  string
		args  args
		want  GroupKey
		want1 int
	}{
		{
			name: "empty data",
			args: args{
				group: GroupData{},
			},
			want:  GroupKey{},
			want1: 0,
		},
		{
			name: "zero snapshots",
			args: args{
				group: GroupData{
					GroupKey:  GroupKey{Hostname: "SK12", Tags: []string{"full-server"}},
					Snapshots: []Snapshot{},
				},
			},
			want:  GroupKey{Hostname: "SK12", Tags: []string{"full-server"}},
			want1: 0,
		},
		{
			name: "two snapshots",
			args: args{
				group: GroupData{
					GroupKey: GroupKey{Hostname: "SK12", Tags: []string{"kuma"}},
					Snapshots: []Snapshot{
						{
							Time:     mustParse(t, "2025-10-08T05:23:10.031203027+02:00"),
							Hostname: "SK12",
							Tags:     []string{"kuma"},
						},
						{
							Time:     mustParse(t, "2025-10-12T05:23:09.346002024+02:00"),
							Hostname: "SK12",
							Tags:     []string{"kuma"},
						},
					},
				},
			},
			want:  GroupKey{Hostname: "SK12", Tags: []string{"kuma"}},
			want1: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := getSnapshotCountByGroup(tt.args.group)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getSnapshotCountByGroup() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("getSnapshotCountByGroup() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_getLastSnapshotByGroup(t *testing.T) {
	type args struct {
		group GroupData
	}
	tests := []struct {
		name    string
		args    args
		want    GroupKey
		want1   Snapshot
		wantErr bool
	}{
		{
			name: "empty data",
			args: args{
				group: GroupData{},
			},
			want:    GroupKey{},
			want1:   Snapshot{},
			wantErr: true,
		},
		{
			name: "zero snapshots",
			args: args{
				group: GroupData{
					GroupKey:  GroupKey{Hostname: "SK12", Tags: []string{"full-server"}},
					Snapshots: []Snapshot{},
				},
			},
			want:    GroupKey{},
			want1:   Snapshot{},
			wantErr: true,
		},
		{
			name: "one snapshot",
			args: args{
				group: GroupData{
					GroupKey: GroupKey{Hostname: "SK12", Tags: []string{"kuma"}},
					Snapshots: []Snapshot{
						{
							Time:     mustParse(t, "2025-10-08T05:23:10.031203027+02:00"),
							Hostname: "SK12",
							Tags:     []string{"kuma"},
						},
					},
				},
			},
			want: GroupKey{Hostname: "SK12", Tags: []string{"kuma"}},
			want1: Snapshot{
				Time:     mustParse(t, "2025-10-08T05:23:10.031203027+02:00"),
				Hostname: "SK12",
				Tags:     []string{"kuma"},
			},
			wantErr: false,
		},
		{
			name: "two snapshots",
			args: args{
				group: GroupData{
					GroupKey: GroupKey{Hostname: "SK12", Tags: []string{"kuma"}},
					Snapshots: []Snapshot{
						{
							Time:     mustParse(t, "2025-10-08T05:23:10.031203027+02:00"),
							Hostname: "SK12",
							Tags:     []string{"kuma"},
						},
						{
							Time:     mustParse(t, "2025-10-12T05:23:09.346002024+02:00"),
							Hostname: "SK12",
							Tags:     []string{"kuma"},
						},
					},
				},
			},
			want: GroupKey{Hostname: "SK12", Tags: []string{"kuma"}},
			want1: Snapshot{
				Time:     mustParse(t, "2025-10-12T05:23:09.346002024+02:00"),
				Hostname: "SK12",
				Tags:     []string{"kuma"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := getLastSnapshotByGroup(tt.args.group)
			if (err != nil) != tt.wantErr {
				t.Errorf("getLastSnapshotByGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getLastSnapshotByGroup() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("getLastSnapshotByGroup() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_getSnapshotMetricsByGroup(t *testing.T) {
	type args struct {
		group GroupData
	}
	tests := []struct {
		name    string
		args    args
		want    GroupKey
		want1   SnapshotMetrics
		wantErr bool
	}{
		{
			name: "empty data",
			args: args{
				group: GroupData{},
			},
			want:    GroupKey{},
			want1:   SnapshotMetrics{},
			wantErr: true,
		},
		{
			name: "zero snapshots",
			args: args{
				group: GroupData{
					GroupKey:  GroupKey{Hostname: "SK12", Tags: []string{"full-server"}},
					Snapshots: []Snapshot{},
				},
			},
			want:    GroupKey{},
			want1:   SnapshotMetrics{},
			wantErr: true,
		},
		{
			name: "two snapshots",
			args: args{
				group: GroupData{
					GroupKey: GroupKey{Hostname: "SK12", Tags: []string{"kuma"}},
					Snapshots: []Snapshot{
						{
							Time:     mustParse(t, "2025-10-08T05:23:10.031203027+02:00"),
							Hostname: "SK12",
							Tags:     []string{"kuma"},
							Summary: Summary{
								BackupStart:         mustParse(t, "2025-10-07T17:56:01.685056163+02:00"),
								BackupEnd:           mustParse(t, "2025-10-07T18:02:13.257197421+02:00"),
								FilesNew:            72799,
								FilesChanged:        0,
								FilesUnmodified:     0,
								DirsNew:             10352,
								DirsChanged:         0,
								DirsUnmodified:      0,
								DataBlobs:           10352,
								TreeBlobs:           9919,
								DataAdded:           3155414417,
								DataAddedPacked:     1241206301,
								TotalFilesProcessed: 72799,
								TotalBytesProcessed: 3771315033,
							},
						},
						{
							Time:     mustParse(t, "2025-10-12T05:23:09.346002024+02:00"),
							Hostname: "SK12",
							Tags:     []string{"kuma"},
							Summary: Summary{
								BackupStart:         mustParse(t, "2025-10-12T00:35:01.347812525+02:00"),
								BackupEnd:           mustParse(t, "2025-10-12T00:36:10.861283523+02:00"),
								FilesNew:            1,
								FilesChanged:        69,
								FilesUnmodified:     73040,
								DirsNew:             0,
								DirsChanged:         80,
								DirsUnmodified:      10419,
								DataBlobs:           230,
								TreeBlobs:           73,
								DataAdded:           138116104,
								DataAddedPacked:     44128842,
								TotalFilesProcessed: 73110,
								TotalBytesProcessed: 3927390997,
							},
						},
					},
				},
			},
			want: GroupKey{Hostname: "SK12", Tags: []string{"kuma"}},
			want1: SnapshotMetrics{
				Time:                mustParse(t, "2025-10-12T05:23:09.346002024+02:00"),
				BackupStart:         mustParse(t, "2025-10-12T00:35:01.347812525+02:00"),
				BackupEnd:           mustParse(t, "2025-10-12T00:36:10.861283523+02:00"),
				FilesNew:            1,
				FilesChanged:        69,
				FilesUnmodified:     73040,
				DirsNew:             0,
				DirsChanged:         80,
				DirsUnmodified:      10419,
				DataBlobs:           230,
				TreeBlobs:           73,
				DataAdded:           138116104,
				DataAddedPacked:     44128842,
				TotalFilesProcessed: 73110,
				TotalBytesProcessed: 3927390997,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := getSnapshotMetricsByGroup(tt.args.group)
			if (err != nil) != tt.wantErr {
				t.Errorf("getSnapshotMetricsByGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getSnapshotMetricsByGroup() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("getSnapshotMetricsByGroup() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
