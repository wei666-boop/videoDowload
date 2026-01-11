package utils

import (
	"fmt"
	"os"
	"os/exec"
)

func AudioAndVideoStart(cmd *exec.Cmd) error {

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return err
	}
	return cmd.Wait()
}

func ThumbnailORSubtitleStart(cmd *exec.Cmd) error {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return err
	}
	return cmd.Wait()
}

//生成mkv格式

func GetMKV(video, subtitle, thumbnail, output string) error {
	var mkvCmd *exec.Cmd

	var args []string

	args = append(args, "-i", video)

	if subtitle != "" && thumbnail == "" {
		args = append(args, "-i", subtitle)
		args = append(args, "-map", "0:v")
		args = append(args, "-map", "0:a?")
		args = append(args, "-map", "1")
		args = append(args, "-c", "copy")
		args = append(args, "-c:s", "srt")
		args = append(args, "-metadata:s:s:0", "language=chi")
		args = append(args, output)
	}

	if thumbnail != "" && subtitle == "" {
		args = append(args, "-i", thumbnail)
		args = append(args, "-map", "0:v")
		args = append(args, "-map", "0:a?")
		args = append(args, "-map", "1")
		args = append(args, "-c", "copy")
		args = append(args, "-disposition:v:1", "attached_pic")
		args = append(args, output)
	}

	if thumbnail != "" && subtitle != "" {
		args = append(args, "-i", subtitle)
		args = append(args, "-i", thumbnail)
		args = append(args, "-map", "0:v")
		args = append(args, "-map", "0:a?")
		args = append(args, "-map", "1")
		args = append(args, "-map", "2")
		args = append(args, "-c", "copy")
		args = append(args, "-c:s", "srt")
		args = append(args, "-disposition:v:1", "attached_pic")
		args = append(args, "-metadata:s:s:0", "language=chi")
		args = append(args, output)
	}

	fmt.Println(args)

	mkvCmd = exec.Command("ffmpeg", args...)

	mkvStdErr, err := mkvCmd.StderrPipe()
	if err != nil {
		return err
	}

	mkvCmd.Start()

	//开启一个新的协程来给ffmpeg输出日志用
	go func() {
		buf := make([]byte, 1024)
		for {
			n, _ := mkvStdErr.Read(buf)
			if n == 0 {
				break
			}
			fmt.Println(string(buf[:n]))
		}
	}()

	return mkvCmd.Wait()
}
