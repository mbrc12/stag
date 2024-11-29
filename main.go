package main

import (
	"flag"
	"io"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	Config *Config

	TaskIndex int

	Restarting bool

	Tasks    map[int]*Task
	UIStates map[int]*UIState
}

type UIState struct {
	Scrolling bool
}

func main() {
	var configPath string

	flag.StringVar(&configPath, "c", "cel.toml", "path to config file")

	flag.Parse()

	file, err := os.Open(configPath)
	if err != nil {
		panic("Config file not found!")
	}

	configData, err := io.ReadAll(file)

	if err != nil {
		panic("Failed to read config file!")
	}

	config := new(Config)
	config.Parse(configData)

	model := Model{
		Config:    config,
		TaskIndex: 0,

		Restarting: false,

		Tasks:    make(map[int]*Task),
		UIStates: make(map[int]*UIState),
	}

	_ = model

	for _, taskConfig := range config.WatchTasks {
		task := &Task{
			Prefix:     config.Prefix,
			Commands:   taskConfig.Run.Commands,
			MaxHistory: uint64(config.Store),
		}

		task.Init()
		task.Watch(taskConfig.Files, taskConfig.Exclude)

		model.Tasks[taskConfig.Id] = task

		var taskEvts chan TaskCmd
		task.Start(taskEvts)
	}

	for _, taskConfig := range config.MenuTasks {
		task := &Task{
			Prefix:     config.Prefix,
			Commands:   taskConfig.Run.Commands,
			MaxHistory: uint64(config.Store),
		}

		task.Init()

		model.Tasks[taskConfig.Id] = task
	}

	program := tea.NewProgram(&model)
	if _, err := program.Run(); err != nil {
		panic(err)
	}
}

func (self *Model) Init() tea.Cmd {
	return nil
}

func (self *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return self, tea.Quit
		}
	}

	return self, nil
}

func (self *Model) View() string {
	print(self.Tasks[0].StatusLong)
	return string(self.Tasks[0].StatusLong)
}
