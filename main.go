package main

import (
	"context"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const fileOpenWaitSeconds = 5
const tailSleepSeconds = 1

func invokeCommand(wg *sync.WaitGroup, ctx context.Context, cancel func()) {
	if len(os.Args) > 1 {
		log.Printf("executing cmd=%s, args=%v", os.Args[1], os.Args[2:])
		cmd := exec.CommandContext(ctx, os.Args[1], os.Args[2:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		go func() {
			wg.Add(1)
			if err := cmd.Run(); err != nil {
				log.Printf("command returned with error: %v", err)
			} else {
				log.Print("command exited with no errors")
			}
			wg.Done()
			cancel()
		}()
	} else {
		log.Print("not invoking any command")
	}
}

func tailFile(wg *sync.WaitGroup, ctx context.Context) {
	filename := os.Getenv("LOG")
	if filename == "" {
		log.Print("no log file specified, will not tail log file")
		return
	}

	go func() {
		wg.Add(1)

		var (
			f   *os.File
			err error
		)

		defer func() {
			wg.Done()
			if f != nil {
				f.Close()
			}
			log.Print("exiting tail")
		}()

		// check if we need to terminate prematurely - close the filehandle
		// so that the io.Copy() can terminate
		go func() {
			wg.Add(1)

			log.Print("close goroutine waiting for context to terminate...")
			<-ctx.Done()
			if f != nil {
				f.Close()
			}

			log.Print("exiting file close goroutine")
			wg.Done()
		}()

		for {
			// try to open file
			for {
				f, err = os.Open(filename)
				if err == nil {
					break
				}
				log.Printf("could not open %s for reading - sleeping...", filename)
				select {
				case <-ctx.Done():
					return
				case <-time.After(fileOpenWaitSeconds * time.Second):
				}
			}

			log.Printf("successfully opened %s for reading", filename)

			for {
				_, err := io.Copy(os.Stdout, f)
				if err != nil {
					log.Printf("tail encountered error: %v", err)
					return
				}
				if _, err := os.Stat(filename); err != nil {
					log.Printf("error encountered while getting file info - file probably got deleted: %v", err)
					f = nil
					break // out of tail loop
				}

				// we have probably reached EOF - sleep to see if more data
				// gets added
				select {
				case <-ctx.Done():
					return
				case <-time.After(tailSleepSeconds * time.Second):
				}
			}
		}
	}()
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	sigCtx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup

	invokeCommand(&wg, ctx, cancel)
	tailFile(&wg, ctx)

	<-sigCtx.Done()
	log.Print("shutting down...")
	cancel()
	wg.Wait()
	log.Print("all goroutines exited")
}
