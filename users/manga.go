package manga

type Manga struct {
	ID          int64
	Title       string
	Description string
}

type Author struct {
	A_ID   int64
	A_name string
}

type User struct {
	User_id  int64
	Name     string
	Email    string
	Password string
}

type Article struct {
	ForumId int64
	Title   string
	Content string
}

type Review struct {
	ReviewId    int64
	Manga_id    int64
	Reviewer_id int64
	Title       string
	Description string
}

type Requests struct {
	Request_id  int64
	Reviewer_id int64
	title       string
}

type Post struct {
	ArticleId  int64
	ForumId    int64
	ReviewerId int64
}

type Forum struct {
	Forum_id    int64
	Title       string
	Description string
}

type Favorite struct {
	Title string
	Count int
}

type Genre struct {
	GId   int64
	GName string
}

type Comment struct {
	ForumId    int64
	ArticleId  int64
	ReviewerId int64
	Content    string
}
