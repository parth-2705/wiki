package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/dgraph-io/badger"
)

// func writeFile(body, filename string) {
// 	file, error := os.Create(filename)
// 	if error != nil {
// 		fmt.Println(error)
// 	}
// 	defer file.Close()
// 	file.WriteString(body)

// }
var counter = 0

func parsePage(url string) {
	fullUrl := ""
	if strings.Contains(url, "https://") {
		if strings.Contains(url, "https://wikipedia.org") || strings.Contains(url, "https://en.wikipedia.org") {
			fullUrl = url
		} else {
			fmt.Println("External Link")
			return
		}
	} else {
		fullUrl = "https://wikipedia.org" + url
	}

	db, err := badger.Open(badger.DefaultOptions("tmp/badger").WithTruncate(true))
	if err != nil {
		log.Fatal("1", err)
	}
	counter++
	if counter > 30 {
		return
	}

	response, error0 := http.Get(fullUrl)

	if error0 != nil {
		fmt.Println("2", error0)
	}
	if response.StatusCode > 400 {
		fmt.Println("Status Code:", response.StatusCode)
	}

	defer response.Body.Close()

	doc, error1 := goquery.NewDocumentFromReader(response.Body)
	if error1 != nil {
		fmt.Println("3", error1)
	}

	body := ""
	var links []string
	doc.Find("div.mw-parser-output").Find("p").Each(func(index int, item *goquery.Selection) {
		para := strings.TrimSpace(item.Text())
		body = body + para
		// fmt.Println(body)
		item.Find("a").Each(func(index1 int, item1 *goquery.Selection) {
			link, _ := item1.Attr("href")
			if len(link) == 0 {
				fmt.Println("Null Link")
				return
			}
			if link[0] != '#' {
				link = strings.TrimRight(link, "#")
				fetcherr := db.View(func(txn *badger.Txn) error {
					fullLink := "https://wikipedia.org" + link
					_, err := txn.Get([]byte(fullLink))

					return err
				})
				if fetcherr != nil {
					links = append(links, link)
					fmt.Println(fetcherr, link)
				} else {
					fmt.Println("already visited : ", link)
				}

			}

		})
	})

	dberr := db.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry([]byte(fullUrl), []byte(body))
		fmt.Println(counter)
		fmt.Println(fullUrl)
		err := txn.SetEntry(e)
		return err
	})
	if dberr != nil {
		fmt.Println("4", dberr)
	}

	db.Close()

	fmt.Println(counter)
	fmt.Println(fullUrl)
	if len(links) == 0 {
		fmt.Println("Oops! No more links to visit")
		return

	} else {
		for _, s := range links {
			parsePage(s)
		}
	}

}

func main() {
	counter = 0
	parsePage("https://en.wikipedia.org/wiki/Charles_IV_of_Spain")

	// writeFile(bottom, "test2.html")

	// response, error0 := http.Get("https://en.wikipedia.org/wiki/Species")

	// if error0 != nil {
	// 	fmt.Println(error0)
	// }
	// if response.StatusCode > 400 {
	// 	fmt.Println("Status Code:", response.StatusCode)
	// }

	// defer response.Body.Close()

	// doc, error1 := goquery.NewDocumentFromReader(response.Body)
	// if error1 != nil {
	// 	fmt.Println(error1)
	// }

	// body := ""
	// // var links []string
	// doc.Find("div.mw-parser-output").Find("p").Each(func(index int, item *goquery.Selection) {
	// 	para := strings.TrimSpace(item.Text())
	// 	body = body + para
	// 	// fmt.Println(body)
	// 	item.Find("a").Each(func(index1 int, item1 *goquery.Selection) {
	// 		link, _ := item1.Attr("href")
	// 		if link[0] != '#' {
	// 			link = strings.TrimRight(link, "#")
	// 			// fetcherr := db.View(func(txn *badger.Txn) error {
	// 			// 	fullLink := "https://wikipedia.org" + link
	// 			// 	_, err := txn.Get([]byte(fullLink))

	// 			// 	return err
	// 			// })
	// 			// if fetcherr != nil {
	// 			// 	links = append(links, link)
	// 			// } else {
	// 			// 	fmt.Println("already visited : ", link)
	// 			// }
	// 			fmt.Println(link)

	// 		}

	// 	})
	// })

	// dberr := db.Update(func(txn *badger.Txn) error {
	// 	e := badger.NewEntry([]byte(url), []byte(bottom))
	// 	err := txn.SetEntry(e)
	// 	return err
	// })
	// if dberr != nil {
	// 	fmt.Println(dberr)
	// }

	// db, err := badger.Open(badger.DefaultOptions("tmp/badger").WithTruncate(true))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// showerr := db.View(func(txn *badger.Txn) error {
	// 	opts := badger.DefaultIteratorOptions
	// 	opts.PrefetchSize = 10
	// 	it := txn.NewIterator(opts)
	// 	defer it.Close()
	// 	dataSize := 0
	// 	for it.Rewind(); it.Valid(); it.Next() {
	// 		dataSize++
	// 		item := it.Item()
	// 		k := item.Key()
	// 		err := item.Value(func(v []byte) error {
	// 			fmt.Printf("%v:key=%s \n", dataSize, k)
	// 			return nil
	// 		})
	// 		if err != nil {
	// 			return err
	// 		}
	// 	}
	// 	return nil
	// })
	// if showerr != nil {
	// 	fmt.Println(showerr)
	// }
	// db.Close()

}
