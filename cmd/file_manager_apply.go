package cmd

import (
	"fmt"
	"time"

	"github.com/bytedance/sonic"
)

// applyFileManagerChunks executes a filemanager operation in chunks.
// It expects the API response body to be compatible with fileManagerOpResp.
func applyFileManagerChunks[T any](
	op string,
	items []T,
	chunkSize int,
	showProgress bool,
	continueOnError bool,
	ignoreErrors bool,
	exec func(filelist string) ([]byte, error),
) error {
	if chunkSize <= 0 {
		return fmt.Errorf("invalid chunk size: %d", chunkSize)
	}
	if len(items) == 0 {
		if showProgress {
			fmt.Printf("%s: no items\n", op)
		}
		return nil
	}
	if ignoreErrors && !continueOnError {
		return fmt.Errorf("--ignore-errors requires --continue-on-error")
	}

	total := len(items)
	appliedChunks := 0
	failedChunks := 0
	failedItems := 0

	for i := 0; i < total; i += chunkSize {
		end := i + chunkSize
		if end > total {
			end = total
		}
		chunk := items[i:end]

		if showProgress {
			fmt.Printf("\n[%d/%d] %s %d item(s)...\n", i+1, total, op, len(chunk))
		}

		filelist, err := sonic.MarshalString(chunk)
		if err != nil {
			if showProgress {
				fmt.Printf("chunk %d-%d/%d marshal failed: %v\n", i+1, end, total, err)
			}
			if continueOnError {
				failedChunks++
				failedItems += len(chunk)
				continue
			}
			return err
		}

		start := time.Now()
		b, err := exec(filelist)
		if err != nil {
			if showProgress {
				fmt.Printf("chunk %d-%d/%d failed: %v\n", i+1, end, total, err)
			}
			if continueOnError {
				failedChunks++
				failedItems += len(chunk)
				continue
			}
			return err
		}

		var resp fileManagerOpResp
		if err := sonic.Unmarshal(b, &resp); err != nil {
			e := fmt.Errorf("%s response parse failed: %w, raw=%s", op, err, string(b))
			if showProgress {
				fmt.Printf("chunk %d-%d/%d failed: %v\n", i+1, end, total, e)
			}
			if continueOnError {
				failedChunks++
				failedItems += len(chunk)
				continue
			}
			return e
		}

		if resp.Errno != 0 {
			itemFails := 0
			for _, info := range resp.Info {
				if info.Errno != 0 {
					itemFails++
					if showProgress {
						fmt.Printf("  - item failed: errno=%d path=%s\n", info.Errno, info.Path)
					}
				}
			}
			if itemFails == 0 {
				itemFails = len(chunk)
			}
			failedItems += itemFails
			failedChunks++

			e := fmt.Errorf("%s chunk failed, errno=%d request_id=%d", op, resp.Errno, resp.RequestId)
			if showProgress {
				fmt.Printf("chunk %d-%d/%d failed: %v\n", i+1, end, total, e)
			}
			if continueOnError {
				continue
			}
			return e
		}

		appliedChunks++
		if showProgress {
			if resp.TaskId != 0 {
				fmt.Printf("submitted chunk %d-%d/%d (taskid=%d, cost=%s)\n", i+1, end, total, resp.TaskId, time.Since(start).Truncate(10*time.Millisecond))
			} else {
				fmt.Printf("applied chunk %d-%d/%d (cost=%s)\n", i+1, end, total, time.Since(start).Truncate(10*time.Millisecond))
			}
		}
	}

	if failedChunks > 0 {
		msg := fmt.Sprintf("completed with failures: applied_chunks=%d failed_chunks=%d failed_items~=%d", appliedChunks, failedChunks, failedItems)
		if ignoreErrors {
			if showProgress {
				fmt.Println(msg)
			}
			return nil
		}
		return fmt.Errorf(msg)
	}
	return nil
}


