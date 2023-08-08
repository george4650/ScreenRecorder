package repository

import (
	"context"
	"fmt"
	"myapp/internal/models"
	"myapp/internal/usecase"

	go_oracle "gitlabnew.nextcontact.ru/r.alfimov/go-oracle"
)

type ScreenRecorderOracle struct {
	ora *go_oracle.Oracle
}

func NewScreenRecorderOracle(ora *go_oracle.Oracle) usecase.ScreenRecorderOracle {
	return &ScreenRecorderOracle{
		ora: ora,
	}
}

func (db *ScreenRecorderOracle) ListProjectsIdName(ctx context.Context) ([]models.ProjectIdName, error) {

	sqlText := `select ID,
					TITLE
					from projects`

	projects, err := go_oracle.SelectMany[models.ProjectIdName](ctx, db.ora, sqlText, []interface{}{})
	if err != nil {
		return nil, fmt.Errorf("ScreenRecorderOracle - ListProjectsIdName - go_oracle.SelectMany: %w", err)
	}

	return projects, nil
}
