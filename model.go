package qishutaLib

type BookInfoModel struct {
	BookId      string `json:"book_id"`
	BookName    string `json:"book_name"`
	Author      string `json:"author"`
	Cover       string `json:"cover"`
	UpdateDate  string `json:"update_date"`
	Status      string `json:"status"`
	Download    string `json:"download"`
	ClickInfo   string `json:"click_info"`
	FileSize    string `json:"file_size"`
	Description string `json:"description"`
}

type TypeListBookInfoModel struct {
	Index    string `json:"index"`
	BookId   string `json:"book_id"`
	BookName string `json:"book_name"`
	Cover    string `json:"cover"`

	Description string `json:"description"`
}

type CatalogModel struct {
	BookID       string `json:"book_id"`
	ChapterIndex string `json:"chapter_index"`
	ChapterTitle string `json:"chapter_title"`
	ChapterId    string `json:"chapter_id"`
}

type ContentModel struct {
	ID           string `json:"_id" bson:"_id"`
	BookId       string `json:"book_id" bson:"book_id"`
	ChapterId    string `json:"chapter_id" bson:"chapter_id"`
	ChapterTitle string `json:"chapter_title" bson:"chapter_title"`
	Content      string `json:"content" bson:"content"`
	ChapterWord  int    `json:"chapter_word" bson:"chapter_word"`
}

type SearchModel struct {
	Index         string `json:"index"`
	BookId        string `json:"book_id"`
	BookName      string `json:"book_name"`
	BookAuthor    string `json:"book_author"`
	LatestChapter string `json:"latest_chapter"`
	Update        string `json:"update"`
}
type BookshelfModel struct {
	BookId        string `json:"book_id"`
	BookName      string `json:"book_name"`
	Index         string `json:"index"`
	LatestChapter string `json:"latest_chapter"`
	UpdateDate    string `json:"update_date"`
}
