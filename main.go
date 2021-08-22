package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"humn/coordinate_mapping"
	"humn/workerpool"
	"log"
	"os"
	"strconv"
)

//var Writer *bufio.Writer

var allTask []*workerpool.Task

func main() {

	accessToken, poolSize := GetCommandLineArgs()

	outputChannel := make(chan coordinate_mapping.CoordinatePostcodeOutput)

	go func() {
		for o := range outputChannel {
			writeOutput(o)

		}
	}()

	pool := workerpool.NewPool(allTask, poolSize)

	go func() {
		coordinateMapper := coordinate_mapping.NewCoordinateMapper(accessToken, outputChannel)
		scanInput(coordinateMapper, pool)
	}()

	pool.RunBackground()

}

func scanInput(coordinateMapper coordinate_mapping.CoordinateMapper, pool *workerpool.Pool) {

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		var coordinate coordinate_mapping.Coordinate
		err := json.Unmarshal(scanner.Bytes(), &coordinate)
		if err != nil {
			log.Fatal(err)
		}

		task := workerpool.NewTask(func(data interface{}) error {
			err = coordinateMapper.GetPostcodeDataForCoordinatesAndWriteToOutput(coordinate)
			return err
		}, nil)
		pool.AddTask(task)
	}
	pool.Stop()
}

func writeOutput(output coordinate_mapping.CoordinatePostcodeOutput) {
	outputJsonBytes, err := json.Marshal(output)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(outputJsonBytes))
}

func GetCommandLineArgs() (string, int) {
	args := os.Args
	if len(args) < 2 {
		log.Fatal("Error: Please pass in 2 command line args. token and pool count")

	}
	accessToken := os.Args[1]
	if accessToken == "" {
		log.Fatal("Error: Please pass in command line arg 1 with your map token")
	}

	poolSizeString := os.Args[2]
	poolSize, err := strconv.Atoi(poolSizeString)
	if err != nil {
		poolSize = 5
		log.Println("Error parsing command line arg 2 'pool size'", err)
		log.Println("Setting pool size to default value", poolSize)
	}
	if poolSize == 0 {
		poolSize = 5
		log.Println("Setting pool size to default value", poolSize)
	}

	return accessToken, poolSize
}
