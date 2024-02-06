
# QiShuTa Library for Go

QiShuTa Library is an unofficial Go client for interacting with the QiShuTa website, which provides a large collection of electronic books that can be downloaded in TXT format for free.

## Features

- Search for books by keywords.
- Retrieve book details including name, author, cover, description, and download link.
- Fetch book covers as byte data.
- Get book catalogs and chapters content.
- Manage user bookshelf after authentication.
- Support for setting proxies, custom user agents, and debug mode.

## Installation

To install QiShuTa Library, use go get:

```bash
go get -u github.com/catnovelapi/qishutaLib
```

Replace `your_username` with your GitHub username if you have forked this project or the original author's username if you are using their package directly.

## Usage

Here is a quick example to get you started:

```go
package main

import (
    "fmt"
    "github.com/catnovelapi/qishutaLib"
)

func main() {
    client := qishutaLib.NewClient().SetBaseURL("https://www.qishuta.org")
    app := client.R()

    // Fetch details for a specific book by its ID
    bookInfo, err := app.GetBookInfo("12345")
    if err != nil {
        fmt.Println("Error fetching book info:", err)
        return
    }
    fmt.Println("Book Info:", bookInfo)

    // Retrieve the book cover image
    coverBytes, err := app.GetCover(bookInfo.Cover)
    if err != nil {
        fmt.Println("Error fetching book cover:", err)
        return
    }
    fmt.Println("Book Cover Byte Length:", len(coverBytes))

    // Search for books by keyword
    searchResults, err := app.GetSearch("keyword")
    if err != nil {
        fmt.Println("Error searching books:", err)
        return
    }
    fmt.Println("Search Results:", searchResults)

    // ... other functionalities
}
```

## Authentication

To access personalized features such as the bookshelf, you need to authenticate:

```go
cookie, err := app.GetCookie("username", "password")
if err != nil {
    fmt.Println("Authentication failed:", err)
    return
}
client.SetCookie(cookie)
```

## Contributing

Contributions are welcome! Feel free to open a pull request to fix a bug or add a new feature. Please follow the existing code style and add unit tests for new logic.

## License

This library is licensed under the MIT License. See the LICENSE file for details.

## Disclaimer

This is an unofficial client and is not endorsed by or affiliated with QiShuTa. Use it at your own risk.

---

For complete documentation on how to use the QiShuTa Library, refer to the GoDoc or check the `examples` folder in the repository.
 