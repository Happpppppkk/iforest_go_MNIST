package main

import (
	"encoding/csv"
	"fmt"
	"image/png"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/e-XpertSolutions/go-iforest/iforest"
	"github.com/montanaflynn/stats"
	"github.com/petar/GoMNIST"
)

func main() {
	// Load the MNIST dataset for train and test
	trainD, testD, err := GoMNIST.Load("MNIST")
	if err != nil {
		panic(err)
	}

	// create image file and print out the first image of MNIST
	file, err := os.Create("image1.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Get the first image and label
	img, label := testD.Get(0)
	err = png.Encode(file, img) // Encode the image to the PNG format and write to the file
	if err != nil {
		panic(err)
	}
	fmt.Printf("image1.png has label %d\n", label)

	//make [][]float64
	trainData := make([][]float64, len(trainD.Images))
	testData := make([][]float64, len(testD.Images))

	// Flatten images
	for i, img := range trainD.Images {
		trainData[i] = make([]float64, len(img))
		for j, pixel := range img {
			trainData[i][j] = float64(pixel)
		}
	}

	for i, img := range testD.Images {
		testData[i] = make([]float64, len(img))
		for j, pixel := range img {
			testData[i][j] = float64(pixel)
		}
	}

	// Dataset information
	fmt.Println("Number of training images:", len(trainData))
	fmt.Println("Number of testing images:", len(testData))
	fmt.Println("First image data:", trainData[0])

	//create
	//model initialization
	forest := iforest.NewForest(100, 256, 0.02)

	//training stage - creating trees
	forest.Train(trainData)
	forest.Test(trainData)

	//iForestAnomalyScores := forest.AnomalyScores
	//threshold := forest.AnomalyBound
	_, goanomalyScores, err := forest.Predict(trainData)
	if err != nil {
		panic(err)
	}
	// labelsTest := forest.Labels
	//fmt.Println("anomalyScores: ", anomalyScores)

	csvFile, err := os.Create("goanomalyscores.csv")
	if err != nil {
		log.Fatalf("Failed to create CSV file: %v", err)
	}
	defer csvFile.Close()
	// Output the anomaly score
	// Create a CSV writer
	csvWriter := csv.NewWriter(csvFile)
	defer csvWriter.Flush()

	if err := csvWriter.Write([]string{"IDX", "iForest_AnomalyScore"}); err != nil {
		fmt.Println("Error writing headers", err)
		return
	}
	for IDX, score := range goanomalyScores {
		record := []string{fmt.Sprintf("%d", IDX), fmt.Sprintf("%f", score)}
		if err := csvWriter.Write(record); err != nil {
			log.Fatalf("Error writing record to CSV file: %v", err)
		}
	}
	//Correlation calcualtion for Python and Go anomaly score
	//read python csv for anomaly score
	pfile, err := os.Open("pythonScores.csv")
	if err != nil {
		fmt.Println("cannot open pythonScores.csvv")
	}
	defer pfile.Close()
	reader := csv.NewReader(pfile)
	// Reading the header
	header, err := reader.Read()
	if err != nil {
		fmt.Println("can't read header of pythonScores.csv")
	}
	var pythonScores []float64

	iforestPythonScoreIdx := -1
	for idx, column := range header {
		if column == "iforestPythonScore" {
			iforestPythonScoreIdx = idx
			break
		}
	}
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break // End of file is reached
		}
		score, _ := strconv.ParseFloat(record[iforestPythonScoreIdx], 64)
		pythonScores = append(pythonScores, score)
	}

	//fmt.Println(pythonScores)
	//correlation of python and go
	corr, err := stats.Correlation(pythonScores, goanomalyScores)
	if err != nil {
		fmt.Println("cannot create correlation of python and go")
	}

	fmt.Printf("The correlation coefficient between Go and Python anomaly scores is: %.2f\n", corr)

	//Correlation between R and Go
	//open file and creaee csv reader
	rfile, err := os.Open("isotreeRScores.csv")
	if err != nil {
		panic(err)
	}
	defer rfile.Close()

	reader1 := csv.NewReader(rfile)
	//skip first element
	if _, err := reader1.Read(); err != nil {
		fmt.Println("issues in r score reading")
	}
	//if _, err := reader1.Read(); err != nil {
	//	fmt.Println("issues in r score reading")
	//}

	var rScores []float64
	//record read
	for {
		record, err := reader1.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("cannot read R score record")
		}
		data, err := strconv.ParseFloat(record[0], 64)
		if err != nil {
			log.Fatalf("Unable to convert record to float: %v", err)
		}

		rScores = append(rScores, data)
	}
	//correlation of r and go
	corrR, err := stats.Correlation(rScores, goanomalyScores)
	if err != nil {
		fmt.Println("cannot create correlation of R and go")
	}
	fmt.Println(len(pythonScores), len(rScores), len(goanomalyScores))
	fmt.Printf("The correlation coefficient between Go and R anomaly scores is: %f\n", corrR)

}
