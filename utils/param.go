package utils

func Audio(path string) []string {
	var args []string
	args = append(args, "-o", path)
	args = append(args, "-x", "--audio-format", "mp3")
	return args
}

func Video(path string) []string {
	var args []string
	args = append(args, "-o", path)
	return args
}

func Thumbnail(path string, args []string) []string {
	args = append(args, "--write-subs", "--write-auto-subs")
	return args
}

func Subtitle(path string, args []string) []string {
	args = append(args, "--write-thumbnail")
	return args
}

func OnlyAudio(path string, url string) []string {
	args := Audio(path)
	args = append(args, url)
	return args
}

func OnlyVideo(path string, url string) []string {
	args := Video(path)
	args = append(args, url)
	return args
}

func ThumbAndVideo(path string, url string) []string {
	args := Video(path)
	args = Thumbnail(path, args)
	args = append(args, url)
	return args
}

func SrtAndVideo(path string, url string) []string {
	args := Video(path)
	args = Subtitle(path, args)
	args = append(args, url)
	return args
}

func Complete(path string, url string) []string {
	args := Video(path)
	args = Subtitle(path, args)
	args = Thumbnail(path, args)
	args = append(args, url)
	return args
}
