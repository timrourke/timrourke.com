package storage

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/timrourke/timrourke.com/model"
	"github.com/timrourke/timrourke.com/query"
	"strconv"
)

// NewPostStorage returns a new instance of PostStorage
func NewPostStorage(DB *sqlx.DB) *PostStorage {
	return &PostStorage{DB}
}

// PostStorage forms SQL queries for posts
type PostStorage struct {
	DB *sqlx.DB
}

// GetAll selects a list of posts
func (s *PostStorage) GetAll(q *query.Query) (uint, []model.Post, error) {
	var (
		posts []model.Post
		count uint
	)

	q.Select("posts.*").From("posts posts")

	sql, boundValues := q.Compile()

	rows, err := s.DB.NamedQuery(sql, boundValues)
	if err != nil {
		return 0, nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var m model.Post
		err = rows.StructScan(&m)
		if err != nil {
			return 0, nil, err
		}

		posts = append(posts, m)
	}

	// Get count of all posts for pagination
	errCount := s.DB.Get(&count, "SELECT COUNT(*) FROM posts")

	if err != nil {
		return 0, nil, err
	} else if errCount != nil {
		return 0, nil, errCount
	}

	return count, posts, nil
}

// GetOne selects a single post
func (s *PostStorage) GetOne(ID string) (*model.Post, error) {
	var post model.Post

	err := s.DB.Get(&post, "SELECT * FROM posts WHERE id=?", ID)

	return &post, err
}

// Insert inserts a single post
func (s *PostStorage) Insert(c model.Post) (*model.Post, error) {
	result, err := s.DB.NamedExec(`INSERT INTO posts (
		title,
		excerpt,
		content,
		permalink,
		user_id
	) VALUES (
		:title,
		:excerpt,
		:content,
		:permalink,
		14
	)`, &c)

	if err != nil {
		fmt.Println("insert error", err)
		return &model.Post{}, err
	}

	insertID, err := result.LastInsertId()
	if err != nil {
		fmt.Println("insert error last insert id", err)
		return &model.Post{}, err
	}

	// Set ID on return struct for rendering to json
	c.SetID(fmt.Sprintf("%d", insertID))

	return s.GetOne(c.GetID())
}

// Delete deletes a single post
func (s *PostStorage) Delete(id string) error {
	_, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return fmt.Errorf("Post id must be integer: %s", id)
	}

	_, err = s.DB.Exec("DELETE FROM posts WHERE id=? LIMIT 1", id)
	if err != nil {
		return err
	}
	return nil
}

// Update updates a single post
func (s *PostStorage) Update(c *model.Post) error {
	_, err := s.DB.NamedExec(`UPDATE posts SET 
		title=:title,
		excerpt=:excerpt,
		content=:content,
		permalink=:permalink,
		WHERE id=:id`, &c)

	if err != nil {
		return err
	}

	return nil
}
