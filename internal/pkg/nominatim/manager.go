package nominatim

import (
	"sync"
	"time"
)

type urlIndex struct {
	sync.Mutex
	index      int
	bannedUrls map[string]time.Time
}

var urlIndexManager = &urlIndex{
	Mutex:      sync.Mutex{},
	index:      0,
	bannedUrls: make(map[string]time.Time),
}

var serverUrls = [3]string{
	"amsterdam.nominatim.openstreetmap.org",
	"slough.nominatim.openstreetmap.org",
	"corvallis.nominatim.openstreetmap.org",
}

const timeToBan = 10 * time.Minute

func getNextUrl() *string {
	urlIndexManager.Lock()
	defer urlIndexManager.Unlock()

	if len(serverUrls) == urlIndexManager.index {
		urlIndexManager.index = 0
	}

	iteration := 0
	for {
		iteration += 1
		if len(serverUrls) < iteration {
			return nil
		}

		url := &serverUrls[urlIndexManager.index]

		if timestamp, found := urlIndexManager.bannedUrls[*url]; found {
			if timestamp.Before(time.Now()) {
				delete(urlIndexManager.bannedUrls, *url)
				urlIndexManager.index += 1

				return url
			}
		} else {
			urlIndexManager.index += 1
			return url
		}
	}
}

func banUrl(url string) {
	urlIndexManager.Lock()
	defer urlIndexManager.Unlock()

	urlIndexManager.bannedUrls[url] = time.Now().Add(timeToBan)
}
