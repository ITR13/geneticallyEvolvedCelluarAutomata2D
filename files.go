/*
    This file is part of InvertoTanks.

    Foobar is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    InvertoTanks is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with InvertoTanks.  If not, see <http://www.gnu.org/licenses/>.
*/
	
package genCell

import (
	"errors"
	"os"

	"github.com/MaxHalford/gago"
)

func Save(filename string, individual gago.Individual) error {
	if SimulationType == Count {
		return nil
	}
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	fitness := int64(individual.Fitness * 1000)
	fitnessB := make([]byte, 8)
	lines := len(individual.Genome)
	linesB := make([]byte, 8)
	for i := 0; i < 8; i++ {
		linesB[i] = byte(lines % 255)
		lines /= 255
		fitnessB[i] = byte(fitness % 255)
		fitness /= 255
	}
	_, err = file.Write(fitnessB)
	if err != nil {
		return err
	}
	_, err = file.Write(linesB)
	if err != nil {
		return err
	}
	for i := range individual.Genome {
		_, err = file.Write([]byte{byte(individual.Genome[i].(int))})
		if err != nil {
			return err
		}
	}
	return nil
}

func Read(filename string) (gago.Individual, error) {
	file, err := os.Open(filename)
	if err != nil {
		return gago.Individual{}, err
	}
	defer file.Close()

	fitnessB := make([]byte, 8)
	linesB := make([]byte, 8)
	n, err := file.Read(fitnessB)
	if err != nil {
		return gago.Individual{}, err
	} else if n != 8 {
		return gago.Individual{}, errors.New("Not enough bytes")
	}
	n, err = file.Read(linesB)
	if err != nil {
		return gago.Individual{}, err
	} else if n != 8 {
		return gago.Individual{}, errors.New("Not enough bytes")
	}
	fitness, lines := int64(0), 0
	for i := 0; i < 8; i++ {
		fitness *= 8
		fitness += int64(fitnessB[i])
		lines *= 8
		lines += int(linesB[i])
	}

	genome := make([]interface{}, lines)

	b := make([]byte, lines)
	n, err = file.Read(b)
	if err != nil {
		return gago.Individual{}, err
	} else if n != lines {
		return gago.Individual{}, errors.New("Not enough bytes")
	}
	for i := range genome {
		genome[i] = int(b[i])
	}
	return gago.Individual{genome, float64(fitness) / 1000, true, "-"}, nil
}
