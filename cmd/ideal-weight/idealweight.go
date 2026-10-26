package main

import "math"

const medicalNote = "Not medical advice. These are population-level estimates; consult a healthcare professional."

type HealthyRange struct {
	MinKg float64 `json:"min_kg"`
	MaxKg float64 `json:"max_kg"`
}

type Estimates struct {
	Devine   float64 `json:"devine"`
	Robinson float64 `json:"robinson"`
	Miller   float64 `json:"miller"`
	Hamwi    float64 `json:"hamwi"`
}

type Result struct {
	HealthyRange HealthyRange `json:"healthy_range_kg"`
	Estimates    Estimates    `json:"estimates_kg"`
	AverageKg    float64      `json:"average_kg"`
	Note         string       `json:"note"`
}

// Calculate — чиста функція: sex має бути вже провалідований ("male" або "female").
func Calculate(heightCm float64, sex string) Result {
	in := heightCm/2.54 - 60 // формули рахують від 5 футів (60 дюймів)

	var devine, robinson, miller, hamwi float64
	if sex == "male" {
		devine = 50 + 2.3*in
		robinson = 52 + 1.9*in
		miller = 56.2 + 1.41*in
		hamwi = 48 + 2.7*in
	} else {
		devine = 45.5 + 2.3*in
		robinson = 49 + 1.7*in
		miller = 53.1 + 1.36*in
		hamwi = 45.5 + 2.2*in
	}

	m := heightCm / 100
	return Result{
		HealthyRange: HealthyRange{MinKg: round1(18.5 * m * m), MaxKg: round1(24.9 * m * m)},
		Estimates: Estimates{
			Devine:   round1(devine),
			Robinson: round1(robinson),
			Miller:   round1(miller),
			Hamwi:    round1(hamwi),
		},
		// Середнє з неокруглених значень, щоб не накопичувати похибку округлення.
		AverageKg: round1((devine + robinson + miller + hamwi) / 4),
		Note:      medicalNote,
	}
}

func round1(x float64) float64 { return math.Round(x*10) / 10 }
