package gorm

import (
	"arcanaeum/logger"
)

// InitData : initialize the db with sample data
func InitData() {
	CreateAppSettings()

	author1 := Author{
		Name: "Haruki Murakami",
		Bio:  "Haruki Murakami (Japanese: 村上 春樹) is a popular contemporary Japanese writer and translator. His work has been described as 'easily accessible, yet profoundly complex'.",
	}
	err := author1.Create()
	if err != nil {
		logger.Error(err)
		return
	}

	author2 := Author{
		Name: "Brian W. Kernighan",
		Bio:  "Brian Kernighan (born January 1, 1942) is a computer scientist who worked at the Bell Labs and contributed to the design of the pioneering AWK and AMPL programming languages. He is most well-known for his co-authorship, with Dennis Ritchie, of the first book on the C programming language.",
	}
	err = author2.Create()
	if err != nil {
		logger.Error(err)
		return
	}

	author3 := Author{
		Name: "Dennis M. Ritchie",
		Bio:  "Dennis Ritchie (September 9, 1941 – c. October 12, 2011) was an American computer scientist and winner, with Kenneth Thompson, of the 1983 Turing Award. He created the C programming language and, with Thompson, the Unix operating system, which have had pervasive and lasting influence on subsequent programming languages and operating systems.",
	}
	err = author3.Create()
	if err != nil {
		logger.Error(err)
		return
	}

	author4 := Author{
		Name: "吉本ばなな",
		Bio:  "吉本 ばなな（本名：吉本 真秀子〈よしもと まほこ〉、旧筆名：よしもと ばなな、は、日本の小説家。事実婚の相手はロルファーの田畑浩良。",
	}
	err = author4.Create()
	if err != nil {
		logger.Error(err)
		return
	}

	author5 := Author{
		Name: "Phan Việt",
		Bio:  "Phan Việt tên thật là Nguyễn Ngọc Hường, sinh năm 1978, tốt nghiệp Đại học Ngoại thương năm 2000, sau đó chị sang Mỹ học Cao học Truyền thông, từ năm 2002 đến nay, chị là nghiên cứu sinh chương trình tiến sĩ về công tác xã hội tại ĐH Chicago, Mỹ. Trong lớp trẻ tuổi viết văn, nếu Nguyễn Ngọc Tư được coi là “hiện tượng” trong suốt những năm 2000, thì Phan Việt được đánh giá là một trong những “đài khí tượng” có khả năng “tiên báo một chiều kích mới cho văn học Việt Nam hiện đại”.",
	}
	err = author5.Create()
	if err != nil {
		logger.Error(err)
		return
	}

	tag1 := Tag{
		Name:        "Japanese Literature",
		Description: "Books, novels, poetry, ... written by Japanese authors",
	}
	err = tag1.Create()
	if err != nil {
		logger.Error(err)
		return
	}

	tag2 := Tag{
		Name:        "Vietnamese Literature",
		Description: "Books, novels, poetry, ... written by Vietnamese authors",
	}
	err = tag2.Create()
	if err != nil {
		logger.Error(err)
		return
	}

	tag3 := Tag{
		Name:        "Nonfiction",
		Description: `Nonfiction means "not a fiction"`,
	}
	err = tag3.Create()
	if err != nil {
		logger.Error(err)
		return
	}

	tag4 := Tag{
		Name:        "Technology",
		Description: "",
	}
	err = tag4.Create()
	if err != nil {
		logger.Error(err)
		return
	}

	language1 := Language{
		Name: "English",
	}
	err = language1.Create()
	if err != nil {
		logger.Error(err)
		return
	}

	language2 := Language{
		Name: "Vietnamese",
	}
	err = language2.Create()
	if err != nil {
		logger.Error(err)
		return
	}

	language3 := Language{
		Name: "Japanese",
	}
	err = language3.Create()
	if err != nil {
		logger.Error(err)
		return
	}

	book1 := Book{
		Title:       "Kafka On The Shore",
		Description: `Kafka on the Shore displays one of the world’s great storytellers at the peak of his powers.Here we meet a teenage boy, Kafka Tamura, who is on the run, and Nakata, an aging simpleton who is drawn to Kafka for reasons that he cannot fathom. As their paths converge, acclaimed author Haruki Murakami enfolds readers in a world where cats talk, fish fall from the sky, and spirits slip out of their bodies to make love or commit murder, in what is a truly remarkable journey.`,
		ISBN10:      "1400079276",
		ISBN13:      "978-1400079278",
		PubMonth:    1,
		PubYear:     2005,
		Shelf:       "A-1",
	}
	err = book1.Create()
	if err != nil {
		logger.Error(err)
		return
	}
	err = book1.AssignLanguage(language1.ID)
	if err != nil {
		logger.Error(err)
		return
	}
	err = book1.AddAuthor(author1.ID)
	if err != nil {
		logger.Error(err)
		return
	}
	err = book1.AddTag(tag1.ID)
	if err != nil {
		logger.Error(err)
		return
	}

	book2 := Book{
		Title:       "The C Programming Language",
		Description: `The C Programming Language (sometimes termed K&R, after its authors' initials) is a computer programming book written by Brian Kernighan and Dennis Ritchie, the latter of whom originally designed and implemented the language, as well as co-designed the Unix operating system with which development of the language was closely intertwined. The book was central to the development and popularization of the C programming language and is still widely read and used today. Because the book was co-authored by the original language designer, and because the first edition of the book served for many years as the de facto standard for the language, the book was regarded by many to be the authoritative reference on C.`,
		ISBN10:      "0131103628",
		ISBN13:      "978-0131103627",
		// No PubMonth intentional
		PubYear: 1988,
		Shelf:   "B-8",
	}
	err = book2.Create()
	if err != nil {
		logger.Error(err)
		return
	}
	err = book2.AssignLanguage(language1.ID)
	if err != nil {
		logger.Error(err)
		return
	}
	err = book2.AddAuthor(author2.ID)
	if err != nil {
		logger.Error(err)
		return
	}
	err = book2.AddAuthor(author3.ID)
	if err != nil {
		logger.Error(err)
		return
	}
	err = book2.AddTag(tag3.ID)
	if err != nil {
		logger.Error(err)
		return
	}
	err = book2.AddTag(tag4.ID)
	if err != nil {
		logger.Error(err)
		return
	}

	book3 := Book{
		Title: "キッチン",
		Description: `ここまで自分の体のような、心のような、手にして開けると懐かしさを感じる本はないと思うほど静かにわたしの中に沈んでいきます。もう知っている物語が愛おしくて大事で、愛犬を愛でるような感情に似ています。可愛らしい内容ではないんですが、ばなな先生がわたしが生まれる前に発表して新人賞などを受賞したこの作品が、旅行から自宅に帰ってきたときのような安心感を与えてくれます。昔のアルバムを開く感覚に近いのかな？他のばなな先生の作品も好きなんですが、強烈なインパクトも大恋愛もしないこの作品が病みつきになって、もう何回読み返したかわかりません。はっきりした理由はありません。でも、わたしはこの作品を一生愛して、自分の子供や将来の孫に与えたいと思っています。今ではなくてはならない本の一冊となっています

(exceprt from https://reviewne.jp/reviews/22607 for placeholder purpose)`,
		ISBN10:   "0802142443",
		ISBN13:   "978-0802142443",
		PubMonth: 4,
		PubYear:  2006,
		Shelf:    "C-3",
	}
	err = book3.Create()
	if err != nil {
		logger.Error(err)
		return
	}
	err = book3.AssignLanguage(language3.ID)
	if err != nil {
		logger.Error(err)
		return
	}
	err = book3.AddAuthor(author4.ID)
	if err != nil {
		logger.Error(err)
		return
	}
	err = book3.AddTag(tag1.ID)
	if err != nil {
		logger.Error(err)
		return
	}

	book4 := Book{
		Title: "Nước Mỹ nước Mỹ",
		Description: `Phải mất tám năm sau ngày cuốn sách ra mắt, tôi mới đọc Nước Mỹ, Nước Mỹ của Phan Việt. Trong lần in lại này, tập truyện được bổ sung thêm hai truyện ngắn ở gần cuối, với không gian tại Hà Nội.

Có khá nhiều lý do cho việc vì sao tôi không đọc cuốn sách này sớm hơn. Thứ nhất tôi không thích tiêu đề của nó. Thứ hai tôi không thích tiêu đề của một cuốn sách khác cùng tác giả – Một mình ở Châu Âu và cuối cùng, tôi không có hứng thú với những tác phẩm văn học của nhà văn nữ Việt Nam trong giai đoạn hiện đại, trừ truyện ngắn của Nguyễn Ngọc Tư và tản văn của Phan Thị Vàng Anh. Nhưng việc đọc sách cũng có những cái gọi là duyên, sau nhiều lần nhìn thấy Nước Mỹ, Nước Mỹ trên kệ sách của Phương Nam, Fahasa và bỏ qua không thương tiếc. Sau nhiều quảng cáo và danh sách Top tác phẩm văn học Việt Nam hiện đại từ Tiki hay Vinabook, và cũng bị bỏ qua không thương tiếc; trong một ngày làm việc của tháng 06, khi tìm mua cho mình một vài cuốn sách trên Tiki, tôi tặc lưỡi bỏ Nước Mỹ, Nước Mỹ vào giỏ hàng vì không còn biết phải mua gì. Và may mắn thay, lựa chọn này không làm tôi thất vọng.

(excerpt from http://duynt.com/2016/07/10/review-nuoc-my-nuoc-my/ for placeholder purpose)`,
		PubMonth: 2,
		PubYear:  2013,
		Shelf:    "C-3",
	}
	err = book4.Create()
	if err != nil {
		logger.Error(err)
		return
	}
	err = book4.AssignLanguage(language2.ID)
	if err != nil {
		logger.Error(err)
		return
	}
	err = book4.AddAuthor(author5.ID)
	if err != nil {
		logger.Error(err)
		return
	}
	err = book4.AddTag(tag2.ID)
	if err != nil {
		logger.Error(err)
		return
	}

	category1 := Category{
		Name:        "News",
		Description: "General news of the library",
	}
	err = category1.Create()
	if err != nil {
		logger.Error(err)
		return
	}

	category2 := Category{
		Name:        "Book Reviews",
		Description: "All the book reviews",
	}
	err = category2.Create()
	if err != nil {
		logger.Error(err)
		return
	}

	return
}
