package stats

import (
	"log"
	"sync"

	"github.com/ivan-jabrony/gitfame/internal/fileparser"
)

func summarizeCommits(commits []fileparser.Commit, isCommitter bool) []fileparser.PersonData {
	summedStatistics := []fileparser.PersonData{}
	personToStats := make(map[string]fileparser.PersonData)

	for _, commit := range commits {
		var curAuthor string

		if isCommitter {
			curAuthor = commit.Commiter
		} else {
			curAuthor = commit.Author
		}

		curPersonData := personToStats[curAuthor]

		if curPersonData.Commits == nil {
			curPersonData.Commits = make(map[string]struct{})
		}
		if curPersonData.Files == nil {
			curPersonData.Files = make(map[string]struct{})
		}
		curPersonData.Name = curAuthor
		curPersonData.Lines += commit.Lines
		curPersonData.Commits[commit.SHA] = struct{}{}
		curPersonData.Files[commit.File] = struct{}{}
		personToStats[curAuthor] = curPersonData

	}

	for _, personStat := range personToStats {
		summedStatistics = append(summedStatistics, personStat)
	}

	return summedStatistics
}

func RunParallel(filePaths []string, maxThreads int, revision, repositoryPath string, isCommitter bool) *[]fileparser.PersonData {
	commitStatsGlobal := []fileparser.Commit{}
	pb := &ProgressBar{total: len(filePaths)}

	mu := make(chan struct{}, 1)
	jobs := make(chan string, len(filePaths))
	wg := sync.WaitGroup{}
	wg.Add(len(filePaths))

	if maxThreads < 0 || maxThreads > 8 {
		// if wasn't specified by user or too big
		maxThreads = 8
	}

	for i := 0; i < min(len(filePaths), maxThreads); i++ {
		go func() {
			for {
				file := <-jobs
				parsedCommits := fileparser.ParseFile(file, repositoryPath, revision)
				mu <- struct{}{}
				commitStatsGlobal = append(commitStatsGlobal, parsedCommits...)
				<-mu
				pb.Increment()
				wg.Done()
			}
		}()
	}

	log.Printf("Started collecting stats, files to parse:%v.\n", len(filePaths))
	for _, file := range filePaths {
		jobs <- file
	}

	wg.Wait()

	sumStats := summarizeCommits(commitStatsGlobal, isCommitter)
	log.Printf("Finished counting statistics for repository with %v personal statistics total.\n", len(sumStats))

	return &sumStats
}
