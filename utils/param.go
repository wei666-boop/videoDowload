package utils

func Audio(path string, url string) []string {
	var args []string
	args = append(args, "-o", path)
	args = append(args, "-x", "--audio-format", "mp3")
	args = append(args, url)
	return args
}

func Video(path string, url string) []string {
	var args []string
	args = append(args, "-o", path)
	args = append(args, url)
	return args
}

func Thumbnail(path string, url string) []string {
	var args []string
	args = append(args, "--skip-download")
	args = append(args, "--write-subs", "--write-auto-subs")
	args = append(args, "-o", path, url)
	return args
}

func Subtitle(path string, url string) []string {
	var args []string
	args = append(args, "--skip-download")
	args = append(args, "--write-thumbnail")
	args = append(args, "-o", path, url)
	return args
}
