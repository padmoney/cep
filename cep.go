package cep

import (
	"database/sql"
	"errors"
	"fmt"
	"regexp"
)

type CEP struct {
	CEP                   string `json:"cep"`
	UF                    string `json:"uf"`
	Localidade            string `json:"localidade"`
	LocalidadeAbrev       string `json:"localidade_abrev"`
	LocalidadeTipo        string `json:"localidade_tipo"`
	Logradouro            string `json:"logradouro"`
	LogradouroAbrev       string `json:"logradouro_abrev"`
	LogradouroTipo        string `json:"localidade_tipo"`
	LogradouroComplemento string `json:"logradouro_complemento"`
	Bairro                string `json:"bairro"`
	BairroAbrev           string `json:"bairro_abrev"`
	CodigoMunicipio       string `json:"codigo_municipio"`
}

type ConnParams struct {
	Driver   string
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	Schema   string
	SSLMode  bool
}

func (c ConnParams) SSLModeAsString() string {
	if c.SSLMode {
		return "enable"
	}
	return "disable"
}

func Search(cep string, connParams ConnParams) (CEP, error) {
	cep = sanitize(cep)
	if !isValid(cep) {
		var entity CEP
		return entity, errors.New("CEP inválido.")
	}

	return searchPostgres(cep, connParams)
}

func isValid(cep string) bool {
	return len(cep) == 8
}

func sanitize(cep string) string {
	reg, err := regexp.Compile("[^0-9]+")
	if err != nil {
		return ""
	}
	return reg.ReplaceAllString(cep, "")
}

func searchPostgres(cep string, connParams ConnParams) (entity CEP, err error) {
	info := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		connParams.Host,
		connParams.Port,
		connParams.User,
		connParams.Password,
		connParams.DBName,
		connParams.SSLModeAsString())
	db, err := sql.Open("postgres", info)
	if err != nil {
		return entity, nil
	}
	defer db.Close()

	// set search_path
	if connParams.Schema != "" && connParams.Schema != "public" {
		sql := fmt.Sprintf("SET search_path TO %s;", connParams.Schema)
		_, err = db.Exec(sql)
		if err != nil {
			return entity, err
		}
	}

	rows, err := db.Query("select * from search_cep($1);", cep)
	if err != nil {
		fmt.Println(2)
		return entity, err
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&entity.CEP,
			&entity.UF,
			&entity.LogradouroTipo,
			&entity.LogradouroAbrev,
			&entity.Logradouro,
			&entity.LogradouroComplemento,
			&entity.Localidade,
			&entity.LocalidadeAbrev,
			&entity.LocalidadeTipo,
			&entity.CodigoMunicipio,
			&entity.Bairro,
			&entity.BairroAbrev,
		)
		if err != nil {
			return entity, err
		}
	}
	return entity, rows.Err()
}
