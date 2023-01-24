/*
Copyright © 2023 Daniel Unverricht (daniel@unverricht.net)
*/
package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/andybalholm/brotli"
	"github.com/klauspost/compress/flate"
	"github.com/klauspost/compress/gzip"
	"go.uber.org/zap"

	"github.com/spf13/cobra"
)

var precompressStaticFilesCmd = &cobra.Command{
	Use:   "precompressStaticFiles",
	Short: "Create precompressed static files",
	RunE: func(cmd *cobra.Command, args []string) error {
		rootDir, _ := cmd.Flags().GetString("rootDir")
		addGzip, _ := cmd.Flags().GetBool("gzip")
		addBrotli, _ := cmd.Flags().GetBool("brotli")
		addDeflate, _ := cmd.Flags().GetBool("deflate")

		fileNames := []string{}
		err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
			if path != rootDir && !info.IsDir() && !strings.HasSuffix(path, ".brotli") && !strings.HasSuffix(path, ".gz") && !strings.HasSuffix(path, "deflate") {
				fileNames = append(fileNames, path)
			}
			return nil
		})
		if err != nil {
			return err
		}

		var buffer bytes.Buffer

		for _, file := range fileNames {
			if addGzip {
				buffer.Reset()
				gzipWriter, err := gzip.NewWriterLevel(&buffer, gzip.BestCompression)
				if err != nil {
					return err
				}

				start := time.Now()
				fileContent, err := os.ReadFile(file)
				if err != nil {
					return err
				}
				if _, err := gzipWriter.Write(fileContent); err != nil {
					return err
				}
				if err := gzipWriter.Flush(); err != nil {
					return err
				}
				if err := gzipWriter.Close(); err != nil {
					return err
				}
				if err := os.WriteFile(fmt.Sprintf("%s.gz", file), buffer.Bytes(), 0664); err != nil {
					return err
				}
				fmt.Printf("Compressed with gzip: %s -> %s.gz (%s)\n", file, file, time.Since(start))
			}
			if addBrotli {
				buffer.Reset()
				brotliWriter := brotli.NewWriterLevel(&buffer, brotli.BestCompression)
				start := time.Now()
				fileContent, err := os.ReadFile(file)
				if err != nil {
					return err
				}
				if _, err := brotliWriter.Write(fileContent); err != nil {
					return err
				}
				if err := brotliWriter.Flush(); err != nil {
					return err
				}
				if err := brotliWriter.Close(); err != nil {
					return err
				}
				if err := os.WriteFile(fmt.Sprintf("%s.brotli", file), buffer.Bytes(), 0664); err != nil {
					return err
				}
				fmt.Printf("Compressed with brotli: %s -> %s.brotli (%s)\n", file, file, time.Since(start))
			}
			if addDeflate {
				buffer.Reset()
				deflateWriter, err := flate.NewWriter(&buffer, flate.BestCompression)
				if err != nil {
					return err
				}
				start := time.Now()
				fileContent, err := os.ReadFile(file)
				if err != nil {
					return err
				}
				if _, err := deflateWriter.Write(fileContent); err != nil {
					return err
				}
				if err := deflateWriter.Flush(); err != nil {
					return err
				}
				if err := deflateWriter.Close(); err != nil {
					return err
				}
				if err := os.WriteFile(fmt.Sprintf("%s.deflate", file), buffer.Bytes(), 0664); err != nil {
					return err
				}
				fmt.Printf("Compressed with deflate : %s -> %s.deflate (%s)\n", file, file, time.Since(start))
			}
		}

		return nil
	},
}

func init() {
	precompressStaticFilesCmd.Flags().StringP("rootDir", "r", ".", "root directory of files")
	if err := precompressStaticFilesCmd.MarkFlagRequired("rootDir"); err != nil {
		zap.S().Error(err)
	}

	precompressStaticFilesCmd.Flags().Bool("gzip", false, "add gzip")
	precompressStaticFilesCmd.Flags().Bool("brotli", false, "add brotli")
	precompressStaticFilesCmd.Flags().Bool("deflate", false, "add deflate")

	rootCmd.AddCommand(precompressStaticFilesCmd)
}
