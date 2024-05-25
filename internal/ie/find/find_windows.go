//go:build windows

package find

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/browserutils/kooky"
	"github.com/browserutils/kooky/internal/cookies"
	"github.com/browserutils/kooky/internal/ie"
)

type IEFinder struct {
	Browser string
}

var _ kooky.CookieStoreFinder = (*IEFinder)(nil)

var registerOnce sync.Once

func init() {
	browser := `ie+edge`
	// don't register multiple times for files shared between ie and edge
	registerOnce.Do(func() {
		kooky.RegisterFinder(browser, &IEFinder{Browser: browser})
	})
}

func (f *IEFinder) FindCookieStores() ([]kooky.CookieStore, error) {
	locApp, _ := os.UserCacheDir()
	home := os.Getenv(`USERPROFILE`)
	windows := os.Getenv(`windir`)
	appData, _ := os.UserConfigDir()

	type pathStruct struct {
		dir   string
		paths [][]string
	}

	// https://tzworks.com/prototypes/index_dat/id.users.guide.pdf
	paths := []pathStruct{
		{
			dir: windows,
			paths: [][]string{
				{`Cookies`}, // IE 4.0
			},
		},
		{
			dir: home,
			paths: [][]string{
				{`Cookies`}, // XP, Vista
			},
		},
		{
			dir: appData,
			paths: [][]string{
				{`Microsoft`, `Windows`, `Cookies`},
				{`Microsoft`, `Windows`, `Cookies`, `Low`},
				{`Microsoft`, `Windows`, `Cookies`, `Low`},
				{`Microsoft`, `Windows`, `Internet Explorer`, `UserData`, `Low`},
			},
		},
	}

	var cookiesFiles []kooky.CookieStore
	for _, p := range paths {
		if len(p.dir) == 0 {
			continue
		}
		for _, path := range p.paths {
			cookiesFiles = append(
				cookiesFiles,
				&cookies.CookieJar{
					CookieStore: &ie.CookieStore{
						CookieStore: &ie.IECacheCookieStore{
							DefaultCookieStore: cookies.DefaultCookieStore{
								BrowserStr:           f.Browser,
								IsDefaultProfileBool: true,
								FileNameStr:          filepath.Join(append(append([]string{p.dir}, path...), `index.dat`)...),
							},
						},
					},
				},
			)
		}
	}

	cookiesFiles = append(
		cookiesFiles,
		&cookies.CookieJar{
			CookieStore: &ie.CookieStore{
				CookieStore: &ie.ESECookieStore{
					DefaultCookieStore: cookies.DefaultCookieStore{
						BrowserStr:           f.Browser,
						IsDefaultProfileBool: true,
						FileNameStr:          filepath.Join(locApp, `Microsoft`, `Windows`, `WebCache`, `WebCacheV01.dat`),
					},
				},
			},
		},
	)

	return cookiesFiles, nil
}
