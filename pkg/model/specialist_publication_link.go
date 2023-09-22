package model

type PublicationLink struct {
	ID        int64  `json:"id"`
	ProfileID int64  `json:"-"`
	Title     string `json:"title"`
	Link      string `json:"link"`
}

func (l *PublicationLink) ToResponse() IResponse {
	l.ProfileID = 0
	return l
}

type PublicationLinkJoin struct {
	ID        *int64  `json:"id"`
	ProfileID *int64  `json:"profile_id"`
	Title     *string `json:"title"`
	Link      *string `json:"link"`
}

func (e PublicationLinkJoin) ConvertToPublicationLink() PublicationLink {
	return PublicationLink{
		ID:        *e.ID,
		ProfileID: *e.ProfileID,
		Title:     *e.Title,
		Link:      *e.Link,
	}
}

type AddPublicationLink struct {
	Title string `json:"title" binding:"required,max=255"`
	Link  string `json:"link" binding:"required,http_url,max=255"`
}

type UpdatePublicationLink AddPublicationLink

type ListPublicationLinks struct {
	PublicationLinks []*PublicationLink `json:"publication_links"`
}

func (l *ListPublicationLinks) ToResponse() IResponse {
	for _, v := range l.PublicationLinks {
		v.ToResponse()
	}
	return l
}

type UpdatePublicationLinkFields map[string]interface{}

func (f UpdatePublicationLinkFields) DBColumns() map[string]interface{} {
	return map[string]interface{}{
		"title": struct{}{},
		"link":  struct{}{},
	}
}

func (f UpdatePublicationLinkFields) Prepare() {
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
