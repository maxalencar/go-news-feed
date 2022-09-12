package model

import "time"

type Articles []Article

type Article struct {
	ID                string     `json:"id,omitempty" bson:"_id,omitempty"`
	Title             string     `json:"title,omitempty" bson:"title,omitempty"`
	Descriptiopn      string     `json:"description,omitempty" bson:"description,omitempty"`
	Link              string     `json:"link,omitempty" bson:"link,omitempty"`
	Source            Source     `json:"source,omitempty" bson:"source,omitempty"`
	PublishedDateTime *time.Time `json:"publishedDateTime,omitempty" bson:"publishedDateTime,omitempty"`
	UpdatedDateTime   *time.Time `json:"updatedDateTime,omitempty" bson:"updatedDateTime,omitempty"`
}

// Len returns the length of Items.
func (a Articles) Len() int {
	return len(a)
}

// Less compares PublishedDateTime of Articles[i], Articles[k]
// and returns true if Articles[i] is less than Articles[k].
func (a Articles) Less(i, k int) bool {
	return a[i].PublishedDateTime.Before(
		*a[k].PublishedDateTime,
	)
}

// Swap swaps Articles[i] and Articles[k].
func (a Articles) Swap(i, k int) {
	a[i], a[k] = a[k], a[i]
}
