# kooky

Package kooky contains routines to reach into cookie stores for Chrome
and Safari, and retrieve the cookies.

You don't usually want to do this, but when you do, you do. Apparently
it's common enough that there is example code scattered around the net
in various languages. If you wanted it in Go, you're in luck.

It aspires to be pure Go (I spent quite a while making
[go-sqlite/sqlite3](https://github.com/go-sqlite/sqlite3) work for
it), but I guess the keychain parts
([keybase/go-keychain](http://github.com/keybase/go-keychain)) mess that up.

## Status

[![No Maintenance Intended](http://unmaintained.tech/badge.svg)](http://unmaintained.tech/)

Basic functionality works, on MacOS. I expect Linux to work too, since
it doesn't encrypt. **The API is currently not expected to be at all stable.**

PRs more than welcome.

TODOs

- [ ] Make it work on Windows. (Look at
      [this](https://play.golang.org/p/fknP9AuLU-) and
      [this](https://github.com/cfstras/chromecsv/blob/master/crypt_windows.go)
      to learn how to decrypt.)
- [ ] Handle rows in Chrome's cookie DB with other than 14 columns (?)
- [ ] Make it work for Firefox.

## Example usage
```go
usr, _ := user.Current()

var cookies []*kooky.Cookie
var err error

chrome := false
if chrome {
	cookiesFile := fmt.Sprintf("%s/Library/Application Support/Google/Chrome/Default/Cookies", usr.HomeDir)
	cookies, err = kooky.ReadChromeCookies(cookiesFile, "", "", time.Time{})
} else {
	cookiesFile := fmt.Sprintf("%s/Library/Cookies/Cookies.binarycookies", usr.HomeDir)
	cookies, err = kooky.ReadSafariCookies(cookiesFile, "", "", time.Time{})
}
if err != nil {
	return err
}
for _, cookie := range cookies {
	fmt.Println(cookie)
}
```

## Thanks/references
- Thanks to [@dacort](http://github.com/dacort) for MacOS cookie decrypting
  code at https://gist.github.com/dacort/bd6a5116224c594b14db.
- Thanks to [@as0ler](http://github.com/as0ler)
  (and originally [@satishb3](http://github.com/satishb3) I believe) for
  Safari cookie-reading Python code at https://github.com/as0ler/BinaryCookieReader.