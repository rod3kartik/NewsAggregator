package main

import ("fmt"
		"net/http"
		"io/ioutil"
		"encoding/xml"
		"strings"
		"html/template"
		"sync")

var wg sync.WaitGroup

func add(x,y float32) float32{
	return x+y
}

const converter float64 = 0.4
type book struct{
	ISBN float64
	copies uint16
	author string
	publisher string
	rack string
}

type SiteMapIndex struct{
	Leagues []string `xml:"sitemap>loc"`
}

type News struct{
	Titles []string `xml:"url>news>title"`
	Keywords []string `xml:"url>news>keywords"`
	Locations []string `xml:"url>loc"`
}

type NewsMap struct{
	Keyword string
	Location string
}

type NewsAggPage struct{
	Title string
	News map[string]NewsMap
}
// type League struct{
// 	Loc string `xml:"loc"`
// }

// func (l League) String() string{
// 	return fmt.Sprintf(l.Loc)
// }

func multiple(a,b string)(string,string){
	return a,b
}

func (b book) noGen() int{
	var y int = int((b.ISBN) * converter)
	return y 
}

func (b *book) new_rack(newrack string){
	b.rack = newrack
}

func index_handler(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Go is tight \n")
	a_book := book{ISBN : 1224546, copies: 10, author: "JK rowling", publisher: "Henry", rack: "A"}
	fmt.Fprintf(w, a_book.author)
	fmt.Fprintf(w, "\n Book serial number is %d", a_book.noGen())
	fmt.Fprintf(w, "\n Changing rack number from " + a_book.rack)
	a_book.new_rack("B")
	fmt.Fprintf(w, "\n new rack is " + a_book.rack)
}

func about_handler(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "<h1>about page under construction</h1>")
	fmt.Fprintf(w, "<p>GO is fast boi</p>")
	fmt.Fprintf(w, "<p>You %s even add %s</p>","can","<strong>variables</strong>")

}

func newsRoutine(c chan News,Location string){
	defer wg.Done()
	resp,_ := http.Get(Location)
	var n News
	bytes,_ := ioutil.ReadAll(resp.Body)
		//fmt.Println("here I am")
	xml.Unmarshal(bytes,&n)
	resp.Body.Close()
	c <- n
}
func newsAggHandler(w http.ResponseWriter, r *http.Request){
	var s SiteMapIndex
	//r.Header.Set("User-Agent", "Mozilla/5.0")
	resp,_ := http.Get("https://www.washingtonpost.com/news-sitemaps/index.xml")
	bytes,_ := ioutil.ReadAll(resp.Body)
	news_map := make(map[string]NewsMap)
	xml.Unmarshal(bytes,&s)
	resp.Body.Close()
	queue := make(chan News, 50)
	for _, Location := range s.Leagues {
		Location = strings.TrimSpace(Location)
		fmt.Println("%s",Location)
		wg.Add(1)
		go newsRoutine(queue, Location)	
	}
	wg.Wait()
	close(queue)
	for elem := range queue{
		for idx,_ := range elem.Keywords{
			news_map[elem.Titles[idx]] = NewsMap{elem.Keywords[idx],elem.Locations[idx]}
	}
	}	
	p:= NewsAggPage{Title:"NewsAggregator", News:news_map}
	t,_ := template.ParseFiles("basictemplating.html")
	fmt.Println(t.Execute(w,p))
}

func main(){
	//fmt.Println("A number from 1-100 is ", rand.Intn(100))
	//num1,num2 := 5.6,9.1
	//w1,w2 := "hey", "there"
	//fmt.Println(multiple(w1,w2))
	//fmt.Println(len(news_map))
	// for idx,data := range news_map {
	// 	fmt.Println("\n\n\n", idx)
	// 	fmt.Println("\n",data.Keyword)
	// 	fmt.Println("\n",data.Location)
	// }

	http.HandleFunc("/",index_handler)
	// http.HandleFunc("/about",about_handler)
	http.HandleFunc("/agg/",newsAggHandler)
	http.ListenAndServe(":8000",nil)
}