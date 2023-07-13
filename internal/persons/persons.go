package persons

import (
	"check-id-api/internal/env"
	"check-id-api/internal/logger"
	"check-id-api/internal/ws"
	"encoding/json"
	"fmt"
)

type Persons struct {
	IdentityNumber string
}

func (p *Persons) GetPersonByIdentityNumber() (*Person, error) {

	resPerson := responsePerson{}
	e := env.NewConfiguration()
	resWs, code, err := ws.ConsumeWS(nil, e.App.UrlPersons+p.IdentityNumber, "GET", "", nil)
	if err != nil || code != 200 {
		logger.Error.Printf("No se pudo obtener la persona por el número de identificación: %v", err)
		return nil, err
	}

	err = json.Unmarshal(resWs, &resPerson)
	if err != nil {
		logger.Error.Printf("No se pudo parsear la respuesta: %v", err)
		return nil, err
	}

	if resPerson.Error {
		logger.Error.Printf("No se pudo obtener los datos de la persona, mensaje: ", resPerson.Msg)
		return nil, fmt.Errorf(resPerson.Msg)
	}

	if resPerson.Data == nil {
		logger.Error.Printf("No se encontro a la persona por su número de identificacion, mensaje: ", resPerson.Msg)
		return nil, fmt.Errorf("no se encontro a la persona por su número de identificacion")
	}

	return resPerson.Data, nil
}
