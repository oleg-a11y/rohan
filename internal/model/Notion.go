package model

type NotionResponse struct {
	Results []Page `json:"results"`
}

type Page struct {
	Properties Properties `json:"properties"`
}

type Properties struct {
	Company TitleProperty  `json:"Company"`
	Date    DateProperty   `json:"Date"`
	Stage   SelectProperty `json:"Stage"`
	Salary  TextProperty   `json:"Salary"`
	Creator TextProperty   `json:"Creator"`
}

type TitleProperty struct {
	Title []Text `json:"title"`
}

type Text struct {
	PlainText string `json:"plain_text"`
}

type DateProperty struct {
	Date DateDetails `json:"date"`
}

type DateDetails struct {
	Start string `json:"start"`
}

type SelectProperty struct {
	Select SelectDetails `json:"select"`
}

type SelectDetails struct {
	Name string `json:"name"`
}

type TextProperty struct {
	RichText []Text `json:"rich_text"`
}