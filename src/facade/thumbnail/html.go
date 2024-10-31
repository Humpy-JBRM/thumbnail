package facade

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"
	"thumbnailer/src/data"
	"time"

	"github.com/chromedp/chromedp"
)

type htmlThumbnailer struct {
	name string
}

func NewHtmlThumbnailer() Thumbnailer {
	return &htmlThumbnailer{
		name: "html",
	}
}

func (t *htmlThumbnailer) GetThumbnail(u *url.URL) (data.Thumbnail, error) {
	// return &data.ThumbnailImpl{
	// 	Content: make([]byte, 0),
	// }, nil
	var err error
	var thumbnail data.ThumbnailImpl

	// This is a two-stage process:
	//
	//	1: use chromedp to screenshot the HTML
	//
	//	2: use the image thumbnailer to generate the "real" thumbnail
	destDir, err := os.MkdirTemp("", "")
	if err != nil {
		log.Printf("HtmlThumbnail(%s, %s): Could not create temp directory: %s", u.String(), destDir, err.Error())
		thumbnail.Error = err
		return &thumbnail, err
	}
	defer os.RemoveAll(destDir)

	// TODO(john): pay the instantiation cost once.  Consider using a pool of instances
	ctx, cancel := chromedp.NewContext(
		context.Background(),
	)
	// release the browser resources when
	// it is no longer needed
	defer cancel()
	var screenshotBuffer []byte
	chromedp.WindowSize(1024, 768)
	err = chromedp.Run(ctx,
		chromedp.Navigate(u.String()),
		chromedp.CaptureScreenshot(&screenshotBuffer),
	)
	if err != nil {
		log.Printf("HtmlThumbnail(%s, %s): Could not create temp directory: %s", u.String(), destDir, err.Error())
		thumbnail.Error = err
		return &thumbnail, err
	}

	// file permissions: 0644 (Owner: read/write, Group: read, Others: read)
	// write the response body to an image file
	destFile := path.Join(destDir, fmt.Sprintf("%s.png", u.Host))
	err = os.WriteFile(destFile, screenshotBuffer, 0644)
	if err != nil {
		log.Printf("HtmlThumbnail(%s, %s): Could not create temp directory: %s", u.String(), destDir, err.Error())
		thumbnail.Error = err
		return &thumbnail, err
	}

	// Make sure that the png file got created
	_, err = os.Stat(destFile)
	if err != nil {
		log.Printf("HtmlThumbnail(%s, %s): %s: %s", u.String(), destDir, destFile, err.Error())
		thumbnail.Error = err
		return &thumbnail, err
	}

	// Now delegate to the image thumbnailer
	imageUrl, err := url.Parse(fmt.Sprintf("file://%s", destFile))
	if err != nil {
		log.Printf("HtmlThumbnail(%s, %s): %s", destFile, destDir, err.Error())
		thumbnail.Error = err
		return &thumbnail, err
	}
	thumb, err := NewImageThumbnailer().GetThumbnail(imageUrl)
	if err != nil {
		log.Printf("HtmlThumbnail(%s, %s): %s", destFile, destDir, err.Error())
		thumbnail.Error = err
		return &thumbnail, err
	}
	return thumb, err
}

func ExecuteHeadless(command []string, wg *sync.WaitGroup, timeout time.Duration) error {
	var err error

	defer func() {
		if wg != nil {
			wg.Done()
		}
	}()

	// Ths is because some html can cause phantomjs to hang
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel() // The cancel should be deferred so resources are cleaned up
	// cmd := exec.CommandContext(ctx, command[0], command[1:]...)
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Stderr = nil
	cmd.Stdout = nil
	//cmd.Env = append(cmd.Env, "QT_QPA_PLATFORM=offscreen")
	commandLine := strings.Join(command, " ")
	log.Printf("Exec(): %s", commandLine)
	out, err := cmd.CombinedOutput()

	// We want to check the context error to see if the timeout was executed.
	// The error returned by cmd.Output() will be OS specific based on what
	// happens when a process is killed.
	// if ctx.Err() == context.DeadlineExceeded {
	//      log.Printf("ExecuteHeadless(%s): timed out", commandLine)
	//      return fmt.Errorf("ExecuteHeadless(%s): timed out", commandLine)
	// }

	output := string(out)
	if output != "" {
		log.Printf("PhantomJS: %s", output)
	}
	if err != nil {
		log.Printf("ExecuteHeadless(%s): cmd.Start() failed with %s", commandLine, strings.ReplaceAll(string(out), "\n", " "))
		return fmt.Errorf("Execute(%s, %s): cmd.Start() failed with '%s'", command[0], command[1:], err)
	}

	return nil
}
