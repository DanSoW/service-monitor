package model

/* Структура основной информации в таблице */
type HeaderInfoModel struct {
	IPv4        string   `json:"ipv4"`
	LoginSSH    []string `json:"property_item"`
	PasswordSSH []string `json:"communicate_variant"`
}

/* Структура идентификатора ячейки таблицы */
type IndexCellModel struct {
	Pos    string `json:"pos"`
	Row    int    `json:"row"`
	Column int    `json:"column"`
	Value  string `json:"string"`
}
