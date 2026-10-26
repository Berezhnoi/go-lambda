package main

import "testing"

// Очікувані значення пораховані вручну:
//
//	180 см = 70.866 in → понад 60: 10.866
//	  Devine 50+2.3×10.866=74.99, Robinson 52+1.9×10.866=72.65,
//	  Miller 56.2+1.41×10.866=71.52, Hamwi 48+2.7×10.866=77.34
//	  середнє 296.50/4=74.12; діапазон 3.24 м² × 18.5/24.9 = 59.94/80.68
//	165 см = 64.961 in → понад 60: 4.961
//	  Devine 45.5+2.3×4.961=56.91, Robinson 49+1.7×4.961=57.43,
//	  Miller 53.1+1.36×4.961=59.85, Hamwi 45.5+2.2×4.961=56.41
//	  середнє 230.61/4=57.65; діапазон 2.7225 м² × 18.5/24.9 = 50.37/67.79
//	152.4 см = рівно 60 in → лише базові коефіцієнти; діапазон 2.3226 м² → 42.97/57.83
func TestCalculate(t *testing.T) {
	tests := []struct {
		name     string
		heightCm float64
		sex      string
		want     Result
	}{
		{"male 180", 180, "male", Result{
			HealthyRange{59.9, 80.7}, Estimates{75.0, 72.6, 71.5, 77.3}, 74.1, medicalNote}},
		{"female 165", 165, "female", Result{
			HealthyRange{50.4, 67.8}, Estimates{56.9, 57.4, 59.8, 56.4}, 57.7, medicalNote}},
		{"male 5 ft", 152.4, "male", Result{
			HealthyRange{43.0, 57.8}, Estimates{50, 52, 56.2, 48}, 51.6, medicalNote}}, // 206.2/4=51.55
		{"female 5 ft", 152.4, "female", Result{
			HealthyRange{43.0, 57.8}, Estimates{45.5, 49, 53.1, 45.5}, 48.3, medicalNote}}, // 193.1/4=48.275
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Calculate(tt.heightCm, tt.sex); got != tt.want {
				t.Errorf("got  %+v\nwant %+v", got, tt.want)
			}
		})
	}
}
