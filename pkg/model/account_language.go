package model

type Language struct {
	AccountID int64  `json:"-"`
	Language  string `json:"language"`
	Level     string `json:"level"`
}

func (l *Language) ToResponse() IResponse {
	l.AccountID = 0
	return l
}

type LanguageJoin struct {
	AccountID *int64  `json:"account_id"`
	Language  *string `json:"language"`
	Level     *string `json:"level"`
}

func (l LanguageJoin) ConvertToLanguage() Language {
	return Language{
		AccountID: *l.AccountID,
		Language:  *l.Language,
		Level:     *l.Level,
	}
}

type AddLanguage struct {
	Language string `json:"language" binding:"required,max=2"`
	Level    string `json:"level" binding:"required,oneof=a1 a2 b1 b2 c1 c2"`
}

type UpdateLanguage struct {
	Level string `json:"level" binding:"required,oneof=a1 a2 b1 b2 c1 c2"`
}

type ListLanguages struct {
	Languages []*Language `json:"languages"`
}

func (l *ListLanguages) ToResponse() IResponse {
	for _, v := range l.Languages {
		v.ToResponse()
	}
	return l
}

type UpdateLanguageFields map[string]interface{}

func (f UpdateLanguageFields) DBColumns() map[string]interface{} {
	return map[string]interface{}{
		//"updated_at":              struct{}{},
		"level": struct{}{},
	}
}

func (f UpdateLanguageFields) Prepare() {
	for k, v := range f {
		columns := f.DBColumns()
		if _, ok := columns[k]; !ok {
			delete(f, k)
			continue
		}

		switch k {
		default:
			f[k] = v
		}
	}
}
