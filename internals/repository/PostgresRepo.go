package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/ZakirAvrora/TechHome/internals/models"
	"github.com/jmoiron/sqlx"
	"log"
	"time"
)

var ErrNoRowAffected = errors.New("bad request, no affect in data")

type PostgresRepo struct {
	db *sqlx.DB
}

func NewPostgresRepo(db *sqlx.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

func (p *PostgresRepo) ListLinks(page int) ([]models.Link, error) {
	var links []models.Link
	query := `SELECT * FROM redirects ORDER BY redirect_id
            	LIMIT 10 OFFSET $1`

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := p.db.SelectContext(ctx, &links, query, page*10); err != nil {
		return nil, err
	}

	if len(links) == 0 {
		return nil, sql.ErrNoRows
	}

	log.Println("LIST:", links, len(links))
	return links, nil
}

func (p *PostgresRepo) GetLink(id int) (models.Link, error) {
	var link models.Link
	query := `SELECT * FROM redirects r where r.redirect_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := p.db.GetContext(ctx, &link, query, id); err != nil {
		return models.Link{}, err
	}

	return link, nil
}

func (p *PostgresRepo) CreateLink(link models.Link) (int64, error) {
	query := `INSERT INTO redirects (active_link, history_link)
			  VALUES (:active_link, :history_link) RETURNING redirect_id`

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	res, err := p.db.NamedExecContext(ctx, query, link)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (p *PostgresRepo) UpdateLink(link models.Link, id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	query := `UPDATE redirects SET active_link = $2, history_link= $3 WHERE redirect_id=$1`

	result, err := p.db.ExecContext(ctx, query, id, link.ActiveLink, link.HistoryLink)

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNoRowAffected
	}

	return nil
}

func (p *PostgresRepo) DeleteLink(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	query := `DELETE FROM redirects WHERE redirect_id = $1`
	result, err := p.db.ExecContext(ctx, query, id)

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNoRowAffected
	}

	return nil
}

func (p *PostgresRepo) FindLink(link string) (models.Link, error) {
	var res models.Link
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	query := `SELECT * FROM redirects WHERE active_link = $1 OR history_link = $1 LIMIT 1`
	if err := p.db.GetContext(ctx, &res, query, link); err != nil {
		return models.Link{}, err
	}

	return res, nil
}
