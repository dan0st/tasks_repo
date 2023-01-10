package main

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"syscall"
	"time"

	"github.com/cheggaaa/pb/v3"
)

func sourceFileValidation(sourceFile string) error {
	sourceFileInfo, err := os.Stat(sourceFile)
	// Проверяем наличие указанного файла-источника
	if err != nil {
		if os.IsNotExist(err) {
			return errors.New("source file is not exists")
		}
		return fmt.Errorf("getting source file info unexpected error: %w", err)
	}

	// Проверяем что указанный файл-источник - это не директория
	if sourceFileInfo.IsDir() {
		return errors.New("can't use source file - it's directory")
	}

	// Проверяем что указанный файл-источник - это регулярный файл
	if !sourceFileInfo.Mode().IsRegular() {
		return errors.New("can't use source file - it's not regular file")
	}

	return nil
}

func destinationFileValidation(destinationFile string) error {
	f, err := os.Stat(destinationFile)
	if err == nil {
		if f.IsDir() {
			return fmt.Errorf(
				"directory with name \"%s\" already exist in destination",
				destinationFile,
			)
		}

		return nil
	}

	var pathError *fs.PathError
	if ok := errors.As(err, &pathError); ok {
		if errors.Is(pathError.Err, syscall.ENOTDIR) {
			return fmt.Errorf("given destination \"%s\" is a directory", destinationFile)
		}
	} else {
		return fmt.Errorf("unexpected error: %w", err)
	}

	return nil
}

func offsetValidation(sourceFile string, offset int64) error {
	sourceFileStat, err := os.Stat(sourceFile)
	if err != nil {
		return fmt.Errorf("getting source file info unexpected error: %w", err)
	}

	if sourceFileStat.Size() < offset {
		return errors.New("offset exceeds file size")
	}

	return nil
}

func validate(
	sourceFile string,
	destinationFile string,
	offset int64,
	limit int64,
) error {
	if sourceFile == "" {
		return errors.New("source file not given")
	}

	if destinationFile == "" {
		return errors.New("destination file not given")
	}

	if offset < 0 || limit < 0 {
		return errors.New("offset and limit can't be less than 0")
	}

	// Валидируем файлы
	err := sourceFileValidation(sourceFile)
	if err != nil {
		return fmt.Errorf("source file validation error: %w", err)
	}

	err = destinationFileValidation(destinationFile)
	if err != nil {
		return fmt.Errorf("destination file validation error: %w", err)
	}

	// Валидируем смещение
	err = offsetValidation(sourceFile, offset)
	if err != nil {
		return fmt.Errorf("offset validation error: %w", err)
	}

	return nil
}

func getEndOfCopying(sourceFile string, offset int64, limit int64) (int64, error) {
	sourceFileInfo, err := os.Stat(sourceFile)
	if err != nil {
		return 0, fmt.Errorf("getting source file info error: %w", err)
	}

	switch {
	case limit == 0 && offset > 0:
		return sourceFileInfo.Size() - offset, nil
	case limit > 0 && offset == 0:
		if limit > sourceFileInfo.Size() {
			return sourceFileInfo.Size(), nil
		}
		return limit, nil
	case limit > 0 && offset > 0:
		restFile := sourceFileInfo.Size() - offset
		if restFile-limit > 0 {
			return limit, nil
		}
		return restFile, nil
	}

	return sourceFileInfo.Size(), nil
}

func Copy(fromPath, toPath string, offset, limit int64) error {
	err := validate(fromPath, toPath, offset, limit)
	if err != nil {
		log.Fatal(fmt.Errorf("validation: %w", err))
	}

	// limit рассматривается как количество копируемых байт либо от начала файла,
	// либо от offset, если он задан
	endOfCopying, err := getEndOfCopying(fromPath, offset, limit)
	if err != nil {
		log.Fatal(fmt.Errorf("end of copy calculation: %w", err))
	}

	fromFile, err := os.Open(fromPath)
	if err != nil {
		return fmt.Errorf("can't open source file while for copying: %w", err)
	}
	defer func() {
		err := fromFile.Close()
		if err != nil {
			fmt.Printf("closing source file error: %v", err)
		}
	}()

	if offset > 0 {
		_, err := fromFile.Seek(offset, 0)
		if err != nil {
			return fmt.Errorf("do offset in file error: %w", err)
		}
	}

	toFile, err := os.Create(toPath)
	if err != nil {
		return fmt.Errorf("can't create destination file for copying: %w", err)
	}
	defer func() {
		err := toFile.Close()
		if err != nil {
			fmt.Printf("closing destination file error: %v", err)
		}
	}()

	bar := pb.Full.Start64(endOfCopying)
	barFromFileProxyReader := bar.NewProxyReader(fromFile)

	_, err = io.CopyN(toFile, barFromFileProxyReader, endOfCopying)
	if err != nil {
		return fmt.Errorf("error while copying: %w", err)
	}
	bar.Finish()

	start := time.Now()
	fmt.Println("Start sync file")
	err = toFile.Sync()
	if err != nil {
		return fmt.Errorf("file sync after copying error: %w", err)
	}
	fmt.Printf("File sync finished in %v \n", time.Since(start))

	return nil
}
