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

type Review struct {
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
