# Aylak
Aylak is a web page scraper application that has javascript rendering feature

# Dependencies
```
Chromedp need Google Chrome so make sure it's installed.
github.com/chromedp/chromedp

github.com/PuerkitoBio/goquery
```

# JSON File
Aylak need a JSON file which has datas of url and selectors. The JSON must contain these values for every website.

``` 
[
    {
        "url": "https://kodilan.com/ilanlar/sayfa/",
        "footerSelector": "footer",
        "count": 10,
        "wrapperSelector": "#page > div.container.job-listing > div.eleven.columns > div > div.listings-container > div",
        "itemSelector": [
            "span.title.tag-post-link"
        ]
    }
]
```

|Url  | The web page's url without page number |
| :------------- | :------------- |
| **FooterSelector** |**Element's tag or class at the bottom of the page** |
| **Count** | **Number of sub pages** |
| **WrapperSelector** | **The biggest required parent element** |
| **ItemSelector** | **The array of elements of desired datas** |
