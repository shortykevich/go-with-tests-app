package fss

import (
	"io"

	"github.com/shortykevich/go-with-tests-app/db/inmem"
	"github.com/shortykevich/go-with-tests-app/pkg/jsontuil"
)

type FileSystemPlayerStorage struct {
	db io.ReadSeeker
}

// TODO: PostPlayerScore(player string) error
// TODO: RecordWin(player string)

func (f *FileSystemPlayerStorage) GetLeagueTable() ([]inmem.Player, error) {
	f.db.Seek(0, io.SeekStart)
	league, err := jsontuil.NewLeagueFromReader(f.db)
	if err != nil {
		return nil, err
	}
	return league, nil
}

func (f *FileSystemPlayerStorage) GetPlayerScore(player string) (int, error) {
	score, err := jsontuil.GetPlayerScoreFromReader(f.db, player)
	if err != nil {
		return 0, err
	}
	return score, nil
}
