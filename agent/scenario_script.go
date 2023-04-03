package agent

import (
	"encoding/json"
	"io"
	"os"
	"strings"
)

type scenarioTrainer struct {
	agent                    *Agent
	scenarioScriptProcessors map[string]func(string)
}

func (s *scenarioTrainer) train(file string) {
	var data *scenarioScriptData
	jsonFile, err := os.Open(file)
	printErr(err)

	byteValue, err := io.ReadAll(jsonFile)
	printErr(err)
	printErr(json.Unmarshal(byteValue, &data))

	for _, fileName := range data.FileNames {
		s.processFile(fileName)
	}
}

func (s *scenarioTrainer) processFile(fileName string) {
	split := strings.Split(fileName, ".")
	extension := split[len(split)-1]
	if processor, seen := s.scenarioScriptProcessors[extension]; seen {
		processor(fileName)
	}
}

func (s *scenarioTrainer) processSntcFile(fileName string) {
	_, sentences, _ := s.agent.language.trainParser.parse(fileName)
	for _, sentence := range sentences {
		sentence.rootNode.build()
	}

	// todo complete this
}

func (s *scenarioTrainer) processFileAction(fileName string) {

}

func (a *Agent) newScenarioTrainer() {
	a.trainer = &scenarioTrainer{agent: a}
	a.trainer.scenarioScriptProcessors = map[string]func(string){
		scenarioScriptFileTypeSntc:   a.trainer.processSntcFile,
		scenarioScriptFileTypeAction: a.trainer.processFileAction,
	}
}

const (
	scenarioScriptFileTypeSntc   = "sss"
	scenarioScriptFileTypeAction = "ssa"
)

type scenarioScriptData struct {
	FileNames []string `json:"fileNames"`
}
