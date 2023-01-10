package main

import (
	"os"
	"syscall"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	sourceTestFile             = "testdata/input.txt"
	minValidTestFileSize int64 = 6000
	maxValidTestFileSize int64 = 7000
)

func validateTestSourceFile(t *testing.T) {
	t.Helper()
	require.FileExists(t, sourceTestFile, "test file not exist")

	sourceTestFileStat, err := os.Stat(sourceTestFile)
	require.NoError(t, err, "cant get test source file stat")

	require.GreaterOrEqualf(
		t,
		sourceTestFileStat.Size(),
		minValidTestFileSize,
		"invalid test source file, it size should be greater or equal %d",
		minValidTestFileSize,
	)
	require.LessOrEqual(
		t,
		sourceTestFileStat.Size(),
		maxValidTestFileSize,
		"invalid test source file, it size should be less or equal %d",
		maxValidTestFileSize,
	)
}

func TestCopy(t *testing.T) {
	validateTestSourceFile(t)

	// Данные кейсы не отлавливаются на этапе валидации так как для это необходима
	// дополнительная более сложная обработка передаваемого пути назначения,
	// поэтому они будут получены непосредственно при инициализации копирования.
	// Вариации позитивных кейсов копирования проверяются в test.sh.

	// Корректно отработает только на UNIX-подобных системах.
	t.Run("destination not created before and it's directory", func(t *testing.T) {
		err := Copy("testdata/input.txt", "test/", 0, 0)

		// Ожидаем ошибку, что переданный destination - это директория.
		require.ErrorIs(t, err, syscall.EISDIR)
	})

	// Корректно отработает только на UNIX-подобных системах.
	t.Run("non-existent destination path", func(t *testing.T) {
		err := Copy("testdata/input.txt", "test_dir/test", 0, 0)

		// Ожидаем ошибку, что переданный destination не существует.
		require.ErrorIs(t, err, syscall.ENOENT)
	})
}

func Test_validate(t *testing.T) {
	validateTestSourceFile(t)

	sourceTestFileStat, err := os.Stat(sourceTestFile)
	require.NoError(t, err, "cant get test source file stat")

	type args struct {
		sourceFile      string
		destinationFile string
		offset          int64
		limit           int64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success - valid fields",
			args: args{
				sourceFile:      sourceTestFile,
				destinationFile: "resFile",
				offset:          0,
				limit:           0,
			},
			wantErr: false,
		},
		{
			name: "fail - empty source",
			args: args{
				sourceFile:      "",
				destinationFile: "resFile",
				offset:          0,
				limit:           0,
			},
			wantErr: true,
		},
		{
			name: "fail - empty dest",
			args: args{
				sourceFile:      sourceTestFile,
				destinationFile: "",
				offset:          0,
				limit:           0,
			},
			wantErr: true,
		},
		{
			name: "fail - negative limit",
			args: args{
				sourceFile:      sourceTestFile,
				destinationFile: "resFile",
				offset:          0,
				limit:           -1,
			},
			wantErr: true,
		},
		{
			name: "fail - negative offset",
			args: args{
				sourceFile:      sourceTestFile,
				destinationFile: "resFile",
				offset:          -1,
				limit:           0,
			},
			wantErr: true,
		},
		{
			name: "fail - offset greater than file size",
			args: args{
				sourceFile:      sourceTestFile,
				destinationFile: "resFile",
				offset:          sourceTestFileStat.Size() + 1000,
				limit:           0,
			},
			wantErr: true,
		},
		{
			name: "fail - non-existent source file",
			args: args{
				sourceFile:      "nonExistentFile",
				destinationFile: "resFile",
				offset:          sourceTestFileStat.Size() + 1000,
				limit:           0,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validate(
				tt.args.sourceFile,
				tt.args.destinationFile,
				tt.args.offset, tt.args.limit,
			); (err != nil) != tt.wantErr {
				t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_getEndOfCopying(t *testing.T) {
	validateTestSourceFile(t)

	sourceTestFileStat, err := os.Stat(sourceTestFile)
	require.NoError(t, err, "cant get test source file stat")

	type args struct {
		sourceFile string
		offset     int64
		limit      int64
	}

	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "success - limit and offset not given",
			args: args{
				sourceFile: sourceTestFile,
				offset:     0,
				limit:      0,
			},
			want:    sourceTestFileStat.Size(),
			wantErr: false,
		},
		{
			name: "success - only offset given",
			args: args{
				sourceFile: sourceTestFile,
				offset:     1000,
				limit:      0,
			},
			want:    sourceTestFileStat.Size() - 1000,
			wantErr: false,
		},
		{
			name: "success - only limit given and limit greater than source file",
			args: args{
				sourceFile: sourceTestFile,
				offset:     0,
				limit:      maxValidTestFileSize + 1000,
			},
			want:    sourceTestFileStat.Size(),
			wantErr: false,
		},
		{
			name: "success - only limit given and limit less than source file",
			args: args{
				sourceFile: sourceTestFile,
				offset:     0,
				limit:      1000,
			},
			want:    1000,
			wantErr: false,
		},
		{
			name: "success - limit and offset given but limit less file after offset",
			args: args{
				sourceFile: sourceTestFile,
				offset:     100,
				limit:      1000,
			},
			want:    1000,
			wantErr: false,
		},
		{
			name: "success - limit and offset given but limit greater file after offset",
			args: args{
				sourceFile: sourceTestFile,
				offset:     5500,
				limit:      2000,
			},
			want:    sourceTestFileStat.Size() - 5500,
			wantErr: false,
		},
		{
			name: "fail - getting source file info",
			args: args{
				sourceFile: "notExistedFile",
				offset:     0,
				limit:      0,
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getEndOfCopying(tt.args.sourceFile, tt.args.offset, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("getEndOfCopying() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getEndOfCopying() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_offsetValidation(t *testing.T) {
	validateTestSourceFile(t)

	type args struct {
		sourceFile string
		offset     int64
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success - given valid file path and offset",
			args: args{
				sourceFile: sourceTestFile,
				offset:     1000,
			},
			wantErr: false,
		},
		{
			name: "fail - getting source file info",
			args: args{
				sourceFile: "notExistedFile",
				offset:     0,
			},
			wantErr: true,
		},
		{
			name: "fail - offset greater than given file",
			args: args{
				sourceFile: sourceTestFile,
				offset:     maxValidTestFileSize + 1000,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := offsetValidation(tt.args.sourceFile, tt.args.offset); (err != nil) != tt.wantErr {
				t.Errorf("offsetValidation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
