package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type PreUrl struct {
	Url    string `json:"datos"`
	Client *http.Client
}

type AutoGenerated []struct {
	Origen     Origen     `json:"origen"`
	Elaborado  string     `json:"elaborado"`
	Nombre     string     `json:"nombre"`
	Provincia  string     `json:"provincia"`
	Prediccion Prediccion `json:"prediccion"`
	ID         int        `json:"id"`
	Version    float64    `json:"version"`
}
type Origen struct {
	Productor string `json:"productor"`
	Web       string `json:"web"`
	Enlace    string `json:"enlace"`
	Language  string `json:"language"`
	Copyright string `json:"copyright"`
	NotaLegal string `json:"notaLegal"`
}
type ProbPrecipitacion struct {
	Value   int    `json:"value"`
	Periodo string `json:"periodo"`
}
type CotaNieveProv struct {
	Value   string `json:"value"`
	Periodo string `json:"periodo"`
}
type EstadoCielo struct {
	Value       string `json:"value"`
	Periodo     string `json:"periodo"`
	Descripcion string `json:"descripcion"`
}
type Viento struct {
	Direccion string `json:"direccion"`
	Velocidad int    `json:"velocidad"`
	Periodo   string `json:"periodo"`
}
type RachaMax struct {
	Value   string `json:"value"`
	Periodo string `json:"periodo"`
}
type Dato struct {
	Value int `json:"value"`
	Hora  int `json:"hora"`
}
type Temperatura struct {
	Maxima int    `json:"maxima"`
	Minima int    `json:"minima"`
	Dato   []Dato `json:"dato"`
}
type SensTermica struct {
	Maxima int    `json:"maxima"`
	Minima int    `json:"minima"`
	Dato   []Dato `json:"dato"`
}
type HumedadRelativa struct {
	Maxima int    `json:"maxima"`
	Minima int    `json:"minima"`
	Dato   []Dato `json:"dato"`
}
type Dia struct {
	ProbPrecipitacion []ProbPrecipitacion `json:"probPrecipitacion"`
	CotaNieveProv     []CotaNieveProv     `json:"cotaNieveProv"`
	EstadoCielo       []EstadoCielo       `json:"estadoCielo"`
	Viento            []Viento            `json:"viento"`
	RachaMax          []RachaMax          `json:"rachaMax"`
	Temperatura       Temperatura         `json:"temperatura"`
	SensTermica       SensTermica         `json:"sensTermica"`
	HumedadRelativa   HumedadRelativa     `json:"humedadRelativa"`
	UvMax             int                 `json:"uvMax,omitempty"`
	Fecha             string              `json:"fecha"`
}
type Prediccion struct {
	Dia []Dia `json:"dia"`
}

func main() {
	//Crida a la funció per invocar la preurl
	result, _ := GetPreUrl()
	//fmt.Println(result)
	valors, _ := GetPrediccio(result)
	for valor := range valors {
		fmt.Println("Valors: ", valors[valor])
	}
}

func GetPreUrl() (string, error) {
	url := "https://opendata.aemet.es/opendata/api/prediccion/especifica/municipio/diaria/08001/?api_key=eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJvZGlnaW9jaW9AZ21haWwuY29tIiwianRpIjoiMzgyYzVjMWUtODcwYy00MmFjLWIwMjUtMzhlODk1OWRjNWVjIiwiaXNzIjoiQUVNRVQiLCJpYXQiOjE2ODM1NjE5MTcsInVzZXJJZCI6IjM4MmM1YzFlLTg3MGMtNDJhYy1iMDI1LTM4ZTg5NTlkYzVlYyIsInJvbGUiOiIifQ.Onfu4qinnfNfHVkBDr8s99pJpFoNRTv7Zed-TqwzAXY"

	//Preparem la peticio
	req, _ := http.NewRequest("GET", url, nil)

	//Afegirem una capçalera
	req.Header.Add("cache-control", "no-cache")
	//Executar la petició
	res, err := http.DefaultClient.Do(req)

	//Controlem errors
	if err != nil {
		log.Println("error conectant amb aemet.es", err)
	}

	defer res.Body.Close() //Diferim el tancament

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("error llegint el json", err)
	}

	preUrl := PreUrl{}
	err = json.Unmarshal(body, &preUrl)
	if err != nil {
		log.Println("error decodificant el json", err)
		return "", err
	}

	//surt be
	return preUrl.Url, err
}

func GetPrediccio(urlRebuda string) ([]int, error) {
	//Definir variable
	url := urlRebuda
	//Preparem la petició
	req, _ := http.NewRequest("GET", url, nil)

	//Afegirem una capçalera
	req.Header.Add("cache-control", "no-cache")
	//Executar la petició
	res, err := http.DefaultClient.Do(req)

	//Controlem errors
	if err != nil {
		log.Println("error conectant amb aemet.es", err)
	}

	defer res.Body.Close() //Diferim el tancament
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("error llegint el json", err)
	}

	prediccio := AutoGenerated{}
	var precipitacio, tempMax, tempMin, humitat int
	err = json.Unmarshal([]byte(body), &prediccio)
	if err != nil {
		log.Println("error decodificant el json", err)
	}
	precipitacio, tempMax, tempMin, humitat = prediccio[0].Prediccion.Dia[0].ProbPrecipitacion[0].Value, prediccio[0].Prediccion.Dia[0].Temperatura.Maxima, prediccio[0].Prediccion.Dia[0].Temperatura.Minima, prediccio[0].Prediccion.Dia[0].HumedadRelativa.Maxima
	valors := []int{precipitacio, tempMax, tempMin, humitat}
	return valors, err
}
