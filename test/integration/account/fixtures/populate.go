package fixtures

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/Hvaekar/med-account/pkg/utils"
	"github.com/Masterminds/squirrel"
	"os"
	"strings"
)

func PopulateDB(c context.Context, db *sql.DB) error {
	files := []string{
		"accounts.json",
		"account_files.json",
		"account_emails.json",
		"account_phones.json",
		"account_addresses.json",
		"account_languages.json",
		"patient_profiles.json",
		"accounts_patient_profiles.json",
		"patient_disability_files.json",
		"patient_metal_components.json",
		"specialist_profiles.json",
		"specialist_specializations.json",
		"specialist_cures_diseases.json",
		"specialist_services.json",
		"specialist_educations.json",
		"specialist_education_files.json",
		"specialist_experiences.json",
		"specialist_experience_specializations.json",
		"specialist_associations.json",
		"specialist_patents.json",
		"specialist_publication_links.json",
	}

	for _, f := range files {
		content, err := os.ReadFile("./fixtures/testdata/" + f)
		if err != nil {
			return err
		}

		data := make([]map[string]any, 0)
		if err := json.Unmarshal(content, &data); err != nil {
			return err
		}

		if len(data) == 0 {
			continue
		}

		columns := make([]string, 0, len(data[0]))
		for k := range data[0] {
			columns = append(columns, k)
		}

		psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).RunWith(db)

		q := psql.Insert(strings.TrimSuffix(f, ".json")).
			Columns(columns...)

		for _, v := range data {
			values := make([]any, 0, len(v))
			for _, col := range columns {
				value := v[col]

				if col == "password" {
					value = utils.HashPassword(value.(string))
				}

				values = append(values, value)
			}

			q = q.Values(values...)
		}

		_, err = q.ExecContext(c)
		if err != nil {
			return err
		}
	}

	return nil
}
