package utils

import pb "gopkg.in/cheggaaa/pb.v1"

func NewProgressBar(prefix string, count int64) *pb.ProgressBar {
	bar := pb.New64(count)
	bar.ShowCounters = false
	bar.ShowSpeed = false

	return bar.Prefix(prefix)
}
